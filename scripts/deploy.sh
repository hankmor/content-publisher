#!/bin/bash

# 部署脚本，用于将 wechat-publisher 部署到远程服务器
# 用法：
# ./deploy.sh [platform] [action]
# 
# 参数：
#   platform: 目标平台（darwin, linux, windows），默认 linux
#   action: 部署动作（build, deploy, restart, logs, status），默认 deploy
#
# 示例：
#   ./deploy.sh              # 构建并部署到 Linux 服务器
#   ./deploy.sh linux build  # 只构建
#   ./deploy.sh linux deploy # 只部署（不构建）
#   ./deploy.sh linux restart # 重启服务
#   ./deploy.sh linux logs   # 查看日志
#   ./deploy.sh linux status # 查看服务状态

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 加载配置文件
CONFIG_FILE="${SCRIPT_DIR}/deploy.config"
if [ -f "${CONFIG_FILE}" ]; then
    source "${CONFIG_FILE}"
else
    echo "错误：配置文件不存在: ${CONFIG_FILE}"
    echo "请复制配置文件模板并修改:"
    echo "  cp ${SCRIPT_DIR}/deploy.config.example ${SCRIPT_DIR}/deploy.config"
    echo "  然后编辑 ${SCRIPT_DIR}/deploy.config 配置您的服务器信息"
    exit 1
fi

# 检查必要的配置项
if [ -z "${SERVER_HOST}" ] || [ "${SERVER_HOST}" = "your-server.com" ]; then
    echo "错误：请在配置文件中设置 SERVER_HOST"
    exit 1
fi

# 设置默认值
SERVER_USER=${SERVER_USER:-"root"}
SERVER_PORT=${SERVER_PORT:-"22"}
DEPLOY_PATH=${DEPLOY_PATH:-"/opt/wechat-publisher"}
SERVICE_NAME=${SERVICE_NAME:-"wechat-publisher"}
BACKUP_DIR=${BACKUP_DIR:-"/opt/wechat-publisher-backup"}

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查参数
PLATFORM=${1:-"linux"}
ACTION=${2:-"deploy"}

log_info "目标平台: ${PLATFORM}"
log_info "部署动作: ${ACTION}"

# 检查 SSH 连接
check_ssh() {
    log_info "检查 SSH 连接..."
    if ssh -p ${SERVER_PORT} -o ConnectTimeout=5 ${SERVER_USER}@${SERVER_HOST} "echo 'SSH 连接成功'" > /dev/null 2>&1; then
        log_info "SSH 连接成功"
        return 0
    else
        log_error "SSH 连接失败，请检查服务器地址和端口"
        exit 1
    fi
}

# 构建项目
build_project() {
    log_info "开始构建项目..."
    
    # 调用构建脚本
    if [ -f "scripts/build.sh" ]; then
        bash scripts/build.sh ${PLATFORM}
    else
        log_error "构建脚本不存在：scripts/build.sh"
        exit 1
    fi
    
    # 检查构建结果
    if [ ! -f "output/wechat-publisher-${PLATFORM}-amd64" ] && [ ! -f "output/wechat-publisher-${PLATFORM}-amd64.exe" ]; then
        log_error "构建失败，未找到可执行文件"
        exit 1
    fi
    
    log_info "构建完成"
}

# 准备部署文件
prepare_deploy() {
    log_info "准备部署文件..."
    
    # 创建临时部署目录
    rm -rf deploy_temp
    mkdir -p deploy_temp
    
    # 复制可执行文件
    if [ "${PLATFORM}" = "windows" ]; then
        cp output/wechat-publisher-${PLATFORM}-amd64.exe deploy_temp/wechat-publisher.exe
    else
        cp output/wechat-publisher-${PLATFORM}-amd64 deploy_temp/wechat-publisher
    fi
    
    # 复制配置文件
    cp -r configs deploy_temp/
    
    # 复制 demo.html
    cp demo.html deploy_temp/
    
    log_info "部署文件准备完成"
}

# 备份远程服务器上的旧版本
backup_remote() {
    log_info "备份远程服务器上的旧版本..."
    
    ssh -p ${SERVER_PORT} ${SERVER_USER}@${SERVER_HOST} << EOF
        # 创建备份目录
        mkdir -p ${BACKUP_DIR}
        
        # 如果旧版本存在，进行备份
        if [ -d "${DEPLOY_PATH}" ]; then
            BACKUP_NAME="${SERVICE_NAME}-\$(date +%Y%m%d-%H%M%S)"
            echo "备份到: ${BACKUP_DIR}/${BACKUP_NAME}"
            cp -r ${DEPLOY_PATH} ${BACKUP_DIR}/${BACKUP_NAME}
            
            # 只保留最近 5 个备份
            cd ${BACKUP_DIR}
            ls -t | tail -n +6 | xargs -r rm -rf
        fi
EOF
    
    log_info "备份完成"
}

# 上传文件到远程服务器
upload_files() {
    log_info "上传文件到远程服务器..."
    
    # 创建远程目录
    ssh -p ${SERVER_PORT} ${SERVER_USER}@${SERVER_HOST} "mkdir -p ${DEPLOY_PATH}"
    
    # 上传文件
    scp -P ${SERVER_PORT} -r deploy_temp/* ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_PATH}/
    
    log_info "文件上传完成"
}

# 设置远程服务器权限
setup_permissions() {
    log_info "设置文件权限..."
    
    ssh -p ${SERVER_PORT} ${SERVER_USER}@${SERVER_HOST} << EOF
        cd ${DEPLOY_PATH}
        
        # 设置可执行文件权限
        chmod +x wechat-publisher
        
        # 创建必要的目录
        mkdir -p logs data
        
        # 设置目录权限
        chmod 755 logs data
        
        echo "权限设置完成"
EOF
    
    log_info "权限设置完成"
}

# 停止远程服务
stop_service() {
    log_info "停止远程服务..."
    
    ssh -p ${SERVER_PORT} ${SERVER_USER}@${SERVER_HOST} << EOF
        # 查找并停止进程
        PID=\$(pgrep -f "${DEPLOY_PATH}/wechat-publisher")
        if [ -n "\$PID" ]; then
            echo "停止进程: \$PID"
            kill \$PID
            sleep 2
            
            # 如果进程还在运行，强制停止
            if pgrep -f "${DEPLOY_PATH}/wechat-publisher" > /dev/null; then
                echo "强制停止进程"
                pkill -9 -f "${DEPLOY_PATH}/wechat-publisher"
            fi
        else
            echo "服务未运行"
        fi
EOF
    
    log_info "服务已停止"
}

# 启动远程服务
start_service() {
    log_info "启动远程服务..."
    
    ssh -p ${SERVER_PORT} ${SERVER_USER}@${SERVER_HOST} << EOF
        cd ${DEPLOY_PATH}
        
        # 检查是否已经在运行
        if pgrep -f "${DEPLOY_PATH}/wechat-publisher" > /dev/null; then
            echo "服务已经在运行"
            exit 0
        fi
        
        # 启动服务
        nohup ./wechat-publisher > logs/startup.log 2>&1 &
        
        # 等待服务启动
        sleep 3
        
        # 检查服务是否启动成功
        if pgrep -f "${DEPLOY_PATH}/wechat-publisher" > /dev/null; then
            echo "服务启动成功"
            echo "PID: \$(pgrep -f '${DEPLOY_PATH}/wechat-publisher')"
        else
            echo "服务启动失败"
            cat logs/startup.log
            exit 1
        fi
EOF
    
    log_info "服务启动完成"
}

# 查看远程日志
view_logs() {
    log_info "查看远程服务日志..."
    
    ssh -p ${SERVER_PORT} ${SERVER_USER}@${SERVER_HOST} << EOF
        cd ${DEPLOY_PATH}
        
        # 查看今天的日志
        TODAY=\$(date +%Y-%m-%d)
        LOG_FILE="logs/\${TODAY}.log"
        
        if [ -f "\${LOG_FILE}" ]; then
            echo "=== 今天的日志 (\${LOG_FILE}) ==="
            tail -n 50 "\${LOG_FILE}"
        else
            echo "未找到今天的日志文件"
        fi
EOF
}

# 查看服务状态
check_status() {
    log_info "检查服务状态..."
    
    ssh -p ${SERVER_PORT} ${SERVER_USER}@${SERVER_HOST} << EOF
        cd ${DEPLOY_PATH}
        
        # 检查进程
        if pgrep -f "${DEPLOY_PATH}/wechat-publisher" > /dev/null; then
            PID=\$(pgrep -f "${DEPLOY_PATH}/wechat-publisher")
            echo "服务状态: 运行中"
            echo "进程 ID: \$PID"
            
            # 显示内存和 CPU 使用情况
            ps -p \$PID -o pid,pcpu,pmem,etime,cmd --no-headers
        else
            echo "服务状态: 未运行"
        fi
        
        # 检查端口
        if netstat -tuln 2>/dev/null | grep -q ':8080 '; then
            echo "端口 8080: 监听中"
        else
            echo "端口 8080: 未监听"
        fi
EOF
}

# 完整部署流程
full_deploy() {
    log_info "开始完整部署流程..."
    
    # 1. 检查 SSH 连接
    check_ssh
    
    # 2. 构建项目
    build_project
    
    # 3. 准备部署文件
    prepare_deploy
    
    # 4. 备份旧版本
    backup_remote
    
    # 5. 上传文件
    upload_files
    
    # 6. 设置权限
    setup_permissions
    
    # 7. 停止旧服务
    stop_service
    
    # 8. 启动新服务
    start_service
    
    # 9. 检查服务状态
    check_status
    
    # 清理临时文件
    rm -rf deploy_temp
    
    log_info "部署完成！"
}

# 主逻辑
case ${ACTION} in
    "build")
        build_project
        ;;
    "deploy")
        check_ssh
        prepare_deploy
        backup_remote
        upload_files
        setup_permissions
        stop_service
        start_service
        check_status
        rm -rf deploy_temp
        log_info "部署完成！"
        ;;
    "restart")
        check_ssh
        stop_service
        start_service
        check_status
        ;;
    "logs")
        check_ssh
        view_logs
        ;;
    "status")
        check_ssh
        check_status
        ;;
    *)
        full_deploy
        ;;
esac

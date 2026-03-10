#!/bin/bash

# 简化的 Linux 部署脚本
# 用法: ./deploy.sh

set -e

# 加载配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/deploy.config"

if [ ! -f "${CONFIG_FILE}" ]; then
    echo "错误：配置文件不存在，请先创建 ${CONFIG_FILE}"
    exit 1
fi

source "${CONFIG_FILE}"

# 默认值
SERVER_USER=${SERVER_USER:-"root"}
SERVER_PORT=${SERVER_PORT:-"22"}
DEPLOY_PATH=${DEPLOY_PATH:-"/opt/wechat-publisher"}

echo "配置信息:"
echo "  服务器: ${SERVER_USER}@${SERVER_HOST}:${SERVER_PORT}"
echo "  部署路径: ${DEPLOY_PATH}"
echo ""

# 测试 SSH 连接
echo "==> 0. 测试 SSH 连接"
if ! ssh -p ${SERVER_PORT} -o ConnectTimeout=10 ${SERVER_USER}@${SERVER_HOST} "echo '连接成功'"; then
    echo "错误：无法连接到服务器"
    exit 1
fi
echo ""

echo "==> 1. 构建 Linux 版本"
bash scripts/build.sh linux-amd64
echo ""

echo "==> 2. 停止远程服务"
ssh -p ${SERVER_PORT} -o ConnectTimeout=10 -o ServerAliveInterval=5 ${SERVER_USER}@${SERVER_HOST} "pkill -f wechat-publisher 2>/dev/null || true" || true
echo "停止服务完成"
echo ""

echo "==> 3. 上传文件"
ssh -p ${SERVER_PORT} ${SERVER_USER}@${SERVER_HOST} "mkdir -p ${DEPLOY_PATH}"
scp -P ${SERVER_PORT} output/wechat-publisher-linux-amd64 ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_PATH}/wechat-publisher
scp -P ${SERVER_PORT} -r configs demo.html API_GUIDE.md ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_PATH}/
echo ""

echo "==> 4. 启动服务"
ssh -p ${SERVER_PORT} ${SERVER_USER}@${SERVER_HOST} << EOF
cd ${DEPLOY_PATH}
chmod +x wechat-publisher
mkdir -p logs data
nohup ./wechat-publisher > logs/startup.log 2>&1 &
sleep 2
if pgrep -f wechat-publisher > /dev/null; then
    echo "服务启动成功"
else
    echo "服务启动失败"
fi
EOF
echo ""

echo "==> 部署完成"
echo "服务地址: http://${SERVER_HOST}:${SERVICE_PORT}"

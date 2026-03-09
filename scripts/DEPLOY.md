# 部署指南

本文档说明如何将 wechat-publisher 部署到远程服务器。

## 前置要求

1. **本地环境**：
   - Go 1.24.0 或更高版本
   - SSH 客户端（Linux/macOS 自带，Windows 需要安装 PuTTY 或 OpenSSH）

2. **服务器环境**：
   - Linux 服务器（推荐 CentOS 7+、Ubuntu 18.04+）
   - 至少 100MB 可用磁盘空间
   - 8080 端口未被占用

3. **网络要求**：
   - 本地可以 SSH 连接到服务器
   - 服务器防火墙允许 8080 端口访问

## 快速开始

### 1. 配置部署参数

复制部署配置文件：

```bash
cp scripts/deploy.config.example scripts/deploy.config
```

编辑 `scripts/deploy.config`，修改以下配置：

```bash
# 服务器配置
SERVER_USER="root"                    # 修改为你的服务器用户名
SERVER_HOST="your-server.com"           # 修改为你的服务器地址
SERVER_PORT="22"                      # SSH 端口，默认 22

# 部署配置
DEPLOY_PATH="/opt/wechat-publisher"      # 部署路径，根据需要修改
SERVICE_NAME="wechat-publisher"           # 服务名称
BACKUP_DIR="/opt/wechat-publisher-backup" # 备份目录
```

### 2. 部署到服务器

**完整部署（构建 + 部署）**：

```bash
./scripts/deploy.sh linux
```

这个命令会：
1. 构建 Linux 版本的可执行文件
2. 准备部署文件
3. 备份服务器上的旧版本
4. 上传新文件到服务器
5. 设置文件权限
6. 停止旧服务
7. 启动新服务
8. 检查服务状态

### 3. 验证部署

部署完成后，访问以下地址验证：

- **API 服务**：`http://your-server.com:8080`
- **Demo 页面**：`http://your-server.com:8080/demo`

## 高级用法

### 只构建不部署

```bash
./scripts/deploy.sh linux build
```

### 只部署不构建（假设已构建）

```bash
./scripts/deploy.sh linux deploy
```

### 重启服务

```bash
./scripts/deploy.sh linux restart
```

### 查看日志

```bash
./scripts/deploy.sh linux logs
```

### 查看服务状态

```bash
./scripts/deploy.sh linux status
```

## 部署到其他平台

### macOS

```bash
./scripts/deploy.sh darwin
```

### Windows

```bash
./scripts/deploy.sh windows
```

## 手动部署

如果自动部署脚本不适用，可以按照以下步骤手动部署：

### 1. 构建项目

```bash
./scripts/build.sh linux
```

### 2. 上传文件

```bash
scp -P 22 output/wechat-publisher-linux-amd64 root@your-server.com:/opt/wechat-publisher/
scp -P 22 -r configs root@your-server.com:/opt/wechat-publisher/
scp -P 22 demo.html root@your-server.com:/opt/wechat-publisher/
```

### 3. SSH 登录服务器

```bash
ssh root@your-server.com
```

### 4. 设置权限

```bash
cd /opt/wechat-publisher
chmod +x wechat-publisher
mkdir -p logs data
chmod 755 logs data
```

### 5. 启动服务

```bash
nohup ./wechat-publisher > logs/startup.log 2>&1 &
```

### 6. 检查服务状态

```bash
# 检查进程
ps aux | grep wechat-publisher

# 检查端口
netstat -tuln | grep 8080

# 查看日志
tail -f logs/$(date +%Y-%m-%d).log
```

## 常见问题

### 1. SSH 连接失败

**问题**：`SSH 连接失败，请检查服务器地址和端口`

**解决方案**：
- 检查服务器地址是否正确
- 检查 SSH 端口是否正确（默认 22）
- 检查服务器防火墙是否允许 SSH 连接
- 尝试手动 SSH 连接：`ssh root@your-server.com`

### 2. 服务启动失败

**问题**：`服务启动失败`

**解决方案**：
- 查看启动日志：`ssh root@your-server.com "tail -f /opt/wechat-publisher/logs/startup.log"`
- 检查配置文件是否正确
- 检查端口 8080 是否被占用：`ssh root@your-server.com "netstat -tuln | grep 8080"`
- 检查是否有足够的磁盘空间：`ssh root@your-server.com "df -h"`

### 3. 无法访问服务

**问题**：部署成功但无法访问 `http://your-server.com:8080`

**解决方案**：
- 检查服务器防火墙是否允许 8080 端口
- 检查云服务商的安全组设置
- 检查服务是否正在运行：`./scripts/deploy.sh linux status`

### 4. 配置文件修改不生效

**问题**：修改配置文件后重启服务，但配置未生效

**解决方案**：
- 确认修改的是服务器上的配置文件，不是本地的
- 重启服务：`./scripts/deploy.sh linux restart`
- 检查配置文件语法是否正确

## 安全建议

1. **使用非 root 用户**：建议创建专门的用户来运行服务
2. **配置防火墙**：只允许必要的端口访问
3. **使用 SSH 密钥**：建议使用 SSH 密钥认证而不是密码
4. **定期备份**：部署脚本会自动备份，但建议定期检查备份
5. **监控日志**：定期查看日志，及时发现异常

## 生产环境配置

部署到生产环境时，建议修改以下配置：

### 1. 修改服务器模式

编辑服务器上的 `configs/config.yaml`：

```yaml
server:
  port: 8080
  mode: "release"  # 改为 release 模式
```

### 2. 配置 IP 白名单

```yaml
api:
  ip_whitelist:
    - "127.0.0.1"
    - "your_server_ip"  # 添加服务器 IP
```

### 3. 启用频率限制

```yaml
api:
  rate_limit:
    enabled: true
    max_draft_per_day: 10  # 根据需要调整
```

### 4. 配置日志级别

```yaml
log:
  level: "info"  # 生产环境使用 info 级别
  console: false  # 生产环境可以关闭控制台日志
```

## 使用 systemd 管理服务（可选）

如果需要使用 systemd 管理服务，可以创建服务文件：

### 1. 创建服务文件

```bash
sudo vi /etc/systemd/system/wechat-publisher.service
```

内容如下：

```ini
[Unit]
Description=WeChat Publisher Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/wechat-publisher
ExecStart=/opt/wechat-publisher/wechat-publisher
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

### 2. 启用并启动服务

```bash
sudo systemctl daemon-reload
sudo systemctl enable wechat-publisher
sudo systemctl start wechat-publisher
```

### 3. 管理服务

```bash
# 查看状态
sudo systemctl status wechat-publisher

# 重启服务
sudo systemctl restart wechat-publisher

# 停止服务
sudo systemctl stop wechat-publisher

# 查看日志
sudo journalctl -u wechat-publisher -f
```

## 更新部署

当需要更新服务时，只需重新运行部署命令：

```bash
./scripts/deploy.sh linux
```

部署脚本会自动：
- 备份当前版本
- 上传新版本
- 重启服务

如果新版本有问题，可以快速回滚：

```bash
ssh root@your-server.com "cd /opt/wechat-publisher-backup && ls -t | head -1"
```

然后手动恢复备份版本。

# 多平台内容发布网关

一个统一的多平台内容发布网关，支持微信公众号、小红书、微博等多个平台的内容发布。目前**已支持微信公众号**，其他平台将陆续接入。

## 功能特性

- 🚀 **统一 API 接口**：一套 API 接口，发布到多个平台
- 📱 **微信公众号**：支持图文消息、图片消息等多种内容类型
- 🔒 **安全认证**：IP 白名单 + API Key 双重验证
- 📊 **日志记录**：完整的操作日志，便于追踪和审计
- ⚙️ **灵活配置**：支持 YAML 配置文件，易于维护

## 已支持平台

| 平台 | 状态 | 支持功能 |
|------|------|----------|
| 微信公众号 | ✅ 已支持 | 素材上传、草稿创建、内容发布 |
| 小红书 | 🚧 规划中 | - |
| 微博 | 🚧 规划中 | - |
| 抖音 | 🚧 规划中 | - |
| B站 | 🚧 规划中 | - |

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/hankmor/wechat-publisher.git
cd wechat-publisher
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境

复制配置文件模板：

```bash
cp configs/config.yaml.example configs/config.yaml
```

编辑 `configs/config.yaml`，填写你的配置信息：

```yaml
# 微信公众号配置
wechat:
  app_id: "your_app_id"
  app_secret: "your_app_secret"
  token: ""

# 服务器配置
server:
  port: 8080
  mode: "debug"

# API 访问控制配置
api:
  ip_whitelist:
    - "127.0.0.1"
  api_keys:
    - "your_api_key"

# 日志配置
log:
  level: "info"
  path: "./logs"
  console: true
```

### 4. 运行服务

```bash
go run ./cmd/publisher
```

服务将在 `http://localhost:8080` 启动。

## API 接口文档

### 微信公众号接口

#### 1. 上传素材库图片

**接口**: `POST /api/upload/material`

**请求方式**: 表单提交，字段名为 `image`，值为图片文件

**响应示例**:

```json
{
  "media_id": "media_id"
}
```

#### 2. 上传图文消息图片

**接口**: `POST /api/upload/news-image`

**请求方式**: 表单提交，字段名为 `image`，值为图片文件

**响应示例**:

```json
{
  "url": "https://mmbiz.qpic.cn/mmbiz_jpg/..."
}
```

#### 3. 创建草稿

**接口**: `POST /api/draft/create`

**请求体示例**:

```json
{
  "type": "news",
  "article": {
    "title": "文章标题",
    "author": "作者",
    "digest": "文章摘要",
    "content": "<p>文章内容</p>",
    "thumb_media_id": "图片素材的 media_id",
    "need_open_comment": 0,
    "only_fans_can_comment": 0
  }
}
```

**响应示例**:

```json
{
  "media_id": "draft_media_id"
}
```

#### 4. 发布草稿

**接口**: `POST /api/draft/publish`

**请求体示例**:

```json
{
  "media_id": "draft_media_id"
}
```

**响应示例**:

```json
{
  "publish_id": 123456789
}
```

## 示例代码

### 上传素材库图片

```bash
curl -X POST http://localhost:8080/api/upload/material \
  -H "X-API-Key: your_api_key" \
  -F "image=@/path/to/image.jpg"
```

### 创建并发布图文消息

```bash
# 1. 上传封面图片
curl -X POST http://localhost:8080/api/upload/material \
  -H "X-API-Key: your_api_key" \
  -F "image=@/path/to/cover.jpg"

# 2. 创建草稿
curl -X POST http://localhost:8080/api/draft/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{
    "type": "news",
    "article": {
      "title": "测试文章",
      "author": "测试作者",
      "digest": "这是一篇测试文章",
      "content": "<p>测试文章内容</p>",
      "thumb_media_id": "your_media_id",
      "need_open_comment": 0,
      "only_fans_can_comment": 0
    }
  }'

# 3. 发布草稿
curl -X POST http://localhost:8080/api/draft/publish \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{
    "media_id": "your_draft_media_id"
  }'
```

## 项目结构

```
wechat-publisher/
├── cmd/
│   └── publisher/          # 应用入口
│       └── main.go
├── internal/
│   ├── api/                # API 接口层
│   ├── config/             # 配置管理
│   ├── log/                # 日志模块
│   └── service/            # 业务逻辑层
├── configs/                # 配置文件
│   ├── config.yaml         # 主配置文件（不提交到 Git）
│   └── config.yaml.example # 配置文件模板
├── logs/                   # 日志文件目录
├── scripts/                # 构建脚本
│   └── build.sh
├── API_GUIDE.md           # API 调用文档
├── demo.html              # 测试页面
└── README.md              # 项目说明
```

## 构建部署

### 构建所有平台

```bash
./scripts/build.sh
```

### 构建指定平台

```bash
./scripts/build.sh darwin   # macOS
./scripts/build.sh linux    # Linux
./scripts/build.sh windows  # Windows
```

## 开发计划

- [x] 微信公众号支持
- [ ] 小红书平台支持（调研中）
- [ ] 微博平台支持
- [ ] 抖音平台支持
- [ ] B站平台支持
- [ ] 统一内容格式转换
- [ ] 批量发布功能
- [ ] 发布状态回调

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

# 微信发布工具

这是一个基于 Go 语言的微信公众号发布工具，提供以下功能：

1. 通过 API 上传图片到素材库
2. 通过 API 上传文章和图片到草稿箱，支持微信贴图、文章等产品
3. 提供最终的 API 接口

## 安装

### 1. 克隆项目

```bash
git clone https://github.com/hankmor/wechat-publisher.git
cd wechat-publisher
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境变量

在运行前，需要设置以下环境变量：

```bash
export WECHAT_APPID=your_appid
export WECHAT_APPSECRET=your_appsecret
export WECHAT_TOKEN=your_token
```

## 运行

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 接口

### 1. 上传素材库图片

**接口**: `POST /api/upload/material`

**请求方式**: 表单提交，字段名为 `image`，值为图片文件

**响应示例**:

```json
{
  "media_id": "media_id"
}
```

### 2. 上传图文消息图片

**接口**: `POST /api/upload/news-image`

**请求方式**: 表单提交，字段名为 `image`，值为图片文件

**响应示例**:

```json
{
  "url": "https://mmbiz.qpic.cn/mmbiz_jpg/..."
}
```

### 3. 创建草稿

**接口**: `POST /api/draft/create`

**请求体示例**:

```json
{
  "title": "文章标题",
  "author": "作者",
  "digest": "文章摘要",
  "content": "文章内容",
  "thumb_media_id": "图片素材的 media_id",
  "need_open_comment": 0,
  "only_fans_can_comment": 0
}
```

**响应示例**:

```json
{
  "media_id": "draft_media_id"
}
```

### 4. 发布草稿

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

### 上传素材库图片示例

```bash
curl -X POST http://localhost:8080/api/upload/material -F "image=@/path/to/image.jpg"
```

### 上传图文消息图片示例

```bash
curl -X POST http://localhost:8080/api/upload/news-image -F "image=@/path/to/image.jpg"
```

### 创建草稿示例

```bash
curl -X POST http://localhost:8080/api/draft/create -H "Content-Type: application/json" -d '{
  "title": "测试文章",
  "author": "测试作者",
  "digest": "这是一篇测试文章",
  "content": "<p>测试文章内容</p>",
  "thumb_media_id": "your_media_id",
  "need_open_comment": 0,
  "only_fans_can_comment": 0
}'
```

### 发布草稿示例

```bash
curl -X POST http://localhost:8080/api/draft/publish -H "Content-Type: application/json" -d '{
  "media_id": "your_draft_media_id"
}'
```

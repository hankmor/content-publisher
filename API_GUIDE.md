# 微信发布工具 API 调用文档

## 1. API 访问控制

所有 API 接口都需要进行访问控制，具体要求如下：

- **IP 白名单**：客户端 IP 必须在配置文件的 `api.ip_whitelist` 中
- **API Key**：客户端必须提供有效的 API Key，API Key 在传输过程中会被自动哈希处理，确保安全传输。可以通过以下方式之一传递：
  - 请求头：`X-API-Key: your_api_key`
  - 查询参数：`?api_key=your_api_key`

**注意**：API Key 在服务器端会使用 SHA256 算法进行哈希处理后再与配置文件中的 API Key 进行比较，确保传输过程中的安全性。

## 2. API 接口列表

### 2.1 上传素材图片（用于图文消息的封面图）

- **接口地址**：`POST /api/upload/material`
- **请求方式**：POST
- **请求参数**：
  - `image`：文件（form-data）
- **响应格式**：
  - 成功：`{"media_id": "素材图片的media_id"}`
  - 失败：`{"error": "错误信息"}`

### 2.2 上传图文消息图片（用于图文消息内容中的图片）

- **接口地址**：`POST /api/upload/news-image`
- **请求方式**：POST
- **请求参数**：
  - `image`：文件（form-data）
- **响应格式**：
  - 成功：`{"url": "图片的URL地址"}`
  - 失败：`{"error": "错误信息"}`

### 2.3 创建草稿

- **接口地址**：`POST /api/draft/create`
- **请求方式**：POST
- **请求参数**：
  - `type`：草稿类型，可选值：`news`（图文消息，默认）、`newspic`（图片消息）
  - `article`：图文消息内容（当 type 为 news 时必填）
    - `title`：标题
    - `author`：作者
    - `digest`：摘要
    - `content`：内容（支持 HTML，图片需要包装到 `<img>` 标签中，使用 UploadNewsImage 返回的 URL）
    - `content_source_url`：原文链接
    - `thumb_media_id`：封面图片的 media_id（通过 UploadMaterialImage 获取）
    - `need_open_comment`：是否打开评论（0：否，1：是）
    - `only_fans_can_comment`：是否仅粉丝可评论（0：否，1：是）
  - `image_article`：图片消息内容（当 type 为 newspic 时必填）
    - `title`：标题
    - `author`：作者
    - `content`：内容
    - `media_id_list`：图片的 media_id 列表（通过 UploadMaterialImage 获取，最多 20 张）
- **响应格式**：
  - 成功：`{"media_id": "草稿的media_id"}`
  - 失败：`{"error": "错误信息"}`

### 2.4 发布草稿

- **接口地址**：`POST /api/draft/publish`
- **请求方式**：POST
- **请求参数**：
  - `media_id`：草稿的 media_id（通过 CreateDraft 获取）
- **响应格式**：
  - 成功：`{"publish_id": "发布任务的ID"}`
  - 失败：`{"error": "错误信息"}`

## 3. API 调用顺序

### 3.1 发布图文消息的流程

1. **上传封面图片**：调用 `UploadMaterialImage` 接口上传封面图片，获取 `media_id`
2. **上传内容图片**：如果图文消息内容中需要图片，调用 `UploadNewsImage` 接口上传图片，获取 `url`
3. **创建图文草稿**：调用 `CreateDraft` 接口，设置 `type` 为 `news`，并提供图文消息内容（包括封面图片的 `media_id` 和内容图片的 `url`，注意图片需要包装到 `<img>` 标签中）
4. **发布草稿**：调用 `PublishDraft` 接口，使用创建草稿时返回的 `media_id`

### 3.2 发布图片消息的流程

1. **上传图片**：调用 `UploadMaterialImage` 接口上传图片，获取 `media_id`（最多 20 张）
2. **创建图片草稿**：调用 `CreateDraft` 接口，设置 `type` 为 `newspic`，并提供图片消息内容（包括图片的 `media_id` 列表）
3. **发布草稿**：调用 `PublishDraft` 接口，使用创建草稿时返回的 `media_id`

## 4. 请求和响应示例

### 4.1 上传封面图片

**请求**：
```bash
curl -X POST http://localhost:8082/api/upload/material \
  -H "X-API-Key: your_api_key" \
  -F "image=@/path/to/image.jpg"
```

**响应**：
```json
{
  "media_id": "MEDIA_ID"
}
```

### 4.2 上传图文消息图片

**请求**：
```bash
curl -X POST http://localhost:8082/api/upload/news-image \
  -H "X-API-Key: your_api_key" \
  -F "image=@/path/to/image.jpg"
```

**响应**：
```json
{
  "url": "http://mmbiz.qpic.cn/mmbiz_jpg/..."
}
```

### 4.3 创建图文草稿（包含图片）

**请求**：
```bash
curl -X POST http://localhost:8082/api/draft/create \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "news",
    "article": {
      "title": "测试标题",
      "author": "测试作者",
      "digest": "测试摘要",
      "content": "<p>测试内容</p><img src=\"http://mmbiz.qpic.cn/mmbiz_jpg/...\" /><p>更多内容</p>",
      "content_source_url": "https://example.com",
      "thumb_media_id": "MEDIA_ID",
      "need_open_comment": 1,
      "only_fans_can_comment": 0
    }
  }'
```

**响应**：
```json
{
  "media_id": "DRAFT_MEDIA_ID"
}
```

### 4.4 创建图片草稿

**请求**：
```bash
curl -X POST http://localhost:8082/api/draft/create \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "newspic",
    "image_article": {
      "title": "测试图片消息",
      "author": "测试作者",
      "content": "测试内容",
      "media_id_list": ["MEDIA_ID1", "MEDIA_ID2"]
    }
  }'
```

**响应**：
```json
{
  "media_id": "DRAFT_MEDIA_ID"
}
```

### 4.5 发布草稿

**请求**：
```bash
curl -X POST http://localhost:8082/api/draft/publish \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "media_id": "DRAFT_MEDIA_ID"
  }'
```

**响应**：
```json
{
  "publish_id": "PUBLISH_ID"
}
```

## 5. 常见错误处理

| 错误信息 | 可能原因 | 解决方案 |
|---------|---------|--------|
| `IP 不在白名单内` | 客户端 IP 没有添加到配置文件的 `api.ip_whitelist` 中 | 将客户端 IP 添加到配置文件的白名单中 |
| `无效的 API Key` | 提供的 API Key 不在配置文件的 `api.api_keys` 中 | 使用配置文件中有效的 API Key |
| `invalid media_id` | 提供的 media_id 无效或已过期 | 重新上传图片获取有效的 media_id |
| `不支持的草稿类型` | 提供的 type 值不是 `news` 或 `newspic` | 使用正确的草稿类型 |
| `bind error` | 请求参数格式错误 | 检查请求参数是否符合要求 |

## 6. 注意事项

1. **图片格式**：支持 JPG、PNG 等常见图片格式
2. **图片大小**：单个图片大小不超过 2M
3. **图片数量**：图片消息最多支持 20 张图片
4. **内容格式**：图文消息内容支持 HTML 格式，图片需要包装到 `<img>` 标签中
5. **内容长度**：图文消息内容长度有限制，建议控制在合理范围内
6. **API 频率**：请遵守微信公众平台的 API 调用频率限制
7. **错误处理**：所有 API 接口都可能返回错误，请在调用时做好错误处理

## 7. 完整调用示例

### 7.1 发布图文消息示例（包含图片）

```bash
# 1. 上传封面图片
COVER_MEDIA_ID=$(curl -X POST http://localhost:8082/api/upload/material \
  -H "X-API-Key: your_api_key" \
  -F "image=@/path/to/cover.jpg" | jq -r '.media_id')

# 2. 上传内容图片
CONTENT_URL=$(curl -X POST http://localhost:8082/api/upload/news-image \
  -H "X-API-Key: your_api_key" \
  -F "image=@/path/to/content.jpg" | jq -r '.url')

# 3. 创建图文草稿（图片需要包装到 img 标签中）
DRAFT_MEDIA_ID=$(curl -X POST http://localhost:8082/api/draft/create \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{  
    "type": "news",
    "article": {
      "title": "测试图文消息",
      "author": "测试作者",
      "digest": "这是一条测试图文消息",
      "content": "<p>测试内容</p><img src=\"'$CONTENT_URL'\" /><p>更多内容</p>",
      "content_source_url": "https://example.com",
      "thumb_media_id": "'$COVER_MEDIA_ID'",
      "need_open_comment": 1,
      "only_fans_can_comment": 0
    }
  }' | jq -r '.media_id')

# 4. 发布草稿
PUBLISH_ID=$(curl -X POST http://localhost:8082/api/draft/publish \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "media_id": "'$DRAFT_MEDIA_ID'"
  }' | jq -r '.publish_id')

echo "发布成功，发布ID：$PUBLISH_ID"
```

### 7.2 发布图片消息示例

```bash
# 1. 上传图片（最多 20 张）
MEDIA_ID1=$(curl -X POST http://localhost:8082/api/upload/material \
  -H "X-API-Key: your_api_key" \
  -F "image=@/path/to/image1.jpg" | jq -r '.media_id')

MEDIA_ID2=$(curl -X POST http://localhost:8082/api/upload/material \
  -H "X-API-Key: your_api_key" \
  -F "image=@/path/to/image2.jpg" | jq -r '.media_id')

# 2. 创建图片草稿
DRAFT_MEDIA_ID=$(curl -X POST http://localhost:8082/api/draft/create \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "newspic",
    "image_article": {
      "title": "测试图片消息",
      "author": "测试作者",
      "content": "这是一条测试图片消息",
      "media_id_list": ["'$MEDIA_ID1'", "'$MEDIA_ID2'"]
    }
  }' | jq -r '.media_id')

# 3. 发布草稿
PUBLISH_ID=$(curl -X POST http://localhost:8082/api/draft/publish \
  -H "X-API-Key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "media_id": "'$DRAFT_MEDIA_ID'"
  }' | jq -r '.publish_id')

echo "发布成功，发布ID：$PUBLISH_ID"
```
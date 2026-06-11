# 后端接口联调说明

当前后端使用 `Gin + GORM + SQLite`，所有数据都会持久化到数据库文件中，而不是只存在内存里。首次启动会自动建表和写入默认数据。

统一返回格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

默认服务地址：

```text
http://localhost:8080
```

## 已持久化的数据模块

- 模型 `models`
- 上游渠道 `channels`
- 模型能力映射 `abilities`
- 访问密钥 `api_keys`
- 系统配置 `system_options`
- 会话 `chat_sessions`
- 消息 `chat_messages`
- 用量日志 `usage_logs`

## 1. 健康检查

```http
GET /health
```

## 2. 模型管理 CRUD

```http
GET    /api/models
GET    /api/models/:id
POST   /api/models
PUT    /api/models/:id
DELETE /api/models/:id
```

模型已预留字段：

- 基础信息：`name`、`display_name`、`provider`、`group_name`
- 能力开关：`enabled`、`supports_vision`、`supports_tools`、`supports_stream`
- 配额能力：`context_window`、`max_output_tokens`
- 计费预留：`input_price`、`output_price`、`request_price`

## 3. 上游渠道 CRUD

```http
GET    /api/channels
GET    /api/channels/:id
POST   /api/channels
PUT    /api/channels/:id
DELETE /api/channels/:id
```

这一层参考 `new-api` 的思路，为真实 API 接入预留：

- `provider_type`
- `base_url`
- `api_key`
- `organization`
- `group_name`
- `model_names`
- `model_mapping`
- `extra_headers`
- `request_template`
- `response_template`
- `priority`
- `weight`
- `rate_limited`
- `max_requests_minute`
- `test_model`

## 4. 模型能力映射 CRUD

```http
GET    /api/abilities
POST   /api/abilities
PUT    /api/abilities/:id
DELETE /api/abilities/:id
```

作用：

- 指定某个分组下，某个模型由哪个渠道提供
- 为后续多渠道路由、优先级、权重调度预留结构

## 5. 密钥管理 CRUD

```http
GET    /api/keys
GET    /api/keys/:id
POST   /api/keys
PUT    /api/keys/:id
DELETE /api/keys/:id
```

查询参数：

- `enabled=true`
- `group_name=codex专用分组`
- `model_id=1`
- `keyword=测试`

说明：

- 密钥会自动生成 `sk-` 前缀
- 生成逻辑改为安全随机字符方案，便于后续真实接口使用
- `model_names` 支持多模型绑定
- `enabled` 使用布尔值控制启用状态
- 已预留完整计费相关字段

## 6. 系统配置 CRUD

```http
GET    /api/options
POST   /api/options
PUT    /api/options/:id
DELETE /api/options/:id
```

## 7. 会话与消息 CRUD

```http
GET    /api/sessions
GET    /api/sessions/:id
POST   /api/sessions
PUT    /api/sessions/:id
DELETE /api/sessions/:id
```

```http
GET    /api/sessions/:session_id/messages
POST   /api/sessions/:session_id/messages
DELETE /api/messages/:id
```

已预留真实聊天接入字段：

- 会话：`model_name`、`channel_id`、`temperature`、`top_p`、`max_tokens`
- 消息：`input_tokens`、`output_tokens`、`finish_reason`、`response_time_ms`

## 8. 用量日志 CRUD

```http
GET    /api/usage-logs
POST   /api/usage-logs
DELETE /api/usage-logs/:id
```

已预留字段：

- `api_key_id`
- `channel_id`
- `model_name`
- `request_id`
- `prompt_tokens`
- `completion_tokens`
- `total_tokens`
- `input_cost`
- `output_cost`
- `request_cost`
- `total_cost`
- `status`
- `error_message`
- `latency_ms`
- `client_ip`
- `user_agent`

## 9. 当前默认初始化内容

后端首次启动后会自动写入：

- 默认管理员 `admin`
- 默认模型 `gpt-5.5`
- 默认模型分组 `codex专用分组`
- 基础系统配置项

## 10. 当前说明

当前后端已经具备“持续化存储 + 增删改查”的基础能力，但还没有真正调用外部 LLM。

下一步接真实 API 时，建议优先接这条链路：

1. 渠道 `channels`
2. 能力映射 `abilities`
3. 会话与消息
4. 用量日志与计费

## 11. DeepSeek 网关准备

当前项目已经增加一条最小可用的 DeepSeek 网关链路：

- 平台密钥鉴权：`Authorization: Bearer sk-...`
- 模型列表：`GET /v1/models`
- 对话转发：`POST /v1/chat/completions`

你需要准备并填写的内容：

- `.env` 里的 `DEEPSEEK_API_KEY`
- `.env` 里的 `DEEPSEEK_BASE_URL`
- 数据库中渠道 `channels` 的 DeepSeek 配置
- 数据库中能力映射 `abilities` 的模型到渠道关系

当前保留为空、等真实接入时填写的字段：

- `channels.api_key`
- `channels.model_mapping`
- `channels.extra_headers`
- `channels.request_template`
- `channels.response_template`
- `api_keys.billing_config`

说明：

- `/v1/chat/completions` 当前只做非流式转发
- `stream=true` 先返回预留提示
- 调用成功后会写入 `usage_logs`
- 会按 `models` 表中的价格字段计算成本并扣减 `api_keys` 额度
## 12. Channel Test / 渠道测试

```http
POST /api/channels/:id/test
```

用途：
- 快速验证某个渠道的 `base_url`、`api_key`、`test_model` 是否可用
- 在不经过客户端 API Key 网关的情况下，直接测试上游渠道

可选请求体：

```json
{
  "model": "deepseek-chat",
  "messages": ["hello from channel test"],
  "max_tokens": 64
}
```

说明：
- 当前仅支持 `deepseek` 类型渠道
- 如果不传 `model`，会优先使用渠道自己的 `test_model`
- 如果不传 `messages`，会发送默认探测消息

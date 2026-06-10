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

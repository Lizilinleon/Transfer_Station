# AI 智能对话平台 PRD

## 1. 当前目标

当前阶段的重点是先把平台后台做成“可持续化存储”的结构，而不是只完成静态页面。

本阶段目标：

- 建立稳定的 SQLite 持久化数据库
- 建立参考 `new-api` 思路的后台数据结构
- 完成核心模块 CRUD
- 为后续真实 LLM API 接入预留足够字段

## 2. 当前已完成范围

### 前端

- 令牌管理页静态 UI
- 创建、搜索、编辑、删除、启用/禁用的前端静态交互
- 可用模型按钮式选择

### 后端

- SQLite 自动建库建表
- 模型 CRUD
- 渠道 CRUD
- 模型能力映射 CRUD
- 密钥 CRUD
- 系统配置 CRUD
- 会话 CRUD
- 消息 CRUD
- 用量日志 CRUD

## 3. 当前数据库设计

### 3.1 users

用于后续用户登录、权限和分组管理。

### 3.2 models

保存模型元数据与计费基础信息。

预留字段包括：

- 模型名称、展示名称、提供商
- 分组
- 是否支持视觉、工具调用、流式输出
- 上下文窗口
- 最大输出
- 输入/输出/请求价格

### 3.3 channels

参考 `new-api` 的渠道层思路，用于真实 API 上游接入。

预留字段包括：

- `provider_type`
- `base_url`
- `api_key`
- `organization`
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

### 3.4 abilities

参考 `new-api` 的能力映射层，用于记录：

- 哪个分组
- 哪个模型
- 使用哪个渠道
- 对应优先级和权重

### 3.5 api_keys

保存平台生成的访问密钥，并预留真实业务变量。

预留字段包括：

- 密钥本身
- 分组
- 单模型绑定
- 多模型名称绑定
- 启用状态
- 总额度、已用额度、剩余额度
- 无限额度开关
- 计费规则与配置
- 输入/输出/请求单价
- 最后使用时间
- 备注与过期时间

### 3.6 system_options

保存站点与系统动态配置。

### 3.7 chat_sessions

保存聊天会话，已预留真实调用参数：

- `model_name`
- `channel_id`
- `temperature`
- `top_p`
- `max_tokens`

### 3.8 chat_messages

保存消息记录，已预留：

- 输入 token
- 输出 token
- 完成原因
- 响应耗时

### 3.9 usage_logs

保存真实 API 调用后的用量日志和计费结果。

预留字段包括：

- 关联密钥
- 关联渠道
- 模型名称
- 请求 ID
- prompt/completion/total tokens
- 输入/输出/请求/总成本
- 状态、错误信息
- 延迟、客户端 IP、User-Agent

## 4. 当前接口范围

### 健康检查

- `GET /health`

### 模型管理

- `GET /api/models`
- `GET /api/models/:id`
- `POST /api/models`
- `PUT /api/models/:id`
- `DELETE /api/models/:id`

### 渠道管理

- `GET /api/channels`
- `GET /api/channels/:id`
- `POST /api/channels`
- `PUT /api/channels/:id`
- `DELETE /api/channels/:id`

### 能力映射

- `GET /api/abilities`
- `POST /api/abilities`
- `PUT /api/abilities/:id`
- `DELETE /api/abilities/:id`

### 密钥管理

- `GET /api/keys`
- `GET /api/keys/:id`
- `POST /api/keys`
- `PUT /api/keys/:id`
- `DELETE /api/keys/:id`

### 系统配置

- `GET /api/options`
- `POST /api/options`
- `PUT /api/options/:id`
- `DELETE /api/options/:id`

### 会话管理

- `GET /api/sessions`
- `GET /api/sessions/:id`
- `POST /api/sessions`
- `PUT /api/sessions/:id`
- `DELETE /api/sessions/:id`

### 消息管理

- `GET /api/sessions/:session_id/messages`
- `POST /api/sessions/:session_id/messages`
- `DELETE /api/messages/:id`

### 用量日志

- `GET /api/usage-logs`
- `POST /api/usage-logs`
- `DELETE /api/usage-logs/:id`

## 5. 验收标准

- SQLite 数据可重复启动并持续保留
- 所有核心业务表能自动建表
- 核心模块具备增删改查能力
- 结构能支撑未来真实 LLM API 接入
- 文档与当前实现保持一致

## 6. 下一步建议

最适合继续推进的顺序：

1. 前端对接 `/api/keys`
2. 前端对接 `/api/models`
3. 增加 `/api/channels` 和 `/api/abilities` 的管理页
4. 接入真实对话转发逻辑
5. 在调用后写入 `usage_logs`

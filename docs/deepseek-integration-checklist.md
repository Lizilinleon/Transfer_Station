# DeepSeek 接入检查清单

本清单参考 New API 的渠道、令牌、用量和计费结构，结合当前项目的第一阶段 SQLite 设计整理。

## 当前已具备

- `provider_channels`：保存上游渠道、Base URL、API Key、模型列表、优先级和权重。
- `model_abilities`：保存 group + model 到 channel 的显式路由规则。
- `api_keys`：保存客户端访问密钥、分组、模型白名单、额度和价格字段。
- `usage_logs`：保存请求模型、token、费用、状态、延迟、客户端信息。
- `/v1/models`：按 API Key 分组和模型白名单返回可用模型。
- `/v1/chat/completions`：非流式转发到 DeepSeek，并写入用量日志。

## 本次补齐

- 兼容更多 OpenAI chat completion 字段：
  - `max_completion_tokens`
  - `presence_penalty`
  - `frequency_penalty`
  - `stop`
  - `response_format`
  - `tools`
  - `tool_choice`
  - `seed`
  - `user`
  - `logprobs`
  - `top_logprobs`
- 支持渠道 `model_mapping`，例如：

```json
{
  "deepseek-chat": "deepseek-reasoner"
}
```

- 支持渠道 `extra_headers` JSON，用于后续扩展供应商自定义请求头。
- `provider_channels` 增加：
  - `used_quota`
  - `last_used_at`
- 成功日志补齐：
  - `request_id`
  - `input_cost`
  - `output_cost`
  - `request_cost`
  - `total_cost`
- 额度扣减改为事务内更新，减少并发请求导致的额度竞争问题。

## 距离正式接入 DeepSeek 还需要

1. 在 `.env` 中配置真实上游密钥：

```env
DEEPSEEK_API_KEY=sk-xxx
```

2. 启动后端，让 seed 写入默认 DeepSeek channel / model / ability。

3. 在前端创建客户端 API Key，确认：
   - `group_name` 为 `codex-group`
   - `model_names` 包含 `deepseek-chat`
   - `enabled` 为 `true`
   - 额度足够，或设置为无限额度

4. 用客户端 API Key 调用：

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer sk-your-client-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "deepseek-chat",
    "messages": [
      { "role": "user", "content": "hello" }
    ]
  }'
```

## 后续建议按小步开发

1. 流式响应：实现 `stream: true` 的 SSE 转发和流式用量统计。
2. 渠道测试：增加 `/api/channels/:id/test`，测试 Base URL、Key、模型是否可用。
3. 失败重试：按优先级和权重选择下一个可用渠道。
4. 自动禁用：连续上游鉴权失败或限流后自动暂停渠道。
5. 计费规则：统一使用模型价格或 API Key 价格，并明确单位。
6. 日志分页：给 `/api/usage-logs` 增加分页和时间范围过滤。
7. 前端页面：继续接入模型服务、渠道管理、能力映射、用量日志。

# Webhook 分发服务

纯 Go 标准库实现的 Webhook 分发平台，零第三方依赖。

## 运行

```bash
cd origin
go build ./...
go run ./cmd/server
```

环境变量：
- `PORT`：监听端口（默认 8080）
- `ADDR`：监听地址（覆盖 PORT）
- `MAX_PAGE_SIZE`：最大分页大小（默认 100）
- `API_KEY`：API 鉴权密钥（空则不鉴权）
- `LOG_LEVEL`：日志级别 debug/info/warn/error

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/endpoints | 创建端点 |
| GET | /api/endpoints | 端点列表 |
| GET | /api/endpoints/{id} | 获取端点 |
| PUT | /api/endpoints/{id} | 更新端点 |
| DELETE | /api/endpoints/{id} | 删除端点 |
| POST | /api/subscriptions | 创建订阅 |
| GET | /api/subscriptions | 订阅列表 |
| GET | /api/subscriptions/{id} | 获取订阅 |
| PUT | /api/subscriptions/{id} | 更新订阅 |
| DELETE | /api/subscriptions/{id} | 删除订阅 |
| POST | /api/event-types | 创建事件类型 |
| GET | /api/event-types | 事件类型列表 |
| GET | /api/event-types/{id} | 获取事件类型 |
| PUT | /api/event-types/{id} | 更新事件类型 |
| DELETE | /api/event-types/{id} | 删除事件类型 |
| POST | /api/events | 创建事件（自动分发） |
| GET | /api/events | 事件列表 |
| GET | /api/events/{id} | 获取事件 |
| PUT | /api/events/{id} | 更新事件 |
| DELETE | /api/events/{id} | 删除事件 |
| POST | /api/deliveries | 创建投递记录 |
| GET | /api/deliveries | 投递记录列表 |
| GET | /api/deliveries/{id} | 获取投递记录 |
| PUT | /api/deliveries/{id} | 更新投递记录 |
| DELETE | /api/deliveries/{id} | 删除投递记录 |
| POST | /api/deliveries/{id}/retry | 手动重试投递 |
| POST | /api/delivery-attempts | 创建投递尝试 |
| GET | /api/delivery-attempts | 投递尝试列表 |
| GET | /api/delivery-attempts/{id} | 获取投递尝试 |
| PUT | /api/delivery-attempts/{id} | 更新投递尝试 |
| DELETE | /api/delivery-attempts/{id} | 删除投递尝试 |
| POST | /api/retry-policies | 创建重试策略 |
| GET | /api/retry-policies | 重试策略列表 |
| GET | /api/retry-policies/{id} | 获取重试策略 |
| PUT | /api/retry-policies/{id} | 更新重试策略 |
| DELETE | /api/retry-policies/{id} | 删除重试策略 |
| POST | /api/signing-keys | 创建签名密钥 |
| GET | /api/signing-keys | 签名密钥列表 |
| GET | /api/signing-keys/{id} | 获取签名密钥 |
| PUT | /api/signing-keys/{id} | 更新签名密钥 |
| DELETE | /api/signing-keys/{id} | 删除签名密钥 |
| POST | /api/signing-keys/verify | HMAC 签名验证 |
| POST | /api/filters | 创建过滤器 |
| GET | /api/filters | 过滤器列表 |
| GET | /api/filters/{id} | 获取过滤器 |
| PUT | /api/filters/{id} | 更新过滤器 |
| DELETE | /api/filters/{id} | 删除过滤器 |
| POST | /api/filters/{id}/match | 测试过滤器匹配 |
| POST | /api/audit-logs | 创建审计日志 |
| GET | /api/audit-logs | 审计日志列表 |
| GET | /api/audit-logs/{id} | 获取审计日志 |
| GET | /api/stats/overview | 投递总览统计 |
| GET | /api/stats/endpoints | 端点性能统计 |
| GET | /api/stats/event-types | 事件类型分布 |
| GET | /api/stats/timeline | 投递时间线 |
| GET | /api/stats/latency | 延迟统计 |
| POST | /api/batch/events | 批量创建事件 |
| POST | /api/batch/endpoints/disable | 批量停用端点 |
| GET | /api/export/snapshot | 全量快照导出 |
| GET | / | 前端页面 |

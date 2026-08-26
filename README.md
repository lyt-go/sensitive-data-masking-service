# Data Masking Service

纯 Go 标准库实现的数据脱敏服务，零第三方依赖。

## 运行

```bash
go run ./cmd/server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/data-classes | 创建数据分类 |
| GET | /api/data-classes | 列表（分页 + level/keyword 筛选） |
| GET | /api/data-classes/{id} | 获取单个 |
| PUT | /api/data-classes/{id} | 更新 |
| DELETE | /api/data-classes/{id} | 删除 |
| POST | /api/field-configs | 创建字段配置 |
| GET | /api/field-configs | 列表（分页 + data_class_id/mask_type/enabled/keyword 筛选） |
| GET | /api/field-configs/{id} | 获取单个 |
| PUT | /api/field-configs/{id} | 更新 |
| DELETE | /api/field-configs/{id} | 删除 |
| POST | /api/mask-rules | 创建脱敏规则 |
| GET | /api/mask-rules | 列表（分页 + mask_type/enabled/keyword 筛选） |
| GET | /api/mask-rules/{id} | 获取单个 |
| PUT | /api/mask-rules/{id} | 更新 |
| DELETE | /api/mask-rules/{id} | 删除 |
| POST | /api/mask-rules/{id}/apply | 对输入应用脱敏规则 |
| POST | /api/mask-tasks | 创建脱敏任务 |
| GET | /api/mask-tasks | 列表（分页 + status/source_type/keyword 筛选） |
| GET | /api/mask-tasks/{id} | 获取单个 |
| PUT | /api/mask-tasks/{id} | 更新 |
| DELETE | /api/mask-tasks/{id} | 删除 |
| POST | /api/mask-tasks/{id}/transition | 状态流转（pending→running→completed/failed） |
| POST | /api/mask-tasks/{id}/advance | 推进处理进度 |
| POST | /api/mask-records | 创建处理记录 |
| GET | /api/mask-records | 列表（分页 + mask_task_id/field_name/rule_id 筛选） |
| GET | /api/mask-records/{id} | 获取单个 |
| PUT | /api/mask-records/{id} | 更新 |
| DELETE | /api/mask-records/{id} | 删除 |
| POST | /api/mask-records/batch | 批量创建处理记录 |
| POST | /api/policies | 创建脱敏策略 |
| GET | /api/policies | 列表（分页 + scope/enabled/keyword 筛选，按优先级排序） |
| GET | /api/policies/{id} | 获取单个 |
| PUT | /api/policies/{id} | 更新 |
| DELETE | /api/policies/{id} | 删除 |
| GET | /api/policies/{id}/evaluate | 评估策略，返回命中的规则列表 |
| GET | /api/stats/overview | 统计概览：分类数、规则数、任务按状态分布、记录按 mask_type 分布、字段记录数 TopN |

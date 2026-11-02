# rpa task get

RPA 自动化

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli rpa task get` |
| HTTP | `GET /open/api/rpa/task/getInfo/{id}` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--id` | integer | 是 |  |  |

## 示例

```bash
pmock-cli rpa task get --id 123
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RGetTaskInfoResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 状态码 |
| `data` | object | 请求返回数据 |
| `msg` | string | 请求返回消息 |
| `status` | string | 状态，表示任务是立即完成还是处于等待状态：finish=已完成，pending=等待中 |
| `traceId` | string | 请求链路 ID |

---

[pmock-rpa](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

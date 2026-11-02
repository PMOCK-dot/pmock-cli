# phone task-status

查询手机任务状态

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone task-status` |
| HTTP | `GET /open/api/phone/task/status` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--phone-id` | integer | 是 |  | 手机id |
| `--operate-type` | string |  |  | 手机操作类型:开机-powerOn，重置-reset、重启-restart、一键新机-renew |

## 示例

```bash
pmock-cli phone task-status --phone-id 123 --operate-type powerOn
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RPhoneTaskStatusResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | object | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

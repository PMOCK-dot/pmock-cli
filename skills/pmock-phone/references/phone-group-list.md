# phone group list

获取手机分组及其手机数量

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone group list` |
| HTTP | `GET /open/api/phone/group/list` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 示例

```bash
pmock-cli phone group list --help
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RListPhoneGroupListResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | []PhoneGroupListResVO | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

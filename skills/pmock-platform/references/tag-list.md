# tag list

按绑定手机数量倒序获取标签列表

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli tag list` |
| HTTP | `GET /open/api/tag/list` |
| 风险 | `read`（只读） |
| 别名 | `tag ls` |
| 异步 | 否 |

## 示例

```bash
pmock-cli tag list --help
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RListTagResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | []TagResVO | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态，表示本次是即时完成还是等待中：finish=已完成，pending=等待中 |
| `traceId` | string | 链路id |

---

[pmock-platform](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

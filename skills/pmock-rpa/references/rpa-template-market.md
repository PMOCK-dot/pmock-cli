# rpa template market

RPA 自动化

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli rpa template market` |
| HTTP | `GET /open/api/rpa/template/getMarkList` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--name-or-code` | string |  |  | 流程模板名称/code |
| `--page-num` | integer |  |  | 当前页数量 |
| `--page-size` | integer |  |  | 每页显示的记录数，默认值为 1，单次查询最大数量为 1 |
| `--tag` | string |  |  | 模板tag |

## 示例

```bash
pmock-cli rpa template market --help
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`TableDataInfoGetMarkTemplatePageResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 消息状态码 |
| `msg` | string | 消息内容 |
| `rows` | []GetMarkTemplatePageResVO | 列表数据 |
| `total` | integer | 记录总数 |
| `traceId` | string | 链路 ID |

---

[pmock-rpa](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

# phone batch-update-tags

批量修改多台手机的标签

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone batch-update-tags` |
| HTTP | `POST /open/api/phone/batchUpdateTags` |
| 风险 | `write`（修改数据） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--mode` | string | 是 |  | 模式：append(追加), replace(替换), clear(清空) |
| `--phone-ids` | integer[] | 是 |  | 手机ID列表 |
| `--tag-ids` | integer[] |  |  | 标签ID列表，mode=clear时为null |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli phone batch-update-tags --mode <string> --phone-ids <csv>
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RBatchUpdateTagsResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | object | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态，表示本次是即时完成还是等待中：finish=已完成，pending=等待中 |
| `traceId` | string | 链路id |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

# tag update

修改指定标签的名称或颜色

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli tag update` |
| HTTP | `POST /open/api/tag/update` |
| 风险 | `write`（修改数据） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--id` | integer | 是 |  | 标签ID |
| `--tag-name` | string | 是 |  | 标签名称 |
| `--color` | string |  |  | 颜色 |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli tag update --id <int> --tag-name <string>
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RVoid`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | null | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-platform](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

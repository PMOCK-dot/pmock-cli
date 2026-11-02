# phone action input

手机操作-输入

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone action input` |
| HTTP | `POST /open/api/phone/action/input` |
| 风险 | `write`（修改数据） |
| 异步 | 支持 `--wait`（`--wait-timeout`、`--poll-interval`） |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--code` | string | 是 |  | 手机编号 |
| `--input-type` | integer |  |  | 输入类型：1=直接输入，2=来源网络请求 |
| `--text` | string |  |  | 输入内容 |
| `--url` | string |  |  | 请求的URL 请求URL |
| `--wait-time` | integer |  |  | 等待结果时长 结果等待时长 |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli phone action input --code <string>
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RLong`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | integer | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

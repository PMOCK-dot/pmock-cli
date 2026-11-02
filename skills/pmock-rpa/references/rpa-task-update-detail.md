# rpa task update-detail

RPA 自动化

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli rpa task update-detail` |
| HTTP | `POST /open/api/rpa/task/updateDetail` |
| 风险 | `write`（修改数据） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--id` | integer | 是 |  | 任务id |
| `--phone-id` | integer |  |  | 手机id |
| `--start-time` | integer |  |  | 任务开始时间 |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli rpa task update-detail --id <int>
```

## 请求体

> schema：`UpdateDetailTaskReqVO`

> 下面这些字段没有对应 flag，只能用 `--json` 传入。

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `params` | object |  | 流程模板参数构建结果 |

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

[pmock-rpa](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

# rpa task add

RPA 自动化

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli rpa task add` |
| HTTP | `POST /open/api/rpa/task/add` |
| 风险 | `write`（修改数据） |
| 异步 | 支持 `--wait`（`--wait-timeout`、`--poll-interval`） |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--name` | string | 是 |  | 任务名称 |
| `--template-id` | integer | 是 |  | 流程模板id |
| `--remark` | string |  |  | 备注 |
| `--stop-phone-flag` | integer |  | `0` | 任务结束是否关机，1=是，0=否 |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli rpa task add --template-id 456 --name daily --json '{"taskList":[{"phoneId":123}]}'
```

## 请求体

> schema：`AddTaskReqVO`

> 下面这些字段没有对应 flag，只能用 `--json` 传入。

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `taskList` | []AddTaskListReqVO |  | 任务集合 |

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RAddTaskResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 状态码 |
| `data` | object | 请求返回数据 |
| `msg` | string | 请求返回消息 |
| `status` | string | 状态，表示任务是立即完成还是处于等待状态：finish=已完成，pending=等待中 |
| `traceId` | string | 请求链路 ID |

---

[pmock-rpa](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

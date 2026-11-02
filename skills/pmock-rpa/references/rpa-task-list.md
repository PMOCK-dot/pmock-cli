# rpa task list

RPA 自动化

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli rpa task list` |
| HTTP | `GET /open/api/rpa/task/getList` |
| 风险 | `read`（只读） |
| 别名 | `rpa task ls` |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--end-time` | integer |  |  | 任务开始执行范围结束时间 |
| `--name` | string |  |  | 任务名称 |
| `--page-num` | integer |  |  | 当前页数量 |
| `--page-size` | integer |  |  | 每页显示的记录数，默认值为 1，单次查询最大数量为 1 |
| `--start-time` | integer |  |  | 任务开始执行范围开始时间 |
| `--status` | integer |  |  | 任务状态 |
| `--template-name` | string |  |  | 关联流程(模板)名称 |

## 示例

```bash
pmock-cli rpa task list --help
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`TableDataInfoGetTaskPageListResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 消息状态码 |
| `msg` | string | 消息内容 |
| `rows` | []GetTaskPageListResVO | 列表数据 |
| `total` | integer | 记录总数 |
| `traceId` | string | 链路 ID |

---

[pmock-rpa](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

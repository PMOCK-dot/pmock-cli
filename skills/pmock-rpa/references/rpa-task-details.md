# rpa task details

RPA 自动化

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli rpa task details` |
| HTTP | `GET /open/api/rpa/task/getDetailList` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--task-id` | integer | 是 |  | 任务ID |
| `--page-num` | integer |  |  | 当前页数量 |
| `--page-size` | integer |  |  | 每页显示的记录数，默认值为 1，单次查询最大数量为 1 |

## 示例

```bash
pmock-cli rpa task details --task-id <int>
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`TableDataInfoGetTaskDetailInfoResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 消息状态码 |
| `msg` | string | 消息内容 |
| `rows` | []GetTaskDetailInfoResVO | 列表数据 |
| `total` | integer | 记录总数 |
| `traceId` | string | 链路 ID |

---

[pmock-rpa](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

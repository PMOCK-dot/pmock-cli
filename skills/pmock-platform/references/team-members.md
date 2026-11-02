# team members

获取团队成员及其角色、已授权手机列表

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli team members` |
| HTTP | `GET /open/api/team/member/list` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 示例

```bash
pmock-cli team members --help
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RListOpenApiTeamMemberResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 状态码 |
| `data` | []OpenApiTeamMemberResVO | 团队成员列表 |
| `msg` | string | 请求返回消息 |
| `status` | string | 状态，表示任务是立即完成还是处于等待状态：finish=已完成，pending=等待中 |
| `traceId` | string | 请求链路 ID |

---

[pmock-platform](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

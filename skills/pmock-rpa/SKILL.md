---
name: pmock-rpa
version: 1.0.0
description: "PMOCK RPA management: flow templates (market, team list, collect, delete) and scheduled tasks (add, list, get, details, update, stop, retry, delete). Use when automating repeated app workflows on PMOCK cloud phones."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli rpa --help"
---

# PMOCK RPA（rpa）

RPA 分两层：`rpa template` 管理流程模板（从市场收集、删除），`rpa task` 管理按模板下发的定时任务。下发、停止、重试都作用于任务 ID，不是手机 ID。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli rpa`
- 收录命令：16 个

## 命令索引

### rpa task

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`rpa task add`](references/rpa-task-add.md) | RPA 新增任务 | `write` |
| [`rpa task delete`](references/rpa-task-delete.md) | RPA 删除任务 | `destructive` |
| [`rpa task delete-detail`](references/rpa-task-delete-detail.md) | RPA 删除任务详情 | `destructive` |
| [`rpa task details`](references/rpa-task-details.md) | RPA 获取任务详情列表 | `read` |
| [`rpa task get`](references/rpa-task-get.md) | RPA 获取任务详情 | `read` |
| [`rpa task list`](references/rpa-task-list.md) | RPA 获取任务列表 | `read` |
| [`rpa task retry`](references/rpa-task-retry.md) | RPA 重试任务 | `write` |
| [`rpa task retry-detail`](references/rpa-task-retry-detail.md) | RPA 重试任务详情 | `write` |
| [`rpa task stop`](references/rpa-task-stop.md) | RPA 停止任务 | `write` |
| [`rpa task stop-detail`](references/rpa-task-stop-detail.md) | RPA 停止任务详情 | `write` |
| [`rpa task update`](references/rpa-task-update.md) | RPA 更新任务 | `write` |
| [`rpa task update-detail`](references/rpa-task-update-detail.md) | RPA 更新任务详情 | `write` |

### rpa template

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`rpa template collect`](references/rpa-template-collect.md) | 添加模板到我的团队 | `write` |
| [`rpa template delete`](references/rpa-template-delete.md) | 从我的团队删除模板 | `destructive` |
| [`rpa template list`](references/rpa-template-list.md) | RPA 获取团队模板列表 | `read` |
| [`rpa template market`](references/rpa-template-market.md) | RPA 获取标记模板列表 | `read` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

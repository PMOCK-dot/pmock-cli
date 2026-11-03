---
name: pmock-watch
version: 1.0.0
description: "PMOCK progress stream: follow RPA tasks, phone lifecycle operations and UI action tasks through one customer-facing contract (running / needs_you / done / failed), with a canonical sentence whose specificity never exceeds the evidence the backend returned. Use when an agent must tell a customer what it is doing right now."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli watch --help"
---

# PMOCK 进度流（watch）

同一条命令覆盖所有长任务：把 RPA 的 `0..4`、手机操作的 `taskStatus+phoneStatus`、动作任务的 `1/2/3` 归一成四个状态，并给出措辞不超过证据的客户话术。加能力是加适配器，不是加状态。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli watch`
- 收录命令：4 个

## 命令索引

### watch

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`watch action`](references/watch-action.md) | 跟踪一次 UI 动作任务（点击、滑动、输入……） | `read` |
| [`watch log`](references/watch-log.md) | 回放本地台账里记录过的进度事件 | `read` |
| [`watch phone`](references/watch-phone.md) | 跟踪开机、关机、重启、重置、一键新机这类手机操作 | `read` |
| [`watch rpa`](references/watch-rpa.md) | 跟踪一个 RPA 任务，把状态变化翻成客户能看懂的一句话 | `read` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

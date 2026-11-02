---
name: pmock-platform
version: 1.0.0
description: "PMOCK platform data: products and SKUs with pricing, phone brands, team members and phone authorization, and background task records. Use when looking up prices, brands, team access or task status."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli product --help"
---

# PMOCK 平台数据（product / brand / team / task）

下单前查商品与价格（`product`）、创建机器时选品牌机型（`brand`）、把机器授权给同事（`team`）、排查异步任务结果（`task`）。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli product`、`pmock-cli brand`、`pmock-cli team`、`pmock-cli task`
- 收录命令：4 个

## 命令索引

### brand

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`brand list`](references/brand-list.md) | 获取品牌列表 | `read` |

### product

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`product list`](references/product-list.md) | 获取商品列表 | `read` |

### task

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`task list`](references/task-list.md) | 获取任务信息列表 | `read` |

### team

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`team members`](references/team-members.md) | 获取团队成员列表 | `read` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

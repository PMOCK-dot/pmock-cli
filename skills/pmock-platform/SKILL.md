---
name: pmock-platform
version: 1.0.0
description: "PMOCK platform data: products and SKUs with pricing, changeable plans, phone brands, tags, team members and phone authorization, and background task records. Use when looking up prices, brands, tags, team access or task status."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli product --help"
---

# PMOCK 平台数据（product / brand / team / task）

下单前查商品与价格（`product`）、创建机器时选品牌机型（`brand`）、维护标签（`tag`）、把机器授权给同事（`team`）、排查异步任务结果（`task`）。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli product`、`pmock-cli brand`、`pmock-cli team`、`pmock-cli task`、`pmock-cli tag`
- 收录命令：9 个

## 命令索引

### brand

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`brand list`](references/brand-list.md) | 获取品牌列表 | `read` |

### product

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`product changeable-sku-list`](references/product-changeable-sku-list.md) | 获取当前套餐及可变更套餐 | `read` |
| [`product list`](references/product-list.md) | 获取商品列表 | `read` |

### tag

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`tag create`](references/tag-create.md) | 创建标签 | `write` |
| [`tag delete`](references/tag-delete.md) | 删除标签 | `destructive` |
| [`tag list`](references/tag-list.md) | 获取标签列表 | `read` |
| [`tag update`](references/tag-update.md) | 修改标签 | `write` |

### task

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`task list`](references/task-list.md) | 获取任务信息列表 | `read` |

### team

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`team members`](references/team-members.md) | 获取团队成员列表 | `read` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

---
name: pmock-billing
version: 1.0.0
description: "PMOCK billing: wallet balance, recharge configuration and orders, wallet top-up, and consumption or recharge records. Use when checking balance, recharging the account or auditing spend."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli payment --help"
---

# PMOCK 钱包与账单（user / payment / order）

钱相关的三组命令：`user balance` 查余额、`payment *` 充值（`payment recharge` 会真实扣款，属 costly，必须先跟用户确认金额）、`order *` 查消费与充值流水。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli user`、`pmock-cli payment`、`pmock-cli order`
- 收录命令：6 个

## 命令索引

### order

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`order consume-record`](references/order-consume-record.md) | 获取消费记录 | `read` |
| [`order recharge-record`](references/order-recharge-record.md) | 获取充值记录 | `read` |

### payment

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`payment recharge`](references/payment-recharge.md) | 钱包充值 | `costly` |
| [`payment recharge-config`](references/payment-recharge-config.md) | 获取充值配置 | `read` |
| [`payment recharge-order`](references/payment-recharge-order.md) | 获取充值订单详情 | `read` |

### user

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`user balance`](references/user-balance.md) | 获取钱包余额 | `read` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

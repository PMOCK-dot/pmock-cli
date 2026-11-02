---
name: pmock-network
version: 1.0.0
description: "PMOCK network management: add, update, remove, inspect and check proxies, organize them in proxy groups, and set the network of a cloud phone. Use when configuring proxies or switching the network of PMOCK cloud phones."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli network --help"
---

# PMOCK 网络（network）

网络是云手机风控的基础：先 `network list` 拿到代理 ID，`network detect` 验证可用性，再用 `phone network-set` 绑定到机器。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli network`
- 收录命令：11 个

## 命令索引

### network

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`network add`](references/network-add.md) | 新增代理 | `write` |
| [`network check`](references/network-check.md) | 检测代理 | `read` |
| [`network detail`](references/network-detail.md) | 获取代理详情 | `read` |
| [`network detect`](references/network-detect.md) | 网络检测 | `read` |
| [`network list`](references/network-list.md) | 获取网络列表 | `read` |
| [`network remove`](references/network-remove.md) | 删除代理 | `destructive` |
| [`network update`](references/network-update.md) | 修改代理 | `write` |

### network group

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`network group create`](references/network-group-create.md) | 创建代理分组 | `write` |
| [`network group delete`](references/network-group-delete.md) | 删除代理分组 | `destructive` |
| [`network group list`](references/network-group-list.md) | 获取代理分组列表 | `read` |
| [`network group update`](references/network-group-update.md) | 修改代理分组 | `write` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

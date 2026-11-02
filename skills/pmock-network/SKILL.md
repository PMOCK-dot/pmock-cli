---
name: pmock-network
version: 1.0.0
description: "PMOCK network management: list proxies, detect a proxy, and set the network of a cloud phone. Use when configuring proxies or switching the network of PMOCK cloud phones."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli network --help"
---

# PMOCK 网络（network）

网络是云手机风控的基础：先 `network list` 拿到代理 ID，`network detect` 验证可用性，再用 `phone network-set` 绑定到机器。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli network`
- 收录命令：2 个

## 命令索引

### network

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`network detect`](references/network-detect.md) | 网络检测 | `read` |
| [`network list`](references/network-list.md) | 获取网络列表 | `read` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

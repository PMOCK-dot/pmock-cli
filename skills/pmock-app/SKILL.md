---
name: pmock-app
version: 1.0.0
description: "PMOCK application management on cloud phones: list installed apps and categories, install an app or APK, check the upgrade list and operate running apps. Use when installing or controlling apps on PMOCK cloud phones."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli app --help"
---

# PMOCK 应用（app）

应用管理围绕云手机上的包：查询已安装、安装（含本地上传的 APK）、查看可升级版本、启动/停止/卸载。安装前先确认目标机器的 `phoneId`。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli app`
- 收录命令：5 个

## 命令索引

### app

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`app categories`](references/app-categories.md) | 获取应用分类列表 | `read` |
| [`app install`](references/app-install.md) | 手机安装应用 | `write` |
| [`app list`](references/app-list.md) | 获取应用列表 | `read` |
| [`app operate`](references/app-operate.md) | 手机应用操作 | `write` |
| [`app upgrades`](references/app-upgrades.md) | 获取可升级应用列表 | `read` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

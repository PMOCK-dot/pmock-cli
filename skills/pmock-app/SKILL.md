---
name: pmock-app
version: 1.0.0
description: "PMOCK application management: list installed apps and categories, install an app or APK, check upgrades, operate running apps, and maintain the shared team app library. Use when installing or controlling apps on PMOCK cloud phones."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli app --help"
---

# PMOCK 应用（app）

应用管理分两层：`app *` 操作单台机器上的包（查询、安装、升级、启停、卸载），`app team-app *` 维护团队共享的应用库与版本。安装前先确认目标机器的 `phoneId`。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli app`
- 收录命令：9 个

## 命令索引

### app

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`app categories`](references/app-categories.md) | 获取应用分类列表 | `read` |
| [`app install`](references/app-install.md) | 手机安装应用 | `write` |
| [`app list`](references/app-list.md) | 获取应用列表 | `read` |
| [`app operate`](references/app-operate.md) | 手机应用操作 | `write` |
| [`app upgrades`](references/app-upgrades.md) | 获取可升级应用列表 | `read` |

### app team-app

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`app team-app add`](references/app-team-app-add.md) | 添加应用到团队应用库 | `write` |
| [`app team-app list`](references/app-team-app-list.md) | 获取团队应用列表 | `read` |
| [`app team-app remove`](references/app-team-app-remove.md) | 从团队应用库删除应用 | `destructive` |
| [`app team-app remove-version`](references/app-team-app-remove-version.md) | 删除应用版本 | `destructive` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

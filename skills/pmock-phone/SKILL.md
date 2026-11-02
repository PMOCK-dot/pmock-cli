---
name: pmock-phone
version: 1.0.0
description: "PMOCK cloud phone management: list, create, start, stop, restart, reset, renew, screenshot, GPS, ADB, network settings, files and UI automation actions. Use when operating PMOCK cloud phones through pmock-cli."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli phone --help"
---

# PMOCK 云手机（phone）

云手机是 PMOCK 的核心资产：创建与续费、开关机、截图与 GPS、ADB 与网络配置、文件上传，以及点击/滑动/输入等 UI 自动化。涉及计费或不可逆的操作，先确认再执行。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli phone`
- 收录命令：32 个

## 命令索引

### phone

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`phone adb`](references/phone-adb.md) | 获取 ADB 信息 | `read` |
| [`phone authorize`](references/phone-authorize.md) | 授权手机给团队成员 | `write` |
| [`phone clear`](references/phone-clear.md) | 手机加速 | `write` |
| [`phone cpy-connect`](references/phone-cpy-connect.md) | 获取 CPY 连接详情。 | `read` |
| [`phone create`](references/phone-create.md) | 创建手机 | `costly` |
| [`phone list`](references/phone-list.md) | 获取手机列表 | `read` |
| [`phone models`](references/phone-models.md) | 获取机型列表 | `read` |
| [`phone network-set`](references/phone-network-set.md) | 手机网络设置 | `write` |
| [`phone new`](references/phone-new.md) | 一键新机 | `destructive` |
| [`phone remote-token`](references/phone-remote-token.md) | 获取远程连接Token | `read` |
| [`phone remove`](references/phone-remove.md) | 删除手机 | `destructive` |
| [`phone renew`](references/phone-renew.md) | 批量续费手机 | `costly` |
| [`phone reset`](references/phone-reset.md) | 恢复出厂设置 | `destructive` |
| [`phone resolutions`](references/phone-resolutions.md) | 获取手机分辨率列表 | `read` |
| [`phone restart`](references/phone-restart.md) | 重启手机 | `write` |
| [`phone screen`](references/phone-screen.md) | 手机截图 | `read` |
| [`phone start`](references/phone-start.md) | 开机 | `write` |
| [`phone stop`](references/phone-stop.md) | 关机 | `write` |
| [`phone stream-token`](references/phone-stream-token.md) | 获取远程推流组件Url | `read` |
| [`phone task-status`](references/phone-task-status.md) | 获取手机任务状态 | `read` |
| [`phone unauthorize`](references/phone-unauthorize.md) | 取消团队成员的手机授权 | `write` |

### phone action

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`phone action app`](references/phone-action-app.md) | 手机操作 - 操作应用 | `write` |
| [`phone action click`](references/phone-action-click.md) | 手机操作 - 单击 | `write` |
| [`phone action double-click`](references/phone-action-double-click.md) | 手机操作 - 双击 | `write` |
| [`phone action drag`](references/phone-action-drag.md) | 手机操作 - 拖动 | `write` |
| [`phone action input`](references/phone-action-input.md) | 手机操作 - 输入 | `write` |
| [`phone action long-click`](references/phone-action-long-click.md) | 手机操作 - 长按 | `write` |
| [`phone action status`](references/phone-action-status.md) | 手机操作 - 查询任务状态 | `read` |
| [`phone action swipe`](references/phone-action-swipe.md) | 手机操作 - 滑动 | `write` |

### phone file

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`phone file download`](references/phone-file-download.md) | 下载手机文件 | `read` |
| [`phone file scan`](references/phone-file-scan.md) | 扫描手机文件 | `read` |
| [`phone file upload`](references/phone-file-upload.md) | 上传本地文件或 APK 到云手机（分片、可续传） | `write` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

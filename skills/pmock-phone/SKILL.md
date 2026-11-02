---
name: pmock-phone
version: 1.0.0
description: "PMOCK cloud phone management: list, create, start, stop, restart, reset, renew, screenshot, GPS, ADB, rename, retag, change plan, groups, ownership transfer, files and UI automation actions. Use when operating PMOCK cloud phones."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli phone --help"
---

# PMOCK 云手机（phone）

云手机是 PMOCK 的核心资产：创建与续费、开关机、截图与 GPS、ADB 与网络配置、改名打标签、分组与转让、文件上传，以及点击/滑动/输入等 UI 自动化。换套餐（update-sku）和转让确认（transfer confirm）影响计费或归属，先确认再执行。

**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。

- 命令前缀：`pmock-cli phone`
- 收录命令：48 个

## 命令索引

### phone

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`phone adb`](references/phone-adb.md) | 获取 ADB 信息 | `read` |
| [`phone authorize`](references/phone-authorize.md) | 授权手机给团队成员 | `write` |
| [`phone batch-update-remark`](references/phone-batch-update-remark.md) | 批量修改手机备注 | `write` |
| [`phone batch-update-tags`](references/phone-batch-update-tags.md) | 批量修改手机标签 | `write` |
| [`phone clear`](references/phone-clear.md) | 手机加速 | `write` |
| [`phone cpy-connect`](references/phone-cpy-connect.md) | 获取 CPY 连接详情。 | `read` |
| [`phone create`](references/phone-create.md) | 创建手机 | `costly` |
| [`phone list`](references/phone-list.md) | 获取手机列表 | `read` |
| [`phone models`](references/phone-models.md) | 获取机型列表 | `read` |
| [`phone move-group`](references/phone-move-group.md) | 移动手机到分组 | `write` |
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
| [`phone update-env-setting`](references/phone-update-env-setting.md) | 修改手机环境设置 | `write` |
| [`phone update-name`](references/phone-update-name.md) | 修改手机名称 | `write` |
| [`phone update-remark`](references/phone-update-remark.md) | 修改手机备注 | `write` |
| [`phone update-sku`](references/phone-update-sku.md) | 修改手机套餐 | `costly` |
| [`phone update-tags`](references/phone-update-tags.md) | 修改手机标签 | `write` |

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
| [`phone file task-status`](references/phone-file-task-status.md) | 获取文件任务状态 | `read` |
| [`phone file upload`](references/phone-file-upload.md) | 上传本地文件或 APK 到云手机（分片、可续传） | `write` |

### phone group

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`phone group create`](references/phone-group-create.md) | 创建手机分组 | `write` |
| [`phone group delete`](references/phone-group-delete.md) | 删除手机分组 | `destructive` |
| [`phone group list`](references/phone-group-list.md) | 获取手机分组列表 | `read` |
| [`phone group update`](references/phone-group-update.md) | 修改手机分组 | `write` |

### phone transfer

| 命令 | 用途 | 风险 |
| --- | --- | --- |
| [`phone transfer cancel`](references/phone-transfer-cancel.md) | 转让手机-取消转让 | `write` |
| [`phone transfer confirm`](references/phone-transfer-confirm.md) | 转让手机-确认转让 | `destructive` |
| [`phone transfer send-code`](references/phone-transfer-send-code.md) | 转让手机-发送验证码 | `write` |

风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · `destructive` 不可逆（需显式 `--yes`）。

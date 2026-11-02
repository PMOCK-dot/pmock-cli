# phone file upload

上传本地文件或 APK 到云手机（分片、可续传）

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone file upload` |
| 风险 | `write`（修改数据） |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--id` | integer | 是 |  | 目标云手机 ID |
| `--file` | string | 是 |  | 本地文件路径 |
| `--type` | string | 否 | `file` | `file` 普通文件；`app` 表示 APK |
| `--dest` | string | 否 | `/sdcard/Download` | 手机上的目标目录 |
| `--chunk-size-mb` | integer | 否 | `30` | 分片大小，30–300 MiB |
| `--md5` | string | 否 |  | 预计算的 MD5，跳过本地计算 |

## 行为

1. 本地计算 MD5（除非传 `--md5`），用 `upload-init` 换取分片清单；
2. 只上传服务端缺失的分片，因此中断后重跑同一命令即可续传；
3. 全部分片就位后 `upload-push` 合并并推送到手机。

`--type app` 时返回的 `data.envFileId` 可以直接传给 `pmock-cli app install --file-id`。

## 示例

```bash
# 普通文件
pmock-cli phone file upload --id 123 --file ./video.mp4

# APK，随后安装
pmock-cli phone file upload --id 123 --file ./app.apk --type app
pmock-cli app install --file-id <envFileId>

# 指定目录，并先看计划不实际传输
pmock-cli phone file upload --id 123 --file ./video.mp4 --dest /sdcard/Movies --dry-run
```

## 注意

- 大文件按 45 分钟预算上传，不受全局 `--timeout` 限制；
- 进度输出到 stderr，仅在 TTY 或 `--verbose` 时出现，不影响 JSON 输出；
- 单文件 MD5 不变时重复执行是幂等的。

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

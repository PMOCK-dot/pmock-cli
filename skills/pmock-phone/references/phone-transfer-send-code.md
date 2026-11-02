# phone transfer send-code

向接收团队的 BOSS 邮箱发送验证码，仅 BOSS 可操作

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone transfer send-code` |
| HTTP | `POST /open/api/phone/transfer/sendCode` |
| 风险 | `write`（修改数据） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--phone-ids` | integer[] | 是 |  | 已选手机ID列表 |
| `--transfer-options` | integer | 是 |  | 转让选项位掩码：1代理 2标签 4备注 |
| `--target-boss-email` | string |  |  | 接收团队BOSS邮箱（未传 targetTeamId 时必填） |
| `--target-team-id` | integer |  |  | 接收团队ID（若传则优先使用，可替代邮箱+团队名） |
| `--target-team-name` | string |  |  | 接收团队名称（未传 targetTeamId 时必填） |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli phone transfer send-code --phone-ids <csv> --transfer-options <int>
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RListPhoneTransferResultItemVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | []PhoneTransferResultItemVO | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态，表示本次是即时完成还是等待中：finish=已完成，pending=等待中 |
| `traceId` | string | 链路id |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

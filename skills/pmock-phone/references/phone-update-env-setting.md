# phone update-env-setting

修改手机的经纬度、归属地、语言、时区等环境设置

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone update-env-setting` |
| HTTP | `POST /open/api/phone/updateEnvSetting` |
| 风险 | `write`（修改数据） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--language-follow-ip` | boolean | 是 |  | 语言是否跟随IP |
| `--lat-lng-follow-ip` | boolean | 是 |  | 经纬度是否跟随 IP：true=跟随，false=自定义 |
| `--location-follow-ip` | boolean | 是 |  | 归属地是否跟随IP |
| `--phone-id` | integer | 是 |  | 手机ID |
| `--timezone-follow-ip` | boolean | 是 |  | 时区是否跟随IP |
| `--language` | string |  |  | 语言 |
| `--latitude` | string |  |  | 纬度 |
| `--location` | string |  |  | 归属地 |
| `--longitude` | string |  |  | 经度 |
| `--next-power-on-effective` | integer |  |  | 是否下次开机生效 0-立即生效，1-下次开机生效 |
| `--timezone` | string |  |  | 时区 |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli phone update-env-setting --language-follow-ip <true|false> --lat-lng-follow-ip <true|false> --location-follow-ip <true|false> --phone-id <int> --timezone-follow-ip <true|false>
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RVoid`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | null | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

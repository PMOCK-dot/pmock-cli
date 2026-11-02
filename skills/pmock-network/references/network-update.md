# network update

修改代理资源

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli network update` |
| HTTP | `POST /open/api/network/update` |
| 风险 | `write`（修改数据） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--account` | string |  |  | 账号 |
| `--check-type` | string |  |  | IP查询渠道: dbip/ip2location/maxmind/ipapi，不传不修改 |
| `--country-code` | string |  |  | 国家代码 |
| `--country-name` | string |  |  | 国家 |
| `--detect-data-type` | integer |  |  | 检测数据类型：1=跟随 IP，2=自定义 |
| `--expire-time` | integer |  |  | 代理到期时间（毫秒时间戳） |
| `--group-id` | integer |  |  | 网络分组ID |
| `--host` | string |  |  | ip/域名 |
| `--id` | integer |  |  | 主键id |
| `--ip` | string |  |  | 实际ip地址 |
| `--language` | string |  |  | 语言 |
| `--latitude` | string |  |  | 纬度 |
| `--longitude` | string |  |  | 经度 |
| `--name` | string |  |  | 网络名称 |
| `--password` | string |  |  | 密码 |
| `--port` | integer |  |  | 端口 |
| `--region` | string |  |  | 地区 |
| `--remark` | string |  |  | 备注 |
| `--status` | string |  |  | 状态 0=未检测 1=检测通过 2=检测失败 |
| `--timezone` | string |  |  | 时区 |
| `--type` | string |  |  | 代理类型 1.socks5 2.https 3.http |
| `--udp` | integer |  |  | 是否开启UDP |
| `--way` | string |  |  | 使用方式 1.静态ip 2.动态ip |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli network update --help
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

[pmock-network](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

# phone network-set

手机网络设置

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone network-set` |
| HTTP | `POST /open/api/phone/setNetwork` |
| 风险 | `write`（修改数据） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--id` | integer | 是 |  | 手机ID |
| `--account` | string |  |  | 代理账号 |
| `--country-code` | string |  |  | 国家代码 |
| `--country-name` | string |  |  | 国家 |
| `--detect-data-type` | integer |  |  | 检测数据类型：1=跟随 IP，2=自定义 |
| `--expire-time` | integer |  |  | 代理过期时间（毫秒时间戳） |
| `--group-id` | integer |  | `-1` | 网络分组ID |
| `--host` | string |  |  | 代理主机 |
| `--ip` | string |  |  | 实际ip地址 |
| `--language` | string |  |  | 语言 |
| `--latitude` | string |  |  | 纬度 |
| `--longitude` | string |  |  | 经度 |
| `--name` | string |  |  | 网络名称 |
| `--network-id` | integer |  |  | 网络ID |
| `--next-power-on-effective` | integer |  |  | 生效时机：0=立即应用，1=下次重启时应用 |
| `--password` | string |  |  | 代理密码 |
| `--port` | integer |  |  | 代理端口 |
| `--region` | string |  |  | 地区 |
| `--remark` | string |  |  | 备注 |
| `--status` | integer |  |  | 代理检测状态；如果设置了 networkId，则以该 networkId 对应的数据为准。 |
| `--tcp` | integer |  |  | 是否开启TCP |
| `--timezone` | string |  |  | 时区 |
| `--type` | string |  |  |  |
| `--udp` | integer |  |  | 是否开启UDP |
| `--way` | string |  |  | 使用方式 1.静态ip 2.动态ip |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli phone network-set --id <int>
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

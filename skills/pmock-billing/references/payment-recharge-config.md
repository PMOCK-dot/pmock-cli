# payment recharge-config

获取充值金额范围及固定金额列表

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli payment recharge-config` |
| HTTP | `GET /open/api/payment/getRechargeConfig` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 示例

```bash
pmock-cli payment recharge-config --help
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RRechargeConfigResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | object | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-billing](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

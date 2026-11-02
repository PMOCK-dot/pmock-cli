# payment recharge-order

预览充值订单，包含折扣、手续费等信息

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli payment recharge-order` |
| HTTP | `POST /open/api/payment/getRechargeOrder` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--recharge-type` | string | 是 |  | 金额类型 |
| `--amount` | number |  |  | 金额 |
| `--recharge-id` | integer |  |  | 充值金额配置id |
| `--team-id` | integer |  |  | 团队id |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli payment recharge-order --recharge-type <string>
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RRechargeOrderInfoResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | object | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-billing](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

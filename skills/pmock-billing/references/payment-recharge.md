# payment recharge

按金额为团队钱包充值

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli payment recharge` |
| HTTP | `POST /open/api/payment/money/recharge` |
| 风险 | `costly`（会产生费用，执行前确认） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--method` | string | 是 |  | 支付方式；取值：air_alipay_cn_qr, air_wechat_qr, wallet_pay, offline_pay, system_gift, alipay_pc_direct, wechat_native, stripe_wechat_qr, stripe_alipay_qr, wechat_global_native, coupon_gift, whop_checkout_configurations, stripe_car_checkout, stripe_googlepay_checkout, stable_pay_checkout, rebate_withdraw, paypal_checkout |
| `--recharge-type` | string | 是 |  | 金额类型 |
| `--team-id` | integer | 是 |  | 团队id |
| `--coupon-id` | integer |  |  | 优惠券id |
| `--money` | number |  |  | 金额 |
| `--recharge-id` | integer |  |  | 充值金额配置id |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli payment recharge --method <string> --recharge-type <string> --team-id <int>
```

## 请求体

> schema：`PaymentReqVO`

> 下面这些字段没有对应 flag，只能用 `--json` 传入。

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `methodExt` | object |  | 流程模板参数构建结果 |

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RWalletRechargeResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | object | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-billing](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

# order recharge-record

获取团队充值账单记录

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli order recharge-record` |
| HTTP | `GET /open/api/order/rechargeRecord` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--channel-code` | string |  |  | 充值支付类型：system=系统，airwallex=空中云汇，alipay=支付宝，wechat=微信，stripe=Stripe，wechatGlobal=微信支付国际版，whop=Whop，stablePay=StablePay，paypal=PayPal；取值：system, airwallex, alipay, wechat, stripe, wechatGlobal, whop, stablePay, paypal |
| `--end-time` | string |  |  | 充值结束时间 |
| `--id` | string |  |  | 订单 ID |
| `--page-num` | integer |  |  | 当前页码 |
| `--page-size` | integer |  |  | 每页显示的记录数 |
| `--pay-method` | string |  |  | 支付方式 |
| `--start-time` | string |  |  | 充值开始时间 |

## 示例

```bash
pmock-cli order recharge-record --help
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`TableDataInfoPayOrderRechargeListResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 消息状态码 |
| `msg` | string | 消息内容 |
| `rows` | []PayOrderRechargeListResVO | 列表数据 |
| `total` | integer | 总记录数 |
| `traceId` | string | 链路id |

---

[pmock-billing](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

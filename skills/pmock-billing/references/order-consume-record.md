# order consume-record

获取团队消费账单记录

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli order consume-record` |
| HTTP | `GET /open/api/order/consumeRecord` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--channel-code` | string |  |  | 支付类型 |
| `--create-user-id` | integer |  |  | 创建人 ID |
| `--end-time` | string |  |  | 扣费结束时间 |
| `--id-or-code` | string |  |  | 订单 ID 或手机编号 |
| `--operate-type` | integer |  |  | 操作类型：RECHARGE_WALLET=钱包充值，USER_ORDER=用户下单，SYSTEM_CHARGING=系统扣费，BUY_BOX=盲盒，AI_REFUND=AI 退款 |
| `--order-ids` | integer[] |  |  | 订单 ID 列表 |
| `--order-type` | integer |  |  | 订单类型 |
| `--page-num` | integer |  |  | 当前页码 |
| `--page-size` | integer |  |  | 每页显示的记录数 |
| `--pay-method` | string |  |  | 支付方式 |
| `--platform` | integer |  |  | 平台来源 |
| `--spu-type` | integer |  |  | 商品类型 |
| `--start-time` | string |  |  | 扣费开始时间 |
| `--status` | integer[] |  |  | 支付状态：0=未支付，1=支付成功，2=已退款，3=支付关闭，4=支付失败，5=待发货，6=待收货，7=完成，8=权益回收且退款，9=仅退款，10=仅关单，11=订单关闭，99=异常订单 |
| `--team-id` | integer |  |  | 团队 ID |
| `--user-ids` | integer[] |  |  | 创建订单的用户 ID 列表 |

## 示例

```bash
pmock-cli order consume-record --help
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`TableDataInfoPayOrderConsumeListResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 消息状态码 |
| `msg` | string | 消息内容 |
| `rows` | []PayOrderConsumeListResVO | 列表数据 |
| `total` | integer | 总记录数 |
| `traceId` | string | 链路id |

---

[pmock-billing](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

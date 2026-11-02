# product list

获取商品信息列表

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli product list` |
| HTTP | `GET /open/api/product/list` |
| 风险 | `read`（只读） |
| 别名 | `product ls` |
| 异步 | 否 |

## 示例

```bash
pmock-cli product list
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RListProductInfoVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | []ProductInfoVO | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-platform](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

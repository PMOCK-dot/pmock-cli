# product changeable-sku-list

获取手机当前套餐（SKU）及可变更套餐列表

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli product changeable-sku-list` |
| HTTP | `POST /open/api/product/getChangeableSkuList` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--phone-id` | integer | 是 |  | 手机id |
| `--sku-id` | integer |  |  | 基准SKU ID，用于计算对比价格，不传则默认以第一个可变更SKU为基准 |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli product changeable-sku-list --phone-id <int>
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RChangeableInfoResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | object | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态，表示本次是即时完成还是等待中：finish=已完成，pending=等待中 |
| `traceId` | string | 链路id |

---

[pmock-platform](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

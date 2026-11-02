# phone create

创建手机

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone create` |
| HTTP | `POST /open/api/phone/create` |
| 风险 | `costly`（会产生费用，执行前确认） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--num` | integer | 是 |  | 购买数量 |
| `--sku-id` | integer | 是 |  | 商品 SKU ID |
| `--spu-id` | integer | 是 |  | 商品 SPU ID |
| `--auto-install-flag` | integer |  | `0` | 自动安装标志：0=关闭，1=开启 |
| `--brand` | string |  |  | 品牌名称 |
| `--config-id` | integer |  |  | 手机机型 ID，根据“获取机型列表”接口返回的 ID 创建对应机型的手机 |
| `--group-id` | integer |  |  | 分组ID，不传则为未分组 |
| `--height` | integer |  |  | 分辨率高度 |
| `--phone-number` | string |  |  | 手机号 |
| `--real-resolution` | boolean |  |  | 是否真机分辨率 |
| `--remark` | string |  |  | 备注，最多200字符 |
| `--tag-ids` | integer[] |  |  | 标签ID列表，最多5个 |
| `--width` | integer |  |  | 分辨率宽度 |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli phone create --spu-id 123 --sku-id 456 --num 1
pmock-cli phone create --spu-id 123 --sku-id 456 --json '{"tagIds":[1,2]}'
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RListPhoneInfoResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 状态码 |
| `data` | []PhoneInfoResVO | 请求返回数据 |
| `msg` | string | 请求返回消息 |
| `status` | string | 状态，表示任务是立即完成还是处于等待状态：finish=已完成，pending=等待中 |
| `traceId` | string | 请求链路 ID |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

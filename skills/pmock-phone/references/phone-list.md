# phone list

分页获取手机列表

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone list` |
| HTTP | `GET /open/api/phone/list` |
| 风险 | `read`（只读） |
| 别名 | `phone ls` |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--authorize-user-ids` | integer[] |  |  | 授权用户 ID 列表，取当前接口返回结果中的 `userId` 值 |
| `--expire-status` | integer |  |  | 计费(过期)状态：0=已经过期，1=正常 |
| `--group-code` | string |  |  | 手机系统分组编号 |
| `--group-id` | integer |  |  | 手机分组ID |
| `--group-type` | string |  |  | 手机分组类型 |
| `--name-or-code` | string |  |  | 手机名称/手机编号 |
| `--network-id` | integer |  |  | 网络ID |
| `--page-num` | integer |  |  | 当前页数 |
| `--page-size` | integer |  |  | 分页大小，每页显示记录数；上游写的「默认1条、最大1条」有误，实测不传即返回全部，建议传 100 一次取回全部 |
| `--real-resolution` | integer |  |  | 是否为真机分辨率：0=否，1=是 |
| `--sku-id` | integer |  |  | 版本ID |
| `--sku-name` | string |  |  | 版本 |
| `--spu-id` | integer |  |  | 商品 SPU ID |
| `--status` | integer |  |  | 手机状态 |
| `--tag-ids` | integer[] |  |  | 标签ID列表，筛选包含这些标签的设备（取并集） |

> **注意**：拉全量列表先用 `--page-size 100` 取一次，只有 `rows` 少于 `total` 时才按 `--page-num` 继续翻页；不要因为上游写了「最大1条」就按 1 条一页逐页扫。

## 示例

```bash
pmock-cli phone list --page-size 100
pmock-cli phone list --name-or-code demo
pmock-cli phone list --status 1 --page-size 100
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`TableDataInfoPhonePageListResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 消息状态码 |
| `msg` | string | 消息内容 |
| `rows` | []PhonePageListResVO | 列表数据 |
| `total` | integer | 总记录数 |
| `traceId` | string | 链路id |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

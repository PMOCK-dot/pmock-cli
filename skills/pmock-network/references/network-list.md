# network list

获取网络列表

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli network list` |
| HTTP | `GET /open/api/network/list` |
| 风险 | `read`（只读） |
| 别名 | `network ls` |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--group-code` | string |  |  | 网络系统分组编码：all_group=全部分组，renew_expired=续费后到期，un_group=未分组 |
| `--group-id` | integer |  |  | 当前接口返回结果中的 `groupId` 值 |
| `--group-type` | string |  |  | 网络分组类型：system_group=系统分组，custom_group=自定义分组 |
| `--host` | string |  |  | 代理网络 IP/域名，取当前接口返回结果中的 `host` 值 |
| `--num` | integer |  |  | 当前接口返回结果中包含的 `num` 值 |
| `--nums` | string |  |  | 代理序号，多个序号使用英文逗号分隔 |
| `--page-num` | integer |  |  | 当前页数量 |
| `--page-size` | integer |  |  | 每页显示的记录数，默认值为 1，单次查询最大数量为 1 |
| `--remark` | string |  |  | 备注 |

## 示例

```bash
pmock-cli network list --help
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`TableDataInfoNetworkListResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 消息状态码 |
| `msg` | string | 消息内容 |
| `rows` | []NetworkListResVO | 列表数据 |
| `total` | integer | 总记录数 |
| `traceId` | string | 链路id |

---

[pmock-network](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

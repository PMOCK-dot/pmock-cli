# app team-app list

分页获取团队应用库列表

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli app team-app list` |
| HTTP | `POST /open/api/app/teamApp/list` |
| 风险 | `read`（只读） |
| 别名 | `app team-app ls` |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--apk-name` | string |  |  | 应用名称 |
| `--apk-pkg` | string |  |  | apk 包名 |
| `--cate-code` | integer |  | `-1` | 应用分类编码 |
| `--id` | integer |  |  |  |
| `--only-show-installed` | boolean |  | `false` | 只显示已经安装的应用, 默认值：false |
| `--page-num` | integer |  |  | 当前页码 |
| `--page-size` | integer |  |  | 每页显示的记录数 |
| `--shelf-status` | integer |  |  | 上架状态：1=已上架，0=已下架 |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli app team-app list --help
```

## 请求体

> schema：`PhoneAppListReqVO`

> 下面这些字段没有对应 flag，只能用 `--json` 传入。

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `pageNum` | integer |  | 当前页数量 |
| `pageSize` | integer |  | 每页显示的记录数，默认值为 1，单次查询最大数量为 1 |

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`TableDataInfoTeamAppListResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 消息状态码 |
| `msg` | string | 消息内容 |
| `rows` | []TeamAppListResVO | 列表数据 |
| `total` | integer | 总记录数 |
| `traceId` | string | 链路id |

---

[pmock-app](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

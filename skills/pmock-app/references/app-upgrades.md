# app upgrades

获取可升级应用列表

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli app upgrades` |
| HTTP | `POST /open/api/app/listUpgradeApp` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--id` | integer | 是 |  |  |
| `--config-id` | integer |  |  | 手机机型 ID，根据“获取机型列表”接口返回的 ID 创建对应机型的手机 |
| `--page-num` | integer |  |  | 当前页数 |
| `--page-size` | integer |  |  | 分页大小 每页显示记录数 默认为1条,单次查询最大条数为:1 |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli app upgrades --id <int>
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`TableDataInfoAppListV2ResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 消息状态码 |
| `msg` | string | 消息内容 |
| `rows` | []AppListV2ResVO | 列表数据 |
| `total` | integer | 总记录数 |
| `traceId` | string | 链路id |

---

[pmock-app](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

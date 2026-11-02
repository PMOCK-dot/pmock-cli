# phone screen

手机截图

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone screen` |
| HTTP | `POST /open/api/phone/screenCap` |
| 风险 | `read`（只读） |
| 别名 | `phone screenshot` |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--id` | integer | 是 |  | 手机 ID，由“获取手机列表”接口返回数据中的 ID。 |
| `--is-remove` | boolean |  | `true` | 执行完成后是否删除手机内容的截图文件(暂时不开放) |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli phone screen --id 123
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RPhoneScreenCapResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 状态码 |
| `data` | object | 请求返回数据 |
| `msg` | string | 请求返回消息 |
| `status` | string | 状态，表示任务是立即完成还是处于等待状态：finish=已完成，pending=等待中 |
| `traceId` | string | 请求链路 ID |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

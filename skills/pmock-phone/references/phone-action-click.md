# phone action click

手机操作-单击

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone action click` |
| HTTP | `POST /open/api/phone/action/click` |
| 风险 | `write`（修改数据） |
| 别名 | `phone action tap` |
| 异步 | 支持 `--wait`（`--wait-timeout`、`--poll-interval`） |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--code` | string | 是 |  | 手机编号 |
| `--pos` | integer[] |  |  | 坐标 X,Y 坐标,格式[X,Y] |
| `--pos-type` | integer |  | `1` | 位置类型：1=节点，2=固定坐标，3=范围位置 |
| `--query` | string[] |  |  | 节点,格式[key=value,key=value,...] |
| `--range` | integer[] |  |  | [x1,y1,x2,y2] 区域,格式[X1,Y1,X2,Y2] |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli phone action click --code ABC123 --pos-type 2 --pos 100,200
pmock-cli phone action click --code ABC123 --pos-type 1 --query "text=登录,index=0"
```

## 请求体

> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；嵌套结构再改用 `--json '<object>'`。

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RLong`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | integer | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

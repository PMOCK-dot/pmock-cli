# phone action swipe

手机操作-滑动

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone action swipe` |
| HTTP | `POST /open/api/phone/action/swipe` |
| 风险 | `write`（修改数据） |
| 异步 | 支持 `--wait`（`--wait-timeout`、`--poll-interval`） |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--code` | string | 是 |  | 手机编号 |
| `--duration` | integer |  |  | 持续时长 |
| `--duration-range` | integer[] |  |  | 持续时长范围 |
| `--pos-type` | integer |  | `1` | 位置类型：1=节点，2=固定坐标，3=范围位置 位置1类型：1=节点，2=固定坐标，3=范围位置 |
| `--pos-type2` | integer |  | `1` | 位置类型：1=节点，2=固定坐标，3=范围位置 位置2类型：1=节点，2=固定坐标，3=范围位置 |
| `--pos1` | integer[] |  |  | 坐标 X,Y 坐标1,格式[X,Y] |
| `--pos2` | integer[] |  |  | x2,y2 坐标2,格式[X,Y] |
| `--query` | string[] |  |  | 节点1,格式[key=value,key=value,...] |
| `--query2` | string[] |  |  | 节点2,格式[key=value,key=value,...] |
| `--range1` | integer[] |  |  | [x1,y1,x2,y2] 区域1,格式[X1,Y1,X2,Y2] |
| `--range2` | integer[] |  |  | [x1,y1,x2,y2] 区域2,格式[X1,Y1,X2,Y2] |

> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，会与上面的 flag 合并。

## 示例

```bash
pmock-cli phone action swipe --code <string>
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

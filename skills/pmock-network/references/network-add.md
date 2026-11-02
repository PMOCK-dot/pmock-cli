# network add

批量新增代理资源

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli network add` |
| HTTP | `POST /open/api/network/add` |
| 风险 | `write`（修改数据） |
| 异步 | 否 |

## 示例

```bash
pmock-cli network add --json '[{"scheme":"socks5","server":"1.2.3.4","port":1080}]'
```

## 请求体

> schema：`NetworkAddBatchReqVO`

> 下面这些字段没有对应 flag，只能用 `--json` 传入。

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `data` | []NetworkAddReqVO | 是 | 网络信息集合 |

> 这些必填字段是嵌套结构，只能用 `--json` 传入：`data`

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RVoid`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 编码 |
| `data` | null | 数据 |
| `msg` | string | 信息 |
| `status` | string | 状态,代表本次是即时完成或者等待中的状态,finish或者pending |
| `traceId` | string | 链路id |

---

[pmock-network](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

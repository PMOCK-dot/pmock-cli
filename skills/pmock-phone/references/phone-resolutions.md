# phone resolutions

获取手机分辨率列表

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli phone resolutions` |
| HTTP | `GET /open/api/phone/resolution/list` |
| 风险 | `read`（只读） |
| 异步 | 否 |

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--spu-id` | integer |  |  | 商品 SPU ID |
| `--status` | integer |  |  | 手机状态：1=正常，2=重启，4=关机，41=开机，6=重置，7=一键新机 |

## 示例

```bash
pmock-cli phone resolutions --help
```

## 响应字段

> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，分页接口把 `total`/`rows` 放在顶层。

> schema：`RListPhoneResolutionListResVO`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | integer | 状态码 |
| `data` | []PhoneResolutionListResVO | 请求返回数据 |
| `msg` | string | 请求返回消息 |
| `status` | string | 状态，表示任务是立即完成还是处于等待状态：finish=已完成，pending=等待中 |
| `traceId` | string | 请求链路 ID |

---

[pmock-phone](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

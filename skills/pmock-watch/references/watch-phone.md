# watch phone

跟踪开机、关机、重启、重置、一键新机这类手机操作

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli watch phone` |
| 风险 | `read`（只读） |

## 参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `--id` | integer | 是 | 云手机 ID |
| `--op` | string | 是 | `powerOn` `powerOff` `restart` `reset` `renew` |

其余参数与 `watch rpa` 相同。

## 行为

接口返回两个轴：`taskStatus`（任务）和 `phoneStatus`（手机本身）。
只有手机真的到了目标状态才算完成——`powerOn` 要求 `phoneStatus == 1`，`powerOff` 要求 `phoneStatus == 4`。
任务状态说成功、但手机还在启动中，就继续报 `running`，不会提前报完成。

## 示例

```bash
pmock-cli phone start --id 123
pmock-cli watch phone --id 123 --op powerOn --human
```

## 注意

- 重启 / 重置 / 一键新机只看 `taskStatus`，手机可能还在过渡态，这是正常的；
- 这个命令只观察，不会替你开机；要真正执行仍然是 `phone start` 等命令。

---

[pmock-watch](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

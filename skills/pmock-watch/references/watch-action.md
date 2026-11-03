# watch action

跟踪一次 UI 动作任务（点击、滑动、输入……）

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli watch action` |
| 风险 | `read`（只读） |

## 参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `--code` | string | 是 | 手机编号（`phone list` 里的 `code`） |
| `--task-id` | integer | 是 | 动作任务 ID |

其余参数与 `watch rpa` 相同。

## 行为

动作通常一秒内结束，这个目标主要是为了让所有能力走同一套状态，而不是各发明一套。
接口文档写的是 `-1/0/1/2`，实际返回 `1/2/3`，两者都认。

## 示例

```bash
pmock-cli watch action --code 133337 --task-id 456 --once
```

## 注意

- 秒级动作没必要播报，`--once` 拿一次快照通常就够了；
- 连续动作编排（点一下、等一下、再输入）应当作为**一个业务步骤**报给客户，而不是逐个动作播报。

---

[pmock-watch](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

# watch log

回放本地台账里记录过的进度事件

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli watch log` |
| 风险 | `read`（只读） |

## 参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `--kind` | string | 是 | `rpa` `phone` `action` |
| `--task-id` | integer | 是 | 任务 ID |
| `--state-dir` | string | 否 | 台账目录，默认 `~/.pmock-cli/watch` |

## 行为

把 `watch --record` 写下的 JSONL 原样打出来，事件格式和实时流完全一致，可以喂给同一个消费方。

## 示例

```bash
pmock-cli watch rpa --task-id 789 --record
pmock-cli watch log --kind rpa --task-id 789 | jq -r .say
```

## 注意

- 台账是事后回答"上次那批检查怎么结束的"的依据，不必为了看结果重跑一遍；
- 没有 `--record` 就没有台账，`watch log` 会明确告诉你没有记录。

---

[pmock-watch](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

# watch rpa

跟踪一个 RPA 任务，把状态变化翻成客户能看懂的一句话

## 概览

| 项 | 值 |
| --- | --- |
| 命令 | `pmock-cli watch rpa` |
| 风险 | `read`（只读） |

## 状态

| 状态 | 含义 | 客户该做什么 |
| --- | --- | --- |
| `running` | 任务在推进 | 等，不需要动作 |
| `needs_you` | 卡在只有客户能做的决定上（费用、授权、登录、选项） | 回一句话 |
| `done` | 有结果可看，可能带"需要关注"的项 | 看结论，决定是否处理 |
| `failed` | 没有可用结果 | 决定重试、换方案，还是算了 |

四个状态覆盖所有能力。加一种能力是加适配器，不是加状态。

## 证据阶梯

每行事件都带 `evidence`，说明后端到底给了什么；`say` 是给客户的那句话，措辞永远不超过证据：

| evidence | 后端给了什么 | 措辞上限 |
| --- | --- | --- |
| `structured` | 结构化结果（比如检查项清单） | 可以报数字 |
| `step` | 业务步骤或失败节点 | 可以报"第几步" |
| `status` | 只有状态码 | 只能说"正在处理" |
| `none` | 心跳：这次没拿到新信号 | 只能说"暂时没有新进度" |

没有步骤信号时 `say` 不会出现"第 N 步"，没有结构化结果时不会报数字。
可以直接把 `say` 转给客户，也可以只把它当作措辞下限。

## 事件格式

每行一个 JSON 对象，关键字段：`seq` `at` `kind` `state` `step` `total` `label` `etaSeconds` `evidence` `reason` `say`。
`--human` 直接打印 `say`，适合给人看；默认 JSON 适合给程序看。

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| `--task-id` | integer | 是 |  | `rpa task add` 返回的任务 ID |
| `--interval` | duration | 否 | `3s` | 轮询间隔 |
| `--heartbeat` | duration | 否 | `30s` | 超过这么久没有新信号就补一条"没有新进度"；负数关闭 |
| `--deadline` | duration | 否 | `20m` | 到点结束，并且先发一条 `failed` 事件 |
| `--once` | bool | 否 | false | 只取一次快照就退出 |
| `--human` | bool | 否 | false | 打印 `say` 而不是 JSON |
| `--record` | bool | 否 | false | 把每条事件追加到本地台账 |
| `--state-dir` | string | 否 | `~/.pmock-cli/watch` | 台账目录，也可用 `PMOCK_CLI_WATCH_DIR` |

## 行为

- 只在变化时输出：状态、步骤、标签或原因变了才发一行，十分钟的任务不会刷出两百行"还在跑"；
- 超过 `--heartbeat` 没有变化，补一条 `evidence: none` 的事件，明确说"没有新进度"，不让流静默；
- 到 `--deadline` 先发一条 `failed` 事件（原因"…内没跑完"）再退出，流不会无声中断；
- RPA 只有拿到 `nodeList` 才报步骤；拿不到就只有状态；再拿不到就只有心跳。

## 示例

```bash
# 下发后立刻开始跟踪，人看的版本
pmock-cli rpa task add --template-id 456 --name daily --json '{"taskList":[{"phoneId":123}]}'
pmock-cli watch rpa --task-id 789 --human

# 给 Agent 解析，并把事件记到台账，事后可回放
pmock-cli watch rpa --task-id 789 --record
pmock-cli watch log --kind rpa --task-id 789
```

## 退出码

0 = 完成 / 待你决定 / `--once` 取到快照；1 = failed；2 = 参数错误；3 = 超时或接口失败。

## 注意

- 失败但后端没给原因时，`say` 只说"这次没做成"，不会编一个原因；这时用 `rpa task details` 查 `failDescription`/`failActionName`；
- `needs_you` 不是终态，但 watch 会停下：只能由人回一句话，再由 Agent 重新开始跟踪。

---

[pmock-watch](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)

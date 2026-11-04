package skillgen

import "github.com/haidaisiganbei/pmock-cli/internal/spec"

// phoneFileUploadPage documents the hand-written upload command. It cannot be
// generated: the CLI merges three OpenAPI operations (upload-init, part PUT,
// upload-push) into one idempotent, resumable command.
func phoneFileUploadPage() ManualPage {
	return ManualPage{
		Skill:   "pmock-phone",
		Command: "phone file upload",
		Summary: "上传本地文件或 APK 到云手机（分片、可续传）",
		Risk:    spec.RiskWrite,
		Body: `## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| ` + "`--id`" + ` | integer | 是 |  | 目标云手机 ID |
| ` + "`--file`" + ` | string | 是 |  | 本地文件路径 |
| ` + "`--type`" + ` | string | 否 | ` + "`file`" + ` | ` + "`file`" + ` 普通文件；` + "`app`" + ` 表示 APK |
| ` + "`--dest`" + ` | string | 否 | ` + "`/sdcard/Download`" + ` | 手机上的目标目录 |
| ` + "`--chunk-size-mb`" + ` | integer | 否 | ` + "`30`" + ` | 分片大小，30–300 MiB |
| ` + "`--md5`" + ` | string | 否 |  | 预计算的 MD5，跳过本地计算 |

## 行为

1. 本地计算 MD5（除非传 ` + "`--md5`" + `），用 ` + "`upload-init`" + ` 换取分片清单；
2. 只上传服务端缺失的分片，因此中断后重跑同一命令即可续传；
3. 全部分片就位后 ` + "`upload-push`" + ` 合并并推送到手机。

` + "`--type app`" + ` 时返回的 ` + "`data.envFileId`" + ` 就是 ` + "`app install --file-id`" + ` 要的值；
安装还需要 ` + "`--id`" + `（目标机器）和 APK 元信息 ` + "`--apk-name`" + `、` + "`--apk-pkg`" + `、
` + "`--apk-version-code`" + `、` + "`--apk-version-name`" + `，见 [app install](../../pmock-app/references/app-install.md)。

## 示例

` + "```bash" + `
# 普通文件
pmock-cli phone file upload --id 123 --file ./video.mp4

# APK，随后安装（--file-id 取上一步的 data.envFileId）
pmock-cli phone file upload --id 123 --file ./app.apk --type app
pmock-cli app install --id 123 --file-id <envFileId> \
  --apk-name TikTok --apk-pkg com.zhiliaoapp.musically \
  --apk-version-code 2024101 --apk-version-name 41.2.1

# 指定目录，并先看计划不实际传输
pmock-cli phone file upload --id 123 --file ./video.mp4 --dest /sdcard/Movies --dry-run
` + "```" + `

## 注意

- 大文件按 45 分钟预算上传，不受全局 ` + "`--timeout`" + ` 限制；
- 进度输出到 stderr，仅在 TTY 或 ` + "`--verbose`" + ` 时出现，不影响 JSON 输出；
- 单文件 MD5 不变时重复执行是幂等的。`,
	}
}

// watchOverview explains the four-state contract every watch target shares.
//
// It is written once because the contract is the product: a new capability adds
// an adapter, never a state.
const watchOverview = `## 状态

| 状态 | 含义 | 客户该做什么 |
| --- | --- | --- |
| ` + "`running`" + ` | 任务在推进 | 等，不需要动作 |
| ` + "`needs_you`" + ` | 卡在只有客户能做的决定上（费用、授权、登录、选项） | 回一句话 |
| ` + "`done`" + ` | 有结果可看，可能带"需要关注"的项 | 看结论，决定是否处理 |
| ` + "`failed`" + ` | 没有可用结果 | 决定重试、换方案，还是算了 |

四个状态覆盖所有能力。加一种能力是加适配器，不是加状态。

## 证据阶梯

每行事件都带 ` + "`evidence`" + `，说明后端到底给了什么；` + "`say`" + ` 是给客户的那句话，措辞永远不超过证据：

| evidence | 后端给了什么 | 措辞上限 |
| --- | --- | --- |
| ` + "`structured`" + ` | 结构化结果（比如检查项清单） | 可以报数字 |
| ` + "`step`" + ` | 业务步骤或失败节点 | 可以报"第几步" |
| ` + "`status`" + ` | 只有状态码 | 只能说"正在处理" |
| ` + "`none`" + ` | 心跳：这次没拿到新信号 | 只能说"暂时没有新进度" |

没有步骤信号时 ` + "`say`" + ` 不会出现"第 N 步"，没有结构化结果时不会报数字。
可以直接把 ` + "`say`" + ` 转给客户，也可以只把它当作措辞下限。

## 事件格式

每行一个 JSON 对象，关键字段：` + "`seq` `at` `kind` `state` `step` `total` `label` `etaSeconds` `evidence` `reason` `say`" + `。
` + "`--human`" + ` 直接打印 ` + "`say`" + `，适合给人看；默认 JSON 适合给程序看。

## 在 Agent 里怎么播报

聊天界面里的进度不是平台推出来的，是 Agent 每次工具调用之间说出来的。所以不要直接跑阻塞的
` + "`watch`" + `（它要等到任务结束才返回，中间用户什么都看不到），用快照加长轮询：

1. 下发任务后先取一次快照：` + "`watch rpa --task-id N --once`" + `，立刻说一句"已接单"；
2. 之后每次用 ` + "`--observe 20s`" + `：有变化立刻返回，没变化最多等 20 秒；
3. ` + "`state`" + ` 或 ` + "`say`" + ` 变了就播报，没变就什么都不用说（不要重复同一句话）；
4. 停在 ` + "`done`" + ` / ` + "`failed`" + ` / ` + "`needs_you`" + `：前两个给结论，最后一个把决定交回客户。

MCP 里对应 ` + "`pmock_watch`" + `（action 取 ` + "`rpa` `phone` `action`" + `），默认就是 20 秒长轮询，
返回一行事件，失败任务也按正常结果返回，不是工具报错。`

// watchPages documents the hand-written watch command tree.
func watchPages() []ManualPage {
	return []ManualPage{
		{
			Skill:   "pmock-watch",
			Command: "watch rpa",
			Summary: "跟踪一个 RPA 任务，把状态变化翻成客户能看懂的一句话",
			Risk:    spec.RiskRead,
			Body: watchOverview + `

## 参数

| 参数 | 类型 | 必填 | 默认 | 说明 |
| --- | --- | --- | --- | --- |
| ` + "`--task-id`" + ` | integer | 是 |  | ` + "`rpa task add`" + ` 返回的任务 ID |
| ` + "`--interval`" + ` | duration | 否 | ` + "`3s`" + ` | 轮询间隔 |
| ` + "`--heartbeat`" + ` | duration | 否 | ` + "`30s`" + ` | 超过这么久没有新信号就补一条"没有新进度"；负数关闭 |
| ` + "`--deadline`" + ` | duration | 否 | ` + "`20m`" + ` | 到点结束，并且先发一条 ` + "`failed`" + ` 事件 |
| ` + "`--once`" + ` | bool | 否 | false | 只取一次快照就退出 |
| ` + "`--human`" + ` | bool | 否 | false | 打印 ` + "`say`" + ` 而不是 JSON |
| ` + "`--record`" + ` | bool | 否 | false | 把每条事件追加到本地台账 |
| ` + "`--state-dir`" + ` | string | 否 | ` + "`~/.pmock-cli/watch`" + ` | 台账目录，也可用 ` + "`PMOCK_CLI_WATCH_DIR`" + ` |

## 行为

- 只在变化时输出：状态、步骤、标签或原因变了才发一行，十分钟的任务不会刷出两百行"还在跑"；
- 超过 ` + "`--heartbeat`" + ` 没有变化，补一条 ` + "`evidence: none`" + ` 的事件，明确说"没有新进度"，不让流静默；
- 到 ` + "`--deadline`" + ` 先发一条 ` + "`failed`" + ` 事件（原因"…内没跑完"）再退出，流不会无声中断；
- RPA 只有拿到 ` + "`nodeList`" + ` 才报步骤；拿不到就只有状态；再拿不到就只有心跳。

## 示例

` + "```bash" + `
# 下发后立刻开始跟踪，人看的版本
pmock-cli rpa task add --template-id 456 --name daily --json '{"taskList":[{"phoneId":123}]}'
pmock-cli watch rpa --task-id 789 --human

# 给 Agent 解析，并把事件记到台账，事后可回放
pmock-cli watch rpa --task-id 789 --record
pmock-cli watch log --kind rpa --task-id 789
` + "```" + `

## 退出码

0 = 完成 / 待你决定 / ` + "`--once`" + ` 取到快照；1 = failed；2 = 参数错误；3 = 超时或接口失败。

## 注意

- 失败但后端没给原因时，` + "`say`" + ` 只说"这次没做成"，不会编一个原因；这时用 ` + "`rpa task details`" + ` 查 ` + "`failDescription`/`failActionName`" + `；
- ` + "`needs_you`" + ` 不是终态，但 watch 会停下：只能由人回一句话，再由 Agent 重新开始跟踪。`,
		},
		{
			Skill:   "pmock-watch",
			Command: "watch phone",
			Summary: "跟踪开机、关机、重启、重置、一键新机这类手机操作",
			Risk:    spec.RiskRead,
			Body: `## 参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| ` + "`--id`" + ` | integer | 是 | 云手机 ID |
| ` + "`--op`" + ` | string | 是 | ` + "`powerOn` `powerOff` `restart` `reset` `renew`" + ` |

其余参数与 ` + "`watch rpa`" + ` 相同。

## 行为

接口返回两个轴：` + "`taskStatus`" + `（任务）和 ` + "`phoneStatus`" + `（手机本身）。
只有手机真的到了目标状态才算完成——` + "`powerOn`" + ` 要求 ` + "`phoneStatus == 1`" + `，` + "`powerOff`" + ` 要求 ` + "`phoneStatus == 4`" + `。
任务状态说成功、但手机还在启动中，就继续报 ` + "`running`" + `，不会提前报完成。

## 示例

` + "```bash" + `
pmock-cli phone start --id 123
pmock-cli watch phone --id 123 --op powerOn --human
` + "```" + `

## 注意

- 重启 / 重置 / 一键新机只看 ` + "`taskStatus`" + `，手机可能还在过渡态，这是正常的；
- 这个命令只观察，不会替你开机；要真正执行仍然是 ` + "`phone start`" + ` 等命令。`,
		},
		{
			Skill:   "pmock-watch",
			Command: "watch action",
			Summary: "跟踪一次 UI 动作任务（点击、滑动、输入……）",
			Risk:    spec.RiskRead,
			Body: `## 参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| ` + "`--code`" + ` | string | 是 | 手机编号（` + "`phone list`" + ` 里的 ` + "`code`" + `） |
| ` + "`--task-id`" + ` | integer | 是 | 动作任务 ID |

其余参数与 ` + "`watch rpa`" + ` 相同。

## 行为

动作通常一秒内结束，这个目标主要是为了让所有能力走同一套状态，而不是各发明一套。
接口文档写的是 ` + "`-1/0/1/2`" + `，实际返回 ` + "`1/2/3`" + `，两者都认。

## 示例

` + "```bash" + `
pmock-cli watch action --code 133337 --task-id 456 --once
` + "```" + `

## 注意

- 秒级动作没必要播报，` + "`--once`" + ` 拿一次快照通常就够了；
- 连续动作编排（点一下、等一下、再输入）应当作为**一个业务步骤**报给客户，而不是逐个动作播报。`,
		},
		{
			Skill:   "pmock-watch",
			Command: "watch log",
			Summary: "回放本地台账里记录过的进度事件",
			Risk:    spec.RiskRead,
			Body: `## 参数

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| ` + "`--kind`" + ` | string | 是 | ` + "`rpa` `phone` `action`" + ` |
| ` + "`--task-id`" + ` | integer | 是 | 任务 ID |
| ` + "`--state-dir`" + ` | string | 否 | 台账目录，默认 ` + "`~/.pmock-cli/watch`" + ` |

## 行为

把 ` + "`watch --record`" + ` 写下的 JSONL 原样打出来，事件格式和实时流完全一致，可以喂给同一个消费方。

## 示例

` + "```bash" + `
pmock-cli watch rpa --task-id 789 --record
pmock-cli watch log --kind rpa --task-id 789 | jq -r .say
` + "```" + `

## 注意

- 台账是事后回答"上次那批检查怎么结束的"的依据，不必为了看结果重跑一遍；
- 没有 ` + "`--record`" + ` 就没有台账，` + "`watch log`" + ` 会明确告诉你没有记录。`,
		},
	}
}

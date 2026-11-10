---
name: pmock-shared
version: 1.0.0
description: "PMOCK CLI shared rules: installation, AK/SK authentication, response envelope, exit codes, risk confirmation, output formats and async tasks. Use when starting with pmock-cli, setting up credentials, or reading its output; read this first before any other PMOCK skill."
metadata:
  requires:
    bins: ["pmock-cli"]
  cliHelp: "pmock-cli --help"
---

# pmock-cli 通用规则

所有 PMOCK 技能共用这一篇。命令本身按领域拆在 `pmock-phone`、`pmock-rpa`、`pmock-app`、`pmock-network`、`pmock-platform` 里。

## 安装与配置

```bash
# 安装（macOS / Linux）
npm install -g pmock-cli@latest

# 写入 AK/SK（交互式，secret 不回显）
pmock-cli config init

# 非交互
pmock-cli config init --access-key AK... --secret-key SK...
```

凭据保存在 `~/.pmock-cli/config.json`（0600）。CI 里优先用环境变量，优先级高于配置文件：

| 变量 | 作用 |
| --- | --- |
| `PMOCK_ACCESS_KEY` | 访问密钥 AK |
| `PMOCK_SECRET_KEY` | 签名密钥 SK |
| `PMOCK_BASE_URL` | 覆盖 API 域名（默认 `https://user.pmock.com/api`） |
| `PMOCK_TIMEOUT` | 单请求超时 |
| `PMOCK_CLI_CONFIG` | 换一份配置文件（测试/多租户） |

## 先验证再执行

```bash
pmock-cli auth status   # 只调用 phone list 取 1 条，不写数据
pmock-cli doctor        # 排查配置、凭据、连通性
```

## 响应信封

所有接口返回同一层信封：

```json
{ "traceId": "...", "code": 0, "msg": "success", "data": { } }
```

分页接口（`phone list`、`app list`、`network list` 等）把分页放在顶层：

```json
{ "traceId": "...", "code": 0, "msg": "ok", "total": 120, "rows": [ ] }
```

判断成功看 `code == 0`（少量接口用 `200`）；`code != 0` 时 `msg` 是失败原因，
排查问题时把 `traceId` 一并带上。

## 退出码

| 退出码 | 含义 | 建议动作 |
| --- | --- | --- |
| 0 | 成功 | — |
| 1 | 业务失败（信封 `code != 0`） | 读 stdout 的 `code`/`msg` 决定是否重试 |
| 2 | 参数或配置错误，请求未发出 | 修正参数，不要原样重试 |
| 3 | 运行时失败（网络、超时、5xx） | 可退避重试 |

退出码 1 时 stdout 仍然是完整信封，`--allow-business-error` 可让它按 0 退出。

## 风险分级与确认

每个命令都有风险等级，技能索引里直接标出：

| 风险 | 含义 | 处理方式 |
| --- | --- | --- |
| `read` | 只读 | 直接执行 |
| `write` | 修改数据 | 明确目标 ID 后执行 |
| `costly` | 会产生费用（创建、续费、升级） | 先向用户确认金额与数量 |
| `destructive` | 不可逆（删除、重置） | 必须 `--yes`，且先确认对象 |

`--dry-run` 只打印将发送的请求，不发出去；改动敏感资源前先用它核对。

## 输出格式

```bash
pmock-cli phone list --format json     # 默认，缩进 JSON
pmock-cli phone list --format pretty   # 带状态摘要的彩色输出
pmock-cli phone list --format table    # 终端表格（不适合再解析）
pmock-cli phone list --compact         # 单行 JSON，省 token
```

给程序或 Agent 解析时用 `--format json`，不要用 `table`。

## 异步任务

部分命令（开机、重置、重启、RPA 下发）返回任务而非最终结果。加 `--wait` 让它轮询到终态：

```bash
pmock-cli phone start --ids 123 --wait
pmock-cli phone reset --id 123 --wait --wait-timeout 10m --poll-interval 3s
```

不加 `--wait` 时用 `pmock-cli task list` 或对应命令的 `task-status` 查进度。

**长任务要给用户播报，不能静默。** 在聊天界面里跑阻塞的 `watch` 没有用——它要等到任务结束才返回，
中间用户什么都看不到。正确做法是先取快照、再长轮询：

```bash
pmock-cli watch rpa --task-id 789 --once          # 下发后立刻说"已接单"
pmock-cli watch rpa --task-id 789 --observe 20s   # 有变化立刻返回，没变化最多等 20 秒
```

`state` 或 `say` 变了就说，没变就不说；停在 `done` / `failed` / `needs_you`。
MCP 侧对应 `pmock_watch`，默认就是 20 秒长轮询。详见 `pmock-watch` 技能。

## 两个易错点

1. **ID 是 int64**：`id`、`phoneId` 等字段超过 2^53 时用字符串传（`--id "1900000000000000123"`），
   否则 JSON 数字会丢精度。
2. **批量接口传逗号分隔**：`--ids 1,2,3`，不是重复传 `--ids`。

## 看屏与客户端窗口

- 看屏优先用 `phone screen`（截图 URL）或 MCP 的 `screenText` / `tapAndSee`；**不要**把
  `phone stream-token` 返回的 `remotePhoneExternal` URL 用浏览器打开。
- 这个网页和 PMOCK 客户端的远程桌面共用同一个远程会话（`forceKick: true`），两边同时开
  会互相踢下线，表现为「已在其他渠道/页面打开」「已断开」。
- 需要让用户看实时画面时，执行任意手机操作，让 PMOCK 客户端自动打开原生窗口；或让用户点
  客户端「Agent 正在操作」浮层的「查看屏幕」。同一台手机只保留一个画面。

## 发现命令

```bash
pmock-cli commands                     # 全部命令的机器可读目录（JSON）
pmock-cli docs phone start             # 单个命令的参数说明
pmock-cli docs phone list --markdown   # 完整 reference 页（示例 + 响应字段）
pmock-cli <命令> --help                # pflag 原始帮助
```

## 安全

- 不要在任何输出里回显 AK/SK；`config show` 会打码。
- 删除、重置、转让等不可逆操作，先复述对象再执行。
- 令牌只从环境变量或 `config.json` 读取，不要写进脚本参数。

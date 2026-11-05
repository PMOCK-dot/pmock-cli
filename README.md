# pmock-cli

![Release](https://github.com/PMOCK-dot/pmock-cli/actions/workflows/release.yml/badge.svg)
![npm](https://img.shields.io/npm/v/pmock-cli.svg)
![Go](https://img.shields.io/badge/go-%3E%3D1.25-blue.svg)
![License: Apache-2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)

PMOCK 开放 API 的命令行客户端。用团队的 **Access Key / Secret Key** 直接管理云手机、
应用、文件、网络和 RPA 任务，不需要打开客户端，也不需要本机登录态。

## 为什么选择 pmock-cli？

- **全 API 覆盖** — 云手机、UI 自动化、文件、应用、代理、RPA、商品、团队、钱包账单
- **单一二进制，零依赖** — Windows / macOS / Linux 同一个可执行文件，AK/SK 直连
- **对脚本与 Agent 友好** — 默认 JSON，`--format table` 给人看， `--json` 覆盖自动化
- **开源** — Apache-2.0，命令、文档、Agent 技能出自同一份 OpenAPI 目录

## 功能概览

| 类别 | 能力 |
|---|---|
| 📱 云手机 | 列表、创建、启停、重启、重置、续费、截图、GPS、分组、标签、改名、转让 |
| 👆 UI 自动化 | 点击、双击、长按、滑动、拖动、输入、应用操作 |
| 📁 文件 | 扫描、上传（自动分片 + MD5）、下载 |
| 🧩 应用 | 列表、安装、升级、启停、团队应用库 |
| 🌐 网络 | 代理与代理组、连通性检测 |
| 🤖 RPA | 模板与定时任务 |
| 💰 账单 | 钱包余额、充值套餐、消费 / 充值流水 |
| 🏢 平台 | 商品与 SKU、品牌、标签、团队成员与授权 |

---

## 安装与快速开始

### 环境要求

- **npm 方式**：Node.js 18+，Windows（x64）、macOS、Linux 通用
- **源码方式**：Go 1.25+

### 快速开始（人类用户）

#### 安装

**方式一 — 从 npm 安装（推荐，Windows / macOS / Linux 通用）：**

```bash
npm install -g pmock-cli@latest
```

国内镜像如果还没同步，指定官方源：

```bash
npm install -g pmock-cli@latest --registry=https://registry.npmjs.org
```

**方式二 — 从源码安装：**

需要 Go 1.25+。

```bash
git clone https://github.com/PMOCK-dot/pmock-cli.git
cd pmock-cli
sudo make install     # 装到 /usr/local/bin；不想用 sudo 就 PREFIX=$HOME/.local make install
```

不想 clone 的话也可以直接：

```bash
go install github.com/PMOCK-dot/pmock-cli@latest
```

产物在 `$(go env GOPATH)/bin/pmock-cli`。

验证安装：

```bash
pmock-cli version
```

> Windows 上用 PowerShell 或 cmd 执行同一条 npm 命令；新开终端后 `where.exe pmock-cli`
> 应指向 `%APPDATA%\npm\pmock-cli.cmd`。如果提示「不是内部或外部命令」，说明 `%APPDATA%\npm`
> 不在 PATH 里，确认后重开终端；版本低于 0.1.12 请重装（那时还没有 Windows 平台包）。

#### 配置与使用

```bash
# 1. 配置 AK/SK（交互式，Secret 不回显）
pmock-cli config init

# 2. 检查配置和连通性
pmock-cli doctor

# 3. 查询设备
pmock-cli phone list --page-num 1 --page-size 100

# 4. 表格查看
pmock-cli phone list --format table
```

### 快速开始（AI Agent）

> 以下步骤面向 AI Agent。部分步骤需要用户提供 AK/SK，可从 PMOCK 客户端的团队设置页获取。

**第一步 — 安装**

```bash
npm install -g pmock-cli@latest
```

**第二步 — 配置 AK/SK**

```bash
pmock-cli config init --access-key "AK..." --secret-key "SK..."
```

也可以只用环境变量，不动配置文件：

```bash
export PMOCK_ACCESS_KEY="AK..."
export PMOCK_SECRET_KEY="SK..."
```

**第三步 — 验证**

```bash
pmock-cli doctor
```

`data.ok` 为 `true` 且 `checks` 里 `api` 通过，说明凭证、地址和网络都正常。

**第四步 — 安装 Skills（可选）**

把技能包装进 Agent，让它知道什么时候该调哪条命令：

```bash
npx skills@latest add PMOCK-dot/pmock-cli -g --all -y
```

---

## 配置与环境变量

`pmock-cli config init` 写入 `~/.pmock-cli/config.json`（Windows 为
`%USERPROFILE%\.pmock-cli\config.json`），macOS / Linux 上权限 `0600`。

CI、容器里用环境变量，优先级 **命令行参数 > 环境变量 > 配置文件 > 默认值**：

| 变量 | 用途 |
|---|---|
| `PMOCK_ACCESS_KEY` / `PMOCK_SECRET_KEY` | 团队凭证，必填 |
| `PMOCK_BASE_URL` | 接口地址，默认 `https://user.pmock.com/api` |
| `PMOCK_TIMEOUT` | 单次请求超时，默认 `30s` |
| `PMOCK_CLI_CONFIG` | 指向另一份配置文件，用于多团队 / 多套配置 |

Windows：PowerShell 用 `$env:PMOCK_ACCESS_KEY = "AK..."`，cmd 用 `set PMOCK_ACCESS_KEY=AK...`；
要长期保留用 `[Environment]::SetEnvironmentVariable("PMOCK_ACCESS_KEY", "AK...", "User")`。

---

## 常用命令

| 命令组 | 用途 |
|---|---|
| `phone` | 云手机列表、创建、启停、重启、重置、续费、截图、GPS |
| `phone action` | 点击、双击、长按、滑动、拖动、输入、应用操作 |
| `phone file` | 文件扫描、下载、上传 |
| `app` | 应用列表、安装、升级、启停 |
| `network` | 代理与网络检测 |
| `rpa` | RPA 模板与任务 |
| `product` | 商品、SKU、价格 |
| `team` | 团队成员 |
| `auth` / `config` | 鉴权与配置 |
| `doctor` | 诊断 |
| `mcp` | MCP Server |

### 云手机

```bash
# 列表
pmock-cli phone list --page-size 100
pmock-cli phone list --name-or-code demo
pmock-cli phone list --status 1
pmock-cli phone list --group-id 123

# 开关机 / 重启
pmock-cli phone start --id 123 --wait
pmock-cli phone stop --id 123 --wait
pmock-cli phone restart --id 123 --wait

# 截图（默认删除手机中的截图文件）
pmock-cli phone screen --id 123
pmock-cli phone screen --id 123 --is-remove=false

# 需要花钱或不可逆的操作
pmock-cli phone create --spu-id 1 --sku-id 2 --num 1
pmock-cli phone remove --id 123 --yes
```

### UI 自动化

动作类命令使用 **code**，不是 id：

```bash
pmock-cli phone action click --code A1 --pos-type 2 --pos 100,200 --wait
pmock-cli phone action input --code A1 --text "hello" --wait
pmock-cli phone action swipe --code A1 --pos1 100,200 --pos2 100,600
```

### 文件上传

```bash
pmock-cli phone file upload --id 123 --file ./video.mp4
pmock-cli phone file upload --id 123 --file ./app.apk --type app
```

- 自动计算 MD5
- 自动分片上传
- 默认推送到 `/sdcard/Download`
- `--dest` 可指定手机目录
- `--type app` 返回的 `envFileId` 可传给 `app install --file-id`

### 应用

```bash
pmock-cli app list --page-size 10 --apk-name tiktok
pmock-cli app install --id 123 --file-id 456 --apk-pkg com.example.app \
  --apk-name Example --apk-version-name 1.0 --apk-version-code 1
```

### RPA

```bash
pmock-cli rpa template list
pmock-cli rpa task list
pmock-cli rpa task add --template-id 9 --name daily --json @task.json --wait
```

### 分组、标签与改名

```bash
pmock-cli phone group list
pmock-cli phone group create --name 投放A
pmock-cli phone move-group --id 7 --phone-ids 1,2        # -1 表示移出分组

pmock-cli tag list
pmock-cli tag create --tag-name tiktok --color '#ff0000'
pmock-cli phone update-tags --phone-id 123 --mode append --tag-ids 1,2
pmock-cli phone batch-update-tags --phone-ids 1,2,3 --mode replace --tag-ids 4

pmock-cli phone update-name --id 123 --name demo-01
pmock-cli phone batch-update-remark --phone-ids 1,2 --mode update --remark "第一批"
```

标签的 `--mode` 取 `append`（追加）/ `replace`（替换）/ `clear`（清空）；备注的 `--mode` 取 `update` / `clear`。

### 代理

```bash
pmock-cli network list
pmock-cli network add --json '[{"scheme":"socks5","server":"1.2.3.4","port":1080}]'
pmock-cli network check --network-id 5
pmock-cli network group list
```

### 钱包与账单

```bash
pmock-cli user balance              # 查余额
pmock-cli order consume-record      # 消费流水
pmock-cli order recharge-record     # 充值流水
pmock-cli payment recharge-config   # 可充值套餐
pmock-cli payment recharge --team-id 1 --method wallet_pay \
  --recharge-type fixed --recharge-id 9      # 真实扣款，属 costly
```

### 转让手机（三步）

```bash
pmock-cli phone transfer send-code --phone-ids 1,2 --transfer-options 7 \
  --target-boss-email boss@team.com --target-team-name 新团队
pmock-cli phone transfer confirm --phone-ids 1,2 --transfer-options 7 \
  --verification-code 123456 --yes                                      # 不可逆
pmock-cli phone transfer cancel --phone-ids 1,2
```

更多命令：

```bash
pmock-cli --help
pmock-cli phone --help
pmock-cli docs phone start              # 纯文本参数说明
pmock-cli docs phone list --markdown    # Agent reference 页（示例 + 响应字段）
pmock-cli commands
```

---

## 输出

默认输出 JSON：

```json
{
  "traceId": "...",
  "code": 200,
  "msg": "查询成功",
  "data": {}
}
```

- `code` 为 `0` 或 `200` 表示成功
- 其他 code 表示失败，`msg` 是原因
- 分页接口统一输出为 `data.total` + `data.rows`

```json
{
  "data": {
    "total": 6,
    "rows": [ ... ]
  }
}
```

- `password`、`adbPassword`、`accessKey`、`token` 等字段会自动替换为 `****`

格式：

```bash
pmock-cli phone list --format json
pmock-cli phone list --format pretty
pmock-cli phone list --format table
```

退出码：

| 码 | 含义 |
|---|---|
| 0 | 成功 |
| 1 | 业务失败 |
| 2 | 参数 / 配置错误 |
| 3 | 网络、超时、鉴权失败 |

`--allow-business-error` 可以让业务失败也返回 0，适合只判断接口可达性的脚本。

---

## MCP（AI Agent 接入）

`pmock-cli` 可以作为 [MCP](https://modelcontextprotocol.io/) 服务器运行，把全部 CLI 能力
（云手机、UI 自动化、文件、应用、网络、RPA、商品、团队、账单）以 MCP 工具的形式暴露。

在 MCP 客户端配置里加入下面这段（Trae 的 `~/.trae/mcp/mcp.json`、Cursor 的 `~/.cursor/mcp.json`，
或 `claude mcp add`）：

```json
{
  "mcpServers": {
    "pmock": {
      "command": "npx",
      "args": ["-y", "pmock-cli@latest", "mcp"],
      "env": {
        "PMOCK_BASE_URL": "https://user.pmock.com/api",
        "PMOCK_ACCESS_KEY": "AK...",
        "PMOCK_SECRET_KEY": "SK..."
      }
    }
  }
}
```

`npx` 写法不需要先全局安装，Windows / macOS / Linux 都一样，版本也始终跟着 `@latest` 走。
本机已经装了二进制的话，把 `command` 直接指向它更省事（启动更快，不依赖网络）：

```json
{
  "mcpServers": {
    "pmock": {
      "command": "pmock-cli",
      "args": ["mcp"]
    }
  }
}
```

说明：

- `PMOCK_ACCESS_KEY` / `PMOCK_SECRET_KEY` 必填，可从 PMOCK 客户端获取
- `PMOCK_BASE_URL` 可选，默认 `https://user.pmock.com/api`
- 环境变量优先级高于配置文件（`~/.pmock-cli/config.json`）；`env` 块可以整个省略，先跑一次
`pmock-cli config init` —— MCP 服务器和 CLI 共享同一份配置
- Windows 上 npm 装出来的是 `pmock-cli.cmd` 垫片，个别客户端不经 shell、无法直接启动它
（Node 18.20.2 起禁止直接 spawn `.cmd`/`.bat`）。这种情况把 `command` 换成 `cmd`、`args` 换成
`["/c", "npx", "-y", "pmock-cli@latest", "mcp"]`，或者把 `command` 指向 `%APPDATA%\npm\node_modules\`
下真实的 `pmock-cli.exe`

工具和每个动作要的参数字段，用 `pmock_help` / `pmock_docs` 查；CLI 侧对应 `pmock-cli docs <command>`。

### 与 PMOCK 客户端联动

`pmock-cli mcp` 发起手机操作、或写命令成功后，如果本机 PMOCK 客户端正在运行，会把事件上报给客户端：

- 云手机列表出现 `AI` 角标，右下角显示「Agent 正在操作 <手机名> · 截图/点击/开机…」
- 同一台手机 5 分钟内首次被操作时自动打开远程桌面窗口（不抢焦点）
- 写操作完成后客户端列表/分组/标签自动刷新

不需要额外配置：客户端会在 `userData/mcp/runtime.json` 写入端口和 token，CLI 自动发现并携带
Bearer token 上报。客户端没启动、端口/token 过期或请求超时都会静默跳过，不影响工具调用。

| 变量 | 作用 |
|---|---|
| `PMOCK_CLIENT_LINK=0` | 关闭联动 |
| `PMOCK_CLIENT_RUNTIME` | 指定 `runtime.json`（或所在目录），用于非默认安装位置 |

需要客户端包含 `POST /agent-event` 端点；没有该端点的旧客户端会自动 404，CLI 忽略。

## 项目结构

```text
cmd/            命令定义（cobra），每个命令组一个文件
internal/
  client/       HTTP 客户端、请求签名、响应信封
  clientbridge/ 向本机 PMOCK 客户端上报 Agent 活动 / 数据变更
  config/       配置文件与环境变量解析
  output/       JSON / pretty / table 三种渲染
  progress/     长任务进度状态归一
  skillgen/     由命令目录生成 skills/ 与 docs/
  spec/         OpenAPI 目录解析
  auth/ build/ exitcode/ phonefile/
openapi/        pmock-openapi.yaml + overlay.yaml（命令的唯一事实来源）
skills/         给外部 Agent 的技能包（make skills 生成）
scripts/        install.sh、打包与 npm 发布
Formula/        Homebrew formula
```

---

## 链接

- 开放 API 文档：https://help.pmock.com/zh/api-docs/api-docs.html
- Releases：https://github.com/PMOCK-dot/pmock-cli/releases
- 问题反馈：https://github.com/PMOCK-dot/pmock-cli/issues

参与开发请见 [CONTRIBUTING.md](CONTRIBUTING.md)。

---

Apache-2.0，见 [LICENSE](LICENSE)。

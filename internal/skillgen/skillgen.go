// Package skillgen renders the committed skills/ tree from the command
// catalog.
//
// OpenAPI plus overlay.yaml stay the single source of truth: one reference
// page is generated per public command, so the agent-facing docs cannot drift
// from the CLI surface. Only skills/pmock-shared is hand-written, because it
// describes conventions no spec field carries (access keys, envelope, exit
// codes, risk confirmation).
package skillgen

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/haidaisiganbei/pmock-cli/internal/spec"
)

const (
	// DefaultBinary is the command name used in every generated example.
	DefaultBinary = "pmock-cli"
	// defaultVersion is the frontmatter version. It is fixed instead of
	// derived from build.Version so a release does not dirty the tree.
	defaultVersion = "1.0.0"
	// maxDescription is the description budget skill loaders enforce.
	maxDescription = 1024
)

var skillNameRe = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// File is one rendered file, relative to the skills root and slash-separated.
type File struct {
	Path    string
	Content []byte
}

// Skill is one generated skill directory.
type Skill struct {
	Name        string     // directory name and frontmatter name
	Title       string     // heading shown in SKILL.md
	Description string     // frontmatter description, at most 1024 characters
	Intent      string     // hand-written paragraph under the heading
	CliHelp     string     // command whose help the skill mirrors
	Prefixes    [][]string // command-path prefixes this skill covers
}

// ManualPage is a page the OpenAPI document cannot describe (a command written
// by hand). The generator still owns the file so pruning stays safe.
type ManualPage struct {
	Skill   string
	Command string // full command path, e.g. "phone file upload"
	Summary string
	Risk    spec.Risk
	Body    string // markdown below the generated heading
}

// Options configures Render.
type Options struct {
	Binary  string // defaults to DefaultBinary
	Version string // frontmatter version, defaults to "1.0.0"
	Skills  []Skill
	Manual  []ManualPage
}

// Render produces every file of the skills tree, sorted by path. The result is
// deterministic: inputs are sorted and nothing depends on the clock.
func Render(cat *spec.Catalog, opts Options) ([]File, error) {
	if cat == nil || cat.Doc == nil {
		return nil, errors.New("skillgen: catalog must carry the OpenAPI document")
	}
	if len(opts.Skills) == 0 {
		return nil, errors.New("skillgen: no skills configured")
	}
	if opts.Binary == "" {
		opts.Binary = DefaultBinary
	}
	if opts.Version == "" {
		opts.Version = defaultVersion
	}

	var files []File
	seen := map[string]bool{}
	add := func(path string, content []byte) error {
		if seen[path] {
			return fmt.Errorf("skillgen: %s rendered twice", path)
		}
		seen[path] = true
		files = append(files, File{Path: path, Content: content})
		return nil
	}

	for _, sk := range opts.Skills {
		if err := sk.validate(); err != nil {
			return nil, err
		}
		cmds := commandsFor(cat, sk)
		manual := manualFor(opts.Manual, sk.Name)
		if len(cmds) == 0 && len(manual) == 0 {
			return nil, fmt.Errorf("skillgen: skill %q matches no command", sk.Name)
		}

		entries := indexEntries(cmds, manual)
		if err := add(sk.Name+"/SKILL.md", renderSkill(sk, entries, opts)); err != nil {
			return nil, err
		}
		for _, c := range cmds {
			path := pagePath(sk.Name, c.FullPath())
			if err := add(path, renderCommand(sk, c, cat.Doc, opts)); err != nil {
				return nil, err
			}
		}
		for _, m := range manual {
			if err := add(pagePath(sk.Name, m.Command), renderManualPage(sk, m)); err != nil {
				return nil, err
			}
		}
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

// DefaultSkills lists the skills this release generates. The phone skill is
// first because it covers the largest part of the API surface; the remaining
// skills follow the product areas a task usually starts from.
func DefaultSkills() []Skill {
	return []Skill{
		{
			Name:  "pmock-phone",
			Title: "PMOCK 云手机（phone）",
			Description: "PMOCK cloud phone management: list, create, start, stop, restart, reset, " +
				"renew, screenshot, GPS, ADB, rename, retag, change plan, groups, ownership transfer, " +
				"files and UI automation actions. Use when operating PMOCK cloud phones.",
			Intent: "云手机是 PMOCK 的核心资产：创建与续费、开关机、截图与 GPS、ADB 与网络配置、" +
				"改名打标签、分组与转让、文件上传，以及点击/滑动/输入等 UI 自动化。" +
				"换套餐（update-sku）和转让确认（transfer confirm）影响计费或归属，先确认再执行。",
			CliHelp:  "pmock-cli phone --help",
			Prefixes: [][]string{{"phone"}},
		},
		{
			Name:  "pmock-rpa",
			Title: "PMOCK RPA（rpa）",
			Description: "PMOCK RPA management: flow templates (market, team list, collect, delete) " +
				"and scheduled tasks (add, list, get, details, update, stop, retry, delete). " +
				"Use when automating repeated app workflows on PMOCK cloud phones.",
			Intent: "RPA 分两层：`rpa template` 管理流程模板（从市场收集、删除），" +
				"`rpa task` 管理按模板下发的定时任务。下发、停止、重试都作用于任务 ID，不是手机 ID。",
			CliHelp:  "pmock-cli rpa --help",
			Prefixes: [][]string{{"rpa"}},
		},
		{
			Name:  "pmock-app",
			Title: "PMOCK 应用（app）",
			Description: "PMOCK application management: list installed apps and categories, install " +
				"an app or APK, check upgrades, operate running apps, and maintain the shared team " +
				"app library. Use when installing or controlling apps on PMOCK cloud phones.",
			Intent: "应用管理分两层：`app *` 操作单台机器上的包（查询、安装、升级、启停、卸载），" +
				"`app team-app *` 维护团队共享的应用库与版本。安装前先确认目标机器的 `phoneId`。",
			CliHelp:  "pmock-cli app --help",
			Prefixes: [][]string{{"app"}},
		},
		{
			Name:  "pmock-network",
			Title: "PMOCK 网络（network）",
			Description: "PMOCK network management: add, update, remove, inspect and check proxies, " +
				"organize them in proxy groups, and set the network of a cloud phone. Use when " +
				"configuring proxies or switching the network of PMOCK cloud phones.",
			Intent: "网络是云手机风控的基础：先 `network list` 拿到代理 ID，" +
				"`network detect` 验证可用性，再用 `phone network-set` 绑定到机器。",
			CliHelp:  "pmock-cli network --help",
			Prefixes: [][]string{{"network"}},
		},
		{
			Name:  "pmock-platform",
			Title: "PMOCK 平台数据（product / brand / team / task）",
			Description: "PMOCK platform data: products and SKUs with pricing, changeable plans, phone " +
				"brands, tags, team members and phone authorization, and background task records. Use " +
				"when looking up prices, brands, tags, team access or task status.",
			Intent: "下单前查商品与价格（`product`）、创建机器时选品牌机型（`brand`）、" +
				"维护标签（`tag`）、把机器授权给同事（`team`）、排查异步任务结果（`task`）。",
			CliHelp:  "pmock-cli product --help",
			Prefixes: [][]string{{"product"}, {"brand"}, {"team"}, {"task"}, {"tag"}},
		},
		{
			Name:  "pmock-billing",
			Title: "PMOCK 钱包与账单（user / payment / order）",
			Description: "PMOCK billing: wallet balance, recharge configuration and orders, wallet " +
				"top-up, and consumption or recharge records. Use when checking balance, recharging " +
				"the account or auditing spend.",
			Intent: "钱相关的三组命令：`user balance` 查余额、`payment *` 充值（`payment recharge` " +
				"会真实扣款，属 costly，必须先跟用户确认金额）、`order *` 查消费与充值流水。",
			CliHelp:  "pmock-cli payment --help",
			Prefixes: [][]string{{"user"}, {"payment"}, {"order"}},
		},
	}
}

// DefaultManual lists hand-written pages that belong to generated skills.
func DefaultManual() []ManualPage {
	return []ManualPage{phoneFileUploadPage()}
}

// validate rejects skills the agent loaders would refuse.
func (s Skill) validate() error {
	if !skillNameRe.MatchString(s.Name) {
		return fmt.Errorf("skillgen: invalid skill name %q", s.Name)
	}
	desc := strings.TrimSpace(s.Description)
	if desc == "" {
		return fmt.Errorf("skillgen: skill %q has no description", s.Name)
	}
	if len([]rune(desc)) > maxDescription {
		return fmt.Errorf("skillgen: skill %q description is %d characters (max %d)", s.Name, len([]rune(desc)), maxDescription)
	}
	if strings.TrimSpace(s.CliHelp) == "" {
		return fmt.Errorf("skillgen: skill %q has no cliHelp", s.Name)
	}
	if len(s.Prefixes) == 0 {
		return fmt.Errorf("skillgen: skill %q has no command prefixes", s.Name)
	}
	return nil
}

// covers reports whether a command path belongs to the skill.
func (s Skill) covers(parts []string) bool {
	for _, prefix := range s.Prefixes {
		if len(parts) < len(prefix) {
			continue
		}
		match := true
		for i := range prefix {
			if parts[i] != prefix[i] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// commandsFor returns the public commands covered by the skill, sorted.
func commandsFor(cat *spec.Catalog, sk Skill) []*spec.Command {
	var out []*spec.Command
	for _, c := range cat.Commands {
		if c.Internal || !sk.covers(c.PathParts) {
			continue
		}
		out = append(out, c)
	}
	return out
}

func manualFor(pages []ManualPage, skill string) []ManualPage {
	var out []ManualPage
	for _, m := range pages {
		if m.Skill == skill {
			out = append(out, m)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Command < out[j].Command })
	return out
}

// indexEntry is one row of the SKILL.md command index.
type indexEntry struct {
	Command string
	Summary string
	Risk    spec.Risk
	Link    string
	Group   string
}

func indexEntries(cmds []*spec.Command, manual []ManualPage) []indexEntry {
	out := make([]indexEntry, 0, len(cmds)+len(manual))
	for _, c := range cmds {
		out = append(out, indexEntry{
			Command: c.FullPath(),
			Summary: c.Summary,
			Risk:    c.Risk,
			Link:    "references/" + slug(c.PathParts) + ".md",
			Group:   groupOf(c.PathParts),
		})
	}
	for _, m := range manual {
		parts := strings.Fields(m.Command)
		out = append(out, indexEntry{
			Command: m.Command,
			Summary: m.Summary,
			Risk:    m.Risk,
			Link:    "references/" + slug(parts) + ".md",
			Group:   groupOf(parts),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Command < out[j].Command })
	return out
}

// pagePath is the reference page of a command, relative to the skills root.
func pagePath(skill, command string) string {
	return skill + "/references/" + slug(strings.Fields(command)) + ".md"
}

// groupOf buckets a command under its first segment, or its first two segments
// when it has a sub-group (phone action click -> "phone action").
func groupOf(parts []string) string {
	if len(parts) >= 3 {
		return parts[0] + " " + parts[1]
	}
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func slug(parts []string) string {
	return strings.ToLower(strings.Join(parts, "-"))
}

package skillgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PMOCK-dot/pmock-cli/internal/spec"
)

// maxResponseFields caps the response tables so one wide schema cannot
// dominate a page; the remainder is counted and left to the OpenAPI document.
const maxResponseFields = 24

func renderSkill(sk Skill, entries []indexEntry, opts Options) []byte {
	var b strings.Builder
	writeFrontmatter(&b, sk, opts)

	fmt.Fprintf(&b, "# %s\n\n", sk.Title)
	if intent := oneLine(sk.Intent); intent != "" {
		b.WriteString(intent + "\n\n")
	}
	b.WriteString("**先读 [`../pmock-shared/SKILL.md`](../pmock-shared/SKILL.md)**：" +
		"认证、响应信封、退出码和风险确认规则都在那里，本文件只描述命令本身。\n\n")
	fmt.Fprintf(&b, "- 命令前缀：%s\n- 收录命令：%d 个\n\n", prefixList(sk), len(entries))

	b.WriteString("## 命令索引\n\n")
	groups := groupEntries(entries)
	for _, g := range groups {
		fmt.Fprintf(&b, "### %s\n\n| 命令 | 用途 | 风险 |\n| --- | --- | --- |\n", g.name)
		for _, e := range g.entries {
			fmt.Fprintf(&b, "| [`%s`](%s) | %s | `%s` |\n", e.Command, e.Link, esc(e.Summary), riskValue(e.Risk))
		}
		b.WriteString("\n")
	}
	b.WriteString("风险等级：`read` 只读 · `write` 修改数据 · `costly` 产生费用 · " +
		"`destructive` 不可逆（需显式 `--yes`）。\n")
	return []byte(b.String())
}

func renderCommand(sk Skill, c *spec.Command, doc *spec.Document, opts Options) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", c.FullPath())
	if d := oneLine(c.Description); d != "" {
		b.WriteString(d + "\n\n")
	} else if s := oneLine(c.Summary); s != "" {
		b.WriteString(s + "\n\n")
	}

	b.WriteString("## 概览\n\n| 项 | 值 |\n| --- | --- |\n")
	fmt.Fprintf(&b, "| 命令 | `%s %s` |\n", opts.Binary, c.FullPath())
	fmt.Fprintf(&b, "| HTTP | `%s %s` |\n", c.Method, c.Path)
	fmt.Fprintf(&b, "| 风险 | %s |\n", riskCell(c.Risk))
	if len(c.Aliases) > 0 {
		fmt.Fprintf(&b, "| 别名 | %s |\n", aliasCell(c))
	}
	if c.CanWait() {
		b.WriteString("| 异步 | 支持 `--wait`（`--wait-timeout`、`--poll-interval`） |\n")
	} else {
		b.WriteString("| 异步 | 否 |\n")
	}
	b.WriteString("\n")

	writeFlags(&b, c)
	writeExamples(&b, c, opts)
	writeBody(&b, c, doc)
	writeResponse(&b, c, doc)
	writeFooter(&b, sk)
	return []byte(b.String())
}

func renderManualPage(sk Skill, m ManualPage) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", m.Command)
	if s := oneLine(m.Summary); s != "" {
		b.WriteString(s + "\n\n")
	}
	b.WriteString("## 概览\n\n| 项 | 值 |\n| --- | --- |\n")
	fmt.Fprintf(&b, "| 命令 | `pmock-cli %s` |\n", m.Command)
	fmt.Fprintf(&b, "| 风险 | %s |\n\n", riskCell(m.Risk))
	b.WriteString(strings.TrimSpace(m.Body))
	b.WriteString("\n\n")
	writeFooter(&b, sk)
	return []byte(b.String())
}

func writeFrontmatter(b *strings.Builder, sk Skill, opts Options) {
	b.WriteString("---\n")
	fmt.Fprintf(b, "name: %s\n", sk.Name)
	fmt.Fprintf(b, "version: %s\n", opts.Version)
	fmt.Fprintf(b, "description: %q\n", oneLine(sk.Description))
	b.WriteString("metadata:\n  requires:\n")
	fmt.Fprintf(b, "    bins: [%q]\n", opts.Binary)
	fmt.Fprintf(b, "  cliHelp: %q\n", sk.CliHelp)
	b.WriteString("---\n\n")
}

func writeFlags(b *strings.Builder, c *spec.Command) {
	if len(c.Flags) == 0 {
		return
	}
	b.WriteString("## 参数\n\n| 参数 | 类型 | 必填 | 默认 | 说明 |\n| --- | --- | --- | --- | --- |\n")
	flags := append([]*spec.Flag(nil), c.Flags...)
	// Required first, then alphabetical: the flags that decide whether the
	// call can succeed at all should not be buried in the middle.
	sort.SliceStable(flags, func(i, j int) bool {
		if flags[i].Required != flags[j].Required {
			return flags[i].Required
		}
		return flags[i].Name < flags[j].Name
	})
	for _, f := range flags {
		required := ""
		if f.Required {
			required = "是"
		}
		def := ""
		if f.Default != nil {
			def = "`" + fmt.Sprint(f.Default) + "`"
		}
		desc := oneLine(f.Description)
		if len(f.Enum) > 0 {
			desc = strings.TrimSpace(desc + "；取值：" + enumList(f.Enum))
		}
		fmt.Fprintf(b, "| `--%s` | %s | %s | %s | %s |\n", f.Name, f.Type, required, def, esc(desc))
	}
	b.WriteString("\n")
	if c.HasJSONInput() {
		b.WriteString("> 嵌套对象或数组字段用 `--json '<object>'` 或 `--json @file.json` 传入，" +
			"会与上面的 flag 合并。\n\n")
	}
}

func writeExamples(b *strings.Builder, c *spec.Command, opts Options) {
	b.WriteString("## 示例\n\n```bash\n")
	if len(c.Examples) > 0 {
		for _, e := range c.Examples {
			b.WriteString(strings.TrimSpace(e) + "\n")
		}
	} else {
		// No curated example: build a runnable skeleton from the required
		// flags, which is more useful than pointing at --help.
		line := opts.Binary + " " + c.FullPath()
		required := 0
		for _, f := range c.Flags {
			if !f.Required {
				continue
			}
			required++
			line += " --" + f.Name + " " + placeholder(f.Type)
		}
		if required == 0 {
			line += " --help"
		}
		b.WriteString(line + "\n")
	}
	b.WriteString("```\n\n")
}

// placeholder is the shell placeholder for a required flag value.
func placeholder(typ string) string {
	switch typ {
	case spec.TypeInteger, spec.TypeNumber:
		return "<int>"
	case spec.TypeBoolean:
		return "<true|false>"
	case spec.TypeStringArray, spec.TypeIntArray:
		return "<csv>"
	default:
		return "<string>"
	}
}

func writeBody(b *strings.Builder, c *spec.Command, doc *spec.Document) {
	if c.Body == nil || len(c.Body.Properties) == 0 {
		return
	}
	// Most body fields become flags. Only the leftovers need a table: listing
	// all of them again would double the page for no extra information.
	withFlag := map[string]bool{}
	for _, f := range c.Flags {
		if f.In == "body" {
			withFlag[f.Field] = true
		}
	}
	var jsonOnly []string
	for _, name := range sortedKeys(c.Body.Properties) {
		if !withFlag[name] {
			jsonOnly = append(jsonOnly, name)
		}
	}

	b.WriteString("## 请求体\n\n")
	if len(jsonOnly) == 0 {
		b.WriteString("> 字段名与上面的 flag 一一对应（`--kebab-case` → `camelCase`），直接用 flag 即可；" +
			"嵌套结构再改用 `--json '<object>'`。\n\n")
	} else {
		if name := oneLine(c.BodySchemaName); name != "" {
			fmt.Fprintf(b, "> schema：`%s`\n\n", name)
		}
		b.WriteString("> 下面这些字段没有对应 flag，只能用 `--json` 传入。\n\n")
		b.WriteString("| 字段 | 类型 | 必填 | 说明 |\n| --- | --- | --- | --- |\n")
		required := stringSet(c.Body.Required)
		for _, name := range jsonOnly {
			s := c.Body.Properties[name]
			mark := ""
			if required[name] {
				mark = "是"
			}
			desc := ""
			if s != nil {
				desc = oneLine(s.Description)
			}
			fmt.Fprintf(b, "| `%s` | %s | %s | %s |\n", name, typeLabel(doc, s), mark, esc(desc))
		}
		b.WriteString("\n")
	}
	if len(c.ComplexRequired) > 0 {
		names := make([]string, 0, len(c.ComplexRequired))
		for _, n := range c.ComplexRequired {
			names = append(names, "`"+n+"`")
		}
		fmt.Fprintf(b, "> 这些必填字段是嵌套结构，只能用 `--json` 传入：%s\n\n", strings.Join(names, "、"))
	}
}

func writeResponse(b *strings.Builder, c *spec.Command, doc *spec.Document) {
	if c.Resp == nil || len(c.Resp.Properties) == 0 {
		return
	}
	b.WriteString("## 响应字段\n\n")
	b.WriteString("> 接口返回的字段；`code`/`msg`/`traceId` 由统一信封提供，" +
		"分页接口把 `total`/`rows` 放在顶层。\n\n")
	if name := oneLine(c.RespSchemaName); name != "" {
		fmt.Fprintf(b, "> schema：`%s`\n\n", name)
	}
	b.WriteString("| 字段 | 类型 | 说明 |\n| --- | --- | --- |\n")
	names := sortedKeys(c.Resp.Properties)
	limit := len(names)
	truncated := 0
	if limit > maxResponseFields {
		limit = maxResponseFields
		truncated = len(names) - limit
	}
	for _, name := range names[:limit] {
		s := c.Resp.Properties[name]
		desc := ""
		if s != nil {
			desc = oneLine(s.Description)
			if s.Format != "" && desc == "" {
				desc = s.Format
			}
		}
		fmt.Fprintf(b, "| `%s` | %s | %s |\n", name, typeLabel(doc, s), esc(desc))
	}
	b.WriteString("\n")
	if truncated > 0 {
		fmt.Fprintf(b, "> 其余 %d 个字段见 OpenAPI 文档。\n\n", truncated)
	}
}

func writeFooter(b *strings.Builder, sk Skill) {
	fmt.Fprintf(b, "---\n\n[%s](../SKILL.md) · [pmock-shared](../../pmock-shared/SKILL.md)\n", sk.Name)
}

// groupEntries buckets index entries for the SKILL.md index.
type entryGroup struct {
	name    string
	entries []indexEntry
}

func groupEntries(entries []indexEntry) []entryGroup {
	byName := map[string][]indexEntry{}
	var order []string
	for _, e := range entries {
		if _, ok := byName[e.Group]; !ok {
			order = append(order, e.Group)
		}
		byName[e.Group] = append(byName[e.Group], e)
	}
	sort.Strings(order)
	out := make([]entryGroup, 0, len(order))
	for _, name := range order {
		out = append(out, entryGroup{name: name, entries: byName[name]})
	}
	return out
}

func riskValue(r spec.Risk) string {
	if r == "" {
		return string(spec.RiskRead)
	}
	return string(r)
}

// riskCell renders the risk plus what it means for the caller.
func riskCell(r spec.Risk) string {
	switch r {
	case spec.RiskCostly:
		return "`costly`（会产生费用，执行前确认）"
	case spec.RiskDestructive:
		return "`destructive`（不可逆，需 `--yes`）"
	case spec.RiskWrite:
		return "`write`（修改数据）"
	case spec.RiskRead:
		return "`read`（只读）"
	default:
		return "`" + string(r) + "`"
	}
}

func aliasCell(c *spec.Command) string {
	parts := make([]string, 0, len(c.Aliases))
	parent := ""
	if len(c.PathParts) > 1 {
		parent = strings.Join(c.PathParts[:len(c.PathParts)-1], " ") + " "
	}
	for _, a := range c.Aliases {
		parts = append(parts, "`"+parent+a+"`")
	}
	return strings.Join(parts, "、")
}

func prefixList(sk Skill) string {
	parts := make([]string, 0, len(sk.Prefixes))
	for _, p := range sk.Prefixes {
		parts = append(parts, "`pmock-cli "+strings.Join(p, " ")+"`")
	}
	return strings.Join(parts, "、")
}

// typeLabel renders a JSON schema as a compact type name.
func typeLabel(doc *spec.Document, s *spec.Schema) string {
	if s == nil {
		return ""
	}
	if doc != nil {
		if resolved := doc.Resolve(s); resolved != nil {
			s = resolved
		}
	}
	if s.Type == "array" {
		inner := typeLabel(doc, s.Items)
		// An array of $ref objects is far more useful when named.
		if s.Items != nil && s.Items.Ref != "" {
			if name := refName(s.Items.Ref); name != "" {
				inner = name
			}
		}
		if inner == "" {
			return "array"
		}
		return "[]" + inner
	}
	if s.Type == "" {
		return "object"
	}
	return s.Type
}

func refName(ref string) string {
	return strings.TrimPrefix(ref, "#/components/schemas/")
}

func enumList(values []any) string {
	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts, fmt.Sprint(v))
	}
	return strings.Join(parts, ", ")
}

func sortedKeys(m map[string]*spec.Schema) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func stringSet(values []string) map[string]bool {
	out := make(map[string]bool, len(values))
	for _, v := range values {
		out[v] = true
	}
	return out
}

// oneLine collapses whitespace so a multi-line spec description fits in a
// markdown table cell.
func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// esc protects table cells from the pipe character.
func esc(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "|", `\|`)
	return strings.ReplaceAll(s, "\n", " ")
}

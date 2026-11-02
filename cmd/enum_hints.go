package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/haidaisiganbei/pmock-cli/internal/spec"
)

// enumPairRe finds "0=已过期" style pairs in flag descriptions. The PMOCK spec
// documents enums as prose rather than JSON Schema enum, so the description is
// the only place the mapping exists.
var enumPairRe = regexp.MustCompile(`(\d+)\s*=\s*([^,，;；/|)）\s]+)`)

// enumHints explains numeric fields of a response using the flag descriptions
// of the command that produced it. It returns "" when nothing matched, so a
// caller can append the result unconditionally.
func enumHints(catalog *spec.Catalog, path []string, payload string) string {
	if catalog == nil || payload == "" {
		return ""
	}
	command := strings.Join(path, " ")
	var target *spec.Command
	for _, c := range catalog.Commands {
		if c.FullPath() == command {
			target = c
			break
		}
	}
	if target == nil {
		return ""
	}

	var obj map[string]any
	if json.Unmarshal([]byte(payload), &obj) != nil {
		return ""
	}
	// Paginated responses put the payload at the top level; everything else
	// nests it under data. Check both, one level deep.
	data, _ := obj["data"].(map[string]any)

	var lines []string
	for _, f := range target.Flags {
		values := enumValues(f)
		if len(values) == 0 {
			continue
		}
		value, ok := obj[f.Field]
		if !ok && data != nil {
			value, ok = data[f.Field]
		}
		if !ok {
			continue
		}
		lines = append(lines, fmt.Sprintf("- %s = %v → %s", f.Field, value, strings.Join(values, ", ")))
	}
	if len(lines) == 0 {
		return ""
	}
	sort.Strings(lines)
	return "[Enum hints] 响应中的数字字段含义：\n" + strings.Join(lines, "\n")
}

// enumValues returns the documented values of a flag, from the schema enum when
// present and otherwise from the description.
func enumValues(f *spec.Flag) []string {
	if len(f.Enum) > 0 {
		out := make([]string, 0, len(f.Enum))
		for _, v := range f.Enum {
			out = append(out, fmt.Sprint(v))
		}
		return out
	}
	matches := enumPairRe.FindAllStringSubmatch(f.Description, -1)
	out := make([]string, 0, len(matches))
	seen := map[string]bool{}
	for _, m := range matches {
		pair := m[1] + "=" + strings.TrimSpace(m[2])
		if seen[pair] {
			continue
		}
		seen[pair] = true
		out = append(out, pair)
	}
	return out
}

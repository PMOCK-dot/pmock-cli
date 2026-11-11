package cmd

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/PMOCK-dot/pmock-cli/internal/spec"
)

// referencePage is the path of a command's reference page inside a skill
// directory, e.g. "references/phone-list.md" for "phone list".
func referencePage(command string) string {
	return "references/" + strings.Join(strings.Fields(command), "-") + ".md"
}

// canonicalCommand maps a user-supplied command path to the catalog's canonical
// path, so an alias ("phone ls") finds the page for "phone list". An empty
// result means the path is not a command the catalog knows.
func canonicalCommand(catalog *spec.Catalog, command string) string {
	if catalog == nil {
		return ""
	}
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}
	want := strings.Join(parts, " ")
	for _, c := range catalog.Commands {
		if c.FullPath() == want {
			return c.FullPath()
		}
	}
	// Aliases only replace the last segment ("phone ls" for "phone list");
	// parent segments are always canonical.
	last := parts[len(parts)-1]
	parent := parts[:len(parts)-1]
	for _, c := range catalog.Commands {
		if len(c.PathParts) != len(parts) {
			continue
		}
		if strings.Join(c.PathParts[:len(parts)-1], " ") != strings.Join(parent, " ") {
			continue
		}
		for _, a := range c.Aliases {
			if a == last {
				return c.FullPath()
			}
		}
	}
	return ""
}

// validCommands lists the commands that extend a prefix, for "did you mean" and
// for rejecting a typo instead of silently answering with the parent's help.
func validCommands(catalog *spec.Catalog, prefix string) []string {
	if catalog == nil {
		return nil
	}
	parts := strings.Fields(prefix)
	var out []string
	for _, c := range catalog.Commands {
		if c.Internal || len(c.PathParts) <= len(parts) {
			continue
		}
		ok := true
		for i, p := range parts {
			if c.PathParts[i] != p {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		out = append(out, c.PathParts[len(parts)])
	}
	sort.Strings(out)
	// The same sub-group appears once per leaf; collapse the duplicates.
	dedup := out[:0]
	for i, s := range out {
		if i == 0 || s != out[i-1] {
			dedup = append(dedup, s)
		}
	}
	return dedup
}

// commandExists reports whether the path names a real command in the CLI tree.
// The tree is broader than the catalog: hand-written commands (doctor, mcp,
// watch, phone file upload) live there without a generated reference page.
func commandExists(catalog *spec.Catalog, command string) bool {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}
	root := newRootWith(catalog, false)
	cmd, rest, err := root.Find(parts)
	return err == nil && len(rest) == 0 && cmd != nil
}

// unknownCommand reports a command path the catalog does not have, naming the
// valid next segments so a typo is fixable in one round.
func unknownCommand(catalog *spec.Catalog, command string) error {
	parts := strings.Fields(command)
	prefix := ""
	if len(parts) > 1 {
		prefix = strings.Join(parts[:len(parts)-1], " ")
	}
	if next := validCommands(catalog, prefix); len(next) > 0 {
		where := "at the top level"
		if prefix != "" {
			where = fmt.Sprintf("under %q", prefix)
		}
		return fmt.Errorf("unknown command %q; %s the commands are: %s",
			command, where, strings.Join(next, ", "))
	}
	return fmt.Errorf("unknown command %q", command)
}

// readReference returns the generated reference page for a command. The page
// lives in exactly one skill directory, so the skill is discovered instead of
// hard-coded: moving a command group between skills keeps this working.
func readReference(skills fs.FS, command string) (string, bool) {
	if skills == nil || strings.TrimSpace(command) == "" {
		return "", false
	}
	page := referencePage(command)
	entries, err := fs.ReadDir(skills, ".")
	if err != nil {
		return "", false
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		data, err := fs.ReadFile(skills, e.Name()+"/"+page)
		if err != nil {
			continue
		}
		return strings.TrimSpace(string(data)), true
	}
	return "", false
}

// skillNames lists the skill directories, for "what can I read" answers.
func skillNames(skills fs.FS) []string {
	if skills == nil {
		return nil
	}
	entries, err := fs.ReadDir(skills, ".")
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}

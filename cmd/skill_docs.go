package cmd

import (
	"io/fs"
	"sort"
	"strings"
)

// referencePage is the path of a command's reference page inside a skill
// directory, e.g. "references/phone-list.md" for "phone list".
func referencePage(command string) string {
	return "references/" + strings.Join(strings.Fields(command), "-") + ".md"
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

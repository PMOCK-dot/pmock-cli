package skillgen

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/PMOCK-dot/pmock-cli/internal/spec"
)

func loadCatalog(t *testing.T) *spec.Catalog {
	t.Helper()
	docFile, err := os.Open("../../openapi/pmock-openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer docFile.Close()
	doc, err := spec.Load(docFile)
	if err != nil {
		t.Fatal(err)
	}
	overlayFile, err := os.Open("../../openapi/overlay.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer overlayFile.Close()
	overlay, err := spec.LoadOverlay(overlayFile)
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := spec.BuildCatalog(doc, overlay)
	if err != nil {
		t.Fatal(err)
	}
	return catalog
}

func render(t *testing.T, cat *spec.Catalog) []File {
	t.Helper()
	files, err := Render(cat, Options{Skills: DefaultSkills(), Manual: DefaultManual()})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func byPath(files []File) map[string]string {
	out := make(map[string]string, len(files))
	for _, f := range files {
		out[f.Path] = string(f.Content)
	}
	return out
}

// Every public command of a covered prefix must have exactly one page: a
// command added to the spec therefore cannot silently miss its documentation.
func TestRenderProducesPagePerPublicCommand(t *testing.T) {
	cat := loadCatalog(t)
	files := byPath(render(t, cat))

	skills := DefaultSkills()
	for _, c := range cat.Commands {
		if c.Internal {
			continue
		}
		var owner *Skill
		for i := range skills {
			if skills[i].covers(c.PathParts) {
				owner = &skills[i]
				break
			}
		}
		if owner == nil {
			continue
		}
		path := pagePath(owner.Name, c.FullPath())
		if _, ok := files[path]; !ok {
			t.Errorf("missing reference page for %q at %s", c.FullPath(), path)
		}
	}
}

func TestManualPageIsIndexed(t *testing.T) {
	files := byPath(render(t, loadCatalog(t)))
	if _, ok := files["pmock-phone/references/phone-file-upload.md"]; !ok {
		t.Fatal("manual upload page is missing")
	}
	index := files["pmock-phone/SKILL.md"]
	if !strings.Contains(index, "references/phone-file-upload.md") {
		t.Error("manual page is not linked from the command index")
	}
}

func TestSkillFrontmatter(t *testing.T) {
	files := byPath(render(t, loadCatalog(t)))
	for _, sk := range DefaultSkills() {
		page, ok := files[sk.Name+"/SKILL.md"]
		if !ok {
			t.Fatalf("%s/SKILL.md is missing", sk.Name)
		}
		if !strings.HasPrefix(page, "---\n") {
			t.Errorf("%s has no YAML frontmatter", sk.Name)
		}
		for _, want := range []string{
			"name: " + sk.Name,
			"description: ",
			`bins: ["pmock-cli"]`,
			"cliHelp: " + quoted(sk.CliHelp),
		} {
			if !strings.Contains(page, want) {
				t.Errorf("%s/SKILL.md is missing %q", sk.Name, want)
			}
		}
		if n := len([]rune(sk.Description)); n > maxDescription {
			t.Errorf("%s description is %d characters", sk.Name, n)
		}
	}
}

func TestLinksResolve(t *testing.T) {
	files := byPath(render(t, loadCatalog(t)))
	// pmock-shared is hand-written, so it is on disk but never rendered.
	if _, err := os.Stat("../../skills/pmock-shared/SKILL.md"); err == nil {
		files["pmock-shared/SKILL.md"] = ""
	}
	linkRe := regexp.MustCompile(`\]\(([^)\s]+)\)`)

	for path, content := range files {
		for _, m := range linkRe.FindAllStringSubmatch(content, -1) {
			target := m[1]
			if strings.HasPrefix(target, "http") || strings.HasPrefix(target, "mailto:") || strings.HasPrefix(target, "#") {
				continue
			}
			resolved := resolveLink(path, target)
			if _, ok := files[resolved]; !ok {
				t.Errorf("%s links to %s, which is not rendered (%s)", path, target, resolved)
			}
		}
	}
}

func TestRenderIsDeterministic(t *testing.T) {
	cat := loadCatalog(t)
	first := render(t, cat)
	second := render(t, cat)
	if len(first) != len(second) {
		t.Fatalf("file count changed: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i].Path != second[i].Path || string(first[i].Content) != string(second[i].Content) {
			t.Fatalf("file %d differs between renders", i)
		}
	}
}

func TestRenderRejectsInvalidSkill(t *testing.T) {
	cat := loadCatalog(t)
	_, err := Render(cat, Options{Skills: []Skill{{
		Name:        "pmock-broken",
		Description: strings.Repeat("x", maxDescription+1),
		CliHelp:     "pmock-cli phone --help",
		Prefixes:    [][]string{{"phone"}},
	}}})
	if err == nil {
		t.Fatal("expected an oversized description to be rejected")
	}
}

// resolveLink turns a relative markdown link into a slash path relative to the
// skills root.
func resolveLink(from, target string) string {
	dir := from
	if i := strings.LastIndexByte(dir, '/'); i >= 0 {
		dir = dir[:i]
	} else {
		dir = ""
	}
	for _, part := range strings.Split(target, "/") {
		switch part {
		case "", ".":
			continue
		case "..":
			if i := strings.LastIndexByte(dir, '/'); i >= 0 {
				dir = dir[:i]
			} else {
				dir = ""
			}
		default:
			if dir == "" {
				dir = part
			} else {
				dir += "/" + part
			}
		}
	}
	return dir
}

func quoted(s string) string { return `"` + s + `"` }

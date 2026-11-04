package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/PMOCK-dot/pmock-cli/internal/spec"
)

func runSkills(t *testing.T, cat *spec.Catalog, args ...string) (string, error) {
	t.Helper()
	root := &cobra.Command{Use: "pmock-cli", SilenceUsage: true, SilenceErrors: true}
	root.AddCommand(newSkillsCmd(cat))
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"skills", "generate"}, args...))
	err := root.Execute()
	return out.String(), err
}

func TestSkillsGenerateAndCheck(t *testing.T) {
	dir := t.TempDir()
	cat := testCatalog(t)

	if _, err := runSkills(t, cat, "--out", dir); err != nil {
		t.Fatalf("generate: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "pmock-phone", "SKILL.md")); err != nil {
		t.Fatalf("SKILL.md not written: %v", err)
	}
	if _, err := runSkills(t, cat, "--out", dir, "--check"); err != nil {
		t.Fatalf("fresh tree reported stale: %v", err)
	}

	page := filepath.Join(dir, "pmock-phone", "references", "phone-list.md")
	if err := os.WriteFile(page, []byte("tampered"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runSkills(t, cat, "--out", dir, "--check"); err == nil {
		t.Fatal("tampered tree passed --check")
	}
}

func TestSkillsPrunesOrphansOnly(t *testing.T) {
	dir := t.TempDir()
	cat := testCatalog(t)
	if _, err := runSkills(t, cat, "--out", dir); err != nil {
		t.Fatal(err)
	}

	orphan := filepath.Join(dir, "pmock-phone", "references", "phone-removed.md")
	if err := os.WriteFile(orphan, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	handWritten := filepath.Join(dir, "pmock-shared", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(handWritten), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(handWritten, []byte("hand written"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := runSkills(t, cat, "--out", dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Error("orphan page survived regeneration")
	}
	if _, err := os.Stat(handWritten); err != nil {
		t.Error("hand-written skill was removed by pruning")
	}
}

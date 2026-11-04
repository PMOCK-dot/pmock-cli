package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/PMOCK-dot/pmock-cli/internal/exitcode"
	"github.com/PMOCK-dot/pmock-cli/internal/skillgen"
	"github.com/PMOCK-dot/pmock-cli/internal/spec"
)

const skillsLong = `Render the agent skills tree from the command catalog.

Every public command becomes one reference page; the catalog comes from
openapi/pmock-openapi.yaml plus openapi/overlay.yaml, so the docs cannot drift
from the CLI surface. The tree is committed, and ` + "`--check`" + ` fails when the
files on disk differ from what would be rendered:

    make skills          # regenerate
    make skills-check    # fail if skills/ is stale`

// newSkillsCmd exposes skill generation. The command is hidden from help: it is
// developer tooling, not part of the agent-facing surface.
func newSkillsCmd(catalog *spec.Catalog) *cobra.Command {
	var (
		out   string
		check bool
	)

	cmd := &cobra.Command{
		Use:    "skills",
		Short:  "Generate the agent skills tree (skills/)",
		Long:   skillsLong,
		Hidden: true,
	}

	gen := &cobra.Command{
		Use:   "generate",
		Short: "Render skills/ (one page per public command)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			files, err := skillgen.Render(catalog, skillgen.Options{
				Skills: skillgen.DefaultSkills(),
				Manual: skillgen.DefaultManual(),
			})
			if err != nil {
				return exitcode.Runtimef("%v", err)
			}
			if check {
				return checkSkills(cmd, out, files)
			}
			return writeSkills(cmd, out, files)
		},
	}
	gen.Flags().StringVar(&out, "out", "skills", "output directory")
	gen.Flags().BoolVar(&check, "check", false, "fail if skills/ is stale instead of writing")

	cmd.AddCommand(gen)
	return cmd
}

// writeSkills writes every rendered file and removes pages the catalog no
// longer produces, so a deleted command cannot leave an orphan page behind.
func writeSkills(cmd *cobra.Command, root string, files []skillgen.File) error {
	if strings.TrimSpace(root) == "" {
		return exitcode.Usagef("--out must not be empty")
	}
	stale, err := staleSkills(root, files)
	if err != nil {
		return exitcode.Runtimef("failed to inspect %s: %v", root, err)
	}
	for _, rel := range stale {
		if err := os.Remove(filepath.Join(root, filepath.FromSlash(rel))); err != nil && !os.IsNotExist(err) {
			return exitcode.Runtimef("failed to remove %s: %v", rel, err)
		}
	}
	for _, f := range files {
		path := filepath.Join(root, filepath.FromSlash(f.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return exitcode.Runtimef("failed to create %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, f.Content, 0o644); err != nil {
			return exitcode.Runtimef("failed to write %s: %v", f.Path, err)
		}
	}
	fmt.Fprintf(cmd.OutOrStdout(), "rendered %d files to %s", len(files), root)
	if len(stale) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), " (removed %d stale)", len(stale))
	}
	fmt.Fprintln(cmd.OutOrStdout())
	return nil
}

// checkSkills reports drift without touching the tree.
func checkSkills(cmd *cobra.Command, root string, files []skillgen.File) error {
	var stale []string
	for _, f := range files {
		path := filepath.Join(root, filepath.FromSlash(f.Path))
		current, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(current, f.Content) {
			stale = append(stale, f.Path)
		}
	}
	orphans, err := staleSkills(root, files)
	if err != nil {
		return exitcode.Runtimef("failed to inspect %s: %v", root, err)
	}
	stale = append(stale, orphans...)
	if len(stale) > 0 {
		sort.Strings(stale)
		return exitcode.Runtimef("%s is stale (%d file(s)); run `make skills`:\n  %s",
			root, len(stale), strings.Join(stale, "\n  "))
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%s is up to date (%d files)\n", root, len(files))
	return nil
}

// staleSkills lists files inside generated skill directories that the
// generator no longer owns. Only directories it writes to are inspected, so
// hand-written skills such as pmock-shared are never touched.
func staleSkills(root string, files []skillgen.File) ([]string, error) {
	owned := make(map[string]bool, len(files))
	dirs := map[string]bool{}
	for _, f := range files {
		owned[f.Path] = true
		if i := strings.IndexByte(f.Path, '/'); i > 0 {
			dirs[f.Path[:i]] = true
		}
	}

	var stale []string
	for dir := range dirs {
		base := filepath.Join(root, dir)
		err := filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			rel = filepath.ToSlash(rel)
			if !owned[rel] {
				stale = append(stale, rel)
			}
			return nil
		})
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	sort.Strings(stale)
	return stale, nil
}

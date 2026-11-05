package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/PMOCK-dot/pmock-cli/internal/build"
	"github.com/PMOCK-dot/pmock-cli/internal/client"
	"github.com/PMOCK-dot/pmock-cli/internal/exitcode"
	"github.com/PMOCK-dot/pmock-cli/internal/spec"
)

// defaultBatchParallel bounds concurrent targets. Six matches the number of
// phones a person usually operates at once; raise it with --max-parallel.
const defaultBatchParallel = 6

// batchResult is one target's outcome inside the aggregate envelope.
type batchResult struct {
	Target   string          `json:"target"`
	OK       bool            `json:"ok"`
	Code     int             `json:"code,omitempty"`
	Msg      string          `json:"msg,omitempty"`
	TraceID  string          `json:"traceId,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
	Error    string          `json:"error,omitempty"`
	ExitCode int             `json:"exitCode,omitempty"`
}

// batchEnvelope is what --codes and --all print: one entry per target instead
// of one process per target.
type batchEnvelope struct {
	OK      bool          `json:"ok"`
	Total   int           `json:"total"`
	Failed  int           `json:"failed"`
	Results []batchResult `json:"results"`
}

// subjectFlag reports which flag addresses a single phone for this command.
// Action commands take --code; lifecycle, app and file commands take --id.
func subjectFlag(sc *spec.Command) string {
	for _, f := range sc.Flags {
		if f.Name == "code" && f.Type == spec.TypeString {
			return "code"
		}
	}
	for _, f := range sc.Flags {
		if f.Name == "id" && f.Type == spec.TypeInteger {
			return "id"
		}
	}
	return ""
}

func addBatchFlags(c *cobra.Command, subject string) {
	if c.Flags().Lookup("codes") == nil {
		label := "phone codes"
		if subject == "id" {
			label = "phone ids"
		}
		c.Flags().StringSlice("codes", nil, "run once per "+label+" and return one aggregate JSON; comma-separated")
	}
	if c.Flags().Lookup("all") == nil {
		c.Flags().Bool("all", false, "run once for every phone in the account (uses the phone list)")
	}
	if c.Flags().Lookup("max-parallel") == nil {
		c.Flags().Int("max-parallel", defaultBatchParallel, "maximum concurrent targets in batch mode")
	}
}

func batchRequested(cmd *cobra.Command) bool {
	if f := cmd.Flags().Lookup("codes"); f != nil {
		if codes, err := cmd.Flags().GetStringSlice("codes"); err == nil && len(codes) > 0 {
			return true
		}
	}
	if f := cmd.Flags().Lookup("all"); f != nil {
		if all, err := cmd.Flags().GetBool("all"); err == nil && all {
			return true
		}
	}
	return false
}

// runBatch executes the same operation once per target concurrently and prints
// a single aggregate envelope. The exit code follows the worst target.
func runBatch(cmd *cobra.Command, sc *spec.Command, byName map[string]*spec.Flag, g *globals) error {
	subject := subjectFlag(sc)
	if subject == "" {
		return exitcode.Usagef("%s has no per-phone subject flag for --codes/--all", sc.FullPath())
	}
	codes, _ := cmd.Flags().GetStringSlice("codes")
	all, _ := cmd.Flags().GetBool("all")
	if len(codes) > 0 && all {
		return exitcode.Usagef("--codes and --all are mutually exclusive")
	}
	if g.dryRun {
		return exitcode.Usagef("--dry-run cannot be combined with --codes/--all; run --dry-run once per target instead")
	}

	targets := append([]string{}, codes...)
	if all {
		var err error
		targets, err = listPhoneTargets(cmd, g, subject)
		if err != nil {
			return err
		}
	}
	targets = dedupeTargets(targets)
	if len(targets) == 0 {
		return exitcode.Usagef("no targets: pass --codes, or --all once the account has phones")
	}

	parallel, _ := cmd.Flags().GetInt("max-parallel")
	if parallel <= 0 {
		parallel = defaultBatchParallel
	}
	if parallel > len(targets) {
		parallel = len(targets)
	}

	snapshot := snapshotFlags(cmd, subject)
	results := make([]batchResult, len(targets))
	sem := make(chan struct{}, parallel)
	var wg sync.WaitGroup
	for i, target := range targets {
		wg.Add(1)
		go func(i int, target string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i] = runBatchTarget(cmd, sc, byName, g, subject, target, snapshot)
		}(i, target)
	}
	wg.Wait()

	out := batchEnvelope{OK: true, Total: len(results), Results: results}
	worst := exitcode.OK
	for _, r := range results {
		if r.OK {
			continue
		}
		out.OK = false
		out.Failed++
		if r.ExitCode > worst {
			worst = r.ExitCode
		}
	}
	if err := renderBatch(cmd, g, out); err != nil {
		return err
	}
	if out.OK {
		return nil
	}
	if worst == exitcode.OK {
		worst = exitcode.Business
	}
	return exitcode.New(worst, fmt.Errorf("%d of %d targets failed", out.Failed, out.Total))
}

// runBatchTarget runs the operation for one target on its own command instance
// so concurrent goroutines never share flag state.
func runBatchTarget(parent *cobra.Command, sc *spec.Command, byName map[string]*spec.Flag, g *globals, subject, target string, snapshot []flagValue) batchResult {
	res := batchResult{Target: target}
	inner := newOperationCmd(sc, g)
	applyFlags(inner, snapshot)
	if err := inner.Flags().Set(subject, target); err != nil {
		res.ExitCode = exitcode.Usage
		res.Error = err.Error()
		return res
	}
	inner.SetContext(commandContext(parent))
	inner.SetOut(io.Discard)
	inner.SetErr(io.Discard)

	env, err := executeOperation(inner, sc, byName, g)
	if env != nil {
		res.Code = env.Code
		res.Msg = env.Msg
		res.TraceID = env.TraceID
		res.Data = env.Data
		res.OK = env.Code == 0 || env.Code == 200
	}
	if err != nil {
		res.OK = false
		res.Error = err.Error()
		res.ExitCode = exitcode.CodeOf(err)
		if res.ExitCode == exitcode.Runtime && env != nil && res.Code != 0 && res.Code != 200 {
			res.ExitCode = exitcode.Business
		}
	}
	if !res.OK && res.ExitCode == exitcode.OK {
		res.ExitCode = exitcode.Business
	}
	return res
}

// flagValue is one flag the caller set, captured before batch goroutines
// start. pflag FlagSets are not safe for concurrent reads, so the parent's
// flags must be snapshot once instead of visited from every goroutine.
type flagValue struct {
	name  string
	value string
}

func snapshotFlags(from *cobra.Command, skip ...string) []flagValue {
	skipSet := map[string]bool{
		"codes": true, "all": true, "max-parallel": true,
		"ids": true, "out": true,
	}
	for _, name := range skip {
		skipSet[name] = true
	}
	var out []flagValue
	from.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed || skipSet[f.Name] {
			return
		}
		value := f.Value.String()
		if sv, ok := f.Value.(pflag.SliceValue); ok {
			value = strings.Join(sv.GetSlice(), ",")
		}
		out = append(out, flagValue{name: f.Name, value: value})
	})
	return out
}

// applyFlags replays the snapshot onto a fresh command instance.
func applyFlags(to *cobra.Command, snapshot []flagValue) {
	for _, fv := range snapshot {
		if to.Flags().Lookup(fv.name) == nil {
			continue
		}
		_ = to.Flags().Set(fv.name, fv.value)
	}
}

func renderBatch(cmd *cobra.Command, g *globals, out batchEnvelope) error {
	b, err := json.Marshal(out)
	if err != nil {
		return exitcode.Runtimef("failed to render batch result: %v", err)
	}
	if !g.compact {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, b, "", "  "); err == nil {
			b = pretty.Bytes()
		}
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", b); err != nil {
		return exitcode.Runtimef("failed to write batch result: %v", err)
	}
	return nil
}

// listPhoneTargets resolves --all through the phone list, following pages until
// the account is exhausted.
func listPhoneTargets(cmd *cobra.Command, g *globals, subject string) ([]string, error) {
	cfg, err := g.config()
	if err != nil {
		return nil, err
	}
	c := client.New(cfg.BaseURL, cfg.AccessKey, cfg.SecretKey, cfg.Timeout)
	c.UserAgent = "pmock-cli/" + build.Version
	ctx, cancel := context.WithTimeout(commandContext(cmd), cfg.Timeout)
	defer cancel()

	const pageSize = 200
	var out []string
	for page := 1; page <= 100; page++ {
		q := url.Values{}
		q.Set("pageNum", strconv.Itoa(page))
		q.Set("pageSize", strconv.Itoa(pageSize))
		env, err := c.Do(ctx, client.Request{Method: "GET", Path: "/open/api/phone/list", Query: q})
		if err != nil {
			return nil, err
		}
		var payload struct {
			Rows []struct {
				ID   json.RawMessage `json:"id"`
				Code string          `json:"code"`
			} `json:"rows"`
		}
		if err := decodePage(env, &payload); err != nil {
			return nil, exitcode.Runtimef("failed to read phone list: %v", err)
		}
		for _, row := range payload.Rows {
			if subject == "code" {
				if row.Code != "" {
					out = append(out, row.Code)
				}
				continue
			}
			if id, ok := parseInt64Loose(row.ID); ok && id > 0 {
				out = append(out, strconv.FormatInt(id, 10))
			}
		}
		if len(payload.Rows) < pageSize {
			break
		}
	}
	return out, nil
}

func dedupeTargets(in []string) []string {
	seen := make(map[string]bool, len(in))
	out := make([]string, 0, len(in))
	for _, t := range in {
		t = strings.TrimSpace(t)
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

// decodePage reads a paginated envelope. Some deployments nest the page under
// data as {total, rows}; others put rows at the envelope top level.
func decodePage(env *client.Envelope, out any) error {
	if env == nil {
		return fmt.Errorf("empty response")
	}
	if len(env.Data) > 0 && string(env.Data) != "null" {
		if err := json.Unmarshal(env.Data, out); err == nil {
			return nil
		}
	}
	if len(env.Rows) > 0 {
		wrapped := make([]byte, 0, len(env.Rows)+10)
		wrapped = append(wrapped, `{"rows":`...)
		wrapped = append(wrapped, env.Rows...)
		wrapped = append(wrapped, '}')
		return json.Unmarshal(wrapped, out)
	}
	return fmt.Errorf("response has no page rows")
}

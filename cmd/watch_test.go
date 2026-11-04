package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/cobra"

	"github.com/PMOCK-dot/pmock-cli/internal/exitcode"
	"github.com/PMOCK-dot/pmock-cli/internal/progress"
)

func watchRoot(t *testing.T) (*cobra.Command, *bytes.Buffer, *bytes.Buffer, *globals) {
	t.Helper()
	g := testGlobals("")
	root := &cobra.Command{Use: "pmock-cli", SilenceErrors: true, SilenceUsage: true}
	root.AddCommand(newWatchCmd(g))
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	return root, &stdout, &stderr, g
}

func TestWatchRPAStreamsUntilDone(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	var mu sync.Mutex
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if want := "/api/open/api/rpa/task/getInfo/7"; r.URL.Path != want {
			t.Errorf("path = %s, want %s", r.URL.Path, want)
		}
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n == 1 {
			_, _ = w.Write([]byte(`{"traceId":"t1","code":0,"msg":"ok","data":{"status":1}}`))
			return
		}
		_, _ = w.Write([]byte(`{"traceId":"t2","code":0,"msg":"ok","data":{"status":2}}`))
	}))
	defer srv.Close()

	root, stdout, _, g := watchRoot(t)
	g.baseURL = srv.URL + "/api"
	root.SetArgs([]string{"watch", "rpa", "--task-id", "7", "--interval", "5ms"})

	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines:\n%s", len(lines), stdout.String())
	}
	var first, last progress.Event
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(lines[1]), &last); err != nil {
		t.Fatal(err)
	}
	if first.State != progress.StateRunning || first.Seq != 1 {
		t.Fatalf("first = %+v", first)
	}
	if last.State != progress.StateDone || last.Seq != 2 {
		t.Fatalf("last = %+v", last)
	}
	if last.Say != "完成了。" {
		t.Fatalf("say = %q", last.Say)
	}
}

func TestWatchFailedTaskExitsOne(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"traceId":"t","code":0,"msg":"ok","data":{"status":4,"taskList":[{"status":3,"remark":"任务创建失败"}]}}`))
	}))
	defer srv.Close()

	root, stdout, stderr, g := watchRoot(t)
	g.baseURL = srv.URL + "/api"
	root.SetArgs([]string{"watch", "rpa", "--task-id", "7", "--interval", "5ms"})

	err := root.Execute()
	if code := exitcode.CodeOf(err); code != exitcode.Business {
		t.Fatalf("exit = %d, want 1 (err=%v)", code, err)
	}
	var ev progress.Event
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout.String())), &ev); err != nil {
		t.Fatal(err)
	}
	if ev.State != progress.StateFailed || !strings.Contains(ev.Say, "任务创建失败") {
		t.Fatalf("event = %+v", ev)
	}
	if stderr.Len() != 0 {
		t.Fatalf("a failed task should not also print an error line: %q", stderr.String())
	}
}

func TestWatchRequiresTarget(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	root, _, _, _ := watchRoot(t)
	root.SetArgs([]string{"watch", "rpa"})
	err := root.Execute()
	if code := exitcode.CodeOf(err); code != exitcode.Usage {
		t.Fatalf("exit = %d, want 2 (err=%v)", code, err)
	}

	root, _, _, _ = watchRoot(t)
	root.SetArgs([]string{"watch", "phone", "--id", "1", "--op", "explode"})
	if code := exitcode.CodeOf(root.Execute()); code != exitcode.Usage {
		t.Fatalf("bad --op should be a usage error")
	}
}

func TestWatchRecordAndReplay(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"traceId":"t","code":0,"msg":"ok","data":{"status":2}}`))
	}))
	defer srv.Close()

	dir := t.TempDir()
	root, stdout, _, g := watchRoot(t)
	g.baseURL = srv.URL + "/api"
	root.SetArgs([]string{"watch", "rpa", "--task-id", "7", "--record", "--state-dir", dir, "--interval", "5ms"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	path := filepath.Join(dir, "rpa-7.jsonl")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ledger: %v", err)
	}
	live := strings.TrimSpace(stdout.String())
	if strings.TrimSpace(string(raw)) != live {
		t.Fatalf("ledger and stream disagree:\nledger=%s\nstream=%s", raw, live)
	}

	root, replay, _, _ := watchRoot(t)
	root.SetArgs([]string{"watch", "log", "--kind", "rpa", "--task-id", "7", "--state-dir", dir})
	if err := root.Execute(); err != nil {
		t.Fatalf("replay: %v", err)
	}
	if strings.TrimSpace(replay.String()) != live {
		t.Fatalf("replay = %q, want %q", replay.String(), live)
	}
}

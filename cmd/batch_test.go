package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/PMOCK-dot/pmock-cli/internal/exitcode"
)

func TestBatchCodesRunsPerTarget(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		_, _ = w.Write([]byte(`{"traceId":"t","code":0,"msg":"ok","data":{"result":"1"},"status":"finish"}`))
	}))
	defer srv.Close()

	sc := testCommand(t, testCatalog(t), "phone action click")
	cmd := newOperationCmd(sc, testGlobals(srv.URL+"/api"))
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	_ = cmd.Flags().Set("codes", "133319,133215")
	_ = cmd.Flags().Set("pos-type", "2")
	_ = cmd.Flags().Set("pos", "10,20")
	_ = cmd.Flags().Set("wait", "true")

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
	out := decodeResult(t, &stdout)
	if out["ok"] != true || out["total"].(float64) != 2 || out["failed"].(float64) != 0 {
		t.Fatalf("out = %#v", out)
	}
	results := out["results"].([]any)
	first := results[0].(map[string]any)
	if first["target"] != "133319" || first["ok"] != true {
		t.Fatalf("first = %#v", first)
	}
}

func TestBatchAllUsesPhoneList(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	var starts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/open/api/phone/list":
			// Real deployments put the page at the envelope top level and
			// encode ids as strings.
			_, _ = w.Write([]byte(`{"code":0,"msg":"ok","rows":[{"id":"11","code":"A"},{"id":"22","code":"B"}]}`))
		case "/api/open/api/phone/start":
			atomic.AddInt32(&starts, 1)
			_, _ = w.Write([]byte(`{"traceId":"t","code":0,"msg":"ok","data":{"taskId":1,"phoneStatus":1},"status":"finish"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	sc := testCommand(t, testCatalog(t), "phone start")
	cmd := newOperationCmd(sc, testGlobals(srv.URL+"/api"))
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	_ = cmd.Flags().Set("all", "true")

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if starts != 2 {
		t.Fatalf("starts = %d, want 2", starts)
	}
	out := decodeResult(t, &stdout)
	results := out["results"].([]any)
	if results[0].(map[string]any)["target"] != "11" || results[1].(map[string]any)["target"] != "22" {
		t.Fatalf("results = %#v", results)
	}
}

func TestBatchFailureExitCode(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "action/click") {
			_, _ = w.Write([]byte(`{"traceId":"t","code":42001,"msg":"手机不存在"}`))
			return
		}
		_, _ = w.Write([]byte(`{"traceId":"t","code":0,"msg":"ok","status":"finish"}`))
	}))
	defer srv.Close()

	sc := testCommand(t, testCatalog(t), "phone action click")
	cmd := newOperationCmd(sc, testGlobals(srv.URL+"/api"))
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	_ = cmd.Flags().Set("codes", "X9")
	_ = cmd.Flags().Set("pos-type", "2")
	_ = cmd.Flags().Set("pos", "10,20")

	err := cmd.RunE(cmd, nil)
	if code := exitcode.CodeOf(err); code != exitcode.Business {
		t.Fatalf("exit code = %d, want %d (%v)", code, exitcode.Business, err)
	}
	out := decodeResult(t, &stdout)
	if out["ok"] != false || out["failed"].(float64) != 1 {
		t.Fatalf("out = %#v", out)
	}
	first := out["results"].([]any)[0].(map[string]any)
	if first["ok"] != false || first["error"] == "" {
		t.Fatalf("first = %#v", first)
	}
}

func TestBatchRejectsDryRun(t *testing.T) {
	sc := testCommand(t, testCatalog(t), "phone action click")
	g := testGlobals("http://127.0.0.1:1/api")
	g.dryRun = true
	cmd := newOperationCmd(sc, g)
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	_ = cmd.Flags().Set("codes", "A")
	_ = cmd.Flags().Set("pos-type", "2")
	_ = cmd.Flags().Set("pos", "1,2")

	err := cmd.RunE(cmd, nil)
	if code := exitcode.CodeOf(err); code != exitcode.Usage {
		t.Fatalf("exit code = %d, want %d (%v)", code, exitcode.Usage, err)
	}
}

func TestBatchUnsupportedTargets(t *testing.T) {
	// phone list has no per-phone subject flag, so it must not grow --codes.
	sc := testCommand(t, testCatalog(t), "phone list")
	if subjectFlag(sc) != "" {
		t.Fatalf("subjectFlag = %q", subjectFlag(sc))
	}
}

func TestCatalogAsyncPolicies(t *testing.T) {
	check := catalogCheck(testCatalog(t))
	if !check.OK {
		t.Fatalf("catalog check = %+v", check)
	}
}

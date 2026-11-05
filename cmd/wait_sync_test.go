package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

// Synchronous endpoints already return a terminal envelope; --wait must accept
// them instead of failing with "response has no taskId".
func TestWaitOnSynchronousEndpoint(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"traceId":"t","code":200,"msg":"操作成功","data":null,"status":"finish"}`))
	}))
	defer srv.Close()

	sc := testCommand(t, testCatalog(t), "app operate")
	cmd := newOperationCmd(sc, testGlobals(srv.URL+"/api"))
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	_ = cmd.Flags().Set("id", "123")
	_ = cmd.Flags().Set("app-op", "start")
	_ = cmd.Flags().Set("pkg", "com.example.app")
	_ = cmd.Flags().Set("wait", "true")

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	out := decodeResult(t, &stdout)
	wait := out["data"].(map[string]any)["wait"].(map[string]any)
	if wait["ok"] != true {
		t.Fatalf("wait = %#v", wait)
	}
}

func TestPkgFailureListsCandidates(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/open/api/app/operate":
			_, _ = w.Write([]byte(`{"traceId":"t","code":42000,"msg":"应用不存在"}`))
		case "/api/open/api/app/list":
			_, _ = w.Write([]byte(`{"traceId":"t","code":200,"msg":"ok","data":{"rows":[{"apkPkg":"com.sec.android.gallery3d"},{"apkPkg":"com.zhiliaoapp.musically"}]}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	sc := testCommand(t, testCatalog(t), "app operate")
	cmd := newOperationCmd(sc, testGlobals(srv.URL+"/api"))
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	_ = cmd.Flags().Set("id", "123")
	_ = cmd.Flags().Set("app-op", "start")
	_ = cmd.Flags().Set("pkg", "com.google.android.apps.photos")

	_ = cmd.RunE(cmd, nil)
	if !strings.Contains(stderr.String(), "Hint") || !strings.Contains(stderr.String(), "com.sec.android.gallery3d") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

// The action task status endpoint is unreliable, so --wait must return on the
// terminal envelope instead of polling it.
func TestActionWaitDoesNotPollStatus(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	var statusCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "action/task/status") {
			atomic.AddInt32(&statusCalls, 1)
		}
		_, _ = w.Write([]byte(`{"traceId":"t","code":200,"msg":"操作成功","data":{"result":"9"},"status":"finish"}`))
	}))
	defer srv.Close()

	sc := testCommand(t, testCatalog(t), "phone action app")
	cmd := newOperationCmd(sc, testGlobals(srv.URL+"/api"))
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	_ = cmd.Flags().Set("code", "133319")
	_ = cmd.Flags().Set("op-type", "1")
	_ = cmd.Flags().Set("pkg", "com.example.app")
	_ = cmd.Flags().Set("wait", "true")

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if statusCalls != 0 {
		t.Fatalf("polled the unreliable status endpoint %d times", statusCalls)
	}
	out := decodeResult(t, &stdout)
	wait := out["data"].(map[string]any)["wait"].(map[string]any)
	if wait["ok"] != true || wait["note"] == nil {
		t.Fatalf("wait = %#v", wait)
	}
}

package clientbridge

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDescribeActivityUsesIDOrCode(t *testing.T) {
	tests := []struct {
		name        string
		path        []string
		params      map[string]any
		wantAction  string
		wantID      string
		wantCode    string
		wantApplied bool
	}{
		{
			name:       "screenshot uses id",
			path:       []string{"phone", "screen"},
			params:     map[string]any{"id": float64(123)},
			wantAction: "screenshot",
			wantID:     "123",
		},
		{
			name:       "click uses code",
			path:       []string{"phone", "action", "click"},
			params:     map[string]any{"code": "133337", "posType": float64(2)},
			wantAction: "click",
			wantCode:   "133337",
		},
		{
			name:        "power on emits activity and applied",
			path:        []string{"phone", "start"},
			params:      map[string]any{"id": "1900000000000000123"},
			wantAction:  "powerOn",
			wantID:      "1900000000000000123",
			wantApplied: true,
		},
		{
			name:        "install app emits installApp and refresh",
			path:        []string{"app", "install"},
			params:      map[string]any{"id": "p1"},
			wantAction:  "installApp",
			wantID:      "p1",
			wantApplied: true,
		},
		{
			name:       "file upload is phone scoped",
			path:       []string{"phone", "file", "upload"},
			params:     map[string]any{"id": "p1"},
			wantAction: "unknown",
			wantID:     "p1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := Describe(tt.path, tt.params)
			if len(events) == 0 {
				t.Fatalf("Describe(%v) returned no events", tt.path)
			}
			activity := events[0]
			if activity.Kind != "activity" {
				t.Fatalf("first event kind = %q, want activity", activity.Kind)
			}
			if activity.Action != tt.wantAction {
				t.Fatalf("action = %q, want %q", activity.Action, tt.wantAction)
			}
			if activity.PhoneID != tt.wantID || activity.PhoneCode != tt.wantCode {
				t.Fatalf("phone = (%q,%q), want (%q,%q)",
					activity.PhoneID, activity.PhoneCode, tt.wantID, tt.wantCode)
			}
			if tt.wantApplied {
				if len(events) != 2 || events[1].Kind != "applied" {
					t.Fatalf("events = %+v, want activity + applied", events)
				}
				if events[1].SettleMS == 0 && tt.path[0] == "phone" {
					t.Fatalf("power-on applied event missing settleMs: %+v", events[1])
				}
			}
		})
	}
}

func TestDescribeAppliedScopesAndPhoneIDs(t *testing.T) {
	events := Describe([]string{"phone", "move-group"}, map[string]any{
		"id":        float64(9),
		"phone-ids": []any{"1", float64(2)},
	})
	if len(events) != 1 || events[0].Kind != "applied" {
		t.Fatalf("events = %+v, want one applied event", events)
	}
	if got := events[0].Scopes; len(got) != 2 || got[0] != "groups" || got[1] != "cloudPhones" {
		t.Fatalf("scopes = %v", got)
	}
	if got := events[0].PhoneIDs; len(got) != 2 || got[0] != "1" || got[1] != "2" {
		t.Fatalf("phoneIds = %v", got)
	}
}

func TestDescribeTagWriteRefreshesTags(t *testing.T) {
	events := Describe([]string{"phone", "batch-update-tags"}, map[string]any{
		"phone-ids": []any{"1", "2"},
		"mode":      "append",
	})
	if len(events) != 1 || events[0].Kind != "applied" {
		t.Fatalf("events = %+v, want one applied event", events)
	}
	if got := events[0].Scopes; len(got) != 2 || got[0] != "tags" || got[1] != "cloudPhones" {
		t.Fatalf("scopes = %v", got)
	}
}

func TestDescribeReadsIDFromJSONPayload(t *testing.T) {
	events := Describe([]string{"phone", "screen"}, map[string]any{
		"json": `{"id":"p1","isRemove":false}`,
	})
	if len(events) != 1 || events[0].PhoneID != "p1" {
		t.Fatalf("events = %+v, want activity for p1", events)
	}
}

func TestDescribeIgnoresReadOnlyCommands(t *testing.T) {
	for _, path := range [][]string{
		{"phone", "list"},
		{"phone", "task-status"},
		{"phone", "remote-token"},
		{"rpa", "task", "list"},
	} {
		if events := Describe(path, map[string]any{"id": "p1"}); len(events) != 0 {
			t.Fatalf("Describe(%v) = %+v, want no events", path, events)
		}
	}
}

func TestReportPostsActivityToRuntime(t *testing.T) {
	resetRuntimeForTest(t)

	var got Event
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode event: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	runtimePath := writeRuntime(t, srv, "token-1")
	t.Setenv(EnvRuntime, runtimePath)
	t.Setenv(EnvLink, "1")

	err := Report(context.Background(), Event{
		Kind:    "activity",
		Tool:    "pmock_cli_phone_screen",
		Action:  "screenshot",
		PhoneID: "p1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if auth != "Bearer token-1" {
		t.Fatalf("authorization = %q", auth)
	}
	if got.Kind != "activity" || got.PhoneID != "p1" || got.Action != "screenshot" {
		t.Fatalf("event = %+v", got)
	}
	if got.At == 0 {
		t.Fatal("event.at should be filled")
	}
}

func TestReportDisabledDoesNotDiscover(t *testing.T) {
	t.Setenv(EnvLink, "0")
	t.Setenv(EnvRuntime, filepath.Join(t.TempDir(), "missing.json"))
	if err := Report(context.Background(), Event{Kind: "activity", PhoneID: "p1", Action: "click"}); err != nil {
		t.Fatalf("disabled bridge returned error: %v", err)
	}
}

func TestReportFailsWhenNoClient(t *testing.T) {
	resetRuntimeForTest(t)
	t.Setenv(EnvLink, "1")
	t.Setenv(EnvRuntime, filepath.Join(t.TempDir(), "missing.json"))
	if err := Report(context.Background(), Event{Kind: "activity", PhoneID: "p1", Action: "click"}); err == nil {
		t.Fatal("Report should fail when runtime.json is missing")
	}
}

func TestReadRuntimeRejectsInvalidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "runtime.json")
	if err := os.WriteFile(path, []byte(`{"port":0,"token":"x"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readRuntime(path); err == nil {
		t.Fatal("expected invalid port error")
	}
	if err := os.WriteFile(path, []byte(`{"port":6874,"token":""}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readRuntime(path); err == nil {
		t.Fatal("expected missing token error")
	}
}

func writeRuntime(t *testing.T, srv *httptest.Server, token string) string {
	t.Helper()
	addr := srv.Listener.Addr().(*net.TCPAddr)
	dir := filepath.Join(t.TempDir(), "mcp")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "runtime.json")
	body := map[string]any{
		"port":      addr.Port,
		"token":     token,
		"pid":       1,
		"startedAt": "2026-11-05T03:00:00.000Z",
	}
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func resetRuntimeForTest(t *testing.T) {
	t.Helper()
	invalidateCache()
	t.Cleanup(invalidateCache)
}

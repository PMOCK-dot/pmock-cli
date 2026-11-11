// Package clientbridge reports pmock-cli operations to a running PMOCK desktop
// client so the client's UI can show the same "Agent is operating" badge and
// data-change refresh as its built-in MCP server.
//
// The bridge is deliberately one-way and fail-open: it only POSTs a small,
// validated event to 127.0.0.1 with the token the client wrote to
// userData/mcp/runtime.json. It never executes a business action through the
// client and never blocks a CLI call for more than a short timeout.
package clientbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Environment variables.
const (
	// EnvLink disables the bridge when set to 0/false/off/no.
	EnvLink = "PMOCK_CLIENT_LINK"
	// EnvRuntime overrides the runtime.json path (or its directory).
	EnvRuntime = "PMOCK_CLIENT_RUNTIME"
)

// eventTimeout keeps the bridge invisible to the caller: a missing or busy
// client must not add noticeable latency to an MCP tool call.
const eventTimeout = 800 * time.Millisecond

// runtimeCacheTTL avoids re-reading runtime.json on every tool call while still
// picking up a client restart (new token/port) within a few seconds.
const runtimeCacheTTL = 3 * time.Second

// Event is the payload accepted by the client's POST /agent-event endpoint.
type Event struct {
	Kind      string   `json:"kind"`
	Tool      string   `json:"tool,omitempty"`
	Action    string   `json:"action,omitempty"`
	PhoneID   string   `json:"phoneId,omitempty"`
	PhoneCode string   `json:"phoneCode,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
	PhoneIDs  []string `json:"phoneIds,omitempty"`
	SettleMS  int      `json:"settleMs,omitempty"`
	At        int64    `json:"at,omitempty"`
}

type runtimeInfo struct {
	Port      int    `json:"port"`
	Token     string `json:"token"`
	PID       int    `json:"pid"`
	StartedAt string `json:"startedAt"`
}

var (
	cacheMu      sync.Mutex
	cachedInfo   runtimeInfo
	cachedErr    error
	cacheExpires time.Time
	cacheValid   bool
	cacheKey     string
)

// Enabled reports whether the client bridge is allowed to report. It is on by
// default so an installed, running client lights up without extra flags.
func Enabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvLink))) {
	case "0", "false", "off", "no", "disabled":
		return false
	default:
		return true
	}
}

// Report sends one event to the running client. Errors are returned for tests
// and diagnostics; callers normally ignore them because the bridge is optional.
func Report(ctx context.Context, event Event) error {
	if !Enabled() {
		return nil
	}
	if event.At == 0 {
		event.At = time.Now().UnixMilli()
	}
	if event.Kind != "activity" && event.Kind != "applied" {
		return fmt.Errorf("clientbridge: invalid event kind %q", event.Kind)
	}

	info, err := currentRuntime()
	if err != nil {
		return err
	}
	if err := postEvent(ctx, info, event); err != nil {
		// The client may have restarted since the cache was filled; rediscover
		// once before giving up.
		invalidateCache()
		if info, freshErr := currentRuntime(); freshErr == nil {
			return postEvent(ctx, info, event)
		}
		return err
	}
	return nil
}

// Describe maps a CLI command path and its params to the client events that
// should be emitted. Activity is emitted before applied so the UI badge shows
// up even if the data refresh is coalesced.
func Describe(path []string, params map[string]any) []Event {
	if len(path) == 0 {
		return nil
	}
	merged := mergeParams(params)
	key := strings.Join(path, " ")
	meta, ok := commandMeta[key]
	if !ok {
		return nil
	}
	tool := "pmock_cli_" + strings.Join(path, "_")
	var events []Event

	action := meta.action
	if key == "app operate" {
		if strings.EqualFold(firstString(merged, "app-op", "appOp"), "uninstall") {
			action = "uninstallApp"
		}
	}
	if action != "" {
		id := firstString(merged, idKeys...)
		code := firstString(merged, codeKeys...)
		if meta.key == "code" {
			id = ""
		} else {
			code = ""
		}
		if id != "" || code != "" {
			events = append(events, Event{
				Kind:      "activity",
				Tool:      tool,
				Action:    action,
				PhoneID:   id,
				PhoneCode: code,
			})
		}
	}

	if len(meta.scopes) > 0 {
		ids := phoneIDs(merged)
		if len(ids) == 0 && meta.key == "id" {
			if s := firstString(merged, idKeys...); s != "" {
				ids = []string{s}
			}
		}
		events = append(events, Event{
			Kind:     "applied",
			Tool:     tool,
			Scopes:   meta.scopes,
			PhoneIDs: ids,
			SettleMS: meta.settleMS,
		})
	}
	return events
}

// DeferredActivity reports whether a command's activity event must wait until
// the command succeeds. Power-state changes (start / restart / reset) are
// deferred: shouting "Agent is powering on" before the request makes the
// desktop client open its remote window while the phone is still off, so the
// window lands on a shutdown screen and never recovers. Every other operation
// keeps the immediate badge semantics.
func DeferredActivity(path []string) bool {
	if len(path) == 0 {
		return false
	}
	meta, ok := commandMeta[strings.Join(path, " ")]
	return ok && meta.afterSuccess
}

// commandMeta is the CLI mirror of the client tool metadata (activity/affects/
// settleMs). Keeping it here means the client does not need to know CLI paths.
type commandMetaEntry struct {
	action   string
	key      string // "id" or "code"; empty means no activity
	scopes   []string
	settleMS int
	// afterSuccess defers the activity event until the API accepted the
	// command. Power-state changes need it: reporting before the request opens
	// the desktop client's remote window onto a phone that is still off.
	afterSuccess bool
}

var commandMeta = map[string]commandMetaEntry{
	// Phone lifecycle / maintenance.
	"phone start":               {action: "powerOn", key: "id", scopes: []string{"cloudPhones"}, settleMS: 5000, afterSuccess: true},
	"phone stop":                {action: "powerOff", key: "id", scopes: []string{"cloudPhones"}, settleMS: 5000},
	"phone restart":             {action: "restart", key: "id", scopes: []string{"cloudPhones"}, settleMS: 5000, afterSuccess: true},
	"phone reset":               {action: "reset", key: "id", scopes: []string{"cloudPhones"}, afterSuccess: true},
	"phone new":                 {action: "oneNew", key: "id", scopes: []string{"cloudPhones"}},
	"phone clear":               {action: "clear", key: "id", scopes: []string{"cloudPhones"}},
	"phone remove":              {action: "delete", key: "id", scopes: []string{"cloudPhones"}},
	"phone create":              {scopes: []string{"cloudPhones"}, settleMS: 8000},
	"phone renew":               {scopes: []string{"cloudPhones"}},
	"phone update-sku":          {scopes: []string{"cloudPhones"}},
	"phone network-set":         {action: "setNetwork", key: "id", scopes: []string{"cloudPhones"}},
	"phone authorize":           {action: "authorize", key: "id", scopes: []string{"cloudPhones"}},
	"phone unauthorize":         {scopes: []string{"cloudPhones"}},
	"phone update-name":         {scopes: []string{"cloudPhones"}},
	"phone update-remark":       {scopes: []string{"cloudPhones"}},
	"phone batch-update-remark": {scopes: []string{"cloudPhones"}},
	"phone update-env-setting":  {scopes: []string{"cloudPhones"}},
	"phone move-group":          {scopes: []string{"groups", "cloudPhones"}},
	"phone update-tags":         {scopes: []string{"tags", "cloudPhones"}},
	"phone batch-update-tags":   {scopes: []string{"tags", "cloudPhones"}},
	"phone transfer send-code":  {scopes: []string{"cloudPhones"}},
	"phone transfer confirm":    {scopes: []string{"cloudPhones"}},
	"phone transfer cancel":     {scopes: []string{"cloudPhones"}},

	// Screen / UI automation.
	"phone screen":              {action: "screenshot", key: "id"},
	"phone screenshot":          {action: "screenshot", key: "id"},
	"phone action click":        {action: "click", key: "code"},
	"phone action double-click": {action: "doubleClick", key: "code"},
	"phone action long-click":   {action: "longClick", key: "code"},
	"phone action swipe":        {action: "swipe", key: "code"},
	"phone action drag":         {action: "drag", key: "code"},
	"phone action input":        {action: "input", key: "code"},
	"phone action app":          {action: "app", key: "code"},

	// App center.
	"app install": {action: "installApp", key: "id", scopes: []string{"cloudPhones"}},
	"app operate": {action: "app", key: "id", scopes: []string{"cloudPhones"}},
	"app upgrade": {scopes: []string{"cloudPhones"}},

	// Inventory (groups / tags).
	"phone group create": {scopes: []string{"groups", "cloudPhones"}},
	"phone group update": {scopes: []string{"groups", "cloudPhones"}},
	"phone group delete": {scopes: []string{"groups", "cloudPhones"}},
	"tag create":         {scopes: []string{"tags", "cloudPhones"}},
	"tag update":         {scopes: []string{"tags", "cloudPhones"}},
	"tag delete":         {scopes: []string{"tags", "cloudPhones"}},

	// File upload moves a file onto the phone; it does not change the list, but
	// it is still a phone-scoped operation worth showing in the activity panel.
	"phone file upload": {action: "unknown", key: "id"},
}

var (
	idKeys   = []string{"id", "phone-id", "phoneId", "phone_id"}
	codeKeys = []string{"code", "phone-code", "phoneCode", "phone_code"}
)

func mergeParams(params map[string]any) map[string]any {
	out := make(map[string]any, len(params)+4)
	for k, v := range params {
		out[k] = v
	}
	raw, ok := out["json"]
	if !ok {
		return out
	}
	var nested map[string]any
	switch v := raw.(type) {
	case string:
		if v == "" || v == "-" || strings.HasPrefix(v, "@") {
			return out
		}
		if json.Unmarshal([]byte(v), &nested) != nil {
			return out
		}
	case map[string]any:
		nested = v
	default:
		return out
	}
	for k, v := range nested {
		if _, exists := out[k]; !exists {
			out[k] = v
		}
	}
	return out
}

func firstString(params map[string]any, keys ...string) string {
	for _, key := range keys {
		if s, ok := valueString(params[key]); ok {
			return s
		}
	}
	return ""
}

func valueString(v any) (string, bool) {
	switch t := v.(type) {
	case string:
		s := strings.TrimSpace(t)
		return s, s != ""
	case float64:
		if math.IsNaN(t) || math.IsInf(t, 0) {
			return "", false
		}
		if t == math.Trunc(t) {
			return strconv.FormatInt(int64(t), 10), true
		}
		return strconv.FormatFloat(t, 'f', -1, 64), true
	case int:
		return strconv.Itoa(t), true
	case int64:
		return strconv.FormatInt(t, 10), true
	case json.Number:
		return t.String(), true
	default:
		return "", false
	}
}

func phoneIDs(params map[string]any) []string {
	var raw any
	for _, key := range []string{"phone-ids", "phoneIds", "phone_ids", "ids", "phone-id", "phoneId", "phone_id"} {
		if v, ok := params[key]; ok {
			raw = v
			break
		}
	}
	if raw == nil {
		return nil
	}
	switch t := raw.(type) {
	case string:
		parts := strings.Split(t, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if s := strings.TrimSpace(p); s != "" {
				out = append(out, s)
			}
		}
		return out
	case []string:
		return t
	case []any:
		out := make([]string, 0, len(t))
		for _, item := range t {
			if s, ok := valueString(item); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		if s, ok := valueString(raw); ok {
			return []string{s}
		}
		return nil
	}
}

func postEvent(ctx context.Context, info runtimeInfo, event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	requestCtx, cancel := context.WithTimeout(ctx, eventTimeout)
	defer cancel()
	url := fmt.Sprintf("http://127.0.0.1:%d/agent-event", info.Port)
	req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+info.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: eventTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("clientbridge: agent-event returned HTTP %d", resp.StatusCode)
	}
	return nil
}

func currentRuntime() (runtimeInfo, error) {
	// The override may change between calls (tests, long-lived processes with
	// a re-pointed env); key the cache so a stale result cannot leak across it.
	key := strings.TrimSpace(os.Getenv(EnvRuntime))
	cacheMu.Lock()
	if cacheValid && cacheKey == key && time.Now().Before(cacheExpires) {
		info := cachedInfo
		err := cachedErr
		cacheMu.Unlock()
		if err != nil {
			return runtimeInfo{}, err
		}
		return info, nil
	}
	cacheMu.Unlock()

	info, err := discoverRuntime()
	if err != nil {
		// Cache the miss too: without it every tool call would rescan the
		// application-support directory tree while the client is closed.
		cacheMu.Lock()
		cachedInfo = runtimeInfo{}
		cachedErr = err
		cacheExpires = time.Now().Add(runtimeCacheTTL)
		cacheValid = true
		cacheKey = key
		cacheMu.Unlock()
		return runtimeInfo{}, err
	}
	cacheMu.Lock()
	cachedInfo = info
	cachedErr = nil
	cacheExpires = time.Now().Add(runtimeCacheTTL)
	cacheValid = true
	cacheKey = key
	cacheMu.Unlock()
	return info, nil
}

func invalidateCache() {
	cacheMu.Lock()
	cacheValid = false
	cachedErr = nil
	cacheKey = ""
	cacheExpires = time.Time{}
	cacheMu.Unlock()
}

func discoverRuntime() (runtimeInfo, error) {
	paths := runtimePaths()
	var best runtimeInfo
	var bestScore int64
	found := false
	var lastErr error
	for _, path := range paths {
		info, score, err := readRuntime(path)
		if err != nil {
			lastErr = err
			continue
		}
		if !found || score > bestScore {
			best = info
			bestScore = score
			found = true
		}
	}
	if !found {
		if lastErr != nil {
			return runtimeInfo{}, fmt.Errorf("clientbridge: no running PMOCK client found: %w", lastErr)
		}
		return runtimeInfo{}, errors.New("clientbridge: no running PMOCK client found")
	}
	return best, nil
}

func readRuntime(path string) (runtimeInfo, int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return runtimeInfo{}, 0, err
	}
	var info runtimeInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return runtimeInfo{}, 0, fmt.Errorf("parse %s: %w", path, err)
	}
	if info.Port <= 0 || info.Port > 65535 {
		return runtimeInfo{}, 0, fmt.Errorf("%s: invalid port %d", path, info.Port)
	}
	if strings.TrimSpace(info.Token) == "" {
		return runtimeInfo{}, 0, fmt.Errorf("%s: missing token", path)
	}
	score := int64(0)
	if t, err := time.Parse(time.RFC3339, info.StartedAt); err == nil {
		score = t.UnixMilli()
	} else if stat, err := os.Stat(path); err == nil {
		score = stat.ModTime().UnixMilli()
	}
	return info, score, nil
}

func runtimePaths() []string {
	if override := strings.TrimSpace(os.Getenv(EnvRuntime)); override != "" {
		if strings.HasSuffix(strings.ToLower(override), ".json") {
			return []string{override}
		}
		return []string{filepath.Join(override, "runtime.json")}
	}

	var roots []string
	switch runtime.GOOS {
	case "darwin":
		if home, err := os.UserHomeDir(); err == nil {
			roots = append(roots, filepath.Join(home, "Library", "Application Support"))
		}
	case "windows":
		for _, env := range []string{"APPDATA", "LOCALAPPDATA"} {
			if dir := strings.TrimSpace(os.Getenv(env)); dir != "" {
				roots = append(roots, dir)
			}
		}
	default:
		if home, err := os.UserHomeDir(); err == nil {
			roots = append(roots, filepath.Join(home, ".config"))
		}
	}

	names := []string{"PMOCK", "PMOCKTest", "MIMO", "MIMOTest"}
	seen := make(map[string]bool)
	var paths []string
	add := func(path string) {
		if path == "" || seen[path] {
			return
		}
		seen[path] = true
		paths = append(paths, path)
	}
	for _, root := range roots {
		if root == "" {
			continue
		}
		for _, name := range names {
			add(filepath.Join(root, name, "mcp", "runtime.json"))
		}
		matches, err := filepath.Glob(filepath.Join(root, "*", "mcp", "runtime.json"))
		if err != nil {
			continue
		}
		for _, match := range matches {
			add(match)
		}
	}
	return paths
}

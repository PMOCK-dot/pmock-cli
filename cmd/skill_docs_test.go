package cmd

import (
	"bytes"
	"context"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func testSkills(t *testing.T) fs.FS {
	t.Helper()
	sub, err := fs.Sub(os.DirFS(".."), "skills")
	if err != nil {
		t.Fatalf("skills tree: %v", err)
	}
	return sub
}

func TestReadReference(t *testing.T) {
	doc, ok := readReference(testSkills(t), "phone list")
	if !ok {
		t.Fatal("no reference page for phone list")
	}
	if !strings.Contains(doc, "# phone list") || !strings.Contains(doc, "## 参数") {
		t.Fatalf("unexpected page:\n%s", doc)
	}
	if _, ok := readReference(testSkills(t), "phone not-a-command"); ok {
		t.Error("unknown command resolved to a page")
	}
	if _, ok := readReference(nil, "phone list"); ok {
		t.Error("nil skills tree resolved to a page")
	}
}

func TestMCPDocsResult(t *testing.T) {
	cat := testCatalog(t)
	s := &mcpServer{catalog: cat, tools: buildMCPTools(cat), skills: testSkills(t)}

	res := s.docsResult("note rpa task add")
	if !res.IsError {
		t.Error("expected an error for an invalid token")
	}
	res = s.docsResult("rpa task add")
	if res.IsError || !strings.Contains(resultText(t, res), "# rpa task add") {
		t.Fatalf("docs page missing: %s", resultText(t, res))
	}
	res = s.docsResult("")
	if res.IsError || !strings.Contains(resultText(t, res), "pmock-phone") {
		t.Fatalf("skill listing missing: %s", resultText(t, res))
	}
}

// A failed call must come back with the reference page attached: examples and
// response fields are what let the model fix its arguments in one round.
func TestMCPFailureAttachesReference(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	t.Setenv("PMOCK_ACCESS_KEY", "AK")
	t.Setenv("PMOCK_SECRET_KEY", "SK")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"traceId":"mcp-fail","code":42001,"msg":"手机不存在"}`))
	}))
	defer srv.Close()
	t.Setenv("PMOCK_BASE_URL", srv.URL+"/api")

	cat := testCatalog(t)
	s := &mcpServer{catalog: cat, tools: buildMCPTools(cat), skills: testSkills(t)}
	tool := mcpToolByName(t, s.tools, "pmock_phone")
	res, _, err := s.callTool(context.Background(), tool, mcpArgs{
		Action: "remove",
		Params: map[string]any{"id": float64(1), "yes": true},
	})
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(t, res)
	if !res.IsError || !strings.Contains(text, "[Reference: pmock-cli phone remove]") {
		t.Fatalf("reference page not attached:\n%s", text)
	}
	if !strings.Contains(text, "## 响应字段") {
		t.Fatalf("reference page looks truncated:\n%s", text)
	}
}

func TestDocsMarkdown(t *testing.T) {
	cmd := newDocsCmd(testCatalog(t), testSkills(t))
	cmd.SilenceErrors, cmd.SilenceUsage = true, true
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"rpa", "task", "add", "--markdown"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "# rpa task add") {
		t.Fatalf("markdown output:\n%s", out.String())
	}

	cmd = newDocsCmd(testCatalog(t), testSkills(t))
	cmd.SilenceErrors, cmd.SilenceUsage = true, true
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"phone", "not-a-command", "--markdown"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected an error for an unknown command")
	}
}

func TestEnumHints(t *testing.T) {
	cat := testCatalog(t)

	hints := enumHints(cat, []string{"phone", "list"}, `{"code":0,"expireStatus":1}`)
	if !strings.Contains(hints, "expireStatus = 1") || !strings.Contains(hints, "1=正常") {
		t.Fatalf("hints = %q", hints)
	}

	// The payload may be nested under data.
	hints = enumHints(cat, []string{"phone", "list"}, `{"code":0,"data":{"expireStatus":0}}`)
	if !strings.Contains(hints, "expireStatus = 0") {
		t.Fatalf("nested hints = %q", hints)
	}

	if hints := enumHints(cat, []string{"phone", "list"}, `not json`); hints != "" {
		t.Errorf("non-JSON payload produced hints: %q", hints)
	}
	if hints := enumHints(cat, []string{"phone", "list"}, `{"code":0}`); hints != "" {
		t.Errorf("payload without enum fields produced hints: %q", hints)
	}
}

// The docs tool must be reachable through a real MCP session, not just as a Go
// function: that is what the agent talks to.
func TestMCPSessionListsDocsTool(t *testing.T) {
	ctx := context.Background()
	cat := testCatalog(t)
	s := &mcpServer{catalog: cat, tools: buildMCPTools(cat), skills: testSkills(t)}

	srv := mcp.NewServer(&mcp.Implementation{Name: "pmock-test", Version: "dev"}, nil)
	s.registerHelp(srv)
	s.registerDocs(srv)
	for _, tool := range s.tools {
		s.registerTool(srv, tool)
	}

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	go func() { _ = srv.Run(ctx, serverTransport) }()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "dev"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer func() { _ = session.Close() }()

	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	names := map[string]bool{}
	for _, tool := range tools.Tools {
		names[tool.Name] = true
	}
	for _, want := range []string{"pmock_docs", "pmock_help", "pmock_phone"} {
		if !names[want] {
			t.Errorf("tool %q is not registered (got %v)", want, names)
		}
	}

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "pmock_docs",
		Arguments: map[string]any{"command": "phone list"},
	})
	if err != nil {
		t.Fatalf("call pmock_docs: %v", err)
	}
	if res.IsError || !strings.Contains(resultText(t, res), "# phone list") {
		t.Fatalf("pmock_docs returned: %s", resultText(t, res))
	}
}

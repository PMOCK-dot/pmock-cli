package cmd

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestComposeScreens(t *testing.T) {
	a := image.NewRGBA(image.Rect(0, 0, 2, 3))
	a.Set(0, 0, color.RGBA{R: 255, A: 255})
	b := image.NewRGBA(image.Rect(0, 0, 4, 2))

	got := composeScreens([]image.Image{a, b})
	if got.Bounds().Dx() != 6 || got.Bounds().Dy() != 3 {
		t.Fatalf("bounds = %v", got.Bounds())
	}
	if px := got.RGBAAt(0, 0); px.R != 255 {
		t.Fatalf("first pixel = %v", px)
	}
}

func TestScreenCompositeCommand(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/open/api/phone/screenCap":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			id := int64(body["id"].(float64))
			_, _ = w.Write([]byte(`{"code":200,"msg":"ok","data":{"url":"` + srv.URL + `/img/` + strconv.FormatInt(id, 10) + `"}}`))
		case strings.HasPrefix(r.URL.Path, "/img/"):
			img := image.NewRGBA(image.Rect(0, 0, 2, 2))
			_ = png.Encode(w, img)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	sc := testCommand(t, testCatalog(t), "phone screen")
	out := filepath.Join(t.TempDir(), "composite.png")
	cmd := newOperationCmd(sc, testGlobals(srv.URL+"/api"))
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	_ = cmd.Flags().Set("ids", "1,2")
	_ = cmd.Flags().Set("out", out)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	res := decodeResult(t, &stdout)
	if res["ok"] != true || res["width"].(float64) != 4 || res["height"].(float64) != 2 {
		t.Fatalf("res = %#v", res)
	}
	info, err := os.Stat(out)
	if err != nil || info.Size() == 0 {
		t.Fatalf("composite not written: %v", err)
	}
}

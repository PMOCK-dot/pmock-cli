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
	"sync"
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

// The composite path builds its own request bodies instead of going through
// buildRequest, so it must reproduce the catalog's isRemove default. Reading
// the flag value alone returns pflag's false default and asks the API to keep
// every capture, which leaves the screenshots in the phone's gallery.
func TestScreenCompositeHonorsRemoveDefault(t *testing.T) {
	t.Setenv("PMOCK_CLI_CONFIG", t.TempDir()+"/none.json")
	var mu sync.Mutex
	var bodies []map[string]any
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/open/api/phone/screenCap" {
			_ = png.Encode(w, image.NewRGBA(image.Rect(0, 0, 2, 2)))
			return
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		mu.Lock()
		bodies = append(bodies, body)
		mu.Unlock()
		_, _ = w.Write([]byte(`{"code":200,"msg":"ok","data":{"url":"` + srv.URL + `/img/1"}}`))
	}))
	defer srv.Close()

	capture := func(setRemove *bool) []map[string]any {
		mu.Lock()
		bodies = nil
		mu.Unlock()

		sc := testCommand(t, testCatalog(t), "phone screen")
		cmd := newOperationCmd(sc, testGlobals(srv.URL+"/api"))
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		_ = cmd.Flags().Set("ids", "1,2")
		_ = cmd.Flags().Set("out", filepath.Join(t.TempDir(), "composite.png"))
		if setRemove != nil {
			_ = cmd.Flags().Set("is-remove", strconv.FormatBool(*setRemove))
		}
		if err := cmd.RunE(cmd, nil); err != nil {
			t.Fatal(err)
		}

		mu.Lock()
		defer mu.Unlock()
		return append([]map[string]any(nil), bodies...)
	}

	defaults := capture(nil)
	if len(defaults) != 2 {
		t.Fatalf("sent %d captures, want 2", len(defaults))
	}
	for _, body := range defaults {
		if body["isRemove"] != true {
			t.Fatalf("default capture sent isRemove = %#v, want true", body["isRemove"])
		}
	}

	no := false
	for _, body := range capture(&no) {
		if body["isRemove"] != false {
			t.Fatalf("explicit capture sent isRemove = %#v, want false", body["isRemove"])
		}
	}
}

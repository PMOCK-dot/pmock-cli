package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/PMOCK-dot/pmock-cli/internal/build"
	"github.com/PMOCK-dot/pmock-cli/internal/client"
	"github.com/PMOCK-dot/pmock-cli/internal/exitcode"
	"github.com/PMOCK-dot/pmock-cli/internal/spec"
)

// screenCapture is one phone's contribution to a composite screenshot.
type screenCapture struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type screenCompositeResult struct {
	OK      bool            `json:"ok"`
	File    string          `json:"file"`
	Width   int             `json:"width"`
	Height  int             `json:"height"`
	Screens []screenCapture `json:"screens"`
}

func addScreenCompositeFlags(c *cobra.Command) {
	if c.Flags().Lookup("ids") == nil {
		c.Flags().Int64Slice("ids", nil, "capture several phones into one composite image; comma-separated")
	}
	if c.Flags().Lookup("out") == nil {
		c.Flags().String("out", "screens.png", "output file for the composite image")
	}
	if c.Flags().Lookup("max-parallel") == nil {
		c.Flags().Int("max-parallel", 4, "maximum concurrent phone captures")
	}
}

// runScreenCommand keeps the single-capture contract of `phone screen --id`
// and adds `--ids`, which stitches several phones into one PNG so a caller can
// verify a whole batch with a single image instead of N downloads.
func runScreenCommand(cmd *cobra.Command, sc *spec.Command, byName map[string]*spec.Flag, g *globals) error {
	ids, _ := cmd.Flags().GetInt64Slice("ids")
	if len(ids) == 0 {
		return runOperation(cmd, sc, byName, g)
	}
	if cmd.Flags().Changed("id") {
		return exitcode.Usagef("--id and --ids are mutually exclusive")
	}
	if g.dryRun {
		return exitcode.Usagef("--dry-run is not supported with --ids; run --dry-run --id once per phone")
	}

	cfg, err := g.config()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL, cfg.AccessKey, cfg.SecretKey, cfg.Timeout)
	c.UserAgent = "pmock-cli/" + build.Version

	parallel, _ := cmd.Flags().GetInt("max-parallel")
	if parallel <= 0 {
		parallel = 4
	}
	if parallel > len(ids) {
		parallel = len(ids)
	}

	ctx, cancel := context.WithTimeout(commandContext(cmd), cfg.Timeout)
	defer cancel()

	isRemove := true
	if f := cmd.Flags().Lookup("is-remove"); f != nil {
		if v, ferr := cmd.Flags().GetBool("is-remove"); ferr == nil {
			isRemove = v
		}
	}
	captures, err := captureScreens(ctx, c, ids, parallel, isRemove)
	if err != nil {
		return err
	}
	imgs := make([]image.Image, 0, len(captures))
	for _, cap := range captures {
		img, derr := downloadImage(ctx, cap.URL)
		if derr != nil {
			return exitcode.Runtimef("phone %d: %v", cap.ID, derr)
		}
		imgs = append(imgs, img)
	}
	composite := composeScreens(imgs)

	out, _ := cmd.Flags().GetString("out")
	if out == "" {
		out = "screens.png"
	}
	if err := writePNG(out, composite); err != nil {
		return err
	}
	abs, err := filepath.Abs(out)
	if err != nil {
		abs = out
	}
	result := screenCompositeResult{
		OK:      true,
		File:    abs,
		Width:   composite.Bounds().Dx(),
		Height:  composite.Bounds().Dy(),
		Screens: captures,
	}
	b, err := json.Marshal(result)
	if err != nil {
		return exitcode.Runtimef("failed to render composite result: %v", err)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", b); err != nil {
		return exitcode.Runtimef("failed to write composite result: %v", err)
	}
	return nil
}

// captureScreens asks the API for each phone's screenshot URL concurrently.
func captureScreens(ctx context.Context, c *client.Client, ids []int64, parallel int, isRemove bool) ([]screenCapture, error) {
	out := make([]screenCapture, len(ids))
	errs := make([]error, len(ids))
	sem := make(chan struct{}, parallel)
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Add(1)
		go func(i int, id int64) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			env, err := c.Do(ctx, client.Request{
				Method: "POST",
				Path:   "/open/api/phone/screenCap",
				Body:   map[string]any{"id": id, "isRemove": isRemove},
			})
			if err != nil {
				errs[i] = fmt.Errorf("phone %d: %w", id, err)
				return
			}
			u, ok := screenURL(env.Data)
			if !ok {
				errs[i] = fmt.Errorf("phone %d: response has no screenshot url", id)
				return
			}
			out[i] = screenCapture{ID: id, URL: u}
		}(i, id)
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func screenURL(raw json.RawMessage) (string, bool) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil && s != "" {
		return s, true
	}
	var obj struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil && obj.URL != "" {
		return obj.URL, true
	}
	return "", false
}

func downloadImage(ctx context.Context, url string) (image.Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download screenshot: HTTP %d", resp.StatusCode)
	}
	img, _, err := image.Decode(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return nil, fmt.Errorf("decode screenshot: %w", err)
	}
	return img, nil
}

// composeScreens lays the screenshots side by side, top aligned, on a black
// canvas so phones with different resolutions still line up.
func composeScreens(imgs []image.Image) *image.RGBA {
	width, height := 0, 0
	for _, img := range imgs {
		b := img.Bounds()
		width += b.Dx()
		if b.Dy() > height {
			height = b.Dy()
		}
	}
	if width <= 0 || height <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	x := 0
	for _, img := range imgs {
		b := img.Bounds()
		draw.Draw(canvas, image.Rect(x, 0, x+b.Dx(), b.Dy()), img, b.Min, draw.Src)
		x += b.Dx()
	}
	return canvas
}

func writePNG(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return exitcode.Runtimef("failed to create %s: %v", filepath.Dir(path), err)
	}
	f, err := os.Create(path)
	if err != nil {
		return exitcode.Runtimef("failed to create %s: %v", path, err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return exitcode.Runtimef("failed to encode %s: %v", path, err)
	}
	return nil
}

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/PMOCK-dot/pmock-cli/internal/client"
	"github.com/PMOCK-dot/pmock-cli/internal/spec"
)

// attachPkgCandidates turns a failed launch of a mistyped package into a
// self-explaining error: business failures stay quiet because the envelope is
// the machine-readable error, so the hint goes to stderr where a human or an
// agent reading the exit code will still see it.
func attachPkgCandidates(ctx context.Context, cmd *cobra.Command, c *client.Client, sc *spec.Command, err error) error {
	if cmd.Flags().Lookup("pkg") == nil || !cmd.Flags().Changed("pkg") {
		return err
	}
	var apiErr *client.APIError
	if !errors.As(err, &apiErr) {
		return err
	}
	pkg, ferr := cmd.Flags().GetString("pkg")
	if ferr != nil || pkg == "" {
		return err
	}
	phoneID, ok := phoneIDFor(ctx, cmd, c)
	if !ok {
		return err
	}
	installed, ierr := installedPackages(ctx, c, phoneID)
	if ierr != nil {
		return err
	}
	cands := matchPackages(pkg, installed)
	if len(cands) == 0 {
		fmt.Fprintf(cmd.ErrOrStderr(),
			"Hint: %q is not installed on phone %d. Run `pmock-cli app list --id %d --only-show-installed --page-size 100` to see packages.\n",
			pkg, phoneID, phoneID)
		return err
	}
	fmt.Fprintf(cmd.ErrOrStderr(),
		"Hint: %q is not installed on phone %d. Installed candidates: %s\n",
		pkg, phoneID, strings.Join(cands, ", "))
	return err
}

// phoneIDFor finds the phone a failing command addresses: --id directly, or
// --code resolved through the phone list.
func phoneIDFor(ctx context.Context, cmd *cobra.Command, c *client.Client) (int64, bool) {
	if cmd.Flags().Lookup("id") != nil {
		if id, err := cmd.Flags().GetInt64("id"); err == nil && id > 0 {
			return id, true
		}
	}
	if cmd.Flags().Lookup("code") == nil {
		return 0, false
	}
	code, err := cmd.Flags().GetString("code")
	if err != nil || code == "" {
		return 0, false
	}
	q := url.Values{}
	q.Set("nameOrCode", code)
	q.Set("pageSize", "10")
	env, err := c.Do(ctx, client.Request{Method: "GET", Path: "/open/api/phone/list", Query: q})
	if err != nil {
		return 0, false
	}
	var payload struct {
		Rows []struct {
			ID   json.RawMessage `json:"id"`
			Code string          `json:"code"`
		} `json:"rows"`
	}
	if err := decodePage(env, &payload); err != nil {
		return 0, false
	}
	for _, row := range payload.Rows {
		if row.Code != code {
			continue
		}
		if id, ok := parseInt64Loose(row.ID); ok && id > 0 {
			return id, true
		}
	}
	return 0, false
}

// installedPackages lists installed packages for one phone.
func installedPackages(ctx context.Context, c *client.Client, phoneID int64) ([]string, error) {
	env, err := c.Do(ctx, client.Request{
		Method: "POST",
		Path:   "/open/api/app/list",
		Body: map[string]any{
			"id":                phoneID,
			"onlyShowInstalled": true,
			"pageNum":           1,
			"pageSize":          strconv.Itoa(100),
		},
	})
	if err != nil {
		return nil, err
	}
	var payload struct {
		Rows []struct {
			Pkg string `json:"apkPkg"`
		} `json:"rows"`
	}
	if err := decodePage(env, &payload); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(payload.Rows))
	for _, row := range payload.Rows {
		if row.Pkg != "" {
			out = append(out, row.Pkg)
		}
	}
	return out, nil
}

// matchPackages narrows the installed list to packages that plausibly answer
// the requested one: same vendor prefix, or the same leaf word (gallery).
func matchPackages(want string, installed []string) []string {
	parts := strings.Split(want, ".")
	prefix := ""
	if len(parts) >= 2 {
		prefix = parts[0] + "." + parts[1]
	}
	leaf := strings.ToLower(parts[len(parts)-1])
	// "photos" vs "gallery" is the classic wrong-guess: treat media/gallery
	// names as one family so the hint lists what is actually installed.
	tokens := map[string]bool{leaf: true}
	for _, group := range [][]string{{"photo", "photos", "gallery", "album"}} {
		for _, alias := range group {
			if alias == leaf {
				for _, other := range group {
					tokens[other] = true
				}
			}
		}
	}
	seen := make(map[string]bool, len(installed))
	out := make([]string, 0, 8)
	for _, p := range installed {
		if p == want || seen[p] {
			continue
		}
		lower := strings.ToLower(p)
		match := prefix != "" && strings.HasPrefix(p, prefix+".")
		for token := range tokens {
			if token != "" && strings.Contains(lower, token) {
				match = true
			}
		}
		if !match {
			continue
		}
		seen[p] = true
		out = append(out, p)
		if len(out) == 8 {
			break
		}
	}
	return out
}

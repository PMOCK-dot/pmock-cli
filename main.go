// pmock-cli is the standalone command line client for the PMOCK open API.
package main

import (
	"embed"
	"os"

	"github.com/PMOCK-dot/pmock-cli/cmd"
)

// assets carries the vendored spec, the product overlay and the generated
// agent skills. Embedding keeps the binary self-contained, which is the main
// reason for choosing Go: the docs an agent reads always match the binary it
// runs, even when the checkout is not around.
//
//go:embed openapi/*.yaml skills
var assetsFS embed.FS

func main() {
	os.Exit(cmd.Execute(assetsFS))
}

// Command gohumanize runs the go-humanize-linter as a standalone
// analysis.Analyzer via singlechecker.Main. This is useful for testing the
// golangci-lint plugin path without golangci-lint itself.
//
// Usage:
//
//	go run ./cmd/gohumanize ./...
package main

import (
	"github.com/larsartmann/go-humanize-linter/plugin"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(plugin.Analyzer)
}

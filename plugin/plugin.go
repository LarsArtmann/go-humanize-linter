// Package plugin exposes go-humanize-linter as a golangci-lint custom linter
// plugin.
//
// # Integration
//
// Add to .golangci.yml:
//
//	plugins:
//	  go-humanize:
//	    path: github.com/larsartmann/go-humanize-linter/plugin
//	    description: Detect hand-rolled reimplementations of go-humanize
//
// Then enable the linter:
//
//	enable:
//	  - gohumanize
//
// # How it works
//
// The Analyzer wraps the same AST detection functions used by the standalone
// CLI. golangci-lint passes a [*analysis.Pass] per package; the Run function
// iterates over pass.Files, finds every [*ast.FuncDecl], applies all seven
// detection rules, and converts [finding.Finding] results to
// [analysis.Diagnostic] via pass.Report.
package plugin

import (
	"go/ast"
	"strings"

	humanizelint "github.com/larsartmann/go-humanize-linter"
	"golang.org/x/tools/go/analysis"
)

// Analyzer is the golangci-lint entry point. Export it so golangci-lint can
// discover it via module plugin loading.
var Analyzer = newAnalyzer()

// newAnalyzer builds the [*analysis.Analyzer] that runs all humanize-lint rules.
func newAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{ //nolint:exhaustruct
		Name: "gohumanize",
		Doc:  "Detect hand-rolled reimplementations of github.com/dustin/go-humanize",
		Run:  run,
	}
}

// run is the analysis.Analyzer.Run implementation. It iterates over every Go
// file in the package, finds function declarations, and applies all detectors.
func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		filePath := pass.Fset.Position(file.Pos()).Filename

		// Skip generated files — they are not hand-written reimplementations.
		if isGenerated(filePath) {
			continue
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			for _, f := range humanizelint.DetectFuncDecl(pass.Fset, file, fn, filePath) {
				pass.Report(analysis.Diagnostic{
					Pos:      fn.Pos(),
					Message:  string(f.Rule) + ": " + f.Message,
					Category: "humanize",
				})
			}
		}
	}

	return nil, nil
}

// isGenerated reports whether a file path looks like generated Go code.
func isGenerated(path string) bool {
	base := path
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		base = path[idx+1:]
	}

	return strings.HasSuffix(base, "_gen.go") ||
		strings.HasSuffix(base, ".gen.go") ||
		strings.HasSuffix(base, "_templ.go")
}

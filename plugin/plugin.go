// Package plugin exposes go-humanize-linter as a golangci-lint v2 custom linter
// module plugin.
//
// # Integration (golangci-lint v2)
//
// Build a custom golangci-lint binary with the plugin compiled in:
//
//	# .custom-gcl.yml
//	version: v2.12.2
//	plugins:
//	  - module: github.com/larsartmann/go-humanize-linter
//	    import: github.com/larsartmann/go-humanize-linter/plugin
//	    path: .
//
// Then run `golangci-lint custom` to produce a binary, and register the
// linter in .golangci.yml:
//
//	linters:
//	  enable:
//	    - gohumanize
//	  settings:
//	    custom:
//	      gohumanize:
//	        type: module
//
// # Standalone usage
//
// The exported [Analyzer] can be used with
// golang.org/x/tools/go/analysis/singlechecker for standalone testing:
//
//	go run ./cmd/gohumanize ./...
package plugin

import (
	"go/ast"
	"strings"

	humanizelint "github.com/larsartmann/go-humanize-linter"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("gohumanize", newPlugin)
}

// pluginSettings holds optional configuration passed via .golangci.yml.
// Currently empty — all 9 rules are enabled by default.
type pluginSettings struct{}

// humanizePlugin implements [register.LinterPlugin] for golangci-lint v2
// module plugin discovery.
type humanizePlugin struct{}

func newPlugin(settings any) (register.LinterPlugin, error) {
	if _, err := register.DecodeSettings[pluginSettings](settings); err != nil {
		return nil, err
	}

	return &humanizePlugin{}, nil
}

func (p *humanizePlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{newAnalyzer()}, nil
}

func (p *humanizePlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

// Analyzer is the standalone entry point. Exported for use with
// golang.org/x/tools/go/analysis/singlechecker (cmd/gohumanize).
var Analyzer = newAnalyzer() //nolint:gochecknoglobals // required by singlechecker

// newAnalyzer builds the [*analysis.Analyzer] that runs all humanize-lint rules.
func newAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{ //nolint:exhaustruct
		Name: "gohumanize",
		Doc:  "Detect hand-rolled reimplementations of github.com/dustin/go-humanize",
		Run:  analyzeHumanize,
	}
}

// analyzeHumanize is the analysis.Analyzer.Run implementation. It iterates over
// every Go file in the package, finds function declarations, and applies all
// detectors.
func analyzeHumanize(pass *analysis.Pass) (any, error) {
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

	return nil, nil //nolint:nilnil // analysis.Analyzer.Run requires (any, error) signature
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

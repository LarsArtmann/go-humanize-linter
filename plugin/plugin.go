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
//	        settings:
//	          enable: "H001,H003"   # only run these rules
//	          disable: "H004"        # skip these rules
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

	"github.com/golangci/plugin-module-register/register"
	humanizelint "github.com/larsartmann/go-humanize-linter"
	"github.com/larsartmann/go-linter-sdk"
	"golang.org/x/tools/go/analysis"
)

// analyzerName is the linter name registered with golangci-lint.
const analyzerName = "gohumanize"

func init() { //nolint:gochecknoinits // required by golangci-lint plugin register API
	register.Plugin(analyzerName, newPlugin)
}

// pluginSettings holds optional configuration passed via .golangci.yml.
// Both fields are comma-separated rule IDs (e.g. "H001,H003").
// When Enable is non-empty, only those rules run.
// When Disable is non-empty, those rules are skipped.
type pluginSettings struct {
	Enable  string `json:"enable"`
	Disable string `json:"disable"`
}

// humanizePlugin implements [register.LinterPlugin] for golangci-lint v2
// module plugin discovery.
type humanizePlugin struct {
	settings pluginSettings
}

func newPlugin(settings any) (register.LinterPlugin, error) { //nolint:ireturn // required by register.NewSettingsPlugin
	s, err := register.DecodeSettings[pluginSettings](settings)
	if err != nil {
		return nil, err
	}

	return &humanizePlugin{settings: s}, nil
}

func (p *humanizePlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	enableSet := parseRuleIDs(p.settings.Enable)
	disableSet := parseRuleIDs(p.settings.Disable)

	rules := filterRules(humanizelint.AllRules(), enableSet, disableSet)
	detector := humanizelint.NewHumanizeDetector(rules...)

	a := &analysis.Analyzer{ //nolint:exhaustruct
		Name: "gohumanize",
		Doc:  "Detect hand-rolled reimplementations of github.com/dustin/go-humanize",
		Run: func(pass *analysis.Pass) (any, error) {
			return runDetector(pass, detector)
		},
	}

	return []*analysis.Analyzer{a}, nil
}

func (p *humanizePlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

// Analyzer is the standalone entry point. Exported for use with
// golang.org/x/tools/go/analysis/singlechecker (cmd/gohumanize).
// All 9 rules are enabled — standalone mode does not support per-rule
// configuration.
var Analyzer = newAnalyzer() //nolint:gochecknoglobals // required by singlechecker

// newAnalyzer builds the [*analysis.Analyzer] that runs all humanize-lint rules.
func newAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{ //nolint:exhaustruct
		Name: "gohumanize",
		Doc:  "Detect hand-rolled reimplementations of github.com/dustin/go-humanize",
		Run:  analyzeHumanize,
	}
}

// analyzeHumanize is the analysis.Analyzer.Run implementation for standalone
// mode (all rules). It iterates over every Go file in the package, finds
// function declarations, and applies all detectors.
func analyzeHumanize(pass *analysis.Pass) (any, error) {
	return runDetector(pass, humanizelint.NewHumanizeDetector())
}

// runDetector applies the given detector to every function declaration in
// every Go file in the package, converting findings to diagnostics.
func runDetector(pass *analysis.Pass, detector *humanizelint.HumanizeDetector) (any, error) {
	for _, file := range pass.Files {
		filePath := pass.Fset.Position(file.Pos()).Filename

		if isGenerated(filePath) {
			continue
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			for _, f := range detector.Run(pass.Fset, file, fn, filePath) {
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

// parseRuleIDs converts a comma-separated string of rule IDs into a set.
// Whitespace around each ID is trimmed; empty entries are ignored.
func parseRuleIDs(spec string) map[string]bool {
	set := make(map[string]bool)

	for _, id := range strings.SplitSeq(spec, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			set[id] = true
		}
	}

	return set
}

// filterRules returns the subset of rules that should run, given the enable
// and disable sets. When enable is non-empty, only those rules are included
// (minus any also disabled). When enable is empty, all rules except disabled
// ones are included.
func filterRules(all []linter.RuleFunc, enable, disable map[string]bool) []linter.RuleFunc {
	if len(enable) == 0 && len(disable) == 0 {
		return all
	}

	var filtered []linter.RuleFunc

	for _, rule := range all {
		if disable[rule.Meta.ID] {
			continue
		}

		if len(enable) > 0 && !enable[rule.Meta.ID] {
			continue
		}

		filtered = append(filtered, rule)
	}

	return filtered
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

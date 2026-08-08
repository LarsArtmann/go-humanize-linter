// Package plugin exposes go-humanize-linter as a golangci-lint v2 custom linter
// module plugin.
//
// # Integration (golangci-lint v2)
//
// 1. Create a .custom-gcl.yml build config:
//
//	# .custom-gcl.yml
//	version: v2.12.2
//	plugins:
//	  - module: github.com/larsartmann/go-humanize-linter
//	    import: github.com/larsartmann/go-humanize-linter/plugin
//	    path: .
//
// 2. Build a custom golangci-lint binary:
//
//	golangci-lint custom   # produces ./custom-gcl
//
// 3. Register the linter in your .golangci.yml (CRITICAL — without this section,
// golangci-lint reports "unknown linters: gohumanize"):
//
//	linters:
//	  enable:
//	    - gohumanize
//	  settings:
//	    custom:
//	      gohumanize:
//	        type: "module"
//	        description: "Detect hand-rolled reimplementations of go-humanize"
//	        settings:
//	          enable: "H001,H003"        # optional: only run these rules
//	          # disable: "H004"          # optional: skip these rules
//	          # minConfidence: "medium"  # optional: low (default), medium, high, full
//	          # verifySuppressions: true # optional: report stale //nolint directives
//
// 4. Run the custom binary:
//
//	./custom-gcl run ./...
//
// See .golangci.custom.yml in the repo root for a complete example config.
//
// # Standalone usage
//
// The exported [Analyzer] can be used with
// golang.org/x/tools/go/analysis/singlechecker for standalone testing:
//
//	go run ./cmd/gohumanize ./...
package plugin

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/gotoken"
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
	Enable             string `json:"enable"`
	Disable            string `json:"disable"`
	MinConfidence      string `json:"minConfidence"`
	VerifySuppressions bool   `json:"verifySuppressions"`
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

	rules := linter.FilterRules(humanizelint.AllRules(), enableSet, disableSet)
	detector := humanizelint.NewHumanizeDetector(rules...)

	minConf, err := finding.ParseConfidence(p.settings.MinConfidence)
	if err != nil {
		return nil, fmt.Errorf("parse min-confidence: %w", err)
	}

	verify := p.settings.VerifySuppressions

	a := &analysis.Analyzer{ //nolint:exhaustruct
		Name: analyzerName,
		Doc:  "Detect hand-rolled reimplementations of github.com/dustin/go-humanize",
		Run: func(pass *analysis.Pass) (any, error) {
			return runDetector(pass, detector, minConf, verify)
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
		Name: analyzerName,
		Doc:  "Detect hand-rolled reimplementations of github.com/dustin/go-humanize",
		Run:  analyzeHumanize,
	}
}

// analyzeHumanize is the analysis.Analyzer.Run implementation for standalone
// mode (all rules). It iterates over every Go file in the package, finds
// function declarations, and applies all detectors.
func analyzeHumanize(pass *analysis.Pass) (any, error) {
	return runDetector(pass, humanizelint.NewHumanizeDetector(), finding.ConfidenceLow, false)
}

// runDetector applies the given detector to every function declaration in every
// Go file in the package. Generated files are skipped via IsGeneratedFile (nil
// content skips the content phase). When verifySuppressions is true, a separate
// pass checks //nolint:gohumanize directives for staleness or misspelled linter
// names. Findings below minConfidence are filtered out before reporting.
func runDetector(
	pass *analysis.Pass,
	detector *humanizelint.HumanizeDetector,
	minConfidence finding.Confidence,
	verifySuppressions bool,
) (any, error) {
	tokenFiles := make(map[string]*token.File)

	for _, file := range pass.Files {
		filePath := pass.Fset.Position(file.Pos()).Filename
		tokenFiles[filePath] = pass.Fset.File(file.Pos())
	}

	var allFindings []finding.Finding

	for _, file := range pass.Files {
		filePath := pass.Fset.Position(file.Pos()).Filename

		if humanizelint.IsGeneratedFile(filePath, nil) {
			continue
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			allFindings = append(allFindings, detector.Run(pass.Fset, file, fn, filePath)...)
		}
	}

	if verifySuppressions {
		report := finding.NewReportFromFindings(
			finding.ToolInfo{Name: "go-humanize-linter"}, //nolint:exhaustruct
			allFindings,
		)
		verifyFindings := humanizelint.VerifySuppressionsInFiles(pass.Fset, pass.Files, report)
		allFindings = append(allFindings, verifyFindings...)
	}

	for _, f := range allFindings {
		// Suppression-verification findings (H0SUP) bypass the confidence
		// filter. They are meta-diagnostics about directive health, not code
		// patterns, and must always be reported when verifySuppressions is on.
		if string(f.Rule) != humanizelint.RuleIDH0SUP && f.Confidence.Compare(minConfidence) < 0 {
			continue
		}

		pass.Report(analysis.Diagnostic{
			Pos:      gotoken.LineColToPos(tokenFiles[string(f.Position.File)], f.Position.Line, f.Position.Column),
			Message:  string(f.Rule) + ": " + f.Message,
			Category: "humanize",
		})
	}

	return nil, nil //nolint:nilnil // analysis.Analyzer.Run requires (any, error) signature
}

// parseRuleIDs converts a comma-separated string of rule IDs into a set.
// Whitespace around each ID is trimmed; empty entries are ignored.
func parseRuleIDs(spec string) map[string]bool {
	set := make(map[string]bool)

	for id := range strings.SplitSeq(spec, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			set[id] = true
		}
	}

	return set
}


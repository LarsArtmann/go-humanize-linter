package plugin

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/gotoken"
	humanizelint "github.com/larsartmann/go-humanize-linter"
	"github.com/larsartmann/go-linter-sdk"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestParseRuleIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty", "", []string{}},
		{"single", "H001", []string{"H001"}},
		{"multiple", "H001,H003", []string{"H001", "H003"}},
		{"whitespace", " H001 , H003 ", []string{"H001", "H003"}},
		{"trailing comma", "H001,", []string{"H001"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := parseRuleIDs(tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf("parseRuleIDs(%q) = %d entries, want %d", tt.input, len(got), len(tt.want))
			}

			for _, id := range tt.want {
				if !got[id] {
					t.Errorf("parseRuleIDs(%q) missing %s", tt.input, id)
				}
			}
		})
	}
}

func TestFilterRules_NoConfigReturnsAll(t *testing.T) {
	t.Parallel()

	all := humanizelint.AllRules()
	got := linter.FilterRules(all, nil, nil)

	if len(got) != len(all) {
		t.Errorf("no config: got %d rules, want %d", len(got), len(all))
	}
}

func TestFilterRules_EnableOnly(t *testing.T) {
	t.Parallel()

	all := humanizelint.AllRules()
	enable := map[string]bool{humanizelint.RuleIDH001: true, humanizelint.RuleIDH003: true}
	got := linter.FilterRules(all, enable, nil)

	if len(got) != 2 {
		t.Fatalf("enable=H001,H003: got %d rules, want 2", len(got))
	}

	for _, rule := range got {
		if !enable[rule.Meta.ID] {
			t.Errorf("unexpected rule %s in filtered set", rule.Meta.ID)
		}
	}
}

func TestFilterRules_DisableOnly(t *testing.T) {
	t.Parallel()

	all := humanizelint.AllRules()
	disable := map[string]bool{humanizelint.RuleIDH001: true}
	got := linter.FilterRules(all, nil, disable)

	if len(got) != len(all)-1 {
		t.Fatalf("disable=H001: got %d rules, want %d", len(got), len(all)-1)
	}

	for _, rule := range got {
		if disable[rule.Meta.ID] {
			t.Errorf("disabled rule %s should not be in filtered set", rule.Meta.ID)
		}
	}
}

func TestFilterRules_EnableOverridesDisable(t *testing.T) {
	t.Parallel()

	all := humanizelint.AllRules()
	enable := map[string]bool{humanizelint.RuleIDH001: true, humanizelint.RuleIDH002: true}
	disable := map[string]bool{humanizelint.RuleIDH001: true}
	got := linter.FilterRules(all, enable, disable)

	if len(got) != 1 {
		t.Fatalf("enable=H001,H002 disable=H001: got %d rules, want 1 (only H002)", len(got))
	}

	if got[0].Meta.ID != humanizelint.RuleIDH002 {
		t.Errorf("expected H002, got %s", got[0].Meta.ID)
	}
}

func TestNewPluginWithSettings(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(map[string]any{
		"enable":             "H001,H002",
		"disable":            "H003",
		"minConfidence":      "high",
		"verifySuppressions": true,
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	impl, ok := plug.(*humanizePlugin)
	if !ok {
		t.Fatalf("expected *humanizePlugin, got %T", plug)
	}

	if impl.settings.Enable != "H001,H002" {
		t.Errorf("expected Enable=%q, got %q", "H001,H002", impl.settings.Enable)
	}

	if impl.settings.Disable != "H003" {
		t.Errorf("expected Disable=%q, got %q", "H003", impl.settings.Disable)
	}

	if impl.settings.MinConfidence != "high" {
		t.Errorf("expected MinConfidence=%q, got %q", "high", impl.settings.MinConfidence)
	}

	if !impl.settings.VerifySuppressions {
		t.Error("expected VerifySuppressions=true")
	}
}

func TestBuildAnalyzers_MinConfidence(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(map[string]any{
		"minConfidence": "medium",
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	analyzers, err := plug.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers with minConfidence=medium failed: %v", err)
	}

	if len(analyzers) != 1 {
		t.Fatalf("expected 1 analyzer, got %d", len(analyzers))
	}
}

func TestBuildAnalyzers_InvalidMinConfidence(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(map[string]any{
		"minConfidence": "bogus",
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	if _, err := plug.BuildAnalyzers(); err == nil {
		t.Fatal("expected error for invalid confidence level, got nil")
	}
}

func TestBuildAnalyzers_VerifySuppressions(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(map[string]any{
		"verifySuppressions": true,
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	analyzers, err := plug.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers with verifySuppressions failed: %v", err)
	}

	if len(analyzers) != 1 {
		t.Fatalf("expected 1 analyzer, got %d", len(analyzers))
	}
}

func TestNewPluginWithEmptySettings_NewFieldsDefault(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(nil)
	if err != nil {
		t.Fatalf("newPlugin(nil) failed: %v", err)
	}

	impl, ok := plug.(*humanizePlugin)
	if !ok {
		t.Fatalf("expected *humanizePlugin, got %T", plug)
	}

	if impl.settings.MinConfidence != "" {
		t.Errorf("empty settings: MinConfidence should be empty, got %q", impl.settings.MinConfidence)
	}

	if impl.settings.VerifySuppressions {
		t.Error("empty settings: VerifySuppressions should be false")
	}
}

func TestNewPluginWithEmptySettings(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(nil)
	if err != nil {
		t.Fatalf("newPlugin(nil) failed: %v", err)
	}

	impl, ok := plug.(*humanizePlugin)
	if !ok {
		t.Fatalf("expected *humanizePlugin, got %T", plug)
	}

	if impl.settings.Enable != "" || impl.settings.Disable != "" {
		t.Errorf("empty settings should produce empty plugin config, got Enable=%q Disable=%q",
			impl.settings.Enable, impl.settings.Disable)
	}
}

func TestBuildAnalyzersProducesConfiguredDetector(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(map[string]any{
		"enable": humanizelint.RuleIDH001,
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	analyzers, err := plug.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers failed: %v", err)
	}

	if len(analyzers) != 1 {
		t.Fatalf("expected 1 analyzer, got %d", len(analyzers))
	}

	if analyzers[0].Name != "gohumanize" {
		t.Errorf("expected analyzer name 'gohumanize', got %q", analyzers[0].Name)
	}

	if analyzers[0].Run == nil {
		t.Fatal("analyzer.Run must not be nil")
	}
}

func TestLineColToPos(t *testing.T) {
	t.Parallel()

	src := "package main\n\nfunc foo() {\n\t_ = 1024\n}\n"
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "test.go", src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parser.ParseFile failed: %v", err)
	}

	tf := fset.File(f.Pos())
	tokenFiles := map[string]*token.File{"test.go": tf}

	line3Start := tf.LineStart(3)

	tests := []struct {
		name      string
		finding   finding.Finding
		wantNoPos bool
		want      token.Pos
	}{
		{
			name:      "valid position line 3 col 1",
			finding:   makeTestFinding("test.go", 3, 1), //nolint:exhaustruct
			want:      line3Start,
			wantNoPos: false,
		},
		{
			name:      "valid position line 3 col 5",
			finding:   makeTestFinding("test.go", 3, 5), //nolint:exhaustruct
			want:      line3Start + 4,
			wantNoPos: false,
		},
		{
			name:      "missing file returns NoPos",
			finding:   makeTestFinding("other.go", 3, 1), //nolint:exhaustruct
			wantNoPos: true,
		},
		{
			name:      "zero line returns NoPos",
			finding:   makeTestFinding("test.go", 0, 0), //nolint:exhaustruct
			wantNoPos: true,
		},
		{
			name:      "negative line returns NoPos",
			finding:   makeTestFinding("test.go", -1, 0), //nolint:exhaustruct
			wantNoPos: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := gotoken.LineColToPos(
				tokenFiles[string(tt.finding.Position.File)],
				tt.finding.Position.Line,
				tt.finding.Position.Column,
			)

			if tt.wantNoPos {
				if got != token.NoPos {
					t.Errorf("expected NoPos, got %d", got)
				}

				return
			}

			if got != tt.want {
				t.Errorf("expected %d, got %d", tt.want, got)
			}
		})
	}
}

func TestLineColToPos_OutOfRangeLine(t *testing.T) {
	t.Parallel()

	src := "package main\n\nfunc foo() {\n\t_ = 1024\n}\n"
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "test.go", src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parser.ParseFile failed: %v", err)
	}

	tf := fset.File(f.Pos())
	tokenFiles := map[string]*token.File{"test.go": tf}

	// Line 1000 is far beyond the file's 5 lines. Without the LineCount
	// guard, token.File.LineStart panics: "illegal line number (past end of
	// file)". A linter must never crash — it must return NoPos instead.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("gotoken.LineColToPos panicked on out-of-range line: %v", r)
		}
	}()

	got := gotoken.LineColToPos(tokenFiles["test.go"], 1000, 1)

	if got != token.NoPos {
		t.Errorf("expected NoPos for out-of-range line, got %d", got)
	}
}

func makeTestFinding(file string, line, col int) finding.Finding {
	return finding.Finding{ //nolint:exhaustruct
		Position: finding.Position{ //nolint:exhaustruct
			File:   finding.FilePath(file),
			Line:   line,
			Column: col,
		},
	}
}

func TestRunDetector_H0SUPBypassesConfidenceFilter(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(map[string]any{
		"minConfidence":      "full",
		"verifySuppressions": true,
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	analyzers, err := plug.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers failed: %v", err)
	}

	testdata := filepath.Join("..", "testdata", "analysistest")
	analysistest.Run(t, testdata, analyzers[0], "./h0supbypass")
}

func TestRunDetector_FiltersByConfidence(t *testing.T) {
	t.Parallel()

	// H001 on h001positive is ConfidenceFull (KMGTPE trick). With
	// minConfidence "full", it should still fire — proving the filter
	// passes full-confidence findings through.
	plug, err := newPlugin(map[string]any{
		"minConfidence": "full",
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	analyzers, err := plug.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers failed: %v", err)
	}

	testdata := filepath.Join("..", "testdata", "analysistest")
	analysistest.Run(t, testdata, analyzers[0], "./h001positive")
}

func TestRunDetector_FiltersOutMediumConfidence(t *testing.T) {
	t.Parallel()

	// h002mediumconfidence triggers H002 via the fallback path, which is
	// ConfidenceMedium. With minConfidence "full" it must be filtered out.
	plug, err := newPlugin(map[string]any{
		"minConfidence": "full",
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	analyzers, err := plug.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers failed: %v", err)
	}

	testdata := filepath.Join("..", "testdata", "analysistest")
	analysistest.Run(t, testdata, analyzers[0], "./h002mediumconfidence")
}

func TestRunDetector_VerifySuppressions(t *testing.T) {
	t.Parallel()

	// With verifySuppressions enabled at default confidence, H0SUP fires
	// on stale directives. The h0supbypass fixture has a stale //nolint
	// directive on a clean function.
	plug, err := newPlugin(map[string]any{
		"verifySuppressions": true,
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	analyzers, err := plug.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers failed: %v", err)
	}

	testdata := filepath.Join("..", "testdata", "analysistest")
	analysistest.Run(t, testdata, analyzers[0], "./h0supbypass")
}

func TestRunDetector_InBodySuppression(t *testing.T) {
	t.Parallel()

	// The h001_inbody_suppressed fixture has an H001 pattern (KMGTPE trick)
	// with a //nolint:gohumanize directive INSIDE the function body. The
	// plugin path must honour it — no diagnostic expected.
	plug, err := newPlugin(map[string]any{})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	analyzers, err := plug.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers failed: %v", err)
	}

	testdata := filepath.Join("..", "testdata", "analysistest")
	analysistest.Run(t, testdata, analyzers[0], "./h001_inbody_suppressed")
}

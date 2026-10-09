package humanizelint //nolint:testpackage // white-box: tests unexported suppression helpers

import (
	"os"
	"testing"
)

func TestSuppressesGohumanize(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		src  string
		want bool
	}{
		{"bare gohumanize", "//nolint:gohumanize", true},
		{"gohumanize scoped", "//nolint:gohumanize:H001", true},
		{"all", "//nolint:all", true},
		{"bare rule id", "//nolint:H001", true},
		{"other linter", "//nolint:other", false},
		{"other scoped to our rule", "//nolint:other:H001", true},
		{"not a directive", "// regular comment", false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			suppressed := suppressedRules(tt.src)
			if suppressed == nil {
				suppressed = []string{}
			}

			if got := suppressesGohumanize(suppressed); got != tt.want {
				t.Errorf("suppressesGohumanize(%q) = %v, want %v", tt.src, got, tt.want)
			}
		})
	}
}

func TestTargetsGohumanize(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		src  string
		want bool
	}{
		{"bare gohumanize", "//nolint:gohumanize", true},
		{"gohumanize scoped", "//nolint:gohumanize:H001", true},
		{"bare rule id", "//nolint:H001", true},
		{"multi-linter with gohumanize", "//nolint:gohumanize,gofmt", true},
		{"all", "//nolint:all", false},
		{"bare nolint", "//nolint", false},
		{"other linter", "//nolint:other", false},
		{"not a directive", "// regular comment", false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			suppressed := suppressedRules(tt.src)
			if suppressed == nil {
				suppressed = []string{}
			}

			if got := targetsGohumanize(suppressed); got != tt.want {
				t.Errorf("targetsGohumanize(%q) = %v, want %v", tt.src, got, tt.want)
			}
		})
	}
}

// TestVerifySuppressions_BareAllNotJudged guards the blanket-suppression
// contract: a //nolint:all directive usually exists for another linter's
// benefit, so H0SUP must not call it stale just because gohumanize itself has
// no finding in that function — while a directive that names gohumanize still
// is judged.
func TestVerifySuppressions_BareAllNotJudged(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	writeFile(t, dir, `package main

//nolint:all // golangci-lint autofix bug workaround
func blanketForOtherLinter() {}

//nolint:gohumanize
func staleGohumanizeDirective() {}

func main() {
	blanketForOtherLinter()
	staleGohumanizeDirective()
}
`)

	findings, err := VerifySuppressions(dir)
	if err != nil {
		t.Fatalf("VerifySuppressions failed: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("expected exactly 1 suppression finding (the stale gohumanize directive), got %d: %+v", len(findings), findings)
	}

	if string(findings[0].Rule) != RuleIDH0SUP {
		t.Errorf("expected rule %s, got %s", RuleIDH0SUP, findings[0].Rule)
	}
}

func TestHasUnknownHumanizeLinterName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		src  string
		want bool
	}{
		{"correct name", "//nolint:gohumanize", false},
		{"module path name", "//nolint:go-humanize-linter", true},
		{"module path scoped", "//nolint:go-humanize-linter/H003", true},
		{"all", "//nolint:all", false},
		{"other linter", "//nolint:other", false},
		{"mixed with correct", "//nolint:gohumanize,go-humanize-linter", true},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			suppressed := suppressedRules(tt.src)
			if suppressed == nil {
				suppressed = []string{}
			}

			if got := hasUnknownHumanizeLinterName(suppressed); got != tt.want {
				t.Errorf("hasUnknownHumanizeLinterName(%q) = %v, want %v", tt.src, got, tt.want)
			}
		})
	}
}

func TestVerifySuppressions_UnknownLinterName(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	writeFile(t, dir, `package main

//nolint:go-humanize-linter/H003
func staleSuppression() {}

func main() {
	staleSuppression()
}
`)

	findings, err := VerifySuppressions(dir)
	if err != nil {
		t.Fatalf("VerifySuppressions failed: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("expected 1 suppression finding, got %d: %+v", len(findings), findings)
	}

	if string(findings[0].Rule) != RuleIDH0SUP {
		t.Errorf("expected rule %s, got %s", RuleIDH0SUP, findings[0].Rule)
	}
}

func TestVerifySuppressions_StaleSuppression(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	writeFile(t, dir, `package main

//nolint:gohumanize:H001
func noByteFormatting() string {
	return "hello"
}

func main() {
	_ = noByteFormatting()
}
`)

	findings, err := VerifySuppressions(dir)
	if err != nil {
		t.Fatalf("VerifySuppressions failed: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("expected 1 stale suppression finding, got %d: %+v", len(findings), findings)
	}

	if string(findings[0].Rule) != RuleIDH0SUP {
		t.Errorf("expected rule %s, got %s", RuleIDH0SUP, findings[0].Rule)
	}
}

// TestVerifySuppressions_UsedSuppression is the regression test for the
// staleness false positive: the detection pipeline drops suppressed findings
// before any report is built, so verification that matched directives
// against a run report flagged every working directive as stale. The
// verdict must come from an unsuppressed re-detection instead — here the
// fixture carries a real H001 pattern, the directive suppresses it in a
// normal run, and verification must accept it.
func TestVerifySuppressions_UsedSuppression(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	writeFile(t, dir, `package main

import "fmt"

//nolint:gohumanize:H001
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func main() {
	_ = formatBytes(1024)
}
`)

	findings, err := VerifySuppressions(dir)
	if err != nil {
		t.Fatalf("VerifySuppressions failed: %v", err)
	}

	if len(findings) != 0 {
		t.Fatalf("expected 0 findings for used suppression, got %d: %+v", len(findings), findings)
	}
}

func TestVerifySuppressions_InBodyStale(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	writeFile(t, dir, `package main

func cleanFunc() string {
	//nolint:gohumanize
	return "hello"
}

func main() {
	_ = cleanFunc()
}
`)

	findings, err := VerifySuppressions(dir)
	if err != nil {
		t.Fatalf("VerifySuppressions failed: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("expected 1 stale suppression finding, got %d: %+v", len(findings), findings)
	}

	if string(findings[0].Rule) != RuleIDH0SUP {
		t.Errorf("expected rule %s, got %s", RuleIDH0SUP, findings[0].Rule)
	}
}

func TestVerifySuppressions_InBodyUsed(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	writeFile(t, dir, `package main

import "fmt"

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp]) //nolint:gohumanize:H001
}

func main() {
	_ = formatBytes(1024)
}
`)

	findings, err := VerifySuppressions(dir)
	if err != nil {
		t.Fatalf("VerifySuppressions failed: %v", err)
	}

	if len(findings) != 0 {
		t.Fatalf("expected 0 findings for used in-body suppression, got %d: %+v", len(findings), findings)
	}
}

func writeFile(t *testing.T, dir, content string) {
	t.Helper()

	if err := os.WriteFile(dir+"/main.go", []byte(content), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}

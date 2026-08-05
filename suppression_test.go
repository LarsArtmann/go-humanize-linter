package humanizelint //nolint:testpackage // white-box: tests unexported suppression helpers

import (
	"os"
	"testing"

	"github.com/larsartmann/go-finding"
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

	writeFile(t, dir, "main.go", `package main

//nolint:go-humanize-linter/H003
func staleSuppression() {}

func main() {
	staleSuppression()
}
`)

	report := finding.NewReport(finding.ToolInfo{Name: "go-humanize-linter"})
	findings, err := VerifySuppressions(dir, report)
	if err != nil {
		t.Fatalf("VerifySuppressions failed: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("expected 1 suppression finding, got %d: %+v", len(findings), findings)
	}

	if string(findings[0].Rule) != suppressionVerificationRuleID {
		t.Errorf("expected rule %s, got %s", suppressionVerificationRuleID, findings[0].Rule)
	}
}

func TestVerifySuppressions_StaleSuppression(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	writeFile(t, dir, "main.go", `package main

//nolint:gohumanize:H001
func noByteFormatting() string {
	return "hello"
}

func main() {
	_ = noByteFormatting()
}
`)

	report := finding.NewReport(finding.ToolInfo{Name: "go-humanize-linter"})
	findings, err := VerifySuppressions(dir, report)
	if err != nil {
		t.Fatalf("VerifySuppressions failed: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("expected 1 stale suppression finding, got %d: %+v", len(findings), findings)
	}

	if string(findings[0].Rule) != suppressionVerificationRuleID {
		t.Errorf("expected rule %s, got %s", suppressionVerificationRuleID, findings[0].Rule)
	}
}

func TestVerifySuppressions_UsedSuppression(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	writeFile(t, dir, "main.go", `package main

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

	// Simulate a report that contains an H001 finding at the function line.
	report := finding.NewReport(finding.ToolInfo{Name: "go-humanize-linter"})
	report.AddFinding(finding.NewBuilder(
		finding.RuleName(RuleIDH001),
		finding.ToolName("go-humanize-linter"),
		"manual byte-size formatting",
		finding.SeverityWarning,
		finding.Pos(finding.FilePath(dir+"/main.go"), 6, 1),
	).MustBuild())
	report.ComputeSummary()

	findings, err := VerifySuppressions(dir, report)
	if err != nil {
		t.Fatalf("VerifySuppressions failed: %v", err)
	}

	if len(findings) != 0 {
		t.Fatalf("expected 0 findings for used suppression, got %d: %+v", len(findings), findings)
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(dir+"/"+name, []byte(content), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}

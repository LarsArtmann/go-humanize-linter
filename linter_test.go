package humanizelint_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsartmann/go-finding"
	humanizelint "github.com/larsartmann/go-humanize-linter"
	"github.com/larsartmann/go-linter-sdk"
)

func testdataDir(t *testing.T, name string) string {
	t.Helper()

	return filepath.Join("testdata", name)
}

func runRule(t *testing.T, rule linter.RuleFunc, dir string) []finding.Finding {
	t.Helper()

	findings, err := rule.Run(context.Background(), dir)
	if err != nil {
		t.Fatalf("rule %s failed: %v", rule.Meta.ID, err)
	}

	return findings
}

func ruleIDs(findings []finding.Finding) []string {
	ids := make([]string, 0, len(findings))
	for _, f := range findings {
		ids = append(ids, string(f.Rule))
	}

	return ids
}

// ---------------------------------------------------------------------------
// H001 — manual-bytes-format
// ---------------------------------------------------------------------------

func TestRuleBytes_KMGTPETrick(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_bytes_kmgtptrick"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}

	if findings[0].Rule != "H001" {
		t.Errorf("expected rule H001, got %s", findings[0].Rule)
	}
}

func TestRuleBytes_SwitchConstants(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_bytes_switch"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

func TestRuleBytes_UnitSlice(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_bytes_unitslice"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

func TestRuleBytes_Negative(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_negative"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings on clean code, got %d: %+v", len(findings), findings)
	}
}

func TestCLI_SuppressedByDirective(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_suppressed"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings when //nolint:gohumanize present, got %d: %+v", len(findings), findings)
	}
}

func TestCLI_SuppressedByDirectiveCommaList(t *testing.T) {
	t.Parallel()

	// Verifies that a comma-list directive like
	// "//nolint:gohumanize,unused" still suppresses gohumanize findings
	// when other linter names are listed alongside it.
	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_suppressed_comma"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings with comma-list //nolint, got %d: %+v", len(findings), findings)
	}
}

// ---------------------------------------------------------------------------
// H002 — manual-comma-format
// ---------------------------------------------------------------------------

func TestRuleComma_Modulo3(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleComma(), testdataDir(t, "h002_comma_mod3"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

func TestRuleComma_StepBy3(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleComma(), testdataDir(t, "h002_comma_step3"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

func TestRuleComma_FallbackNamedConstant(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleComma(), testdataDir(t, "h002_comma_fallback"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
	// Fallback path should have medium confidence.
	if findings[0].Confidence > finding.ConfidenceMedium {
		t.Errorf("expected medium confidence, got %v", findings[0].Confidence)
	}
}

// ---------------------------------------------------------------------------
// H003 — manual-reltime-format
// ---------------------------------------------------------------------------

func TestRuleRelTime_Positive(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleRelTime(), testdataDir(t, "h003_reltime"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

func TestRuleRelTime_Negative(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleRelTime(), testdataDir(t, "h003_negative"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d: %+v", len(findings), findings)
	}
}

// ---------------------------------------------------------------------------
// H004 — manual-plural
// ---------------------------------------------------------------------------

func TestRulePlural_NamedParams(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RulePlural(), testdataDir(t, "h004_plural_params"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

func TestRuleComma_Negative(t *testing.T) {
	t.Parallel()

	// formatSSN writes separators but lacks mod-3 / step-by-3 / for-loop +
	// digit-conversion patterns. H002 must not fire.
	findings := runRule(t, humanizelint.RuleComma(), testdataDir(t, "h002_negative"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings on H002 negative fixture, got %d: %+v", len(findings), findings)
	}
}

func TestRulePlural_IfOne(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RulePlural(), testdataDir(t, "h004_plural_if"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

func TestRulePlural_Negative(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RulePlural(), testdataDir(t, "h004_negative"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings for non-string-returning if==1, got %d: %+v", len(findings), findings)
	}
}

// ---------------------------------------------------------------------------
// H005 — manual-si-format
// ---------------------------------------------------------------------------

func TestRuleSI_Positive(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleSI(), testdataDir(t, "h005_si"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

func TestRuleSI_Negative(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleSI(), testdataDir(t, "h005_negative"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d: %+v", len(findings), findings)
	}
}

// ---------------------------------------------------------------------------
// H006 — manual-ftoa
// ---------------------------------------------------------------------------

func TestRuleFtoa_Positive(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleFtoa(), testdataDir(t, "h006_ftoa"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

// ---------------------------------------------------------------------------
// H009 — manual-commaf
// ---------------------------------------------------------------------------

func TestRuleCommaf_Positive(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleCommaf(), testdataDir(t, "h009_commaf"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}

	if findings[0].Rule != "H009" {
		t.Errorf("expected rule H009, got %s", findings[0].Rule)
	}
}

func TestRuleCommaf_Negative(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleCommaf(), testdataDir(t, "h009_negative"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings on H009 negative fixture, got %d: %+v", len(findings), findings)
	}
}

// ---------------------------------------------------------------------------
// H008 — manual-ordinal
// ---------------------------------------------------------------------------

func TestRuleOrdinal_Positive(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleOrdinal(), testdataDir(t, "h008_ordinal"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}

	if findings[0].Rule != "H008" {
		t.Errorf("expected rule H008, got %s", findings[0].Rule)
	}
}

func TestRuleOrdinal_Negative(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleOrdinal(), testdataDir(t, "h008_negative"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings on H008 negative fixture, got %d: %+v", len(findings), findings)
	}
}

// ---------------------------------------------------------------------------
// H007 — manual-parse-bytes
// ---------------------------------------------------------------------------

func TestRuleParseBytes_SuffixChecks(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleParseBytes(), testdataDir(t, "h007_parsebytes_suffix"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

func TestRuleParseBytes_MultiplierMap(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleParseBytes(), testdataDir(t, "h007_parsebytes_map"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
}

func TestRuleParseBytes_PackageVarMultiplierMap(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleParseBytes(), testdataDir(t, "h007_package_var"))

	found := false

	for _, f := range findings {
		if string(f.Rule) == humanizelint.RuleIDH007 && strings.Contains(f.Message, "package-level") {
			found = true
		}
	}

	if !found {
		t.Fatalf("expected package-level multiplier map finding, got %d: %+v", len(findings), findings)
	}
}

// ---------------------------------------------------------------------------
// Cross-rule: clean code produces no findings from any rule
// ---------------------------------------------------------------------------

func TestAllRules_CleanCode(t *testing.T) {
	t.Parallel()
	dir := testdataDir(t, "clean")

	for _, rule := range humanizelint.AllRules() {
		findings := runRule(t, rule, dir)
		if len(findings) != 0 {
			t.Errorf("rule %s produced %d findings on clean code: %+v",
				rule.Meta.ID, len(findings), findings)
		}
	}
}

// ---------------------------------------------------------------------------
// Registry integration
// ---------------------------------------------------------------------------

func TestDefaultRegistry_AllRules(t *testing.T) {
	t.Parallel()

	r := humanizelint.DefaultRegistry()
	if len(r.All()) != 9 {
		t.Fatalf("expected 9 rules, got %d", len(r.All()))
	}
}

func TestDefaultRegistry_RunOnRelTime(t *testing.T) {
	t.Parallel()

	r := humanizelint.DefaultRegistry()

	report, err := r.Run(context.Background(), testdataDir(t, "h003_reltime"))
	if err != nil {
		t.Fatal(err)
	}

	if report.Len() != 1 {
		t.Fatalf("expected 1 finding, got %d", report.Len())
	}
}

func TestDefaultRegistry_RunOnClean(t *testing.T) {
	t.Parallel()

	r := humanizelint.DefaultRegistry()

	report, err := r.Run(context.Background(), testdataDir(t, "clean"))
	if err != nil {
		t.Fatal(err)
	}

	if report.Len() != 0 {
		t.Fatalf("expected 0 findings on clean code, got %d", report.Len())
	}
}

func TestDefaultRegistry_RunOnKMGTPETrick(t *testing.T) {
	t.Parallel()

	r := humanizelint.DefaultRegistry()

	report, err := r.Run(context.Background(), testdataDir(t, "h001_bytes_kmgtptrick"))
	if err != nil {
		t.Fatal(err)
	}

	ids := ruleIDs(report.FindingsSnapshot())
	if len(ids) != 1 || ids[0] != "H001" {
		t.Fatalf("expected exactly H001, got %v", ids)
	}
}

// ---------------------------------------------------------------------------
// Exit code
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Self-scan: the linter must produce zero findings on its own source
// ---------------------------------------------------------------------------

func TestLintsItself_Clean(t *testing.T) {
	t.Parallel()

	r := humanizelint.DefaultRegistry()

	report, err := r.Run(context.Background(), ".")
	if err != nil {
		t.Fatalf("self-scan failed: %v", err)
	}

	if report.Len() != 0 {
		t.Fatalf("expected 0 findings on own source, got %d: %+v",
			report.Len(), report.FindingsSnapshot())
	}
}

func TestExitCode_CleanIsZero(t *testing.T) {
	t.Parallel()

	r := humanizelint.DefaultRegistry()

	report, err := r.Run(context.Background(), testdataDir(t, "clean"))
	if err != nil {
		t.Fatal(err)
	}

	if code := linter.ExitCodeFromReport(report); code != 0 {
		t.Errorf("expected exit 0, got %d", code)
	}
}

func TestExitCode_FindingsIsOne(t *testing.T) {
	t.Parallel()

	r := humanizelint.DefaultRegistry()

	report, err := r.Run(context.Background(), testdataDir(t, "h002_comma_mod3"))
	if err != nil {
		t.Fatal(err)
	}

	if code := linter.ExitCodeFromReport(report); code != 1 {
		t.Errorf("expected exit 1, got %d", code)
	}
}

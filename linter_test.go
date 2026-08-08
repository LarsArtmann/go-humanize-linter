package humanizelint_test

import (
	"context"
	"path/filepath"
	"slices"
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

func TestRuleBytes_SizeBucket_NoFalsePositive(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_sizebucket"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings on size-bucket lookup code, got %d: %+v", len(findings), findings)
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

func TestCLI_SuppressedByInBodyDirective(t *testing.T) {
	t.Parallel()

	// Verifies that a //nolint:gohumanize directive placed INSIDE the function
	// body (on a specific statement line) is honoured as a suppression.
	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_suppressed_inbody"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings with in-body //nolint, got %d: %+v", len(findings), findings)
	}
}

func TestCLI_ScopedSuppression_H001Only(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_scoped_h001"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings with //nolint:gohumanize:H001, got %d: %+v", len(findings), findings)
	}
}

func TestCLI_ScopedSuppression_DoesNotSuppressOtherRule(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_scoped_h002_only"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding (H002 scope must NOT suppress H001), got %d: %+v", len(findings), findings)
	}
}

func TestCLI_InBodyDirectiveIsolation(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleBytes(), testdataDir(t, "h001_inbody_isolation"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding (func B only; func A suppressed by in-body //nolint), got %d: %+v", len(findings), findings)
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

func TestRuleRelTime_DotImport(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleRelTime(), testdataDir(t, "h003_dot_import"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
	}
	// Dot-imported time constants should still produce full confidence.
	if findings[0].Confidence != finding.ConfidenceFull {
		t.Errorf("expected full confidence for dot-imported time threshold, got %v", findings[0].Confidence)
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

// TestRuleComma_StringsJoinSpace_NoFalsePositive is a regression test for the
// issue where H002/H009 fired on functions that joined CLI args with a single
// space (" "). A space is not a thousands separator, so this pattern must
// stay clean even when paired with a for-loop and strconv.Itoa (the fallback
// heuristic's other triggers). See pattern_comma.go and the feedback doc at
// docs/feedback/new/2026-08-05_h002-h009-strings-join-space-false-positive.md.
func TestRuleComma_StringsJoinSpace_NoFalsePositive(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleComma(), testdataDir(t, "h002_strings_join_space"))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings on strings.Join(args, \" \") fixture, got %d: %+v", len(findings), findings)
	}
}

func TestRuleCommaf_StringsJoinSpace_NoFalsePositive(t *testing.T) {
	t.Parallel()

	findings := runRule(t, humanizelint.RuleCommaf(), testdataDir(t, "h002_strings_join_space"))
	if len(findings) != 0 {
		t.Fatalf(
			"expected 0 findings on strings.Join(args, \" \") fixture for H009, got %d: %+v",
			len(findings),
			findings,
		)
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

func TestRuleParseBytes_AliasedImport(t *testing.T) {
	t.Parallel()

	// Verifies that H007 detects byte-unit suffix checks through an aliased
	// import (e.g. `import str "strings"` → str.HasSuffix).
	findings := runRule(t, humanizelint.RuleParseBytes(), testdataDir(t, "h007_aliased_import"))
	if len(findings) == 0 {
		t.Fatal("expected H007 finding for aliased import, got 0")
	}
}

func TestRuleParseBytes_DotImport(t *testing.T) {
	t.Parallel()

	// Verifies that H007 detects byte-unit suffix checks through a dot
	// import (e.g. `import . "strings"` → bare HasSuffix calls).
	findings := runRule(t, humanizelint.RuleParseBytes(), testdataDir(t, "h007_dot_import"))
	if len(findings) == 0 {
		t.Fatal("expected H007 finding for dot import, got 0")
	}
}

func TestH009H002_NoOverlap(t *testing.T) {
	t.Parallel()

	dir := testdataDir(t, "h009_h002_overlap")
	r := humanizelint.DefaultRegistry()

	report, err := r.Run(context.Background(), dir)
	if err != nil {
		t.Fatalf("registry.Run failed: %v", err)
	}

	var ids []string

	for f := range report.All() {
		ids = append(ids, string(f.Rule))
	}

	hasH009 := slices.Contains(ids, "H009")
	hasH002 := slices.Contains(ids, "H002")

	if !hasH009 {
		t.Errorf("expected H009 finding (float+comma function), got rules: %v", ids)
	}

	if !hasH002 {
		t.Errorf("expected H002 finding (integer-comma function), got rules: %v", ids)
	}

	// Verify H009 and H002 don't fire on the SAME function. The overlap
	// fixture has two functions: one float+comma (should fire H009 only)
	// and one integer+comma (should fire H002 only). Each should fire exactly once.
	countRule := func(ruleID string) int {
		matches := 0

		for _, id := range ids {
			if id == ruleID {
				matches++
			}
		}

		return matches
	}

	h009Count := countRule("H009")
	h002Count := countRule("H002")

	if h009Count != 1 {
		t.Errorf("expected exactly 1 H009 finding, got %d (rules: %v)", h009Count, ids)
	}

	if h002Count != 1 {
		t.Errorf("expected exactly 1 H002 finding, got %d (rules: %v)", h002Count, ids)
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

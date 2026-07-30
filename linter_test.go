package humanizelint_test

import (
	"context"
	"path/filepath"
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
	ids := make([]string, len(findings))
	for i, f := range findings {
		ids[i] = string(f.Rule)
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

func TestRulePlural_IfOne(t *testing.T) {
	t.Parallel()
	findings := runRule(t, humanizelint.RulePlural(), testdataDir(t, "h004_plural_if"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %+v", len(findings), findings)
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
	if len(r.All()) != 6 {
		t.Fatalf("expected 6 rules, got %d", len(r.All()))
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

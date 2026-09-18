package humanizelint //nolint:testpackage // white-box: tests collectWordSeriesEvidence

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// wordSeriesEvidenceFor parses src (a full Go file) and returns the H012
// evidence collected for its first function declaration.
func wordSeriesEvidenceFor(t *testing.T, src string) wordSeriesEvidence {
	t.Helper()

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "fixture.go", src, 0)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		return collectWordSeriesEvidence(fn, buildImportAliases(file))
	}

	t.Fatal("fixture has no function declaration")

	return wordSeriesEvidence{}
}

// TestCollectWordSeriesEvidence_ConjunctionPosition guards the conjunction
// position filter at the unit level: a conjunction literal (" and ") counts
// as word-series evidence only in a join position — a strings.Join separator
// or a concatenation operand. The same literal in prose (a Sprintf format
// string, a bare string variable) must not count even when a comma join is
// present. This is the discrimination proof behind the h012_negative
// fixtures; loosening the filter fails this test before the corpus false
// positive returns.
func TestCollectWordSeriesEvidence_ConjunctionPosition(t *testing.T) {
	t.Parallel()

	t.Run("prose in Sprintf format next to comma join stays clean", func(t *testing.T) {
		t.Parallel()

		evidence := wordSeriesEvidenceFor(t, `package p

import (
	"fmt"
	"strings"
)

func summary(names []string) string {
	return fmt.Sprintf("open and closed issues: %s", strings.Join(names, ", "))
}`)

		if !evidence.anyCommaJoin {
			t.Fatal("fixture must contain a comma join for the discrimination to have teeth")
		}

		if evidence.conjunctionLiteral {
			t.Error(`prose "and" in a Sprintf format string must not count as conjunction evidence`)
		}
	})

	t.Run("prose literal in variable next to comma join stays clean", func(t *testing.T) {
		t.Parallel()

		evidence := wordSeriesEvidenceFor(t, `package p

import (
	"strings"
)

func heading(names []string) string {
	title := "cause and effect"
	return title + ": " + strings.Join(names, ", ")
}`)

		if !evidence.anyCommaJoin {
			t.Fatal("fixture must contain a comma join for the discrimination to have teeth")
		}

		if evidence.conjunctionLiteral {
			t.Error(`prose "and" in a bare string variable must not count as conjunction evidence`)
		}
	})

	t.Run("Join separator counts as conjunction evidence", func(t *testing.T) {
		t.Parallel()

		evidence := wordSeriesEvidenceFor(t, `package p

import (
	"strings"
)

func pair(parts []string) string {
	return strings.Join(parts, " and ")
}`)

		if !evidence.conjunctionLiteral {
			t.Error(`strings.Join(parts, " and ") must count as conjunction evidence`)
		}
	})

	t.Run("concatenation operand counts as conjunction evidence", func(t *testing.T) {
		t.Parallel()

		evidence := wordSeriesEvidenceFor(t, `package p

func pair(a, b string) string {
	return a + " and " + b
}`)

		if !evidence.conjunctionLiteral {
			t.Error(`+ " and " + must count as conjunction evidence`)
		}
	})
}

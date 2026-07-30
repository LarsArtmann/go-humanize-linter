package humanizelint_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	humanizelint "github.com/larsartmann/go-humanize-linter"
)

// ExampleDefaultRegistry shows building a registry with every rule enabled by
// default — the typical starting point for running the full linter.
func ExampleDefaultRegistry() {
	registry := humanizelint.DefaultRegistry()

	fmt.Println("registered rules:", len(registry.All()))

	// Output: registered rules: 9
}

// ExampleAllRules lists the stable rule IDs in registration order.
func ExampleAllRules() {
	for _, rule := range humanizelint.AllRules() {
		fmt.Println(rule.Meta.ID)
	}

	// Output:
	// H001
	// H002
	// H003
	// H004
	// H005
	// H006
	// H007
	// H008
	// H009
}

// ExampleRuleBytes inspects a single rule factory's metadata.
func ExampleRuleBytes() {
	rule := humanizelint.RuleBytes()

	fmt.Println(rule.Meta.ID, rule.Meta.Name)

	// Output: H001 manual-bytes-format
}

// ExampleRuleComma shows the H002 manual-comma-format rule.
func ExampleRuleComma() {
	rule := humanizelint.RuleComma()

	fmt.Println(rule.Meta.ID, rule.Meta.Name)

	// Output: H002 manual-comma-format
}

// ExampleRuleRelTime shows the H003 manual-reltime-format rule.
func ExampleRuleRelTime() {
	rule := humanizelint.RuleRelTime()

	fmt.Println(rule.Meta.ID, rule.Meta.Name)

	// Output: H003 manual-reltime-format
}

// ExampleRulePlural shows the H004 manual-plural rule.
func ExampleRulePlural() {
	rule := humanizelint.RulePlural()

	fmt.Println(rule.Meta.ID, rule.Meta.Name)

	// Output: H004 manual-plural
}

// ExampleRuleSI shows the H005 manual-si-format rule.
func ExampleRuleSI() {
	rule := humanizelint.RuleSI()

	fmt.Println(rule.Meta.ID, rule.Meta.Name)

	// Output: H005 manual-si-format
}

// ExampleRuleFtoa shows the H006 manual-ftoa rule.
func ExampleRuleFtoa() {
	rule := humanizelint.RuleFtoa()

	fmt.Println(rule.Meta.ID, rule.Meta.Name)

	// Output: H006 manual-ftoa
}

// ExampleRuleParseBytes shows the H007 manual-parse-bytes rule.
func ExampleRuleParseBytes() {
	rule := humanizelint.RuleParseBytes()

	fmt.Println(rule.Meta.ID, rule.Meta.Name)

	// Output: H007 manual-parse-bytes
}

// ExampleRuleOrdinal shows the H008 manual-ordinal rule.
func ExampleRuleOrdinal() {
	rule := humanizelint.RuleOrdinal()

	fmt.Println(rule.Meta.ID, rule.Meta.Name)

	// Output: H008 manual-ordinal
}

// ExampleRuleCommaf shows the H009 manual-commaf rule.
func ExampleRuleCommaf() {
	rule := humanizelint.RuleCommaf()

	fmt.Println(rule.Meta.ID, rule.Meta.Name)

	// Output: H009 manual-commaf
}

// ExampleDetectFuncDecl runs the detectors against a hand-rolled byte formatter
// parsed from source and reports how many findings the KMGTPE trick produces.
func ExampleDetectFuncDecl() {
	const src = `package main

import "fmt"

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
`

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "demo.go", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	findings := 0

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		findings += len(humanizelint.DetectFuncDecl(fset, file, fn, "demo.go"))
	}

	fmt.Println("findings:", findings)

	// Output: findings: 1
}

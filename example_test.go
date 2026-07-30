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

	// Output: registered rules: 7
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
}

// ExampleRuleBytes inspects a single rule factory's metadata.
func ExampleRuleBytes() {
	rule := humanizelint.RuleBytes()

	fmt.Println(rule.Meta.ID, rule.Meta.Name)

	// Output: H001 manual-bytes-format
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

package humanizelint

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func mustParse(b *testing.B, src string) *ast.File {
	b.Helper()

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "bench.go", src, parser.ParseComments)
	if err != nil {
		b.Fatalf("parse error: %v", err)
	}

	return file
}

const aliasedSrc = `
package main

import (
	str "strings"
	"fmt"
)

func parse(s string) bool {
	return str.HasSuffix(s, "KB") || str.HasPrefix(s, "MB")
}

func main() { fmt.Println(parse("x")) }
`

const standardSrc = `
package main

import (
	"strings"
	"fmt"
)

func parse(s string) bool {
	return strings.HasSuffix(s, "KB") || strings.HasPrefix(s, "MB")
}

func main() { fmt.Println(parse("x")) }
`

func BenchmarkIsPackageCall_WithoutAliases(b *testing.B) {
	file := mustParse(b, standardSrc)

	var count int

	calls := collectCalls(file)

	b.ResetTimer()

	for b.Loop() {
		for _, call := range calls {
			if isPackageCall(call, "strings", "HasSuffix") {
				count++
			}
		}
	}

	b.StopTimer()

	if count == 0 {
		b.Fatal("expected at least one match")
	}
}

func BenchmarkIsPackageCall_WithAliases(b *testing.B) {
	file := mustParse(b, aliasedSrc)

	aliases := buildImportAliases(file)
	calls := collectCalls(file)

	var count int

	b.ResetTimer()

	for b.Loop() {
		for _, call := range calls {
			if isPackageCall(call, "strings", "HasSuffix", aliases) {
				count++
			}
		}
	}

	b.StopTimer()

	if count == 0 {
		b.Fatal("expected at least one match")
	}
}

func BenchmarkBuildImportAliases(b *testing.B) {
	file := mustParse(b, aliasedSrc)

	var result map[string]string

	b.ResetTimer()

	for b.Loop() {
		result = buildImportAliases(file)
	}

	if len(result) == 0 {
		b.Fatal("expected non-empty aliases map")
	}
}

func collectCalls(file *ast.File) []*ast.CallExpr {
	var calls []*ast.CallExpr

	ast.Inspect(file, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			calls = append(calls, call)
		}

		return true
	})

	return calls
}

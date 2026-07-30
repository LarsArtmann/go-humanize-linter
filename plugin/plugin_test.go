package plugin_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	humanizelint "github.com/larsartmann/go-humanize-linter"
	"github.com/larsartmann/go-humanize-linter/plugin"
)

func TestAnalyzerNotNil(t *testing.T) {
	if plugin.Analyzer == nil {
		t.Fatal("plugin.Analyzer must not be nil")
	}

	if plugin.Analyzer.Name != "gohumanize" {
		t.Fatalf("expected analyzer name 'gohumanize', got %q", plugin.Analyzer.Name)
	}

	if plugin.Analyzer.Run == nil {
		t.Fatal("plugin.Analyzer.Run must not be nil")
	}

	if plugin.Analyzer.Doc == "" {
		t.Fatal("plugin.Analyzer.Doc must not be empty")
	}
}

func TestDetectFuncDeclBytesPositive(t *testing.T) {
	src := `package main

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

func main() {}
`

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "main.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		findings := humanizelint.DetectFuncDecl(fset, file, fn, "main.go")
		if len(findings) == 0 {
			continue
		}

		// Expect at least one H001 byte-format finding.
		foundH001 := false

		for _, f := range findings {
			if string(f.Rule) == "H001" {
				foundH001 = true
			}
		}

		if !foundH001 {
			t.Fatalf("expected H001 byte-format finding, got: %v", findings)
		}

		return
	}

	t.Fatal("no findings at all — expected H001")
}

func TestDetectFuncDeclCleanNegative(t *testing.T) {
	src := `package main

import "fmt"

func hashID(data []byte) string {
	return fmt.Sprintf("%x", data[:8])
}

func main() {}
`

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "main.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		findings := humanizelint.DetectFuncDecl(fset, file, fn, "main.go")
		if len(findings) != 0 {
			t.Fatalf("expected 0 findings for clean code, got %d: %v", len(findings), findings)
		}
	}
}

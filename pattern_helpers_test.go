package humanizelint //nolint:testpackage // white-box: tests unexported detection helpers

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestNoLintMatches(t *testing.T) {
	t.Parallel()

	cases := []struct {
		comment string
		want    []string
	}{
		{"//nolint", []string{"all"}},
		{"//nolint:all", []string{"all"}},
		{"//nolint:gohumanize", []string{"gohumanize"}},
		{"//nolint:gohumanize,other", []string{"gohumanize", "other"}},
		{"//nolint:other,gohumanize", []string{"other", "gohumanize"}},
		{"//nolint:gohumanize // trailing reason text", []string{"gohumanize"}},
		{"//nolint:other", []string{"other"}},
		{"//nolint :", []string{}},
		{"//lint:ignore gohumanize", []string{"gohumanize"}},
		{"//lint:ignore gohumanize:H001", []string{"gohumanize", "H001"}},
		{"//lint:ignore gohumanize reason text", []string{"gohumanize"}},
		{"// regular comment", nil},
		{"//notanolint", nil},
		{"", nil},
	}

	for _, tt := range cases {
		if got := suppressedRules(tt.comment); !equalSlices(got, tt.want) {
			t.Errorf("suppressedRules(%q) = %v, want %v", tt.comment, got, tt.want)
		}
	}
}

// equalSlices compares two []string for equality (order-independent since
// suppression checks are order-independent).
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestHasNoLintDirective(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		src  string
		want bool
	}{
		{
			name: "leading directive",
			src: `package main

//nolint:gohumanize
func f() {}
`,
			want: true,
		},
		{
			name: "trailing directive",
			src: `package main

func f() {} //nolint:gohumanize
`,
			want: true,
		},
		{
			name: "multi-line doc comment with directive",
			src: `package main

// f does things.
//
//nolint:gohumanize
func f() {}
`,
			want: true,
		},
		{
			name: "bare nolint",
			src: `package main

//nolint
func f() {}
`,
			want: true,
		},
		{
			name: "no directive",
			src: `package main

// just a doc comment
func f() {}
`,
			want: false,
		},
		{
			name: "directive for other linter",
			src: `package main

//nolint:other
func f() {}
`,
			want: false,
		},
		{
			name: "directive separated by blank line",
			src: `package main

//nolint:gohumanize

func f() {}
`,
			want: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fset, file, fn := parseFirstFunc(t, tt.src)
			if got := hasNoLintDirective(fset, file, fn); got != tt.want {
				t.Errorf("hasNoLintDirective = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:gohumanize // test fixtures legitimately contain byte-unit names (KB, MB) to verify H001 detection
func TestHasConst1024(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		src  string
		want bool
	}{
		{
			name: "in-function const",
			src: `package main

func f(b int64) int64 {
	const unit = 1024

	return b / unit
}
`,
			want: true,
		},
		{
			name: "file-level const",
			src: `package main

const kb = 1024

func f(b int64) int64 {
	return b / kb
}
`,
			want: true,
		},
		{
			name: "mul chain 1024*1024",
			src: `package main

func f(b int64) int64 {
	const mb = 1024 * 1024

	return b / mb
}
`,
			want: true,
		},
		{
			name: "no power-of-1024 const",
			src: `package main

func f(b int64) int64 {
	const unit = 1000

	return b / unit
}
`,
			want: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, file, fn := parseFirstFunc(t, tt.src)
			if got := hasConst1024(file, fn); got != tt.want {
				t.Errorf("hasConst1024 = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasStepBy3(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		src  string
		want bool
	}{
		{
			name: "add-assign i += 3",
			src: `package main

func f(n int) int {
	for i := 0; i < n; i += 3 {
		_ = i
	}

	return 0
}
`,
			want: true,
		},
		{
			name: "add binary i = i + 3",
			src: `package main

func f(n int) {
	for i := 0; i < n; i = i + 3 {
		_ = i
	}
}
`,
			want: true,
		},
		{
			name: "step by 5 no match",
			src: `package main

func f(n int) {
	for i := 0; i < n; i += 5 {
		_ = i
	}
}
`,
			want: false,
		},
		{
			name: "no loop",
			src: `package main

func f() int { return 3 }
`,
			want: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, _, fn := parseFirstFunc(t, tt.src)
			if got := hasStepBy3(fn); got != tt.want {
				t.Errorf("hasStepBy3 = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasEqualsOneBranch(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		src  string
		want bool
	}{
		{
			name: "if n == 1 with string literal in branch",
			src: `package main

func f(n int) string {
	if n == 1 {
		return "one"
	}

	return "many"
}
`,
			want: true,
		},
		{
			name: "if 1 == n (literal on left)",
			src: `package main

func f(n int) string {
	if 1 == n {
		return "one"
	}

	return "many"
}
`,
			want: true,
		},
		{
			name: "if n != 1 (NEQ operator)",
			src: `package main

func f(n int) string {
	if n != 1 {
		return "many"
	}

	return "one"
}
`,
			want: true,
		},
		{
			name: "if n == 1 with string concat",
			src: `package main

func f(n int) string {
	if n == 1 {
		return "item-" + "1"
	}

	return "items-many"
}
`,
			want: true,
		},
		{
			name: "if n == 1 no string in branch (FP filter)",
			src: `package main

func f(n int) int {
	if n == 1 {
		return 42
	}

	return 100
}
`,
			want: false,
		},
		{
			name: "if n == 2 (not equal to 1)",
			src: `package main

func f(n int) string {
	if n == 2 {
		return "two"
	}

	return "other"
}
`,
			want: false,
		},
		{
			name: "if len(s) == 1 (function call, not simple ident)",
			src: `package main

func f(s string) string {
	if len(s) == 1 {
		return "single"
	}

	return "many"
}
`,
			want: false,
		},
		{
			name: "no if at all",
			src: `package main

func f(n int) string { return "x" }
`,
			want: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, _, fn := parseFirstFunc(t, tt.src)
			if got := hasEqualsOneBranch(fn); got != tt.want {
				t.Errorf("hasEqualsOneBranch = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasCommaOrSeparator(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		src  string
		want bool
	}{
		{
			name: "WriteString comma",
			src: `package main

import "strings"

func f() {
	var b strings.Builder
	b.WriteString(",")
}
`,
			want: true,
		},
		{
			name: "WriteString dot",
			src: `package main

import "strings"

func f() {
	var b strings.Builder
	b.WriteString(".")
}
`,
			want: true,
		},
		{
			name: "WriteByte rune comma",
			src: `package main

import "strings"

func f() {
	var b strings.Builder
	b.WriteByte(',')
}
`,
			want: true,
		},
		{
			name: "WriteRune rune comma",
			src: `package main

import "strings"

func f() {
	var b strings.Builder
	b.WriteRune(',')
}
`,
			want: true,
		},
		{
			name: "strings.Join with comma",
			src: `package main

import "strings"

func f() string {
	return strings.Join([]string{"a", "b"}, ",")
}
`,
			want: true,
		},
		{
			name: "strings.Join with space",
			src: `package main

import "strings"

func f() string {
	return strings.Join([]string{"a", "b"}, " ")
}
`,
			want: true,
		},
		{
			name: "WriteString no separator",
			src: `package main

import "strings"

func f() {
	var b strings.Builder
	b.WriteString("hello")
}
`,
			want: false,
		},
		{
			name: "no WriteString at all",
			src: `package main

func f() string { return "x" }
`,
			want: false,
		},
		{
			name: "strings.Join with non-separator",
			src: `package main

import "strings"

func f() string {
	return strings.Join([]string{"a", "b"}, "X")
}
`,
			want: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, _, fn := parseFirstFunc(t, tt.src)
			if got := hasCommaOrSeparator(fn); got != tt.want {
				t.Errorf("hasCommaOrSeparator = %v, want %v", got, tt.want)
			}
		})
	}
}

func parseFirstFunc(t *testing.T, src string) (*token.FileSet, *ast.File, *ast.FuncDecl) {
	t.Helper()

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "main.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			return fset, file, fn
		}
	}

	t.Fatal("no FuncDecl in source")

	return nil, nil, nil
}

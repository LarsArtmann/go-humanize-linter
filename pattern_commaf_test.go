package humanizelint //nolint:testpackage // white-box: tests unexported format-string helpers

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestWalkFormatFloatVerbs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		val  string
		want bool
	}{
		{"dotted precision %.2f", "%.2f", true},
		{"dotted precision %.0f", "%.0f", true},
		{"dotted precision %.10f", "%.10f", true},
		{"embedded in string", "value: %.3f units", true},
		{"bare width %3f", "%3f", true},
		{"escaped percent then real verb", "100%% done %.2f", true},
		{"plain %f no precision", "%f", false},
		{"escaped percent only", "100%%", false},
		{"empty string", "", false},
		{"no percent at all", "hello world", false},
		{"integer verb %d", "%d", false},
		{"g verb %.2g", "%.2g", false},
		{"s verb %s", "%s", false},
		{"trailing percent", "50%", false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := walkFormatFloatVerbs(tt.val); got != tt.want {
				t.Errorf("walkFormatFloatVerbs(%q) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}

func TestScanDottedPercentFloat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		val  string
		i    int
		want bool
	}{
		{"%.2f at 0", "%.2f", 0, true},
		{"%.0f at 0", "%.0f", 0, true},
		{"%.10f at 0", "%.10f", 0, true},
		{"%.22f at 0", "%.22f", 0, true},
		{"%.2f embedded at offset", "x%.2f", 1, true},
		{"%.f no digit after dot", "%.f", 0, false},
		{"%.2g wrong verb", "%.2g", 0, false},
		{"%.2 no verb", "%.2", 0, false},
		{"%.2d wrong verb", "%.2d", 0, false},
		{"too short", "%.", 0, false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.i >= len(tt.val) {
				t.Fatalf("test index %d out of range for %q", tt.i, tt.val)
			}

			if got := scanDottedPercentFloat(tt.val, tt.i); got != tt.want {
				t.Errorf("scanDottedPercentFloat(%q, %d) = %v, want %v", tt.val, tt.i, got, tt.want)
			}
		})
	}
}

func TestScanBarePercentFloat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		val  string
		i    int
		want bool
	}{
		{"%3f at 0", "%3f", 0, true},
		{"%0f at 0", "%0f", 0, true},
		{"%33f repeated digit", "%33f", 0, true},
		{"%34f multi-digit width", "%34f", 0, true},
		{"%10f two-digit width", "%10f", 0, true},
		{"%3f embedded at offset", "x%3f", 1, true},
		{"%3d wrong verb", "%3d", 0, false},
		{"%3 no verb", "%3", 0, false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.i >= len(tt.val) {
				t.Fatalf("test index %d out of range for %q", tt.i, tt.val)
			}

			if got := scanBarePercentFloat(tt.val, tt.i); got != tt.want {
				t.Errorf("scanBarePercentFloat(%q, %d) = %v, want %v", tt.val, tt.i, got, tt.want)
			}
		})
	}
}

func TestHasFormatFloatPrecision(t *testing.T) {
	t.Parallel()

	intLit := func(v string) ast.Expr {
		return &ast.BasicLit{Kind: token.INT, Value: v}
	}

	cases := []struct {
		name string
		call *ast.CallExpr
		want bool
	}{
		{
			name: "precision 2",
			call: &ast.CallExpr{Args: []ast.Expr{nil, nil, intLit("2")}},
			want: true,
		},
		{
			name: "precision 10",
			call: &ast.CallExpr{Args: []ast.Expr{nil, nil, intLit("10")}},
			want: true,
		},
		{
			name: "precision 0 excluded",
			call: &ast.CallExpr{Args: []ast.Expr{nil, nil, intLit("0")}},
			want: false,
		},
		{
			name: "precision -1 excluded",
			call: &ast.CallExpr{Args: []ast.Expr{nil, nil, intLit("-1")}},
			want: false,
		},
		{
			name: "string literal not int",
			call: &ast.CallExpr{Args: []ast.Expr{nil, nil, &ast.BasicLit{Kind: token.STRING, Value: `"2"`}}},
			want: false,
		},
		{
			name: "nil arg",
			call: &ast.CallExpr{Args: []ast.Expr{nil, nil, nil}},
			want: false,
		},
		{
			name: "too few args",
			call: &ast.CallExpr{Args: []ast.Expr{nil}},
			want: false,
		},
		{
			name: "no args",
			call: &ast.CallExpr{Args: nil},
			want: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := hasFormatFloatPrecision(tt.call); got != tt.want {
				t.Errorf("hasFormatFloatPrecision = %v, want %v", got, tt.want)
			}
		})
	}
}

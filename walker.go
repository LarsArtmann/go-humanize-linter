package humanizelint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/larsartmann/go-finding"
)

// ParsedFile holds a parsed Go source file together with its token file set.
type ParsedFile struct {
	Path string
	Fset *token.FileSet
	File *ast.File
}

// WalkError is returned by WalkGoDir and checkFuncDecls when the directory walk
// itself fails. It carries the directory the walker was asked to scan so callers
// can produce actionable error messages without parsing the error string.
type WalkError struct {
	Dir string
	Err error
}

// Error implements the error interface.
func (e *WalkError) Error() string {
	return fmt.Sprintf("walking %s: %v", e.Dir, e.Err)
}

// Unwrap returns the underlying error so errors.Is / errors.AsType work
// through the chain to the original fs.PathError, syscall.Errno, etc.
func (e *WalkError) Unwrap() error {
	return e.Err
}

// skipDirs are directory basenames that WalkGoDir never descends into.
var skipDirs = map[string]bool{ //nolint:gochecknoglobals // package-level lookup table
	"vendor":       true,
	".git":         true,
	"testdata":     true,
	"node_modules": true,
	".idea":        true,
	"__debug":      true,
}

// WalkGoDir walks dir recursively, parses every non-test .go file, and returns
// them. Files that fail to parse are silently skipped — syntax errors are the
// compiler's job, not the linter's. On walk failure it returns a *WalkError
// wrapping the underlying fs error. Callers should errors.AsType[*WalkError](err)
// to read .Dir.
//
//nolint:erraudit // always returns *WalkError; signature stays error for v0.1.x compat
func WalkGoDir(dir string) ([]ParsedFile, error) { //nolint:erraudit
	var files []ParsedFile

	fset := token.NewFileSet()

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			//nolint:erraudit // err is from filepath.WalkDir; fset/base/parseErr flagged are out-of-scope variables (erraudit false positive)
			return fmt.Errorf("walk %s: %w", path, err)
		}

		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Skip generated files — they are not hand-written reimplementations.
		base := filepath.Base(path)
		if strings.HasSuffix(base, "_gen.go") || strings.HasSuffix(base, ".gen.go") ||
			strings.HasSuffix(base, "_templ.go") {
			return nil
		}

		file, parseErr := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if parseErr != nil {
			return nil //nolint:nilerr // syntax errors are the compiler's job, not the linter's
		}

		files = append(files, ParsedFile{Path: path, Fset: fset, File: file})

		return nil
	})
	if err != nil {
		return nil, &WalkError{Dir: dir, Err: err}
	}

	return files, nil
}

// detectorFunc inspects a single function declaration and returns findings for
// any reimplementation patterns it detects.
type detectorFunc func(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl, filePath string) []finding.Finding

// checkFuncDecls walks dir, parses Go files, and applies detect to every
// function declaration. This is the shared execution path used by every rule.
// On walk failure it returns a *WalkError wrapping the underlying fs error so
// callers can read Dir via errors.AsType[*WalkError].
func checkFuncDecls(dir string, detect detectorFunc) ([]finding.Finding, error) {
	files, err := WalkGoDir(dir)
	if err != nil {
		return nil, &WalkError{Dir: dir, Err: err}
	}

	var all []finding.Finding

	for _, pf := range files {
		for _, decl := range pf.File.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if hasNoLintDirective(pf.Fset, pf.File, fn) {
				continue
			}

			all = append(all, detect(pf.Fset, pf.File, fn, pf.Path)...)
		}
	}

	return all, nil
}

// posOf extracts the line and column from an AST position within fset.
func posOf(fset *token.FileSet, pos token.Pos) (int, int) {
	p := fset.Position(pos)

	return p.Line, p.Column
}
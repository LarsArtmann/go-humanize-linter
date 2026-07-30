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
// compiler's job, not the linter's.
func WalkGoDir(dir string) ([]ParsedFile, error) {
	var files []ParsedFile

	fset := token.NewFileSet()

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
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
		return nil, fmt.Errorf("walking %s: %w", dir, err)
	}

	return files, nil
}

// detectorFunc inspects a single function declaration and returns findings for
// any reimplementation patterns it detects.
type detectorFunc func(fset *token.FileSet, file *ast.File, fn *ast.FuncDecl, filePath string) []finding.Finding

// checkFuncDecls walks dir, parses Go files, and applies detect to every
// function declaration. This is the shared execution path used by every rule.
func checkFuncDecls(dir string, detect detectorFunc) ([]finding.Finding, error) {
	files, err := WalkGoDir(dir)
	if err != nil {
		return nil, err
	}

	var all []finding.Finding

	for _, pf := range files {
		for _, decl := range pf.File.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
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

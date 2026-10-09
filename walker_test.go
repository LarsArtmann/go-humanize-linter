package humanizelint_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	humanizelint "github.com/larsartmann/go-humanize-linter"
)

// TestWalkGoDir_NonexistentDirectory verifies that WalkGoDir returns a typed
// *WalkError when given a path that doesn't exist. The typed error lets callers
// read .Dir without parsing the error string, while Unwrap exposes the
// underlying *fs.PathError for syscall-aware handling.
func TestWalkGoDir_NonexistentDirectory(t *testing.T) {
	t.Parallel()

	nonexistent := filepath.Join(t.TempDir(), "this-path-does-not-exist", "nope")

	files, err := humanizelint.WalkGoDir(nonexistent)
	if err == nil {
		t.Fatalf("expected error for nonexistent dir %q, got files=%v", nonexistent, files)
	}

	if files != nil {
		t.Errorf("expected nil files slice on error, got %v", files)
	}

	walkErr, ok := errors.AsType[*humanizelint.WalkError](err)
	if !ok {
		t.Fatalf("expected error chain to contain *WalkError, got %T: %v", err, err)
	}

	if walkErr.Dir != nonexistent {
		t.Errorf("expected WalkError.Dir=%q, got %q", nonexistent, walkErr.Dir)
	}

	if !strings.Contains(err.Error(), nonexistent) {
		t.Errorf("expected error string to mention the missing dir %q, got: %v", nonexistent, err)
	}

	if _, ok := errors.AsType[*fs.PathError](err); !ok {
		t.Errorf("expected unwrapped error chain to contain *fs.PathError, got %T: %v", err, err)
	}
}

// TestWalkGoDir_SkipsHiddenBackupArchivedForks verifies the walker's skip
// rules: hidden directories (leading "."), backup trees (trailing ".bak"),
// and the archived/forks basenames are never descended into, so frozen or
// upstream-owned copies cannot double-report findings from live code.
func TestWalkGoDir_SkipsHiddenBackupArchivedForks(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	const violating = `package main

func formatBytes(bytes int64) string {
	const unit = 1024
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return "x"
}
`

	skipped := []string{
		".hidden-state/skipped.go",
		"proj.vendor.bak/skipped.go",
		"archived/skipped.go",
		"forks/skipped.go",
	}

	kept := []string{"main.go", filepath.Join("sub", "kept.go")}

	write := func(rel, content string) {
		path := filepath.Join(root, rel)
		if mkdirErr := os.MkdirAll(filepath.Dir(path), 0o750); mkdirErr != nil {
			t.Fatalf("mkdir for %s: %v", rel, mkdirErr)
		}

		if writeErr := os.WriteFile(path, []byte(content), 0o600); writeErr != nil {
			t.Fatalf("write %s: %v", rel, writeErr)
		}
	}

	for _, rel := range skipped {
		write(rel, violating)
	}

	for _, rel := range kept {
		write(rel, violating)
	}

	files, err := humanizelint.WalkGoDir(root)
	if err != nil {
		t.Fatalf("WalkGoDir(%s): %v", root, err)
	}

	got := make([]string, 0, len(files))
	for _, pf := range files {
		rel, relErr := filepath.Rel(root, pf.Path)
		if relErr != nil {
			t.Fatalf("Rel(%s): %v", pf.Path, relErr)
		}

		got = append(got, rel)
	}

	sort.Strings(got)

	want := append([]string{}, kept...)
	sort.Strings(want)

	if !slices.Equal(got, want) {
		t.Errorf("expected walker to return exactly %v, got %v", want, got)
	}
}

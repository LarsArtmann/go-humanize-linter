package humanizelint_test

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	humanizelint "github.com/larsartmann/go-humanize-linter"
)

// TestWalkGoDir_NonexistentDirectory verifies that WalkGoDir returns a wrapped
// error when given a path that doesn't exist. The wrap must include the
// directory so the error chain tells the user where the walk failed.
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

	if !strings.Contains(err.Error(), nonexistent) {
		t.Errorf("expected error to mention the missing dir %q, got: %v", nonexistent, err)
	}

	// WalkGoDir wraps the underlying fs error with %w. WalkDir reports
	// ENOENT as a *fs.PathError on Linux/macOS; we assert the chain
	// contains one so callers can use errors.AsType for syscall-aware
	// handling.
	if _, ok := errors.AsType[*fs.PathError](err); !ok {
		t.Errorf("expected error chain to contain *fs.PathError, got %T: %v", err, err)
	}
}
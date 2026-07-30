package humanizelint_test

import (
	"errors"
	"io/fs"
	"path/filepath"
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

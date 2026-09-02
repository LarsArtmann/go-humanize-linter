package main

import (
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestSinglechecker_CleanCode builds the singlechecker binary and runs it
// against testdata/clean. The singlechecker entry point is a thin
// analysis.SingleMain wrapper around plugin.Analyzer; this test exercises that
// wrapper end-to-end and closes the 0% coverage gap.
func TestSinglechecker_CleanCode(t *testing.T) {
	t.Parallel()

	binary := filepath.Join(t.TempDir(), "gohumanize")

	//nolint:gosec // test build command
	build := exec.CommandContext(context.Background(), "go", "build", "-o", binary, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("failed to build singlechecker: %v\n%s", err, out)
	}

	cleanDir, _ := filepath.Abs(filepath.Join("..", "..", "testdata", "clean"))

	//nolint:gosec // test binary path is trusted
	run := exec.CommandContext(context.Background(), binary, cleanDir)

	out, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 on clean code, got error: %v\n%s", err, out)
	}
}

// TestSinglechecker_PositiveFinding verifies the singlechecker emits an H001
// diagnostic on the KMGTPE testdata. The plugin path uses pass.Report which
// routes through singlechecker.Main's standard exit codes (3 = diagnostics
// reported).
func TestSinglechecker_PositiveFinding(t *testing.T) {
	t.Parallel()

	binary := filepath.Join(t.TempDir(), "gohumanize")

	//nolint:gosec // test build command
	build := exec.CommandContext(context.Background(), "go", "build", "-o", binary, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("failed to build singlechecker: %v\n%s", err, out)
	}

	kmgtpeDir, _ := filepath.Abs(filepath.Join("..", "..", "testdata", "h001_bytes_kmgtptrick"))

	//nolint:gosec // test binary path is trusted
	run := exec.CommandContext(context.Background(), binary, kmgtpeDir)

	out, runErr := run.CombinedOutput()

	// singlechecker uses exit 3 for "diagnostics reported" — extract via
	// errors.As to handle wrapped errors cleanly.

	if runErr == nil {
		t.Fatalf("expected non-zero exit on positive finding, got success\n%s", out)
	}

	exitErr, ok := errors.AsType[*exec.ExitError](runErr)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got %T: %v", runErr, runErr)
	}

	if exitErr.ExitCode() != 3 {
		t.Errorf("expected exit 3 for diagnostics reported, got %d\n%s",
			exitErr.ExitCode(), out)
	}

	if !strings.Contains(string(out), "H001") {
		t.Errorf("expected H001 in singlechecker output, got: %s", out)
	}
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// binaryOnce ensures the CLI is built only once across all tests.
var (
	binaryOnce sync.Once
	binaryPath string
	binaryErr  error
)

// buildCLI builds the CLI binary to a path in the project directory (Nix
// temp dirs have noexec).
func buildCLI(t *testing.T) string {
	t.Helper()

	binaryOnce.Do(func() {
		binaryPath = filepath.Join(os.Getenv("HOME"), ".cache", "go-humanize-linter-test")

		cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-humanize-linter/")
		cmd.Dir = "../.."
		cmd.Env = append(os.Environ(),
			"GOEXPERIMENT=jsonv2",
			"GOPRIVATE=github.com/larsartmann/*",
			"CGO_ENABLED=0",
		)

		var out strings.Builder
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			binaryErr = fmt.Errorf("failed to build CLI: %w\n%s", err, out.String())
			return
		}

		// Nix Go toolchain doesn't set execute permission on build output.
		if err := os.Chmod(binaryPath, 0o755); err != nil {
			binaryErr = fmt.Errorf("failed to chmod binary: %w", err)
		}
	})

	if binaryErr != nil {
		t.Fatal(binaryErr)
	}

	return binaryPath
}

func TestCLI_RunOnCleanCode(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)
	testdata, _ := filepath.Abs(filepath.Join("..", "..", "testdata", "clean"))

	cmd := exec.Command(binary, testdata)
	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 on clean code, got error: %v\n%s", err, output)
	}

	if !strings.Contains(string(output), "0 findings") {
		t.Errorf("expected '0 findings' in output, got: %s", output)
	}
}

func TestCLI_RunOnPositiveFinding(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)
	testdata, _ := filepath.Abs(filepath.Join("..", "..", "testdata", "h001_bytes_kmgtptrick"))

	cmd := exec.Command(binary, "--quiet", testdata)
	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	output, _ := cmd.CombinedOutput()

	if !strings.Contains(string(output), "H001") {
		t.Errorf("expected H001 in output, got: %s", output)
	}
}

func TestCLI_EnableFilter(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)
	testdata, _ := filepath.Abs(filepath.Join("..", "..", "testdata", "h001_bytes_kmgtptrick"))

	cmd := exec.Command(binary, "--quiet", "--enable", "H003", testdata)
	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 when only H003 enabled on H001 testdata: %v\n%s", err, output)
	}
}

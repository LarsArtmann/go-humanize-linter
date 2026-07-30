package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/larsartmann/go-finding"
	humanizelint "github.com/larsartmann/go-humanize-linter"
)

// binaryOnce ensures the CLI is built only once across all tests.
var (
	binaryOnce sync.Once //nolint:gochecknoglobals // test-only build cache
	binaryPath string    //nolint:gochecknoglobals // test-only build cache
	errBinary  error     //nolint:gochecknoglobals // test-only build error
)

// buildCLI builds the CLI binary to a path in the project directory (Nix
// temp dirs have noexec).
func buildCLI(t *testing.T) string {
	t.Helper()

	binaryOnce.Do(func() {
		binaryPath = filepath.Join(os.Getenv("HOME"), ".cache", "go-humanize-linter-test")

		cmd := exec.CommandContext( //nolint:gosec // test build command
			context.Background(),
			"go", "build", "-o", binaryPath, "./cmd/go-humanize-linter/",
		)
		cmd.Dir = "../.."

		cmd.Env = append(
			os.Environ(),
			"GOEXPERIMENT=jsonv2",
			"GOPRIVATE=github.com/larsartmann/*",
			"CGO_ENABLED=0",
		)

		var out strings.Builder

		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			errBinary = fmt.Errorf("failed to build CLI: %w\n%s", err, out.String())

			return
		}

		// Nix Go toolchain doesn't set execute permission on build output.
		if err := os.Chmod(binaryPath, 0o755); err != nil { //nolint:gosec // executable binary needs 0o755
			errBinary = fmt.Errorf("failed to chmod binary: %w", err)
		}
	})

	if errBinary != nil {
		t.Fatal(errBinary)
	}

	return binaryPath
}

func TestCLI_RunOnCleanCode(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)
	testdata, _ := filepath.Abs(filepath.Join("..", "..", "testdata", "clean"))

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, testdata,
	)

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

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--quiet", testdata,
	)

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

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--quiet", "--enable", "H003", testdata,
	)

	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 when only H003 enabled on H001 testdata: %v\n%s", err, output)
	}
}

// ---------------------------------------------------------------------------
// Unit tests for buildRegistry / output / stringList (no subprocess)
// ---------------------------------------------------------------------------

func TestStringList(t *testing.T) {
	t.Parallel()

	var s stringList

	if s.String() != "" {
		t.Errorf("empty stringList.String() = %q, want %q", s.String(), "")
	}

	if err := s.Set("H001"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if err := s.Set("H002"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if got, want := s.String(), "H001,H002"; got != want {
		t.Errorf("stringList.String() = %q, want %q", got, want)
	}
}

func TestBuildRegistry_AllEnabledByDefault(t *testing.T) {
	t.Parallel()

	reg := buildRegistry(nil, nil)
	if got, want := len(reg.All()), len(humanizelint.AllRules()); got != want {
		t.Errorf("default registry has %d rules, want %d", got, want)
	}
}

func TestBuildRegistry_EnableOnly(t *testing.T) {
	t.Parallel()

	reg := buildRegistry([]string{"H001"}, nil)

	rules := reg.All()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}

	if rules[0].ID() != "H001" {
		t.Errorf("expected H001, got %s", rules[0].ID())
	}
}

func TestBuildRegistry_Disable(t *testing.T) {
	t.Parallel()

	reg := buildRegistry(nil, []string{"H001"})
	if got, want := len(reg.All()), len(humanizelint.AllRules())-1; got != want {
		t.Errorf("after disabling H001, registry has %d rules, want %d", got, want)
	}
}

func TestBuildRegistry_UnknownIDIgnored(t *testing.T) {
	t.Parallel()

	// Unknown enable id registers nothing.
	reg := buildRegistry([]string{"H999"}, nil)
	if len(reg.All()) != 0 {
		t.Errorf("expected 0 rules for unknown enable id, got %d", len(reg.All()))
	}
}

// sampleReport runs the full registry against a known-positive fixture so the
// output tests operate on a real *finding.Report.
func sampleReport(t *testing.T) *finding.Report {
	t.Helper()

	reg := buildRegistry(nil, nil)

	dir, err := filepath.Abs(filepath.Join("..", "..", "testdata", "h001_bytes_kmgtptrick"))
	if err != nil {
		t.Fatal(err)
	}

	report, err := reg.Run(context.Background(), dir)
	if err != nil {
		t.Fatalf("registry.Run failed: %v", err)
	}

	if report.Len() == 0 {
		t.Fatal("expected at least one finding in sample report")
	}

	return report
}

func TestOutput_TextContainsFinding(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	output(&buf, sampleReport(t), "text", true)

	out := buf.String()
	if !strings.Contains(out, "H001") {
		t.Errorf("text output missing H001: %q", out)
	}

	if strings.Contains(out, "findings") {
		t.Errorf("quiet=true should suppress the summary line: %q", out)
	}
}

func TestOutput_TextSummaryWhenNotQuiet(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	output(&buf, sampleReport(t), "text", false)

	if !strings.Contains(buf.String(), "findings") {
		t.Errorf("non-quiet text output should contain summary: %q", buf.String())
	}
}

func TestOutput_JSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	output(&buf, sampleReport(t), "json", true)

	out := buf.String()
	if !strings.Contains(out, "H001") {
		t.Errorf("json output missing H001: %q", out)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Errorf("json output is not valid JSON: %v", err)
	}
}

func TestOutput_SARIF(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	output(&buf, sampleReport(t), "sarif", true)

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("sarif output is not valid JSON: %v", err)
	}

	if _, ok := parsed["runs"]; !ok {
		t.Errorf("sarif output missing 'runs' key: %v", parsed)
	}
}

// ---------------------------------------------------------------------------
// CLI flag tests (subprocess)
// ---------------------------------------------------------------------------

func TestCLI_VersionFlag(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--version",
	)

	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 for --version, got error: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "go-humanize-linter") {
		t.Errorf("version output missing program name: %q", out)
	}
}

func TestCLI_RulesFlag(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--rules",
	)

	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 for --rules, got error: %v\n%s", err, out)
	}

	for _, id := range []string{"H001", "H002", "H003", "H004", "H005", "H006", "H007"} {
		if !strings.Contains(string(out), id) {
			t.Errorf("--rules output missing %s: %q", id, out)
		}
	}
}

func TestCLI_SARIFOutput(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)
	testdata, _ := filepath.Abs(filepath.Join("..", "..", "testdata", "h001_bytes_kmgtptrick"))

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--quiet", "--format", "sarif", testdata,
	)

	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	out, _ := cmd.CombinedOutput()

	var sarif map[string]any
	if err := json.Unmarshal(out, &sarif); err != nil {
		t.Fatalf("sarif CLI output is not valid JSON: %v\n%s", err, out)
	}

	if _, ok := sarif["runs"]; !ok {
		t.Errorf("sarif CLI output missing 'runs' key: %s", out)
	}
}

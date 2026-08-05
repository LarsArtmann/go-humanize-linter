package main

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"errors"
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

	err := output(&buf, sampleReport(t), "text", true)
	if err != nil {
		t.Fatalf("text output failed: %v", err)
	}

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

	err := output(&buf, sampleReport(t), "text", false)
	if err != nil {
		t.Fatalf("text output failed: %v", err)
	}

	if !strings.Contains(buf.String(), "findings") {
		t.Errorf("non-quiet text output should contain summary: %q", buf.String())
	}
}

func TestOutput_JSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	err := output(&buf, sampleReport(t), "json", true)
	if err != nil {
		t.Fatalf("json output failed: %v", err)
	}

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

	err := output(&buf, sampleReport(t), "sarif", true)
	if err != nil {
		t.Fatalf("sarif output failed: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("sarif output is not valid JSON: %v", err)
	}

	if _, ok := parsed["runs"]; !ok {
		t.Errorf("sarif output missing 'runs' key: %v", parsed)
	}
}

func TestOutput_SARIFCleanReport(t *testing.T) {
	t.Parallel()

	reg := buildRegistry(nil, nil)

	dir, err := filepath.Abs(filepath.Join("..", "..", "testdata", "clean"))
	if err != nil {
		t.Fatal(err)
	}

	report, err := reg.Run(context.Background(), dir)
	if err != nil {
		t.Fatalf("registry.Run on clean testdata failed: %v", err)
	}

	if report.Len() != 0 {
		t.Fatalf("expected 0 findings on clean testdata, got %d", report.Len())
	}

	var buf bytes.Buffer

	if err := output(&buf, report, "sarif", true); err != nil {
		t.Fatalf("sarif output on clean report failed: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("clean sarif output is not valid JSON: %v\n%s", err, buf.String())
	}

	if _, ok := parsed["runs"]; !ok {
		t.Errorf("clean sarif output missing 'runs' key: %v", parsed)
	}
}

// failingWriter is an io.Writer that always returns an error. Used to exercise
// the failure paths in output().
type failingWriter struct {
	err error
}

func (f *failingWriter) Write(_ []byte) (int, error) { return 0, f.err }

// errWrite is the sentinel error returned by failingWriter.
var errWrite = errors.New("disk full: write failed")

func TestOutput_JSONWriterFailure(t *testing.T) {
	t.Parallel()

	err := output(&failingWriter{err: errWrite}, sampleReport(t), "json", true)
	if err == nil {
		t.Fatal("expected error when writer fails, got nil")
	}

	if !errors.Is(err, errWrite) {
		t.Errorf("expected error to wrap errWrite, got: %v", err)
	}

	outErr, ok := errors.AsType[*OutputError](err)
	if !ok {
		t.Fatalf("expected *OutputError in chain, got %T: %v", err, err)
	}

	if outErr.Format != "json" {
		t.Errorf("expected Format=%q, got %q", "json", outErr.Format)
	}

	if outErr.Stage != "write" {
		t.Errorf("expected Stage=%q, got %q", "write", outErr.Stage)
	}
}

func TestOutput_SARIFWriterFailure(t *testing.T) {
	t.Parallel()

	err := output(&failingWriter{err: errWrite}, sampleReport(t), "sarif", true)
	if err == nil {
		t.Fatal("expected error when writer fails, got nil")
	}

	if !errors.Is(err, errWrite) {
		t.Errorf("expected error to wrap errWrite, got: %v", err)
	}

	outErr, ok := errors.AsType[*OutputError](err)
	if !ok {
		t.Fatalf("expected *OutputError in chain, got %T: %v", err, err)
	}

	if outErr.Format != "sarif" {
		t.Errorf("expected Format=%q, got %q", "sarif", outErr.Format)
	}

	if outErr.Stage != "render" {
		t.Errorf("expected Stage=%q, got %q", "render", outErr.Stage)
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

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("expected exit 0 for --version, got error: %v\nstdout: %s\nstderr: %s",
			err, stdout.String(), stderr.String())
	}

	if !strings.Contains(stdout.String(), "go-humanize-linter") {
		t.Errorf("version output missing program name: %q", stdout.String())
	}

	// A "dev" build must emit a stderr warning so users notice they are
	// running an unversioned binary.
	if !strings.Contains(stderr.String(), "version is unset") {
		t.Errorf("--version with dev build should warn on stderr, got: %q", stderr.String())
	}
}

func TestCLI_RulesFlag(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--rules",
	)

	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("expected exit 0 for --rules, got error: %v\nstdout: %s\nstderr: %s",
			err, stdout.String(), stderr.String())
	}

	// --rules must write to stdout (so it pipes cleanly) and be silent on stderr.
	if stderr.Len() != 0 {
		t.Errorf("--rules wrote to stderr (should be stdout only): %q", stderr.String())
	}

	for _, id := range []string{"H001", "H002", "H003", "H004", "H005", "H006", "H007", "H008", "H009"} {
		if !strings.Contains(stdout.String(), id) {
			t.Errorf("--rules stdout missing %s: %q", id, stdout.String())
		}
	}
}

func TestCLI_ListFiles(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)
	cleanDir, _ := filepath.Abs(filepath.Join("..", "..", "testdata", "clean"))

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--list-files", cleanDir,
	)

	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 for --list-files, got error: %v\n%s", err, out)
	}

	// Clean testdata contains main.go — must appear in the listing.
	if !strings.Contains(string(out), "main.go") {
		t.Errorf("--list-files output missing main.go: %s", out)
	}

	// Listing must be limited to the requested directory; testdata/h001_bytes_kmgtptrick
	// should not appear when listing testdata/clean.
	if strings.Contains(string(out), "h001_bytes_kmgtptrick") {
		t.Errorf("--list-files should not list sibling testdata dirs: %s", out)
	}
}

func TestCLI_ExplainFlag(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--explain", "H001",
	)

	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 for --explain, got error: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "H001") {
		t.Errorf("--explain H001 output missing rule ID: %s", out)
	}

	if !strings.Contains(string(out), "humanize.Bytes") {
		t.Errorf("--explain H001 output missing suggested replacement: %s", out)
	}
}

func TestCLI_ExplainFlagUnknownRule(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--explain", "H999",
	)

	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for unknown rule, got success")
	}

	var exitErr *exec.ExitError

	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *exec.ExitError, got %T: %v", err, err)
	}

	if exitErr.ExitCode() != 2 {
		t.Errorf("expected exit 2 for unknown rule, got %d", exitErr.ExitCode())
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

// ---------------------------------------------------------------------------
// writeReport file-output tests
// ---------------------------------------------------------------------------

func TestWriteReport_TextFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "report.txt")

	writeReport(path, sampleReport(t), "text", true)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !strings.Contains(string(data), "H001") {
		t.Errorf("text file output missing H001: %q", data)
	}

	if strings.Contains(string(data), "findings") {
		t.Errorf("quiet=true should suppress summary in file: %q", data)
	}
}

func TestWriteReport_JSONFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "report.json")

	writeReport(path, sampleReport(t), "json", true)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !strings.Contains(string(data), "H001") {
		t.Errorf("json file output missing H001: %q", data)
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Errorf("json file output is not valid JSON: %v", err)
	}
}

func TestWriteReport_FileHandleClosed(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "report.txt")

	writeReport(path, sampleReport(t), "text", true)

	// If the file handle were still open, Remove would fail on some OSes.
	if err := os.Remove(path); err != nil {
		t.Fatalf("could not remove file (handle not closed?): %v", err)
	}
}

// ---------------------------------------------------------------------------
// CLI --output flag tests (subprocess)
// ---------------------------------------------------------------------------

func TestCLI_OutputToFile(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)
	testdata, _ := filepath.Abs(filepath.Join("..", "..", "testdata", "h001_bytes_kmgtptrick"))
	outFile := filepath.Join(t.TempDir(), "results.txt")

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--quiet", "--output", outFile, testdata,
	)

	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	// Findings exist in this testdata → exit code 1 is expected.
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit (findings present), got success")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *exec.ExitError, got %T: %v", err, err)
	}

	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit 1 (findings), got %d", exitErr.ExitCode())
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}

	if !strings.Contains(string(data), "H001") {
		t.Errorf("output file missing H001 finding: %q", data)
	}
}

// ---------------------------------------------------------------------------
// Config file tests
// ---------------------------------------------------------------------------

func TestLoadConfig_Valid(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.yaml")

	err := os.WriteFile(path, []byte("enable:\n  - H001\n  - H003\ndisable:\n  - H004\n"), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}

	if len(cfg.Enable) != 2 || cfg.Enable[0] != "H001" || cfg.Enable[1] != "H003" {
		t.Errorf("expected Enable=[H001,H003], got %v", cfg.Enable)
	}

	if len(cfg.Disable) != 1 || cfg.Disable[0] != "H004" {
		t.Errorf("expected Disable=[H004], got %v", cfg.Disable)
	}
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.yaml")
	_ = os.WriteFile(path, []byte(""), 0o600)

	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig on empty file failed: %v", err)
	}

	if cfg.Enable != nil || cfg.Disable != nil {
		t.Errorf("empty config should produce nil slices, got Enable=%v Disable=%v", cfg.Enable, cfg.Disable)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	t.Parallel()

	_, err := loadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.yaml")
	_ = os.WriteFile(path, []byte("enable: [H001\ninvalid: yaml: content"), 0o600)

	_, err := loadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestParseConfidenceLevel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		want  finding.Confidence
		err   bool
	}{
		{"low", finding.ConfidenceLow, false},
		{"medium", finding.ConfidenceMedium, false},
		{"high", finding.ConfidenceHigh, false},
		{"full", finding.ConfidenceFull, false},
		{"", finding.ConfidenceNone, true},
		{"invalid", finding.ConfidenceNone, true},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := parseConfidenceLevel(tt.input)
			if (err != nil) != tt.err {
				t.Fatalf("parseConfidenceLevel(%q) error = %v, wantErr %v", tt.input, err, tt.err)
			}

			if got != tt.want {
				t.Errorf("parseConfidenceLevel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestExitCodeFromReport(t *testing.T) {
	t.Parallel()

	t.Run("nil report", func(t *testing.T) {
		t.Parallel()
		if got := exitCodeFromReport(nil); got != 0 {
			t.Errorf("exitCodeFromReport(nil) = %d, want 0", got)
		}
	})

	t.Run("empty report", func(t *testing.T) {
		t.Parallel()
		report := finding.NewReport(finding.ToolInfo{Name: "test"})
		if got := exitCodeFromReport(report); got != 0 {
			t.Errorf("exitCodeFromReport(empty) = %d, want 0", got)
		}
	})

	t.Run("medium only", func(t *testing.T) {
		t.Parallel()
		report := finding.NewReportFromFindings(finding.ToolInfo{Name: "test"}, []finding.Finding{
			makeFinding(finding.ConfidenceMedium),
		})
		if got := exitCodeFromReport(report); got != 2 {
			t.Errorf("exitCodeFromReport(medium) = %d, want 2", got)
		}
	})

	t.Run("high only", func(t *testing.T) {
		t.Parallel()
		report := finding.NewReportFromFindings(finding.ToolInfo{Name: "test"}, []finding.Finding{
			makeFinding(finding.ConfidenceHigh),
		})
		if got := exitCodeFromReport(report); got != 1 {
			t.Errorf("exitCodeFromReport(high) = %d, want 1", got)
		}
	})

	t.Run("mixed medium and high", func(t *testing.T) {
		t.Parallel()
		report := finding.NewReportFromFindings(finding.ToolInfo{Name: "test"}, []finding.Finding{
			makeFinding(finding.ConfidenceMedium),
			makeFinding(finding.ConfidenceHigh),
		})
		if got := exitCodeFromReport(report); got != 1 {
			t.Errorf("exitCodeFromReport(mixed) = %d, want 1", got)
		}
	})
}

func makeFinding(conf finding.Confidence) finding.Finding {
	return finding.NewBuilder(
		finding.RuleName("H001"),
		finding.ToolName("test"),
		"test finding",
		finding.SeverityWarning,
		finding.Pos(finding.FilePath("test.go"), 1, 1),
	).
		WithConfidence(conf).
		MustBuild()
}

func TestFilterReportByConfidence(t *testing.T) {
	t.Parallel()

	report := finding.NewReportFromFindings(finding.ToolInfo{Name: "test"}, []finding.Finding{
		makeFinding(finding.ConfidenceMedium),
		makeFinding(finding.ConfidenceHigh),
	})

	filtered := filterReportByConfidence(report, finding.ConfidenceHigh)
	if filtered.Len() != 1 {
		t.Fatalf("expected 1 high finding, got %d", filtered.Len())
	}
}

func TestCLI_MinConfidence(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)
	dir := t.TempDir()

	err := os.WriteFile(
		filepath.Join(dir, "main.go"),
		[]byte(`package main

func labels() string {
	return "KB and MB"
}

func main() {
	_ = labels()
}
`),
		0o644,
	)
	if err != nil {
		t.Fatal(err)
	}

	// With the default (low) threshold the two unit strings produce a finding.
	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--quiet", dir,
	)
	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit with default threshold, got success: %s", out)
	}

	// With --min-confidence high the medium-confidence finding is filtered out.
	cmd = exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--quiet", "--min-confidence", "high", dir,
	)
	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 with --min-confidence high, got error: %v\n%s", err, out)
	}
}

func TestCLI_VerifySuppressions_UnknownLinterName(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)
	dir := t.TempDir()

	err := os.WriteFile(
		filepath.Join(dir, "main.go"),
		[]byte(`package main

//nolint:go-humanize-linter/H003
func staleSuppression() {}

func main() {
	staleSuppression()
}
`),
		0o644,
	)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--quiet", "--verify-suppressions", dir,
	)
	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit with stale suppression, got success: %s", out)
	}

	if !strings.Contains(string(out), "H0SUP") {
		t.Errorf("expected H0SUP verification finding, got: %s", out)
	}

	if !strings.Contains(string(out), "gohumanize") {
		t.Errorf("expected hint about gohumanize, got: %s", out)
	}
}

func TestCLI_VerifySuppressions_StaleSuppression(t *testing.T) {
	t.Parallel()

	binary := buildCLI(t)
	dir := t.TempDir()

	err := os.WriteFile(
		filepath.Join(dir, "main.go"),
		[]byte(`package main

//nolint:gohumanize:H001
func noByteFormatting() string {
	return "hello"
}

func main() {
	_ = noByteFormatting()
}
`),
		0o644,
	)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.CommandContext( //nolint:gosec // test binary path is trusted
		context.Background(), binary, "--quiet", "--verify-suppressions", dir,
	)
	cmd.Env = append(os.Environ(), "GOEXPERIMENT=jsonv2")

	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit with stale suppression, got success: %s", out)
	}

	if !strings.Contains(string(out), "H0SUP") {
		t.Errorf("expected H0SUP verification finding, got: %s", out)
	}
}

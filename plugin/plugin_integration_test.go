package plugin_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golangci/plugin-module-register/register"
)

// TestPluginRegisteredWithGolangciLint verifies that the init() function in
// plugin.go correctly registers the gohumanize linter with the
// plugin-module-register package. This is the same lookup golangci-lint v2
// performs at runtime when discovering module plugins.
//
// If this test passes, the custom-gcl binary built via `golangci-lint custom`
// will discover the gohumanize linter when the .golangci.yml includes:
//
//	linters:
//	  enable:
//	    - gohumanize
//	  settings:
//	    custom:
//	      gohumanize:
//	        type: "module"
func TestPluginRegisteredWithGolangciLint(t *testing.T) {
	t.Parallel()

	constructor, err := register.GetPlugin("gohumanize")
	if err != nil {
		t.Fatalf("gohumanize not registered with plugin-module-register: %v", err)
	}

	if constructor == nil {
		t.Fatal("constructor is nil")
	}

	plugin, err := constructor(nil)
	if err != nil {
		t.Fatalf("constructor(nil) failed: %v", err)
	}

	if plugin == nil {
		t.Fatal("plugin is nil")
	}

	loadMode := plugin.GetLoadMode()
	if loadMode != register.LoadModeSyntax {
		t.Errorf("expected LoadMode %q, got %q", register.LoadModeSyntax, loadMode)
	}

	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers failed: %v", err)
	}

	if len(analyzers) != 1 {
		t.Fatalf("expected 1 analyzer, got %d", len(analyzers))
	}

	if analyzers[0].Name != "gohumanize" {
		t.Errorf("expected analyzer name 'gohumanize', got %q", analyzers[0].Name)
	}

	if analyzers[0].Run == nil {
		t.Error("analyzer Run function is nil")
	}
}

// TestPluginRegisteredWithSettings verifies that the plugin constructor
// accepts and applies enable/disable settings passed from .golangci.yml.
func TestPluginRegisteredWithSettings(t *testing.T) {
	t.Parallel()

	constructor, err := register.GetPlugin("gohumanize")
	if err != nil {
		t.Fatalf("gohumanize not registered: %v", err)
	}

	plugin, err := constructor(map[string]any{
		"enable": "H001,H003",
	})
	if err != nil {
		t.Fatalf("constructor with settings failed: %v", err)
	}

	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers failed: %v", err)
	}

	if len(analyzers) != 1 {
		t.Fatalf("expected 1 analyzer, got %d", len(analyzers))
	}

	if analyzers[0].Run == nil {
		t.Error("analyzer Run function is nil")
	}
}

// TestCustomGCLIntegration builds the custom-gcl binary via
// `golangci-lint custom`, runs it on a temp module containing a known H001
// pattern, and asserts the gohumanize diagnostic appears.
//
// This test requires network access (golangci-lint custom clones the
// golangci-lint source) and is therefore gated behind testing.Short().
func TestCustomGCLIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping custom-gcl integration test in short mode (requires network + git clone)")
	}

	if _, err := exec.LookPath("golangci-lint"); err != nil {
		t.Skip("golangci-lint not found in PATH")
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}

	projectRoot := filepath.Dir(wd)

	// Build the custom-gcl binary.
	t.Log("building custom-gcl binary (this may take a while)...")

	buildCmd := exec.Command("golangci-lint", "custom")
	buildCmd.Dir = projectRoot

	buildCmd.Env = append(os.Environ(),
		"GOEXPERIMENT=jsonv2",
		"GOPRIVATE=github.com/larsartmann/*,github.com/LarsArtmann/*",
		"GONOSUMDB=github.com/larsartmann/*,github.com/LarsArtmann/*",
	)

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("golangci-lint custom failed: %v\n%s", err, output)
	}

	customGCL := filepath.Join(projectRoot, "custom-gcl")

	if _, err := os.Stat(customGCL); err != nil {
		t.Fatalf("custom-gcl binary not found after build: %v", err)
	}

	t.Cleanup(func() { _ = os.Remove(customGCL) })

	// Create a temp module with a known H001 pattern (KMGTPE index trick).
	tmpDir := t.TempDir()

	if err := os.WriteFile(
		filepath.Join(tmpDir, "go.mod"),
		[]byte("module testexample\n\ngo 1.26\n"),
		0o644,
	); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	h001Source := `package main

import "fmt"

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	println(formatBytes(1073741824))
}
`

	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(h001Source), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	gclConfig := `version: "2"
run:
  timeout: 5m
linters:
  default: none
  enable:
    - gohumanize
  settings:
    custom:
      gohumanize:
        type: "module"
        description: "test"
        original-url: github.com/LarsArtmann/go-humanize-linter
`

	if err := os.WriteFile(filepath.Join(tmpDir, ".golangci.yml"), []byte(gclConfig), 0o644); err != nil {
		t.Fatalf("write .golangci.yml: %v", err)
	}

	// Run custom-gcl on the temp module.
	runCmd := exec.Command(customGCL, "run", "-c", filepath.Join(tmpDir, ".golangci.yml"), "./...")
	runCmd.Dir = tmpDir

	runCmd.Env = append(os.Environ(),
		"GOEXPERIMENT=jsonv2",
		"GOPRIVATE=github.com/larsartmann/*,github.com/LarsArtmann/*",
		"GONOSUMDB=github.com/larsartmann/*,github.com/LarsArtmann/*",
	)

	output, err := runCmd.CombinedOutput()
	// golangci-lint exits 1 when findings are reported — that's expected.
	if err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			t.Fatalf("custom-gcl run failed unexpectedly: %v\n%s", err, output)
		}
	}

	out := string(output)

	if !strings.Contains(out, "H001") {
		t.Errorf("expected H001 diagnostic in custom-gcl output, got:\n%s", out)
	}

	if !strings.Contains(out, "gohumanize") {
		t.Errorf("expected 'gohumanize' linter name in custom-gcl output, got:\n%s", out)
	}
}

package plugin_test

import (
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
// `golangci-lint custom`, runs it on temp modules with various plugin
// settings, and asserts the gohumanize diagnostics appear correctly.
//
// Subtests:
//   - basic_detection: default settings, H001 fires
//   - min_confidence_full: minConfidence "full", H001 (ConfidenceFull) still fires
//   - verify_suppressions: stale //nolint directive produces H0SUP
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

	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}

	projectRoot := filepath.Dir(wd)

	t.Log("building custom-gcl binary (this may take a while)...")

	buildCmd := exec.CommandContext(t.Context(), "golangci-lint", "custom")
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

	t.Run("basic_detection", func(t *testing.T) {
		t.Parallel()

		out := runCustomGCL(t, customGCL, h001Source, "")

		if !strings.Contains(out, "H001") {
			t.Errorf("expected H001 diagnostic in output, got:\n%s", out)
		}

		if !strings.Contains(out, "gohumanize") {
			t.Errorf("expected 'gohumanize' linter name in output, got:\n%s", out)
		}
	})

	t.Run("min_confidence_full", func(t *testing.T) {
		t.Parallel()

		settings := `        settings:
          minConfidence: "full"
`

		out := runCustomGCL(t, customGCL, h001Source, settings)

		if !strings.Contains(out, "H001") {
			t.Errorf("expected H001 (ConfidenceFull) to survive minConfidence=full, got:\n%s", out)
		}
	})

	t.Run("verify_suppressions", func(t *testing.T) {
		t.Parallel()

		staleSource := `package main

//nolint:gohumanize // stale: no humanize pattern here
func cleanFunc() string {
	return "hello"
}

func main() {
	_ = cleanFunc()
}
`
		settings := `        settings:
          verifySuppressions: true
`

		out := runCustomGCL(t, customGCL, staleSource, settings)

		if !strings.Contains(out, "H0SUP") {
			t.Errorf("expected H0SUP diagnostic for stale directive, got:\n%s", out)
		}
	})
}

func runCustomGCL(t *testing.T, customGCL, source, settingsYAML string) string {
	t.Helper()

	tmpDir := t.TempDir()

	if err := os.WriteFile(
		filepath.Join(tmpDir, "go.mod"),
		[]byte("module testexample\n\ngo 1.26\n"),
		0o600,
	); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(source), 0o600); err != nil {
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
` + settingsYAML

	if err := os.WriteFile(filepath.Join(tmpDir, ".golangci.yml"), []byte(gclConfig), 0o600); err != nil {
		t.Fatalf("write .golangci.yml: %v", err)
	}

	runCmd := exec.CommandContext(
		t.Context(),
		customGCL,
		"run",
		"-c",
		filepath.Join(tmpDir, ".golangci.yml"),
		"./...",
	) //nolint:gosec // customGCL is built by this test into the project dir
	runCmd.Dir = tmpDir

	runCmd.Env = append(os.Environ(),
		"GOEXPERIMENT=jsonv2",
		"GOPRIVATE=github.com/larsartmann/*,github.com/LarsArtmann/*",
		"GONOSUMDB=github.com/larsartmann/*,github.com/LarsArtmann/*",
	)

	output, _ := runCmd.CombinedOutput()

	return string(output)
}

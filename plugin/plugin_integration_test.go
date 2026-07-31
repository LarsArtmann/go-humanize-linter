package plugin_test

import (
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

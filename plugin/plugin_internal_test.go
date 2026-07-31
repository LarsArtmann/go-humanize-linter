package plugin

import (
	"testing"

	humanizelint "github.com/larsartmann/go-humanize-linter"
)

func TestParseRuleIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty", "", []string{}},
		{"single", "H001", []string{"H001"}},
		{"multiple", "H001,H003", []string{"H001", "H003"}},
		{"whitespace", " H001 , H003 ", []string{"H001", "H003"}},
		{"trailing comma", "H001,", []string{"H001"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := parseRuleIDs(tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf("parseRuleIDs(%q) = %d entries, want %d", tt.input, len(got), len(tt.want))
			}

			for _, id := range tt.want {
				if !got[id] {
					t.Errorf("parseRuleIDs(%q) missing %s", tt.input, id)
				}
			}
		})
	}
}

func TestFilterRules_NoConfigReturnsAll(t *testing.T) {
	t.Parallel()

	all := humanizelint.AllRules()
	got := filterRules(all, nil, nil)

	if len(got) != len(all) {
		t.Errorf("no config: got %d rules, want %d", len(got), len(all))
	}
}

func TestFilterRules_EnableOnly(t *testing.T) {
	t.Parallel()

	all := humanizelint.AllRules()
	enable := map[string]bool{humanizelint.RuleIDH001: true, humanizelint.RuleIDH003: true}
	got := filterRules(all, enable, nil)

	if len(got) != 2 {
		t.Fatalf("enable=H001,H003: got %d rules, want 2", len(got))
	}

	for _, rule := range got {
		if !enable[rule.Meta.ID] {
			t.Errorf("unexpected rule %s in filtered set", rule.Meta.ID)
		}
	}
}

func TestFilterRules_DisableOnly(t *testing.T) {
	t.Parallel()

	all := humanizelint.AllRules()
	disable := map[string]bool{humanizelint.RuleIDH001: true}
	got := filterRules(all, nil, disable)

	if len(got) != len(all)-1 {
		t.Fatalf("disable=H001: got %d rules, want %d", len(got), len(all)-1)
	}

	for _, rule := range got {
		if disable[rule.Meta.ID] {
			t.Errorf("disabled rule %s should not be in filtered set", rule.Meta.ID)
		}
	}
}

func TestFilterRules_EnableOverridesDisable(t *testing.T) {
	t.Parallel()

	all := humanizelint.AllRules()
	enable := map[string]bool{humanizelint.RuleIDH001: true, humanizelint.RuleIDH002: true}
	disable := map[string]bool{humanizelint.RuleIDH001: true}
	got := filterRules(all, enable, disable)

	if len(got) != 1 {
		t.Fatalf("enable=H001,H002 disable=H001: got %d rules, want 1 (only H002)", len(got))
	}

	if got[0].Meta.ID != humanizelint.RuleIDH002 {
		t.Errorf("expected H002, got %s", got[0].Meta.ID)
	}
}

func TestNewPluginWithSettings(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(map[string]any{
		"enable":  "H001,H002",
		"disable": "H003",
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	impl, ok := plug.(*humanizePlugin)
	if !ok {
		t.Fatalf("expected *humanizePlugin, got %T", plug)
	}

	if impl.settings.Enable != "H001,H002" {
		t.Errorf("expected Enable=%q, got %q", "H001,H002", impl.settings.Enable)
	}

	if impl.settings.Disable != "H003" {
		t.Errorf("expected Disable=%q, got %q", "H003", impl.settings.Disable)
	}
}

func TestNewPluginWithEmptySettings(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(nil)
	if err != nil {
		t.Fatalf("newPlugin(nil) failed: %v", err)
	}

	impl, ok := plug.(*humanizePlugin)
	if !ok {
		t.Fatalf("expected *humanizePlugin, got %T", plug)
	}

	if impl.settings.Enable != "" || impl.settings.Disable != "" {
		t.Errorf("empty settings should produce empty plugin config, got Enable=%q Disable=%q",
			impl.settings.Enable, impl.settings.Disable)
	}
}

func TestBuildAnalyzersProducesConfiguredDetector(t *testing.T) {
	t.Parallel()

	plug, err := newPlugin(map[string]any{
		"enable": humanizelint.RuleIDH001,
	})
	if err != nil {
		t.Fatalf("newPlugin failed: %v", err)
	}

	analyzers, err := plug.BuildAnalyzers()
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
		t.Fatal("analyzer.Run must not be nil")
	}
}

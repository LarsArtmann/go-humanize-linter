package humanizelint

import "github.com/larsartmann/go-linter-sdk"

// DefaultRegistry returns a Registry pre-loaded with all humanize-lint rules,
// all enabled by default.
func DefaultRegistry() *linter.Registry {
	r := linter.NewRegistry()

	for _, rule := range AllRules() {
		r.Register(rule)
	}

	return r
}

// AllRules returns every rule in this linter as a slice. Useful for consumers
// that want to cherry-pick rules into their own registry.
func AllRules() []linter.RuleFunc {
	return []linter.RuleFunc{
		RuleBytes(),
		RuleComma(),
		RuleRelTime(),
		RulePlural(),
		RuleSI(),
		RuleFtoa(),
		RuleParseBytes(),
	}
}

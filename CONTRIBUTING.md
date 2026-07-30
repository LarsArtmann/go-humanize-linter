# Contributing

Thanks for your interest in contributing to go-humanize-linter!

## How to Contribute

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## Development Setup

### Prerequisites

- **Go 1.26+**
- **Nix** (recommended) — all build commands set the required env vars automatically

### Required Environment Variables

All `go` commands require these env vars (Nix sets them automatically):

```bash
export GOEXPERIMENT=jsonv2
export GOPRIVATE='github.com/larsartmann/*'
export GONOSUMDB='github.com/larsartmann/*'
```

### Build Commands (via Nix)

```bash
nix run .#test          # go test ./... -count=1
nix run .#test-race     # go test ./... -race -count=1
nix run .#build         # go build ./...
nix run .#lint          # golangci-lint run ./...
nix run .#vet           # go vet ./...
nix run .#coverage      # go test with coverage report
```

### Build Commands (direct Go)

```bash
export GOEXPERIMENT=jsonv2 GOPRIVATE='github.com/larsartmann/*' GONOSUMDB='github.com/larsartmann/*'
go build ./...
go test ./... -race -count=1
go vet ./...
```

## Testing Strategy

- **Unit tests** (`linter_test.go`) — run each rule against `testdata/` fixtures
- **White-box tests** (`pattern_helpers_test.go`, `pattern_commaf_test.go`) — test AST helpers directly
- **Analysistest** (`plugin/plugin_test.go`) — end-to-end plugin path via `analysistest.Run`
- **Anti-regression** (`TestRuleCountConsistency`) — guards against ghost rules (registered but never fire)
- **Self-scan** (`TestLintsItself_Clean`) — the linter must produce 0 findings on its own source

### Testdata Layout

```
testdata/
  h008_ordinal/main.go        # positive fixture (should trigger H008)
  h008_negative/main.go       # negative fixture (must NOT trigger H008)
  analysistest/
    h008positive/main.go      # analysistest fixture with // want "H008"
    clean/main.go             # must produce no diagnostics
```

## Adding a New Rule (H0NN)

1. **`pattern_h0NN.go`** — AST pattern helpers (`hasPattern`, etc.). Each helper takes `*ast.FuncDecl` and returns `bool`.

2. **`rule_h0NN.go`** — Rule factory (`RuleName()`) returning `linter.RuleFunc`, and detector function (`detectName`). The detector calls the pattern helpers and returns `[]finding.Finding`.

3. **`rules.go`** — Register the rule in **two** places:
   - `AllRules()` — add `RuleName()` to the slice
   - `allRuleDetectors()` — add `{"H0NN", detectName}` to the slice
   
   `TestRuleCountConsistency` will fail if you forget either one.

4. **`testdata/h0NN_name/main.go`** — Positive fixture that triggers the rule.

5. **`testdata/h0NN_negative/main.go`** — Negative fixture with similar-looking code that must NOT trigger.

6. **`testdata/analysistest/h0NNpositive/main.go`** — Analysistest fixture with `// want "H0NN"` on the function line.

7. **`plugin/plugin_test.go`** — Add `"./h0NNpositive"` to the `analysistest.Run` call in `TestAnalyzerAnalysistest`.

8. **`linter_test.go`** — Add `TestRuleName_Positive` and `TestRuleName_Negative`.

9. **`cmd/go-humanize-linter/main.go`** — Add to `ruleExplanations` map and the local `h0NN` constant.

10. **`docs/rules/H0NN.md`** — Rule documentation (follow `H001.md` format).

11. **`README.md`** — Add to the rules table.

## Adding a CLI Flag

1. **`cmd/go-humanize-linter/main.go`** — Add the flag definition, wire it to the processing logic, and update `flag.Usage`.
2. **`cmd/go-humanize-linter/main_test.go`** — Add test coverage.
3. **`README.md`** — Document the flag in the usage section.

## Reporting Issues

Please use GitHub Issues to report bugs or request features.

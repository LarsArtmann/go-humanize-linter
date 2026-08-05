# Roadmap

> Long-term direction and raw ideas not yet refined into actionable tasks.
>
> Bounded, estimable work lives in `TODO_LIST.md`. This file is vision: themes and
> unrefined concepts. Items graduate to `TODO_LIST.md` when they become scoped.

## Themes

### Precision & ergonomics

- **Per-line diagnostics** — report at the actual matched pattern, not just the
  func-decl position. Today every detector returns a finding at `fn.Pos()`; per-line
  reporting requires each detector to return a specific `token.Pos`.
- **Type-aware detection** — optional `go/types` / `pass.TypesInfo` integration to
  resolve import aliases (`s "strings"`) and typed values. Would cut false negatives
  on aliased and generic code.
- **Per-statement suppression** — today `//nolint` directives are matched at the
  function level; per-statement or per-line suppression would let developers
  silence individual findings within a function (see TODO T19).
- **Rule-overlap disambiguation** — H009 (manual-commaf) and H002 (manual-comma)
  both fire on the same comma-loop pattern, producing double diagnostics for one
  mistake. H009 should require a signal H002 cannot match (e.g. `strconv.FormatFloat`
  or a dot-split on the formatted float) so the two rules partition the space.

### Detection breadth

- **H010 and beyond** — more `go-humanize` coverage candidates:
  - `humanize.LookupMenuItem` (Kubernetes-style suffix lookup)
  - `time.Round` detection for H003
  - `fmt.Sprintf("%.1f", x)` + `strings.TrimRight` combined detection for H006
  - Inline `.String() + " ago"` patterns for H003
- **Package-scope scanning** — today only `FuncDecl` bodies are scanned;
  package-level `var` initializers are invisible.

### Ecosystem & distribution

- **golangci-lint plugin index** — publish once a tagged version is `go install`-able
  and `go.mod` `replace` directives are removed (depends on `go-linter-sdk` first tag).
- **Stable rule IDs** — freeze H001–H0xx (no renames) after v1.0.
- **Benchmark against large codebases** — k8s, cockroach; track scan-speed regressions
  across releases with `benchstat`.

### Integration

- **Auto-fix** — rewrite detected code in-place to use `humanize.X` via the
  `go-finding` pipeline `FixEngine`.
- **Editor integration** — LSP server mode (`humanize-lint server`) and/or a VSCode
  extension via `gopls` analyzer.
- **Better SARIF** — include rule descriptions, help URIs, and fix suggestions.

## Non-goals

Deliberately NOT pursued:

- **Multi-language support** — this is a Go linter; detecting `humanize`-style
  reimplementations in Python/Rust/JS is out of scope.
- **ML-based scoring** — the multi-signal heuristic is hand-tuned and auditable;
  replacing it with a trained model trades interpretability for marginal accuracy.
- **Hosted/SaaS version** — PR commenting, team dashboards, and custom rule packs
  are not on the path; the linter is a local CLI/library/plugin.

# Roadmap

> Long-term direction and raw ideas not yet refined into actionable tasks.
>
> Bounded, estimable work lives in `TODO_LIST.md`. This file is vision: themes and
> unrefined concepts. Items graduate to `TODO_LIST.md` when they become scoped.

## Themes

### Precision & ergonomics

- **Per-line diagnostics** — report at the actual matched pattern, not just the
  func-decl position. Today every detector returns a finding at `fn.Pos()`; per-line
  reporting requires each detector to return a specific `token.Pos`. The plugin path
  already maps findings to positions via `findingToTokenPos`; the detectors need to
  produce more precise positions.
- **Type-aware detection** — optional `go/types` / `pass.TypesInfo` integration to
  resolve import aliases (`s "strings"`) and typed values. Would cut false negatives
  on aliased and generic code. ADR 0001 documents the current syntactic-only approach.
- **Per-statement suppression** — today `//nolint` directives are matched at the
  function level; per-statement or per-line suppression would let developers
  silence individual findings within a function (see TODO T19).
- ~~**Rule-overlap disambiguation**~~ — DONE (ADR 0005): H009 suppresses H002
  on shared matches; H010/H002 are inverse operations and cannot co-fire.

### Detection breadth

- **H011+ — more go-humanize coverage candidates**:
  - `humanize.ParseSI` hand-rolls (suffix parsing → value + unit; upstream
    v1.1.0 even added a µ/mu alias — parsing rules are timely)
  - `humanize.WordSeries` / `OxfordWordSeries` ("a, b, and c" joining)
  - `humanize.LookupMenuItem` (Kubernetes-style suffix lookup)
  - `time.Round` detection for H003
  - `fmt.Sprintf("%.1f", x)` + `strings.TrimRight` combined detection for H006
  - Inline `.String() + " ago"` patterns for H003
- **Project-level consistency check** — flag mixed SI/IEC byte-formatting
  conventions within a module (e.g., `humanize.Bytes` in one package, custom
  IEC formatting in another).

### Ecosystem & distribution

- **golangci-lint plugin index** — v0.3.0 is tagged and `go get`-verified
  (clean-module check passed); submission is TODO T18.
- **Stable rule IDs** — freeze H001–H0xx (no renames) after v1.0.
- **Benchmark against large codebases** — k8s, cockroach; track scan-speed regressions
  across releases with `benchstat`.

### Integration

- **Auto-fix** — rewrite detected code in-place to use `humanize.X` via the
  `go-finding` pipeline `FixEngine`.
- **Editor integration** — LSP server mode (`humanize-lint server`) and/or a VSCode
  extension via `gopls` analyzer.
- **Better SARIF** — include rule descriptions, help URIs, and fix suggestions.
- **`--behavior-delta`** — warn when replacing hand-rolled code with `humanize.X`
  would change visible output (SI vs IEC byte formatting, rounding differences).

## Non-goals

Deliberately NOT pursued:

- **Multi-language support** — this is a Go linter; detecting `humanize`-style
  reimplementations in Python/Rust/JS is out of scope.
- **ML-based scoring** — the multi-signal heuristic is hand-tuned and auditable;
  replacing it with a trained model trades interpretability for marginal accuracy.
- **Hosted/SaaS version** — PR commenting, team dashboards, and custom rule packs
  are not on the path; the linter is a local CLI/library/plugin.

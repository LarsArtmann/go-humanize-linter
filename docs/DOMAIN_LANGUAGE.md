# Domain Language

Ubiquitous language for the go-humanize-linter project.

## Core Concepts

### Humanize Reimplementation

A hand-written function that duplicates functionality already provided by [`github.com/dustin/go-humanize`](https://github.com/dustin/go-humanize). The linter detects these and suggests the equivalent `go-humanize` API call.

### Rule

A named detection pattern that identifies a specific class of reimplementation. Each rule has:

- **ID**: A stable identifier (H001-H009) used for suppression matching, filter configuration, and reporting.
- **Detector function**: Walks the AST of a function declaration and returns findings when the pattern matches.
- **Confidence level**: How certain the detector is that this is a true positive (not a false positive).

### Corroborating Signals

Each rule requires **multiple independent signals** in the same function before firing. No single signal is enough. For example, H001 (byte formatting) requires byte-unit strings AND division by 1024, or the KMGTPE index trick. This multi-signal approach keeps the false-positive rate near zero.

### Finding

A diagnostic produced by a rule detector. Contains:

- **Rule ID**: Which rule fired (H001-H009, or H0SUP for suppression verification).
- **Message**: Human-readable description of the detected pattern.
- **Suggestion**: The recommended `go-humanize` API replacement.
- **Confidence**: How certain the detector is (None, Low, Medium, High, Full).
- **Position**: File path, line, and column where the finding was detected.

### Confidence Level

A finding's certainty score, from the `go-finding` library:

| Level | Value | Meaning |
|-------|-------|---------|
| `None` | 0.0 | No confidence (unused) |
| `Low` | 0.25 | Minimal signal — show by default |
| `Medium` | 0.5 | Moderate signal — worth triage |
| `High` | 0.75 | Strong signal — likely a true positive |
| `Full` | 1.0 | Near-certain — the canonical anti-pattern |

The CLI's `--min-confidence` flag and the plugin's `minConfidence` setting filter findings below the threshold. Exit codes reflect confidence: 0 = clean, 1 = high/full finding, 2 = only medium/low.

### Suppression Directive

A `//nolint:gohumanize` or `//nolint:gohumanize:Hxxx` comment that silences findings for a specific function. The linter namespace is `gohumanize` (the analyzer name), NOT `go-humanize-linter` (the module path).

### Suppression Verification

A post-detection pass (`--verify-suppressions`) that checks every `//nolint` directive in the scanned code for two problems:

1. **Unknown linter name**: Directives containing "humanize" but not exactly "gohumanize" (e.g., `//nolint:go-humanize-linter`). These silently do nothing.
2. **Stale suppression**: `//nolint:gohumanize` directives that didn't suppress any finding in the current run. These accumulate when code is refactored.

Suppression-verification findings use pseudo-rule ID **H0SUP**.

### Pseudo-Rule

A rule ID that is NOT registered in `AllRules()` or `DefaultRegistry()`. It cannot be enabled or disabled via `--enable`/`--disable`. Currently, `H0SUP` is the only pseudo-rule. It bypasses confidence filtering in the plugin path so stale-directive diagnostics are always surfaced.

### Ghost Rule

A rule that is registered in the registry but never fires on any test fixture. The `TestRuleCountConsistency` anti-regression test guards against ghost rules by verifying that every registered rule has at least one positive test fixture that triggers it.

## Rule IDs

| ID | Name | Detects | Suggests |
|----|------|---------|----------|
| H001 | manual-bytes-format | Byte-size formatting (KMGTPE trick, unit strings, 1024 division) | `humanize.Bytes`, `humanize.IBytes` |
| H002 | manual-comma-format | Thousands-separator insertion (mod-3 grouping, digit conversion) | `humanize.Comma` |
| H003 | manual-reltime | Relative-time formatting (time.Since/Sub + "ago" + threshold) | `humanize.RelTime`, `humanize.Time` |
| H004 | manual-plural | English pluralization (`if n == 1` switches, singular/plural params) | `english.PluralWord`, `english.Plural` |
| H005 | manual-si-format | SI-prefix formatting (division by 1000, K/M/G suffix) | `humanize.SI` |
| H006 | manual-ftoa | Trailing-zero stripping (nested `strings.TrimRight`) | `humanize.Ftoa` |
| H007 | manual-parse-bytes | Byte-size string parsing (HasSuffix/CutSuffix, multiplier maps) | `humanize.ParseBytes` |
| H008 | manual-ordinal | Ordinal formatting (`switch n%10`, st/nd/rd/th cases) | `humanize.Ordinal` |
| H009 | manual-commaf | Float-with-comma formatting (`%.Nf` + manual separator grouping) | `humanize.Commaf` |
| H0SUP | suppression-verification | Stale or misspelled `//nolint:gohumanize` directives | Fix or remove the directive |

## Detection Architecture

### Two Execution Paths

| Path | Entry point | How it works |
|------|-------------|--------------|
| **CLI** | `cmd/go-humanize-linter` | `WalkGoDir` walks directories, `checkFuncDecls` finds functions, detectors run per-function |
| **Plugin** | `plugin/plugin.go` | `analysis.Pass` provides pre-parsed files, `runDetector` iterates declarations, findings reported as `analysis.Diagnostic` |

### Import-Alias Resolution

`buildImportAliases(file)` maps import aliases to canonical package paths (e.g., `str → strings`). Supports named aliases (`str "strings"`) and dot imports (`. "strings"` → `aliases["."] = "strings"`). Threaded through pattern helpers via variadic `aliases ...map[string]string`.

### Generated-File Detection

`IsGeneratedFile` is the single source of truth for skipping generated files. Uses `gogenfilter` two-phase detection (filename-only + content) with a legacy suffix fallback (`_gen.go`, `.gen.go`, `_templ.go`) for generators gogenfilter doesn't enumerate.

## Configuration Terms

| Term | Meaning |
|------|---------|
| `enable` | Comma-separated rule IDs to run (default: all) |
| `disable` | Comma-separated rule IDs to skip |
| `minConfidence` | Confidence threshold (`low`, `medium`, `high`, `full`) |
| `verifySuppressions` | Boolean: run stale-directive verification pass |
| `gohumanize` | The analyzer name registered with golangci-lint (NOT `go-humanize-linter`) |

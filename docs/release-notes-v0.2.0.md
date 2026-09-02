# v0.2.0 Release Notes

_Published 2026-09-02. See the full [CHANGELOG](../CHANGELOG.md) for details._

## go-humanize-linter v0.2.0

The detection-completeness and production-readiness release. All 9 rules
are now registered, the CLI has confidence-aware exit codes and behavior
delta tracking, and the golangci-lint v2 module plugin supports configurable
rules, confidence filtering, and suppression verification.

### Highlights

**Two new rules:**

- **H008** (manual-ordinal): Detects `switch n%10` with st/nd/rd/th cases. Suggests `humanize.Ordinal`.
- **H009** (manual-commaf): Detects `%.Nf` + manual comma/separator grouping. Suggests `humanize.Commaf`.

**CLI improvements:**

- `--min-confidence` — Filter findings by confidence (low/medium/high/full). Exit 1 for high/full, exit 2 for medium/low only.
- `--verify-suppressions` — Report stale and misspelled `//nolint` directives.
- `--save-baseline` / `--behavior-delta` — Track finding changes across refactors by (rule, file, line) tuple.
- `--config` — Load rule enable/disable settings from YAML.
- `--output` — Write reports to a file.
- `--explain` / `--list-files` — Inspection and debugging flags.

**Plugin improvements:**

- Configurable rules via `.golangci.yml` (`enable`/`disable` settings).
- `minConfidence` and `verifySuppressions` settings.
- golangci-lint v2 module plugin registration via `plugin-module-register`.

**Suppression system:**

- Scoped `//nolint:gohumanize:H001` directives — suppress individual rules.
- In-body `//nolint` — place the directive on the specific statement, not just the function declaration.
- Go-style `//lint:ignore gohumanize` alternative syntax.
- Consistent behavior across CLI, plugin, and verification paths.

**Detection improvements:**

- H001 size-bucket false-positive filter (switch + slice lookups without div1024 excluded).
- H004 false-positive filter (requires string return type + string in branch).
- H009/H002 overlap disambiguation (H002 suppressed when H009 fires).
- Import-alias and dot-import aware detection across all rules.
- Package-level `var` detection for H007.
- gogenfilter-driven generated-file detection (sqlc, templ, protobuf, wire, moq, etc.).

**Quality:**

- 9 rules, 6 ADRs, 83.9% overall test coverage (core 90.9%, CLI 54.9%, plugin 94.5%).
- CI with test+vet, lint, govulncheck, coverage reporting, and self-scan.
- GitHub Action with all flags exposed.

### Fixed Late in the Cycle

- **Plugin-path H0SUP diagnostics were silently dropped by golangci-lint's nolint filter** — stale-directive findings were anchored at the `//nolint:gohumanize` directive itself, which golangci suppresses (self-referential). The plugin path now re-anchors them to the file position with the directive coordinates in the message. The CLI path keeps precise positions.
- **`go-linter-sdk` bumped to v0.2.0** — provides `RuleMeta.ToolName`, unblocking the custom-gcl plugin build and its end-to-end test.
- **Local absolute-path `replace` directives removed from `go.mod`** — they only worked on the dev machine; local development now resolves siblings exclusively via `go.work`.

### Breaking Changes

None. All v0.1.0 APIs, flags, and directives continue to work unchanged.

### Installation

```bash
go install github.com/larsartmann/go-humanize-linter/cmd/go-humanize-linter@v0.2.0
```

### golangci-lint v2 Module Plugin

```bash
golangci-lint custom  # builds custom-gcl with gohumanize compiled in
./custom-gcl run ./...
```

See `.golangci.custom.yml` for the required config section.

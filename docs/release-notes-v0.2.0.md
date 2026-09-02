# v0.2.0 Release Notes Draft

_Draft for GitHub release. Tag pending user approval._

---

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

- 9 rules, 9 ADRs, 83.9% overall test coverage.
- CI with test+vet, lint, govulncheck, coverage reporting, and self-scan.
- GitHub Action with all flags exposed.

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

# Roadmap

> Long-term direction and raw ideas not yet refined into actionable tasks.

## v0.2 — Suppression and Precision

- `//nolint:gohumanize` directive support (inline suppression)
- Per-line diagnostics (report at the actual pattern, not just func-decl)
- Self-exclusion (don't flag own source)
- Configurable confidence threshold (`--threshold=high`)

## v0.3 — Type-aware Detection

- Optional `go/types` integration for H001 (type-based byte detection)
- Detect `fmt.Sprintf` with `%d` + " byte(s)" pattern using type info
- Reduce false negatives on generic code

## v0.4 — More Rules

- H008: humanize.Ordinal (`1st`, `2nd`, `3rd`)
- H009: humanize.Commaf (float comma formatting)
- H010: humanize.LookupMenuItem (Kubernetes-style suffix lookup)
- Detect `time.Round` usage for H003

## v1.0 — Ecosystem

- Publish to golangci-lint plugin index
- Stable rule IDs (no more changes after v1.0)
- Full analysistest suite for plugin path
- Benchmark against large codebases (k8s, cockroach)
- Auto-fix support via go-finding pipeline FixEngine

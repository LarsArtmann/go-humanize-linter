# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                     | Tier   | Effort | Status  |
| --- | -------------------------------------------------------- | ------ | ------ | ------- |
| T21 | Push H010 session work + green CI on main                | High   | S      | ready   |
| T22 | H010 real-world corpus sweep + validation doc            | High   | M      | ready   |
| T23 | Release v0.3.0 (H010) after sweep                        | High   | S      | blocked by T22 |
| T24 | H010 fixture pack (alias/dot/scoped-nolint/rune-loop/NewReplacer) | Medium | S | ready |
| T25 | Investigate why stale gogenfilter test never failed CI   | High   | S      | ready   |
| T26 | Fix dependabot actions-group PR (red since 2026-09-17)   | Medium | S      | ready   |
| T27 | Align doc.go/README examples to ExitCodeByConfidence     | Medium | S      | ready   |
| T28 | Suggestion enrichments (BytesN/BigComma/CommafWithDigits) | Medium | S     | ready   |
| T29 | H011 rule: manual SI-string parsing (ParseSI)            | Medium | M      | planned |
| T30 | H012 rule: WordSeries/Oxford hand-rolls                  | Low    | M      | planned |
| T31 | Delete DEPLOY_KEY_* secrets + dep-repo keys (post-green-CI) | Low  | S      | blocked by T21 |
| T18 | Publish to golangci-lint plugin index                    | High   | S      | blocked by T23 |
| T20 | Release gogenfilter (sibling ahead of pinned v3.6.0)     | Medium | S      | planned |

Full breakdown with 27 medium tasks and 70 fine-grained steps:
`docs/planning/2026-09-18_21-09_SUPERB-pareto-ship-h010-and-beyond.md`.

---

## Distribution

### T18 — Publish to golangci-lint plugin index · High · _blocked by T23 (v0.3.0 release)_

v0.2.0 is tagged, the GitHub Release is published (with linux-amd64 binaries),
and `go get github.com/larsartmann/go-humanize-linter@v0.2.0` was verified in a
clean module (the `plugin` package compiles from the published version).

- [ ] Submit to the golangci-lint plugin index

---

## Upstream

### T20 — Release gogenfilter with the new sqlc detection · Medium · _planned_

The workspace sibling `../gogenfilter` is ahead of the pinned v3.6.0
(2026-09-18: local sibling carries further commits; v3.6.0 already shipped the
config-aware sqlc change this repo adopted). The `go.work` `use` entry serves
the sibling locally, so local dev can run newer detection behavior than CI —
a latent split-brain.

- [ ] Tag + push a gogenfilter release (requires explicit user approval)
- [ ] `go get github.com/LarsArtmann/gogenfilter/v3@<new>` + `go mod tidy`
- [ ] Re-run the testdata sweep to confirm no detection behavior changed

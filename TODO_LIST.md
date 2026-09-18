# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                          | Tier   | Effort | Status                                              |
| --- | ------------------------------------------------------------- | ------ | ------ | --------------------------------------------------- |
| T18 | Publish to golangci-lint plugin index                         | High   | S      | ready (v0.3.0 released)                             |
| T20 | Release gogenfilter (sibling ahead of pinned v3.6.0)          | Medium | S      | planned, needs approval                             |
| T29 | H011 rule: manual SI-string parsing (ParseSI)                 | Medium | M      | deferred — zero corpus demand (2026-09-18 research) |
| T32 | Remove read-only deploy keys from 4 dep repos (foreign repos) | Low    | S      | ready                                               |

Done 2026-09-18 (recorded in CHANGELOG): T21 push+CI, T22 H010 sweep,
T23 v0.3.0 release, T24 H010 fixture pack, T25 CI-miss investigation
(main was red since 09-13, unnoticed), T26 dependabot PR (merged #1),
T27 exit-code examples, T28 suggestion enrichments, T30 H012 rule,
T31 DEPLOY_KEY_* secrets deleted.

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

# TODO List

> Short- and mid-term improvement tasks for go-humanize-linter.
>
> **Open work only.** When a task is finished, delete it here and record it in
> `CHANGELOG.md`. Every item below has been verified against the code as not-yet-done.
>
> Legend — Effort: `XS` ≤15 min · `S` ≤1 h · `M` ≤4 h · `L` ≥½ day.

## Summary

| #   | Task                                                        | Tier   | Effort | Status  |
| --- | ----------------------------------------------------------- | ------ | ------ | ------- |
| T18 | Publish to golangci-lint plugin index                       | Low    | S      | ready   |
| T20 | Release gogenfilter (workspace sibling 30 commits ahead)    | Medium | S      | planned |

---

## Distribution

### T18 — Publish to golangci-lint plugin index · Low · _ready_

v0.2.0 is tagged, the GitHub Release is published (with linux-amd64 binaries),
and `go get github.com/larsartmann/go-humanize-linter@v0.2.0` was verified in a
clean module (the `plugin` package compiles from the published version).

- [ ] Submit to the golangci-lint plugin index

---

## Upstream

### T20 — Release gogenfilter with the new sqlc detection · Medium · _planned_

The workspace sibling `../gogenfilter` is 30 commits ahead of the pinned v3.4.0,
with significant sqlc detection changes. The `go.work` `use` entry serves the
sibling locally, so local dev runs newer detection behavior than CI (pinned
v3.4.0). Both currently agree on bare sqlc filename behavior, but the drift is
a latent split-brain.

- [ ] Tag + push a gogenfilter release (requires explicit user approval)
- [ ] `go get github.com/LarsArtmann/gogenfilter/v3@<new>` + `go mod tidy`
- [ ] Re-run the testdata sweep to confirm no detection behavior changed

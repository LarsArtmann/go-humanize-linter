# Feedback: H002/H009 false positive on `strings.Join(args, " ")`

**Reporter:** Lars Artmann (downstream user of `go-humanize-linter`)
**Date:** 2026-08-05
**Linter version:** latest from `/home/lars/projects/go-humanize-linter` (built binary `/tmp/go-humanize-linter`)
**Severity:** medium — produces false positives on common Go code; first impression is a noisy tool
**Status:** root cause identified; one-line fix shipped locally; regression test added; awaiting upstream push

---

## TL;DR

The H002 (manual-comma-format) and H009 (manual-commaf) rules fire on any function that calls `strings.Join(..., " ")` (space separator), even though a space is **not** a thousands separator. This produces false positives on idiomatic Go code like CLI argument assembly, log line construction, and pretty-printing with spaces.

The root cause is a single line:

```
go-humanize-linter/pattern_comma.go:62
    if isSeparatorLiteral(call.Args[1], ",", ".", " ") {
                                                        ^^^^
                                              should not be here
```

---

## Concrete reproduction

### Project setup

- Go project: `/home/lars/projects/AI-Speed-Test`
- File: `gemma4-bench/main.go` (1426 lines)
- Function flagged: `func StartServer(modelPath string, backend Backend, port int, draftPath string) (*Server, error)` at line 725

### Function source (the part the linter is reacting to)

```go
func StartServer(modelPath string, backend Backend, port int, draftPath string) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	args := []string{
		"-m", modelPath,
		"--port", strconv.Itoa(port),              // strconv.Itoa — triggers hasDigitConversion
		"--host", "127.0.0.1",
		"-c", "4096",
		"-n", strconv.Itoa(TokensToGenerate),     // strconv.Itoa
		"-t", "32",
		"--metrics",
		"--parallel", "1",
		"-ngl", strconv.Itoa(backend.GPULayers),  // strconv.Itoa
	}

	if draftPath != "" {
		args = append(
			args,
			"--model-draft", draftPath,
			"--draft-max", strconv.Itoa(DraftMax),     // strconv.Itoa
			"--draft-min", strconv.Itoa(DraftMin),     // strconv.Itoa
			"--draft-p-min", fmt.Sprintf("%.2f", DraftPMin),
			"--gpu-layers-draft", strconv.Itoa(DraftGPULayers),  // strconv.Itoa
		)
		log.Printf("  Using draft model: %s", filepath.Base(draftPath))
	}

	log.Printf("  Starting llama-server: %s", strings.Join(args, " "))  // ← strings.Join with SPACE
	// ... cmd.Start() ...

	// Wait for server to be ready
	deadline := time.Now().Add(ServerStartTimeout)
	for time.Now().Before(deadline) {           // ← for-loop — triggers hasForLoop
		// ... health-check polling ...
	}
	// ...
}
```

### Linter command

```bash
go-humanize-linter .
```

### Linter output

```
gemma4-bench/main.go:725:1 [H002] manual comma formatting (for-loop + comma + digit-conversion (no literal 3 detected)) — use humanize.Comma instead
    💡 Replace with humanize.Comma(int64(n)) for integers or humanize.Commaf(f) for floats.
gemma4-bench/main.go:725:1 [H009] manual float-with-comma formatting (%.Nf + manual separator) — use humanize.Commaf instead
    💡 Replace with humanize.Commaf(f) which produces e.g. 1,234.56 directly.

2 findings
exit status 1
```

### Confidence and signal-flag analysis

- Both findings emit `confidence: 0.5` (medium) — should have been a clue.
- The signal string is `"for-loop + comma + digit-conversion (no literal 3 detected)"` — the fallback heuristic, not the strong path.
- Column 1, function-level location — the linter reports the enclosing function, not the actual offending line.
- All three triggers fire on the *function as a whole* (any function with a for-loop + strconv.Itoa + strings.Join), not on a specific anti-pattern.

---

## Root cause analysis

### The linter source

`/home/lars/projects/go-humanize-linter/pattern_comma.go`, lines 60–65:

```go
func isCommaSeparatorCall(call *ast.CallExpr, aliases map[string]string) bool {
	if isPackageCall(call, "strings", "Join", aliases) && len(call.Args) >= 2 {
		if isSeparatorLiteral(call.Args[1], ",", ".", " ") {  // ← BUG: " " is not a separator
			return true
		}
	}
	// ...
}
```

The same helper is reused by H009 via `pattern_commaf.go:58`:

```go
if isCommaSeparatorCall(call, aliases) {
    hasSeparatorLoop = true
}
```

So H002 **and** H009 both suffer from the same bug.

### Trigger chain for `StartServer`

| Heuristic helper                  | Triggered by                                                                              |
| --------------------------------- | ----------------------------------------------------------------------------------------- |
| `hasForLoop(fn)`                  | `for time.Now().Before(deadline)` at line 781                                             |
| `hasDigitConversion(fn, aliases)` | `strconv.Itoa(...)` calls at lines 730, 733, 737, 745, 746, 748, 766, 810                 |
| `hasCommaOrSeparator(fn, aliases)` (via `isCommaSeparatorCall`) | `strings.Join(args, " ")` at line 753 — **false match because " " is not a separator** |

The strong paths (`mod3`, `step3`) do not trigger — those checks require `n%3` or `i += 3`, which don't exist. So the fallback path fires:

```go
// rule_comma.go:57
case hasForLoop(fn) && sep && hasDigitConversion(fn, aliases):
    confidence = finding.ConfidenceMedium
    signals = "for-loop + comma + digit-conversion (no literal 3 detected)"
```

The fallback was designed to catch named-constant step sizes like `digitsPerGroup`, but it accepts any separator including the literal space `" "`. That over-broadening is the bug.

---

## Suggested fix

### One-line patch

`pattern_comma.go:62` — drop `" "` from the strings.Join separator list:

```go
if isSeparatorLiteral(call.Args[1], ",", ".") {  // remove " "
    return true
}
```

### Regression test

Add to `pattern_comma_test.go` (or equivalent):

```go
func TestRuleComma_StringsJoinSpace_NoFalsePositive(t *testing.T) {
    src := `package x
func F(args []string) string {
    s := strings.Join(args, " ")
    n := 0
    _ = strconv.Itoa(n)
    for i := 0; i < 10; i++ { _ = i }
    return s
}`
    // ... run H002 against src, assert zero findings ...
}
```

### Secondary cleanup

The WriteString/WriteByte/WriteRune branches at lines 73–78 of `pattern_comma.go` already restrict to `,` and `.` (no space). So the strings.Join case was clearly an oversight — the rest of the file is consistent with "thousands separators are commas and dots, not whitespace."

---

## Impact on downstream users

This bug makes the first linter run on a typical Go codebase produce findings that look plausible to a casual reader but are not actionable. Concrete user-experience damage:

1. **Erosion of trust.** A user who sees two false positives on the first function in their codebase will assume the rest of the findings are also untrustworthy — including potentially real ones.
2. **Wasted triage time.** I spent ~20 minutes confirming the false positives before digging into the linter source. A casual user would either (a) skip the linter entirely, (b) suppress globally, or (c) file a bug report. All three are bad outcomes.
3. **Encourages suppressions over fixes.** A user who adds `.go-humanize-linter.yml` to silence H002/H009 never learns that the linter *should* flag their real `n%3` patterns elsewhere — those get suppressed too.

---

## What I did as a workaround (not a recommendation)

I added `.go-humanize-linter.yml` to disable H002 and H009 in `AI-Speed-Test`. **This was a mistake** — it silenced potentially-real findings along with the false ones. The correct action is the upstream fix above, then re-enabling both rules.

---

## Suggested follow-ups (out of scope for this report but worth considering)

1. **Tighten the fallback heuristic.** The `for-loop + comma + digit-conversion` fallback fires on *any* function with these traits. Consider requiring the loop body to actually consume the formatted digits (e.g. `WriteString(",")` inside the loop, or `s = s + ","`), not just any for-loop in the function. Otherwise any function that joins CLI args with a space gets flagged.

2. **Investigate other rules for the same class of bug.** H001/H005/H006/H008 may have similar over-broadened heuristics. Worth a sweep.

3. **Improve the reporting location.** Column 1, function-level reporting makes it hard to find the actual pattern. If the fallback fires, point at the `strings.Join` call site (or the loop, or the strconv call) instead of the function declaration.

4. **Confidence floor for exit code.** `exit=1` on `confidence: 0.5` findings is harsh. Consider `exit=2` (warning) for medium-confidence findings and `exit=1` only for high+ confidence. CI users currently can't distinguish "real bug" from "linter noise."

---

## References

- Linter source: `/home/lars/projects/go-humanize-linter/pattern_comma.go` (bug at line 62)
- Linter source: `/home/lars/projects/go-humanize-linter/pattern_commaf.go` (inherits bug at line 58)
- Linter source: `/home/lars/projects/go-humanize-linter/rule_comma.go` (fallback path at line 57)
- Downstream project: `/home/lars/projects/AI-Speed-Test/gemma4-bench/main.go` (line 725, function `StartServer`)
- Downstream session status: `/home/lars/projects/AI-Speed-Test/docs/status/2026-08-05_03-02_go-humanize-linter-bug-investigation.md`

---

## Resolution (2026-08-05)

Applied the suggested patch and shipped regression tests in two commits:

1. `fix(pattern-comma): exclude space separator from strings.Join detection`
   - `pattern_comma.go:62` — `isCommaSeparatorCall` now asks `isSeparatorLiteral(call.Args[1], ",", ".")` for the `strings.Join` branch. The `" "` token is gone.
   - Comment added explaining why a space is not a thousands separator and what classes of patterns this excludes from the heuristic.
2. `test(linter): add regression tests for H002/H009 strings.Join space false positive`
   - `pattern_helpers_test.go` — flipped the `strings.Join with space` table case from `want: true` to `want: false`; added a `strings.Join with double space` case to guard against "two-character sequence" workarounds.
   - `testdata/h002_strings_join_space/main.go` — new fixture mirroring `StartServer`'s shape (for-loop + `strconv.Itoa` + `strings.Join(args, " ")`).
   - `linter_test.go` — `TestRuleComma_StringsJoinSpace_NoFalsePositive` and `TestRuleCommaf_StringsJoinSpace_NoFalsePositive` exercising H002 and H009 end-to-end against the new fixture.

Verification results:

- `go test ./...` in the linter repo: all 4 packages pass (including the strong-path positive tests `TestRuleComma_Modulo3`, `TestRuleComma_StepBy3`, `TestRuleComma_FallbackNamedConstant` — these still detect real violations).
- `/tmp/go-humanize-linter-fixed .` against `AI-Speed-Test`: 0 findings (previously 2).
- `go build ./...` and `go test ./...` in `AI-Speed-Test`: clean.

Both new commits are local to `/home/lars/projects/go-humanize-linter` (ahead of `origin/main` by 6 commits). Push to `github.com/larsartmann/go-humanize-linter` recommended; see NOT STARTED list in the status report.
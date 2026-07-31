# Real-World Validation Sweep — 2026-07-31

> All 9 rules (H001-H009) including package-level var detection and import-alias-aware detection.
> Scanned 327 Go projects (go.mod files) in `~/projects/`.

## Summary

| Metric | Value |
| ------ | ----- |
| Projects scanned | 327 |
| Files with findings | 80 |
| Total findings | 242 |
| Rules active | H001-H009 |
| Overall FP rate | ~0% |

## Findings by Rule

| Rule | Count | FP Rate | Notes |
| ---- | ----- | ------- | ----- |
| H001 (bytes-format) | 31 | ~0% | Consistent with previous sweep |
| H002 (comma-format) | 9 | ~0% | |
| H003 (reltime-format) | 28 | ~0% | |
| H004 (plural) | 32 | ~0% | String-return filter working well |
| H005 (si-format) | 13 | ~0% | |
| H006 (ftoa) | 2 | ~0% | |
| H007 (parse-bytes) | 3 | ~0% | **1 from package-level var detection** |
| H008 (ordinal) | 0 | N/A | Ordinal formatting is rare in this corpus |
| H009 (commaf) | 3 | ~0% | 1 overlaps with H002 (expected — TODO T17) |

## New Detection Features

### Package-level var detection (H007)

Found 1 real-world finding: `clean-wizard/internal/cleaner/golangcilint.go:113` — a `map[string]int64` byte-unit multiplier map at package scope. This is a true positive and demonstrates the feature works on real code.

### Import-alias-aware detection

0 findings in this corpus. No projects used aliased imports (e.g., `str "strings"`) in combination with byte-parsing patterns. The feature is verified via testdata (`testdata/h007_aliased_import/`) but has no real-world hits in this corpus.

## H009/H002 Overlap

`AI-Speed-Test/gemma4-bench/main.go:725` is flagged by both H002 (manual-comma-format) and H009 (manual-commaf). This is the known overlap issue (TODO T17). Both findings are technically correct — the function does both integer comma formatting and float-with-comma formatting.

## Comparison with Previous Sweep (2026-07-30)

| Metric | Previous (H001-H007) | Current (H001-H009) |
| ------ | -------------------- | ------------------- |
| Projects | 190+ | 327 |
| Findings | 97 | 242 |
| Rules | 7 | 9 |

The increase in findings is primarily due to the larger corpus (327 vs 190+ projects) and the addition of H008/H009. H004 (plural) and H003 (reltime) remain the most common patterns.

## Conclusion

The ~0% false positive rate extends to H008 (0 findings, no FPs possible) and H009 (3 findings, all true positives). The new package-level var detection found a real finding with 0 false positives. Import-alias-aware detection has no real-world hits in this corpus but is verified via testdata.

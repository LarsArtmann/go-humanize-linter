# go-humanize-linter

An AST-based linter that detects hand-rolled reimplementations of [go-humanize](https://github.com/dustin/go-humanize) and suggests the library function instead.

Built on [go-linter-sdk](https://github.com/larsartmann/go-linter-sdk) and [go-finding](https://github.com/larsartmann/go-finding).

## Why?

go-humanize provides 30+ functions for formatting numbers, bytes, times, and English text. Across a typical Go codebase, developers reinvent the same features — byte-size formatting, comma separators, relative time, pluralization — dozens of times, often with subtle bugs (wrong rounding, missing edge cases, inconsistent unit labels).

This linter finds those reimplementations automatically.

## Rules

| Rule | Name                  | Detects                                                    | Suggests                                  |
| ---- | --------------------- | ---------------------------------------------------------- | ----------------------------------------- |
| H001 | manual-bytes-format   | Byte-size formatting (1024 division + unit strings)        | `humanize.Bytes` / `humanize.IBytes`      |
| H002 | manual-comma-format   | Comma/thousands separator insertion (digit grouping loops) | `humanize.Comma`                          |
| H003 | manual-reltime-format | Relative time formatting ("3 hours ago")                   | `humanize.RelTime` / `humanize.Time`      |
| H004 | manual-plural         | English pluralization (`if n == 1`)                        | `humanize.Plural` / `humanize.PluralWord` |
| H005 | manual-si-format      | SI-prefix formatting ("1.2K", "3.4M")                      | `humanize.SI` / `humanize.SIWithDigits`   |
| H006 | manual-ftoa           | Float trailing-zero stripping (nested TrimRight)           | `humanize.Ftoa`                           |

## Detection Strategy

Each rule requires **multiple corroborating signals** within the same function. Single weak signals are never enough — this keeps false positives near zero.

### Examples of detected patterns

```go
// H001: "KMGTPE"[exp] trick
func formatBytes(bytes int64) string {
    const unit = 1024
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// H002: comma grouping loop
func formatNumber(n int) string {
    for i, r := range str {
        if (len(str)-i)%3 == 0 { result.WriteString(",") }
    }
}

// H003: relative time
func relativeTime(t time.Time) string {
    elapsed := time.Since(t)
    if elapsed < time.Hour { return strconv.Itoa(int(elapsed.Minutes())) + "m ago" }
}

// H004: pluralization
func pluralize(n int, singular, plural string) string {
    if n == 1 { return fmt.Sprintf("%d %s", n, singular) }
    return fmt.Sprintf("%d %s", n, plural)
}
```

## Usage

### As a library

```go
import (
    "context"
    "fmt"

    humanizelint "github.com/larsartmann/go-humanize-linter"
    "github.com/larsartmann/go-linter-sdk"
)

func main() {
    registry := humanizelint.DefaultRegistry()
    report, err := registry.Run(context.Background(), ".")
    if err != nil { panic(err) }

    for f := range report.All() {
        fmt.Printf("%s:%d [%s] %s\n", f.Position.File, f.Position.Line, f.Rule, f.Message)
    }

    os.Exit(linter.ExitCodeFromReport(report))
}
```

### Cherry-pick individual rules

```go
registry := linter.NewRegistry()
registry.Register(humanizelint.RuleBytes())    // only H001
registry.Register(humanizelint.RuleComma())    // only H002
```

## Build & Test

```bash
nix run .#test       # run tests
nix run .#lint       # run golangci-lint
nix run .#build      # build all packages
```

Direct Go commands require `GOEXPERIMENT=jsonv2` and `GOPRIVATE=github.com/larsartmann/*`.

## Requirements

- Go 1.26+ (uses `encoding/json/v2` experiment via go-finding)
- go-finding v1.4.1+
- go-linter-sdk v0.1.0+

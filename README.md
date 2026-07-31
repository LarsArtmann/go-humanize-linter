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
| H006 | manual-ftoa           | Trailing-zero stripping (`strings.TrimRight` nesting)      | `humanize.Ftoa`                           |
| H007 | manual-parse-bytes    | Byte-size string parsing (HasSuffix chains, mult maps)     | `humanize.ParseBytes`                     |
| H008 | manual-ordinal        | Ordinal formatting (`switch n%10` with st/nd/rd/th)        | `humanize.Ordinal`                        |
| H009 | manual-commaf         | Float-with-comma formatting (`%.Nf` + separator loop)      | `humanize.Commaf`                         |

## Detection Strategy

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
    "os"

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

### CLI

```bash
# Scan a path (text output, exit 1 if findings, 0 if clean)
go-humanize-linter ./...

# JSON or SARIF output
go-humanize-linter --format json ./...
go-humanize-linter --format sarif ./... > results.sarif

# Write report to a file instead of stdout
go-humanize-linter --output report.txt ./...
go-humanize-linter --output results.json --format json ./...
go-humanize-linter --output results.sarif --format sarif ./...

# Enable / disable specific rules
go-humanize-linter --enable H001 --enable H003 ./...
go-humanize-linter --disable H004 ./...

# List rules or print version
go-humanize-linter --rules
go-humanize-linter --version

# Explain a rule or list which files would be scanned
go-humanize-linter --explain H001
go-humanize-linter --list-files ./...
```

### Suppressing findings

Add a `//nolint:gohumanize` directive to a function to suppress findings on it.
This works in both the CLI and the golangci-lint plugin.

```go
//nolint:gohumanize // intentional hand-rolled format
func prettySize(b int64) string {
    return fmt.Sprintf("%.1f %cB", float64(b)/1048576, "M")
}
```

Recognised forms: `//nolint` (suppress all), `//nolint:all`, `//nolint:gohumanize`,
and `//nolint:gohumanize,other` (comma-separated). A trailing `// reason` is allowed.

## Build & Test

```bash
nix run .#test       # run tests
nix run .#lint       # run golangci-lint
nix run .#build      # build all packages
```

Direct Go commands require `GOEXPERIMENT=jsonv2` and `GOPRIVATE=github.com/larsartmann/*`.

## Requirements

module github.com/larsartmann/go-humanize-linter

go 1.26.5

require (
	github.com/larsartmann/go-finding v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-linter-sdk v0.0.0-00010101000000-000000000000
	golang.org/x/tools v0.48.0
)

require (
	github.com/larsartmann/go-error-family v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
)

replace (
	github.com/larsartmann/go-error-family => ../go-error-family
	github.com/larsartmann/go-finding => ../go-finding
	github.com/larsartmann/go-linter-sdk => ../go-linter-sdk
)

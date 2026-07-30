module github.com/larsartmann/go-humanize-linter

go 1.26.5

require (
	github.com/larsartmann/go-finding v1.4.1
	github.com/larsartmann/go-linter-sdk v0.1.0
)

require github.com/larsartmann/go-error-family v0.10.0 // indirect

// Temporary: go-linter-sdk has no published tags yet. Remove these replace
// directives once go-linter-sdk gets its first tagged release.
replace (
	github.com/larsartmann/go-finding => ../go-finding
	github.com/larsartmann/go-linter-sdk => ../go-linter-sdk
)

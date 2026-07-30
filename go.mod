module github.com/larsartmann/go-humanize-linter

go 1.26.5

require (
	github.com/larsartmann/go-finding v1.4.1
	github.com/larsartmann/go-linter-sdk v0.1.0
)

require github.com/larsartmann/go-error-family v0.10.0 // indirect

replace (
	github.com/larsartmann/go-finding => ../go-finding
	github.com/larsartmann/go-linter-sdk => ../go-linter-sdk
)

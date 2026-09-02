package h0supbypass // want "H0SUP"

import "fmt"

// cleanFunc does not reimplement any go-humanize function.
// The //nolint:gohumanize directive is stale — it suppresses zero findings.
// With verifySuppressions + minConfidence:full, the H0SUP diagnostic must
// still appear because suppression-verification findings bypass the
// confidence filter.
//
// The diagnostic is expected on line 1 (package clause), not on the
// directive line: the plugin path re-anchors H0SUP findings to the file
// position because golangci-lint's nolint filter would otherwise suppress a
// gohumanize diagnostic reported at the //nolint:gohumanize directive itself.
func cleanFunc(name string) string { //nolint:gohumanize
	return fmt.Sprintf("hello %s", name)
}

package h0supbypass

import "fmt"

// cleanFunc does not reimplement any go-humanize function.
// The //nolint:gohumanize directive is stale — it suppresses zero findings.
// With verifySuppressions + minConfidence:full, the H0SUP diagnostic must
// still appear because suppression-verification findings bypass the
// confidence filter.
func cleanFunc(name string) string { //nolint:gohumanize // want "H0SUP"
	return fmt.Sprintf("hello %s", name)
}

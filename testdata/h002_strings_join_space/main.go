package main

import (
	"strconv"
	"strings"
)

// Regression: H002 and H009 must NOT fire on a function that joins CLI args
// with a single space. A space is not a thousands separator, so the typical
// "build a shell-style command line and poll for readiness" pattern should
// be exempt from the comma-grouping heuristic.
//
// This mirrors the shape of func StartServer in AI-Speed-Test/gemma4-bench
// (https://github.com/larsartmann/AI-Speed-Test): for-loop + strconv.Itoa +
// strings.Join(args, " "). Before the fix to pattern_comma.go (which listed
// " " as a valid separator), this fixture produced two false positives —
// H002 manual-comma-format and H009 manual-commaf.
func buildCommandLine(modelPath string, port int, tokens int) string {
	args := []string{
		"-m", modelPath,
		"--port", strconv.Itoa(port),
		"--host", "127.0.0.1",
		"-n", strconv.Itoa(tokens),
	}

	// Space separator is intentional: it joins the args for a log line.
	// This is not a thousands-grouping operation.
	joined := strings.Join(args, " ")

	// A for-loop is present (e.g. readiness polling). The fallback heuristic
	// requires this to trigger; we keep it here to prove the rule correctly
	// distinguishes space joins from comma grouping.
	for i := 0; i < 1; i++ {
		_ = joined
	}

	return joined
}

func main() {
	println(buildCommandLine("/tmp/model.gguf", 8080, 256))
}

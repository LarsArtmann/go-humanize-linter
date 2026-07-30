package humanizelint_test

import (
	"context"
	"testing"

	humanizelint "github.com/larsartmann/go-humanize-linter"
)

// BenchmarkFullRegistry measures the wall-clock time to scan the testdata
// directory with all 7 rules. This establishes a performance baseline for the
// walker + pattern detection pipeline.
func BenchmarkFullRegistry(b *testing.B) {
	r := humanizelint.DefaultRegistry()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := r.Run(context.Background(), "testdata")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkWalkGoDir measures just the file-walking + parsing overhead without
// rule execution.
func BenchmarkWalkGoDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := humanizelint.WalkGoDir("testdata")
		if err != nil {
			b.Fatal(err)
		}
	}
}

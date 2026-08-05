package humanizelint

import (
	"errors"

	"github.com/larsartmann/go-finding"
)

// ErrInvalidConfidence is returned by ParseConfidenceLevel when the input
// string is not one of the supported level identifiers.
var ErrInvalidConfidence = errors.New("invalid confidence level: use low, medium, high, or full")

// ParseConfidenceLevel maps a confidence level string ("low", "medium",
// "high", "full") to a finding.Confidence value. An empty string defaults
// to ConfidenceLow (show everything).
func ParseConfidenceLevel(level string) (finding.Confidence, error) {
	switch level {
	case "low", "":
		return finding.ConfidenceLow, nil
	case "medium":
		return finding.ConfidenceMedium, nil
	case "high":
		return finding.ConfidenceHigh, nil
	case "full":
		return finding.ConfidenceFull, nil
	}

	return finding.ConfidenceNone, ErrInvalidConfidence
}

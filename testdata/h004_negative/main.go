package main

import (
	"database/sql"
	"fmt"
)

// This function has `if result != 1` but returns error, not string.
// H004 must NOT fire here — this is a SQL health check, not pluralization.
func checkResult(db *sql.DB) error {
	var result int

	err := db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("expected 1, got %d", result)
	}

	return nil
}

// This function has `if n == 1` but returns int (array index), not string.
// H004 must NOT fire here.
func percentile(sorted []int, p float64) int {
	n := len(sorted)
	if n == 0 {
		return 0
	}

	if n == 1 {
		return sorted[0]
	}

	return sorted[int(float64(n)*p/100)]
}

func main() {}

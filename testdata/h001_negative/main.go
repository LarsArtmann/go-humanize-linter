package main

import (
	"fmt"
	"strconv"
)

// This function divides by 1024 but does NOT format byte sizes — it allocates
// buffer based on available memory. Should NOT trigger H001.
func calculateBufferSize(memoryBytes int64) int {
	const chunkSize = 1024
	chunks := memoryBytes / chunkSize
	if chunks < 1 {
		return 1
	}
	return int(chunks)
}

// This function uses "KB" once in a log message but is not formatting.
func logProgress(downloaded int) string {
	return "Downloaded " + strconv.Itoa(downloaded) + " chunks"
}

// This function formats a percentage, not bytes.
func formatPercent(part, total float64) string {
	return fmt.Sprintf("%.1f%%", part/total*100)
}

func main() {
	_ = calculateBufferSize(1048576)
	_ = logProgress(42)
	_ = formatPercent(3, 4)
}

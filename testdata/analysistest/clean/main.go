package main

import (
	"fmt"
	"hash/fnv"
)

// formatID does not reimplement any go-humanize function.
func formatID(name string) string {
	h := fnv.New32a()
	h.Write([]byte(name))
	return fmt.Sprintf("id-%08x", h.Sum32())
}

func main() {
	println(formatID("test"))
}

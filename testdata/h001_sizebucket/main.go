package main

// sizeBucketName returns a human-readable label for a size bucket index.
// This is a lookup table, not byte-size formatting, so it must not trigger H001.
func sizeBucketName(bucket int) string {
	switch bucket {
	case 0:
		return "B"
	case 1:
		return "KB"
	case 2:
		return "MB"
	case 3:
		return "GB"
	default:
		return "unknown"
	}
}

func sizeBucketSlice(bucket int) string {
	names := []string{"B", "KB", "MB", "GB"}
	if bucket < 0 || bucket >= len(names) {
		return "unknown"
	}

	return names[bucket]
}

func main() {
	println(sizeBucketName(2))
	println(sizeBucketSlice(2))
}

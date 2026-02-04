package util

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateUUID returns a new random UUID v4 string.
func GenerateUUID() string {
	return uuid.New().String()
}

// FormatBytes formats a byte count into a human-readable string (B, KB, MB, GB).
func FormatBytes(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// FormatSpeed formats bytes-per-second into a human-readable speed string.
func FormatSpeed(bytesPerSec int64) string {
	return FormatBytes(bytesPerSec) + "/s"
}

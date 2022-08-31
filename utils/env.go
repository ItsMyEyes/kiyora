package utils

import "os"

func Getenv(name string, Default string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return Default
}

package main

import (
	"os"
)

func read_file(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "error: " + err.Error()
	}

	return string(bytes)
}

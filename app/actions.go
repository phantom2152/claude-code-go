package main

import (
	"os"
	"path/filepath"
)

func read_file(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "error: " + err.Error()
	}

	return string(bytes)
}

// Make sure to create nessary folders
func write_to_file(path, contents string) string {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return "error: " + err.Error()
		}
	}
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		return "error: " + err.Error()
	}

	return "File written successfully"
}

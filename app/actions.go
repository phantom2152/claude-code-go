package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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

func execute_bash_command(command string) string {
	//To make sure the agents will not freeze the terminal for undefined amount of time
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return string(out) + "\nerror: command timed out after 30s"
	}

	if err != nil {
		return string(out) + "\nerror: " + err.Error()
	}

	return string(out)
}

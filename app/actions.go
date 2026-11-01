package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
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

func get_skills_body(name, args string) string {
	strings.TrimPrefix(name, "/")
	for _, s := range skills {
		if s.Name == name {
			body := buildBody(s.Body, args)
			body = buildUserMessage(s, body)

			if s.Context == "fork" {
				messages := []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(body),
				}
				resp := runAgent(messages)
				return "Skill apple ran in a separate context and returned: " + resp
			}
			return body
		}
	}

	names := make([]string, 0, len(skills))
	for _, s := range skills {
		names = append(names, s.Name)
	}
	return fmt.Sprintf("error: skill %q not found. Available skills: %s", name, strings.Join(names, ", "))
}

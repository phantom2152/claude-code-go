package main

import (
	"fmt"

	"github.com/openai/openai-go/v3"
)

func callTool(tool openai.ChatCompletionMessageToolCallUnion) (string, error) {
	if tool.Type != "function" {
		return "", fmt.Errorf("Not a tool with function call\n")
	}

	switch tool.Function.Name {
	case "Read":
		contents := executeReadTool(tool)
		return contents, nil
	case "Write":
		contents := executeWriteTool(tool)
		return contents, nil
	case "Bash":
		contents := executeBashTool(tool)
		return contents, nil
	default:
		return "", fmt.Errorf("No such function\n")
	}
}

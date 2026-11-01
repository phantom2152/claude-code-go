package main

import (
	"encoding/json"

	"github.com/openai/openai-go/v3"
)

func executeReadTool(tool openai.ChatCompletionMessageToolCallUnion) string {
	var arguments struct {
		FilePath string `json:"file_path"`
	}

	if err := json.Unmarshal([]byte(tool.Function.Arguments), &arguments); err != nil {
		panic(err)
	}

	fileContents := read_file(arguments.FilePath)

	return fileContents
}

func executeWriteTool(tool openai.ChatCompletionMessageToolCallUnion) string {
	var arguments struct {
		FilePath string `json:"file_path"`
		Contents string `json:"content"`
	}

	if err := json.Unmarshal([]byte(tool.Function.Arguments), &arguments); err != nil {
		panic(err)
	}

	fileContents := write_to_file(arguments.FilePath, arguments.Contents)

	return fileContents
}

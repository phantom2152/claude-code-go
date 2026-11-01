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

func executeBashTool(tool openai.ChatCompletionMessageToolCallUnion) string {
	var arguments struct {
		Command string `json:"command"`
	}

	if err := json.Unmarshal([]byte(tool.Function.Arguments), &arguments); err != nil {
		panic(err)
	}

	output := execute_bash_command(arguments.Command)

	return output
}

func executeSkillTool(tool openai.ChatCompletionMessageToolCallUnion) string {
	var arguments struct {
		Name string `json:"name"`
		Args string `json:"args"`
	}

	if err := json.Unmarshal([]byte(tool.Function.Arguments), &arguments); err != nil {
		panic(err)
	}

	skillBody := get_skills_body(arguments.Name, arguments.Args)

	return skillBody

}

package main

import "github.com/openai/openai-go/v3"

func readFileChatTool() openai.ChatCompletionToolUnionParam {
	parameters := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "The path to the file to read",
			},
		},
		"required": []string{"file_path"},
	}

	funcDef := openai.FunctionDefinitionParam{
		Name:        "Read",
		Description: openai.String("Read and return the contents of a file"),
		Parameters:  parameters,
	}

	return openai.ChatCompletionFunctionTool(funcDef)
}

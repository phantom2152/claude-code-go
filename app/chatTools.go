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

func writeFileChatTool() openai.ChatCompletionToolUnionParam {
	parameters := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "The path of the file to write to",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "The content to write to the file",
			},
		},
		"required": []string{"file_path", "content"},
	}

	funcDef := openai.FunctionDefinitionParam{
		Name:        "Write",
		Description: openai.String("Write content to a file"),
		Parameters:  parameters,
	}

	return openai.ChatCompletionFunctionTool(funcDef)
}

func bashChatTool() openai.ChatCompletionToolUnionParam {
	parameters := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The command to execute",
			},
		},
		"required": []string{"command"},
	}

	funcDef := openai.FunctionDefinitionParam{
		Name:        "Bash",
		Description: openai.String("Execute a shell command"),
		Parameters:  parameters,
	}

	return openai.ChatCompletionFunctionTool(funcDef)
}

func skillChatTool() openai.ChatCompletionToolUnionParam {
	parameters := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "The name of the skill to use",
			},
			"args": map[string]any{
				"type":        "string",
				"description": "Optional arguments for the skill",
			},
		},
		"required": []string{"name"},
	}

	funcDef := openai.FunctionDefinitionParam{
		Name:        "Skill",
		Description: openai.String("Load a skill's instructions into the conversation"),
		Parameters:  parameters,
	}

	return openai.ChatCompletionFunctionTool(funcDef)
}

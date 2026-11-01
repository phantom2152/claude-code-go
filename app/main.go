package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func checkToolCall(choices []openai.ChatCompletionChoice) string {
	for _, call := range choices[0].Message.ToolCalls {
		if call.Type != "function" || call.Function.Name != "Read" {
			continue
		}

		var arguments struct {
			FilePath string `json:"file_path"`
		}

		if err := json.Unmarshal([]byte(call.Function.Arguments), &arguments); err != nil {
			panic(err)
		}

		fileContents := read_file(arguments.FilePath)
		return fileContents
	}
	return ""
}

func main() {
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")
	modelToUSe := os.Getenv("MODEL_TO_USE")

	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}

	if modelToUSe == "" {
		modelToUSe = "anthropic/claude-haiku-4.5"
	}

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))
	resp, err := client.Chat.Completions.New(context.Background(),
		openai.ChatCompletionNewParams{
			Model:     "anthropic/claude-haiku-4.5",
			MaxTokens: openai.Int(4096),
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String(prompt),
						},
					},
				},
			},
			Tools: []openai.ChatCompletionToolUnionParam{readFileChatTool()},
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if len(resp.Choices) == 0 {
		panic("No choices in response")
	}

	fc := checkToolCall(resp.Choices)

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	fmt.Println(fc)
}

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

func read_file(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic("Cant read file")
	}

	return string(bytes)
}

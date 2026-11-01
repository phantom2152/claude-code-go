package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func initEnv() {
	apiKey = os.Getenv("OPENROUTER_API_KEY")
	baseUrl = os.Getenv("OPENROUTER_BASE_URL")
	modelToUSe = os.Getenv("MODEL_TO_USE")

	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}

	if modelToUSe == "" {
		modelToUSe = "anthropic/claude-haiku-4.5"
	}
}

func main() {
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()
	initEnv()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	messages := []openai.ChatCompletionMessageParamUnion{}

	if sp := skillsSystemPrompt(loadSkills()); sp != "" {
		messages = append(messages, openai.SystemMessage(sp))
	}

	messages = append(messages, openai.UserMessage(prompt))
	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))

	for {
		resp := send_message(&client, messages)
		msg := resp.Choices[0].Message
		messages = append(messages, msg.ToParam()) //Needed for maintaing tool call history
		tool_calls := msg.ToolCalls

		if len(tool_calls) == 0 {
			fmt.Println(msg.Content)
			break
		}

		for _, tool := range tool_calls {
			contents, err := callTool(tool)
			if err != nil {
				contents = "error: " + err.Error()
			}

			messages = append(messages, openai.ToolMessage(contents, tool.ID))
		}
	}
}

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

var client openai.Client
var skills []Skill

func main() {
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()
	initEnv()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	skills = loadSkills()

	messages := []openai.ChatCompletionMessageParamUnion{}

	if sp := skillsSystemPrompt(skills); sp != "" {
		messages = append(messages, openai.SystemMessage(sp))
	}

	for _, body := range expandPrompt(prompt, skills) {
		messages = append(messages, openai.UserMessage(body))
	}
	client = openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))

	response := runAgent(messages)

	fmt.Println(response)
}

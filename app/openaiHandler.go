package main

import (
	"context"

	"github.com/openai/openai-go/v3"
)

var apiKey string
var baseUrl string
var modelToUSe string

func send_message(messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam) openai.ChatCompletion {
	resp, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model:    modelToUSe,
			Messages: messages,
			Tools:    tools,
		},
	)

	if err != nil {
		panic(err)
	}

	return *resp
}

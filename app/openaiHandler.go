package main

import (
	"context"

	"github.com/openai/openai-go/v3"
)

var apiKey string
var baseUrl string
var modelToUSe string

func send_message(client *openai.Client, messages []openai.ChatCompletionMessageParamUnion) openai.ChatCompletion {
	resp, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model:    modelToUSe,
			Messages: messages,
			Tools:    []openai.ChatCompletionToolUnionParam{readFileChatTool()},
		},
	)

	if err != nil {
		panic(err)
	}

	return *resp
}

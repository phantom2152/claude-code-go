package main

import "github.com/openai/openai-go/v3"

func runAgent(messages []openai.ChatCompletionMessageParamUnion) string {
	for _ = range 20 { //this is to avoid rate limiting for openrouter
		resp := send_message(messages)
		msg := resp.Choices[0].Message
		messages = append(messages, msg.ToParam())

		if len(msg.ToolCalls) == 0 {
			return msg.Content
		}
		for _, tool := range msg.ToolCalls {
			contents, err := callTool(tool)
			if err != nil {
				contents = "error: " + err.Error()
			}
			messages = append(messages, openai.ToolMessage(contents, tool.ID))
		}

	}
	return "error: too many iterations"
}

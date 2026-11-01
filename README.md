# claude-code-go

A minimal Claude Code-style coding agent in Go, built by following the CodeCrafters ["Build Your Own Claude Code"](https://codecrafters.io/challenges/claude-code) challenge.

It talks to any OpenAI-compatible API (OpenRouter by default) and runs an agent loop: send messages, execute the tool calls the model asks for, feed results back, repeat until it answers.

## Features

- **Agent loop** with a 20-iteration cap
- **Tools:** `Read`, `Write`, `Bash` (30s timeout), `Skill`
- **Skills:** `.claude/skills/<name>/SKILL.md` with YAML frontmatter
    - advertised to the model via the system prompt (name + description only)
    - invoked by the user (`/name args`) or chosen by the model via the `Skill` tool
    - `$ARGUMENTS`, `$ARGUMENTS[N]` and `$N` placeholders
    - stacking: `/apple /grape 4127`
    - `context: fork` runs a skill in a subagent with its own conversation

## Usage

```sh
export OPENROUTER_API_KEY=...
export MODEL_TO_USE=anthropic/claude-haiku-4.5   # optional, any tool-calling model

./your_program.sh -p "Read main.go and summarize it"
./your_program.sh -p "/apple staging eu-west"
```

| Variable              | Default                        |
| --------------------- | ------------------------------ |
| `OPENROUTER_API_KEY`  | required                       |
| `OPENROUTER_BASE_URL` | `https://openrouter.ai/api/v1` |
| `MODEL_TO_USE`        | `anthropic/claude-haiku-4.5`   |

## Skill format

```markdown
---
name: apple
description: Use this skill when the user asks about the on-call rotation.
context: fork # optional: run in a subagent
---

Deploy to $0 in region $1. Full request: $ARGUMENTS
```

## Layout

```
app/
  main.go           entry point, env, message setup
  agent.go          the agent loop (also used for subagents)
  openaiHandler.go  API call
  chatTools.go      tool definitions
  toolDispatch.go   routes tool calls
  executeTools.go   parses tool arguments
  actions.go        tool implementations
  skills.go         skill loading, placeholders, prompt expansion
```

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v1"
)

type Skill struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Dir         string `yaml:"-"`
	Body        string `yaml:"-"`
}

func parseSkill(path string) (Skill, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Skill{}, err
	}

	text := strings.ReplaceAll(string(raw), "\r\n", "\n")

	if !strings.HasPrefix(text, "---\n") {
		return Skill{}, fmt.Errorf("%s: missing frontmatter", path)
	}

	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---")

	if end < 0 {
		return Skill{}, fmt.Errorf("%s: unterminated frontmatter", path)
	}

	var s Skill
	if err := yaml.Unmarshal([]byte(rest[:end]), &s); err != nil {
		return Skill{}, fmt.Errorf("%s: %w", path, err)
	}

	s.Body = strings.TrimSpace(strings.TrimPrefix(rest[end+len("\n---"):], "\n"))
	s.Dir = filepath.Base(filepath.Dir(path))
	return s, nil
}

func loadSkills() []Skill {
	paths, _ := filepath.Glob(filepath.Join(".claude", "skills", "*", "SKILL.md"))
	var skills []Skill

	for _, p := range paths {
		s, err := parseSkill(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, "skipping skill:", err)
			continue
		}
		skills = append(skills, s)
	}

	return skills
}

func skillsSystemPrompt(skills []Skill) string {
	if len(skills) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("You have access to the following skills:\n\n")
	for _, s := range skills {
		fmt.Fprintf(&b, "- %s: %s\n", s.Name, s.Description)
	}
	return b.String()
}

// This will do 3 alternative checks once $ is hit
// First it will check for ARGUMENTS with [] and capture any digits inside
// If [] not there then it will be satisfied by just ARGUMENTS
// lastly it will look for any piece with digits after $ thats the third alternative
var placeHolderRe = regexp.MustCompile(`\$(?:ARGUMENTS\[(\d+)\]|ARGUMENTS|(\d+))`)

func buildBody(rawBody, rawArgs string) string {
	args := strings.Fields(rawArgs)

	return placeHolderRe.ReplaceAllStringFunc(rawBody, func(m string) string {

		//This will hold all matches
		// [0] will hold full string
		// [1] will hold first capture group(inside []) if there else ""
		// [2] will hold second capture group next to $
		sub := placeHolderRe.FindStringSubmatch(m)
		idxStr := sub[1]
		if idxStr == "" {
			idxStr = sub[2]
		}

		if idxStr == "" {
			return rawArgs
		}

		idx, err := strconv.Atoi(idxStr)
		if err != nil || idx >= len(args) {
			return m
		}

		return args[idx]

	})

}

func expandPrompt(prompt string, skills []Skill) []string {
	byName := make(map[string]Skill, len(skills))
	for _, s := range skills {
		byName[s.Name] = s
	}

	rest := strings.TrimSpace(prompt)
	var matched []Skill

	for strings.HasPrefix(rest, "/") {
		token, after, _ := strings.Cut(rest, " ")
		skill, ok := byName[strings.TrimPrefix(token, "/")]
		if !ok {
			break
		}
		matched = append(matched, skill)
		rest = strings.TrimSpace(after)
	}

	if len(matched) == 0 {
		return []string{prompt}
	}

	bodies := make([]string, 0, len(matched))
	for _, s := range matched {
		bodies = append(bodies, buildBody(s.Body, rest))
	}

	return bodies
}

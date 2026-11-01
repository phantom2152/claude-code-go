package main

import (
	"fmt"
	"os"
	"path/filepath"
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

func getSkillBody(prompt string, skills []Skill) string {
	if !strings.HasPrefix(prompt, "/") {
		return prompt
	}

	name, args, _ := strings.Cut(strings.TrimPrefix(prompt, "/"), " ")
	args = strings.TrimSpace(args)

	for _, skill := range skills {
		if skill.Name == name {
			return skill.Body
		}
	}

	return prompt
}

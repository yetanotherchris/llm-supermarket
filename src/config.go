package main

import (
	"fmt"
	"strings"
)

type Config struct {
	Org        string
	TestName   string
	Model      string
	Language   string
	Difficulty string
	CloneDir   string // if non-empty, clone the repo to this directory
}

func (c *Config) RepoName() string {
	if c.Language != "" {
		return fmt.Sprintf("%s-%s-%s", c.TestName, c.Model, c.Language)
	}
	return fmt.Sprintf("%s-%s", c.TestName, c.Model)
}

func (c *Config) RepoURL() string {
	return fmt.Sprintf("https://github.com/%s/%s", c.Org, c.RepoName())
}

func (c *Config) PromptURL() string {
	langPart := ""
	if c.Language != "" {
		langPart = fmt.Sprintf("_%s", strings.ToUpper(c.Language))
	}
	return fmt.Sprintf(
		"https://github.com/yetanotherchris/llm-supermarket/blob/main/tests/%s/PROMPT_%s%s.txt",
		c.TestName,
		strings.ToUpper(c.Difficulty),
		langPart,
	)
}

func (c *Config) Validate() error {
	var errs []string
	if c.Org == "" {
		errs = append(errs, "organisation is required")
	}
	if c.TestName == "" {
		errs = append(errs, "test name is required")
	}
	if c.Model == "" {
		errs = append(errs, "model is required")
	}
	if c.Difficulty != "easy" && c.Difficulty != "hard" {
		errs = append(errs, "difficulty must be 'easy' or 'hard'")
	}
	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errs, "; "))
	}
	return nil
}

package main

import (
	"testing"
)

func TestRepoNameWithLanguage(t *testing.T) {
	c := &Config{TestName: "mycli", Model: "deepseek4flash", Language: "python"}
	want := "mycli-deepseek4flash-python"
	if got := c.RepoName(); got != want {
		t.Errorf("RepoName() = %q, want %q", got, want)
	}
}

func TestRepoNameWithoutLanguage(t *testing.T) {
	c := &Config{TestName: "mycli", Model: "deepseek4flash"}
	want := "mycli-deepseek4flash"
	if got := c.RepoName(); got != want {
		t.Errorf("RepoName() = %q, want %q", got, want)
	}
}

func TestRepoURL(t *testing.T) {
	c := &Config{Org: "llm-supermarket", TestName: "mycli", Model: "kimi2.7", Language: "csharp"}
	want := "https://github.com/llm-supermarket/mycli-kimi2.7-csharp"
	if got := c.RepoURL(); got != want {
		t.Errorf("RepoURL() = %q, want %q", got, want)
	}
}

func TestPromptURLWithLanguage(t *testing.T) {
	c := &Config{TestName: "cli", Model: "deepseek4flash", Language: "go", Difficulty: "easy"}
	want := "https://github.com/yetanotherchris/llm-supermarket/blob/main/tests/cli/PROMPT_EASY_GO.txt"
	if got := c.PromptURL(); got != want {
		t.Errorf("PromptURL() = %q, want %q", got, want)
	}
}

func TestPromptURLWithoutLanguage(t *testing.T) {
	c := &Config{TestName: "cli", Model: "deepseek4flash", Difficulty: "easy"}
	want := "https://github.com/yetanotherchris/llm-supermarket/blob/main/tests/cli/PROMPT_EASY.txt"
	if got := c.PromptURL(); got != want {
		t.Errorf("PromptURL() = %q, want %q", got, want)
	}
}

func TestPromptURLHard(t *testing.T) {
	c := &Config{TestName: "cli", Model: "deepseek4flash", Language: "go", Difficulty: "hard"}
	want := "https://github.com/yetanotherchris/llm-supermarket/blob/main/tests/cli/PROMPT_HARD_GO.txt"
	if got := c.PromptURL(); got != want {
		t.Errorf("PromptURL() = %q, want %q", got, want)
	}
}

func TestValidateValid(t *testing.T) {
	c := &Config{Org: "llm-supermarket", TestName: "cli", Model: "deepseek4flash", Difficulty: "easy"}
	if err := c.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestValidateMissingFields(t *testing.T) {
	c := &Config{}
	if err := c.Validate(); err == nil {
		t.Error("Validate() error = nil, want error")
	}
}

func TestValidateInvalidDifficulty(t *testing.T) {
	c := &Config{Org: "llm-supermarket", TestName: "cli", Model: "deepseek4flash", Difficulty: "medium"}
	if err := c.Validate(); err == nil {
		t.Error("Validate() error = nil, want error")
	}
}

func TestValidateAllowsEasyAndHard(t *testing.T) {
	for _, d := range []string{"easy", "hard"} {
		c := &Config{Org: "o", TestName: "t", Model: "m", Difficulty: d}
		if err := c.Validate(); err != nil {
			t.Errorf("Validate() with difficulty %q error = %v", d, err)
		}
	}
}

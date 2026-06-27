package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func prompt(scanner *bufio.Scanner, question, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", question, defaultVal)
	} else {
		fmt.Printf("%s: ", question)
	}
	if !scanner.Scan() {
		return defaultVal
	}
	val := strings.TrimSpace(scanner.Text())
	if val == "" {
		return defaultVal
	}
	return val
}

func main() {
	cfg := &Config{}
	scanner := bufio.NewScanner(os.Stdin)

	cfg.Org = prompt(scanner, "GitHub organisation", "llm-supermarket")
	cfg.TestName = prompt(scanner, "Test name (e.g. cli, markdown-table-formatting)", "")
	cfg.Model = prompt(scanner, "Model", "")
	cfg.Language = prompt(scanner, "Programming language", "")
	cfg.Difficulty = prompt(scanner, "Difficulty", "")

	root, err := findRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defaultClone := filepath.Dir(root)
	clonePath := prompt(scanner, "Clone to (leave empty to skip)", defaultClone)
	if clonePath != "" {
		cfg.CloneDir = clonePath
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n--- Confirmation ---")
	fmt.Printf("Organisation:  %s\n", cfg.Org)
	fmt.Printf("Test name:     %s\n", cfg.TestName)
	fmt.Printf("Model:         %s\n", cfg.Model)
	if cfg.Language != "" {
		fmt.Printf("Language:      %s\n", cfg.Language)
	}
	fmt.Printf("Difficulty:    %s\n", cfg.Difficulty)
	if cfg.CloneDir != "" {
		fmt.Printf("Clone to:      %s\n", cfg.CloneDir)
	} else {
		fmt.Printf("Clone:         no\n")
	}
	fmt.Printf("Repository:    %s\n", cfg.RepoURL())
	fmt.Printf("Prompt URL:    %s\n", cfg.PromptURL())

	confirm := prompt(scanner, "Proceed?", "n")
	if !strings.EqualFold(confirm, "y") && !strings.EqualFold(confirm, "yes") {
		fmt.Println("Aborted.")
		return
	}

	s := &Scaffold{
		Runner: &RealCommandRunner{},
		Stdout: os.Stdout,
	}

	if err := s.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nRepository URL: %s\n", cfg.RepoURL())
	fmt.Printf("Prompt URL:     %s\n", cfg.PromptURL())
}

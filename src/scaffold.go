package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CommandRunner interface {
	Run(dir, name string, args ...string) (string, error)
}

type RealCommandRunner struct{}

func (r *RealCommandRunner) Run(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

type Scaffold struct {
	Runner   CommandRunner
	Stdout   io.Writer
	RepoRoot string // root of the llm-supermarket repo; if empty, auto-detected
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if info, err := os.Stat(filepath.Join(dir, ".git")); err == nil && info.IsDir() {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repository root (no .git directory found)")
		}
		dir = parent
	}
}

func (s *Scaffold) repoRoot() string {
	if s.RepoRoot != "" {
		return s.RepoRoot
	}
	root, err := findRepoRoot()
	if err != nil {
		return "."
	}
	return root
}

func (s *Scaffold) CreateRepo(cfg *Config) (string, error) {
	repoFull := fmt.Sprintf("%s/%s", cfg.Org, cfg.RepoName())

	if cfg.CloneDir != "" {
		cloneAbs, err := filepath.Abs(cfg.CloneDir)
		if err != nil {
			return "", fmt.Errorf("failed to resolve clone path: %w", err)
		}
		if err := os.MkdirAll(cloneAbs, 0755); err != nil {
			return "", fmt.Errorf("failed to create clone directory: %w", err)
		}
		fmt.Fprintf(s.Stdout, "Creating repository %s (cloning to %s)...\n", repoFull, cloneAbs)
		out, err := s.Runner.Run(cloneAbs, "gh", "repo", "create", repoFull, "--public", "--clone")
		if err != nil {
			return "", fmt.Errorf("failed to create repo: %s: %w", strings.TrimSpace(out), err)
		}
		return filepath.Join(cloneAbs, cfg.RepoName()), nil
	}

	fmt.Fprintf(s.Stdout, "Creating repository %s (remote only)...\n", repoFull)
	out, err := s.Runner.Run(".", "gh", "repo", "create", repoFull, "--public")
	if err != nil {
		return "", fmt.Errorf("failed to create repo: %s: %w", strings.TrimSpace(out), err)
	}

	tmpDir, err := os.MkdirTemp("", "llm-supermarket-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	repoDir := filepath.Join(tmpDir, cfg.RepoName())

	out, err = s.Runner.Run(tmpDir, "git", "init", cfg.RepoName())
	if err != nil {
		return "", fmt.Errorf("git init failed: %s: %w", strings.TrimSpace(out), err)
	}

	remoteURL := fmt.Sprintf("https://github.com/%s/%s.git", cfg.Org, cfg.RepoName())
	out, err = s.Runner.Run(repoDir, "git", "remote", "add", "origin", remoteURL)
	if err != nil {
		return "", fmt.Errorf("git remote add failed: %s: %w", strings.TrimSpace(out), err)
	}

	return repoDir, nil
}

func (s *Scaffold) WriteReadme(repoDir string, cfg *Config) error {
	readmePath := filepath.Join(s.repoRoot(), "tests", cfg.TestName, fmt.Sprintf("README_%s.MD", strings.ToUpper(cfg.Difficulty)))
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("failed to read README template: %w", err)
	}

	readme := strings.ReplaceAll(string(content), "{MODEL NAME}", cfg.Model)
	return os.WriteFile(filepath.Join(repoDir, "README.md"), []byte(readme), 0644)
}

func (s *Scaffold) CopyRepoFiles(repoDir string, cfg *Config) error {
	srcDir := filepath.Join(s.repoRoot(), "tests", cfg.TestName, "repo-files")
	_, err := os.Stat(srcDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat repo-files: %w", err)
	}

	return filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		dest := filepath.Join(repoDir, rel)
		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dest, data, 0644)
	})
}

func (s *Scaffold) InitAndPush(repoDir string) error {
	fmt.Fprintln(s.Stdout, "Committing and pushing files...")

	if out, err := s.Runner.Run(repoDir, "git", "add", "."); err != nil {
		return fmt.Errorf("git add failed: %s: %w", strings.TrimSpace(out), err)
	}
	if out, err := s.Runner.Run(repoDir, "git", "branch", "-M", "main"); err != nil {
		return fmt.Errorf("git branch failed: %s: %w", strings.TrimSpace(out), err)
	}
	if out, err := s.Runner.Run(repoDir, "git", "commit", "-m", "Initial commit"); err != nil {
		return fmt.Errorf("git commit failed: %s: %w", strings.TrimSpace(out), err)
	}
	if out, err := s.Runner.Run(repoDir, "git", "push", "-u", "origin", "HEAD"); err != nil {
		return fmt.Errorf("git push failed: %s: %w", strings.TrimSpace(out), err)
	}
	return nil
}

func (s *Scaffold) Run(cfg *Config) error {
	repoDir, err := s.CreateRepo(cfg)
	if err != nil {
		return err
	}
	if err := s.WriteReadme(repoDir, cfg); err != nil {
		return err
	}
	if err := s.CopyRepoFiles(repoDir, cfg); err != nil {
		return err
	}
	return s.InitAndPush(repoDir)
}

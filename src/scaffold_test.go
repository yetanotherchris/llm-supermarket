package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type mockRunner struct {
	commands []struct {
		Dir  string
		Name string
		Args []string
	}
	err error
}

func (m *mockRunner) Run(dir, name string, args ...string) (string, error) {
	m.commands = append(m.commands, struct {
		Dir  string
		Name string
		Args []string
	}{dir, name, args})
	return "", m.err
}

func TestWriteReadmeEasy(t *testing.T) {
	root := t.TempDir()
	repoDir := t.TempDir()

	testDir := filepath.Join(root, "tests", "mycli")
	mkdirAll(t, testDir)
	writeFile(t, filepath.Join(testDir, "README_EASY.MD"), "# {MODEL NAME} Encryptor\n\nReference content.")

	chdir(t, root)

	s := &Scaffold{Runner: &mockRunner{}, Stdout: os.Stdout, RepoRoot: root}
	cfg := &Config{TestName: "mycli", Model: "kimi2.7", Difficulty: "easy"}

	if err := s.WriteReadme(repoDir, cfg); err != nil {
		t.Fatal(err)
	}

	data := readFile(t, filepath.Join(repoDir, "README.md"))
	want := "# kimi2.7 Encryptor\n\nReference content."
	if got := string(data); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteReadmeHard(t *testing.T) {
	root := t.TempDir()
	repoDir := t.TempDir()

	testDir := filepath.Join(root, "tests", "mycli")
	mkdirAll(t, testDir)
	writeFile(t, filepath.Join(testDir, "README_HARD.MD"), "# {MODEL NAME} Encryptor")

	chdir(t, root)

	s := &Scaffold{Runner: &mockRunner{}, Stdout: os.Stdout, RepoRoot: root}
	cfg := &Config{TestName: "mycli", Model: "deepseek4flash", Difficulty: "hard"}

	if err := s.WriteReadme(repoDir, cfg); err != nil {
		t.Fatal(err)
	}

	data := readFile(t, filepath.Join(repoDir, "README.md"))
	want := "# deepseek4flash Encryptor"
	if got := string(data); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteReadmeMissingTemplate(t *testing.T) {
	root := t.TempDir()
	repoDir := t.TempDir()

	chdir(t, root)

	s := &Scaffold{Runner: &mockRunner{}, Stdout: os.Stdout, RepoRoot: root}
	cfg := &Config{TestName: "nonexistent", Model: "m1", Difficulty: "easy"}

	if err := s.WriteReadme(repoDir, cfg); err == nil {
		t.Fatal("expected error for missing README template, got nil")
	}
}

func TestCopyRepoFiles(t *testing.T) {
	root := t.TempDir()
	repoDir := t.TempDir()

	repoFiles := filepath.Join(root, "tests", "mycli", "repo-files")
	mkdirAll(t, repoFiles)
	writeFile(t, filepath.Join(repoFiles, "file1.txt"), "alpha")
	writeFile(t, filepath.Join(repoFiles, "file2.bin"), string([]byte{0, 1, 2}))

	chdir(t, root)

	s := &Scaffold{Runner: &mockRunner{}, Stdout: os.Stdout, RepoRoot: root}
	cfg := &Config{TestName: "mycli"}

	if err := s.CopyRepoFiles(repoDir, cfg); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"file1.txt", "file2.bin"} {
		data, err := os.ReadFile(filepath.Join(repoDir, name))
		if err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
		if len(data) == 0 {
			t.Errorf("empty %s", name)
		}
	}
}

func TestCopyRepoFilesSubdirectories(t *testing.T) {
	root := t.TempDir()
	repoDir := t.TempDir()

	repoFiles := filepath.Join(root, "tests", "mycli", "repo-files")
	mkdirAll(t, filepath.Join(repoFiles, "subdir"))
	writeFile(t, filepath.Join(repoFiles, "root.txt"), "root")
	writeFile(t, filepath.Join(repoFiles, "subdir", "nested.txt"), "nested")

	chdir(t, root)

	s := &Scaffold{Runner: &mockRunner{}, Stdout: os.Stdout, RepoRoot: root}
	cfg := &Config{TestName: "mycli"}

	if err := s.CopyRepoFiles(repoDir, cfg); err != nil {
		t.Fatal(err)
	}

	checkFile(t, repoDir, "root.txt", "root")
	checkFile(t, repoDir, filepath.Join("subdir", "nested.txt"), "nested")
}

func TestCopyRepoFilesMissing(t *testing.T) {
	root := t.TempDir()
	repoDir := t.TempDir()

	chdir(t, root)

	s := &Scaffold{Runner: &mockRunner{}, Stdout: os.Stdout, RepoRoot: root}
	cfg := &Config{TestName: "nonexistent"}

	if err := s.CopyRepoFiles(repoDir, cfg); err != nil {
		t.Fatalf("expected no error for missing repo-files, got: %v", err)
	}
}

func TestInitAndPushCallsGitCommands(t *testing.T) {
	repoDir := t.TempDir()
	m := &mockRunner{}
	s := &Scaffold{Runner: m, Stdout: os.Stdout}

	if err := s.InitAndPush(repoDir); err != nil {
		t.Fatal(err)
	}

	if len(m.commands) < 4 {
		t.Fatalf("expected >=4 commands, got %d", len(m.commands))
	}

	assertCmd(t, m.commands[0], "git", "add", ".")
	assertCmd(t, m.commands[1], "git", "branch", "-M", "main")
	assertCmd(t, m.commands[2], "git", "commit", "-m", "Initial commit")
	assertCmd(t, m.commands[3], "git", "push", "-u", "origin", "HEAD")
}

func TestInitAndPushPropagatesGitAddError(t *testing.T) {
	m := &mockRunner{err: fmt.Errorf("git add failed")}
	s := &Scaffold{Runner: m, Stdout: os.Stdout}

	if err := s.InitAndPush(t.TempDir()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateRepoRemoteOnly(t *testing.T) {
	m := &mockRunner{}
	var buf bytes.Buffer
	s := &Scaffold{Runner: m, Stdout: &buf}
	cfg := &Config{Org: "myorg", TestName: "mycli", Model: "m1"}

	repoDir, err := s.CreateRepo(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(m.commands) != 3 {
		t.Fatalf("expected 3 commands, got %d: %v", len(m.commands), m.commands)
	}

	// gh repo create (no --clone)
	cmd := m.commands[0]
	if cmd.Name != "gh" {
		t.Errorf("expected gh, got %s", cmd.Name)
	}
	if len(cmd.Args) < 3 || cmd.Args[0] != "repo" || cmd.Args[1] != "create" {
		t.Errorf("expected 'repo create', got %v", cmd.Args)
	}
	if cmd.Args[2] != "myorg/mycli-m1" {
		t.Errorf("expected 'myorg/mycli-m1', got %q", cmd.Args[2])
	}
	if cmd.Args[3] != "--public" {
		t.Errorf("expected --public, got %q", cmd.Args[3])
	}
	if len(cmd.Args) > 4 && cmd.Args[4] == "--clone" {
		t.Errorf("unexpected --clone flag for remote-only create")
	}

	// git init
	assertCmd(t, m.commands[1], "git", "init", "mycli-m1")

	// git remote add
	assertCmd(t, m.commands[2], "git", "remote", "add", "origin", "https://github.com/myorg/mycli-m1.git")

	wantRepoDirSuffix := "mycli-m1"
	if !samePathSuffix(repoDir, wantRepoDirSuffix) {
		t.Errorf("repoDir should end with %q, got %q", wantRepoDirSuffix, repoDir)
	}
}

func TestCreateRepoWithClone(t *testing.T) {
	m := &mockRunner{}
	var buf bytes.Buffer
	s := &Scaffold{Runner: m, Stdout: &buf}
	cfg := &Config{Org: "myorg", TestName: "mycli", Model: "m1", CloneDir: t.TempDir()}

	repoDir, err := s.CreateRepo(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(m.commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(m.commands))
	}

	cmd := m.commands[0]
	if cmd.Name != "gh" {
		t.Errorf("expected gh, got %s", cmd.Name)
	}
	if len(cmd.Args) < 3 || cmd.Args[0] != "repo" || cmd.Args[1] != "create" {
		t.Errorf("expected 'repo create', got %v", cmd.Args)
	}
	if cmd.Args[2] != "myorg/mycli-m1" {
		t.Errorf("expected 'myorg/mycli-m1', got %q", cmd.Args[2])
	}
	if cmd.Args[3] != "--public" {
		t.Errorf("expected --public, got %q", cmd.Args[3])
	}
	if cmd.Args[4] != "--clone" {
		t.Errorf("expected --clone, got %q", cmd.Args[4])
	}

	wantRepoDirSuffix := "mycli-m1"
	if !samePathSuffix(repoDir, wantRepoDirSuffix) {
		t.Errorf("repoDir should end with %q, got %q", wantRepoDirSuffix, repoDir)
	}
}

func TestRunFullFlow(t *testing.T) {
	root := t.TempDir()

	testDir := filepath.Join(root, "tests", "mycli")
	mkdirAll(t, testDir)
	writeFile(t, filepath.Join(testDir, "README_EASY.MD"), "# {MODEL NAME}")

	chdir(t, root)

	m := &mockRunner{}
	var buf bytes.Buffer
	s := &Scaffold{Runner: m, Stdout: &buf, RepoRoot: root}
	cfg := &Config{Org: "test-org", TestName: "mycli", Model: "m1", Difficulty: "easy"}

	_ = s.Run(cfg)

	if len(m.commands) < 1 {
		t.Fatal("expected at least one command")
	}
	if m.commands[0].Args[2] != "test-org/mycli-m1" {
		t.Errorf("expected 'test-org/mycli-m1', got %q", m.commands[0].Args[2])
	}
}

func assertCmd(t *testing.T, cmd struct {
	Dir  string
	Name string
	Args []string
}, name string, args ...string) {
	t.Helper()
	if cmd.Name != name {
		t.Errorf("expected command %q, got %q", name, cmd.Name)
	}
	if len(cmd.Args) != len(args) {
		t.Errorf("expected %d args, got %d: %v", len(args), len(cmd.Args), cmd.Args)
		return
	}
	for i := range args {
		if cmd.Args[i] != args[i] {
			t.Errorf("arg[%d]: expected %q, got %q", i, args[i], cmd.Args[i])
		}
	}
}

func checkFile(t *testing.T, dir, name, want string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Errorf("file %s: %v", name, err)
		return
	}
	if got := string(data); got != want {
		t.Errorf("file %s: got %q, want %q", name, got, want)
	}
}

func mkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	prev, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(prev) })
}

func samePathSuffix(path, suffix string) bool {
	path = filepath.ToSlash(path)
	suffix = filepath.ToSlash(suffix)
	return len(path) >= len(suffix) && path[len(path)-len(suffix):] == suffix
}

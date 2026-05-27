package tests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func findRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("pwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			t.Fatalf("could not find repo root")
		}
		wd = parent
	}
}

func buildCLI(t *testing.T) string {
	t.Helper()
	repo := findRepoRoot(t)
	tmp, err := os.MkdirTemp("", "pgi-build-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	bin := filepath.Join(tmp, "pgi")
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags=-s -w", "-o", bin, "./cmd/pgi")
	cmd.Dir = repo
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, string(out))
	}
	return bin
}

func TestHelpOutputShowsUsageAndExamples(t *testing.T) {
	bin := buildCLI(t)
	stdout, stderr, err := runBin(t, bin, t.TempDir(), nil, "--help")
	if err != nil {
		t.Fatalf("help failed: %v, %s", err, stderr)
	}
	if stderr != "" {
		t.Fatalf("expected no stderr, got: %s", stderr)
	}
	if !strings.Contains(stdout, "pgi [--local|--global] [--help] <command> [pattern]") {
		t.Fatalf("usage line missing from help output: %s", stdout)
	}
	if !strings.Contains(stdout, "Examples:") {
		t.Fatalf("examples section missing from help output: %s", stdout)
	}
	if !strings.Contains(stdout, "pgi list \"*.log\"") {
		t.Fatalf("list glob example missing from help output: %s", stdout)
	}
	if !strings.Contains(stdout, "pgi --global add \"*.env\"") {
		t.Fatalf("example missing from help output: %s", stdout)
	}
}

func runBin(t *testing.T, bin string, cwd string, env map[string]string, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func buildEditorHelper(t *testing.T, dir string) string {
	t.Helper()
	source := filepath.Join(dir, "editor-helper.go")
	program := `package main

import "os"

func main() {
	if len(os.Args) != 3 {
		os.Exit(2)
	}
	if os.Args[1] != "--write" {
		os.Exit(3)
	}
	if err := os.WriteFile(os.Args[2], []byte("from-editor\n"), 0o644); err != nil {
		os.Exit(4)
	}
}
`
	if err := os.WriteFile(source, []byte(program), 0o644); err != nil {
		t.Fatalf("write editor source: %v", err)
	}

	bin := filepath.Join(dir, "editor-helper")
	cmd := exec.Command("go", "build", "-o", bin, source)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build editor helper: %v\n%s", err, string(out))
	}
	return bin
}

func TestLocalInlineCRUDAutoInitializesGitInfoExclude(t *testing.T) {
	bin := buildCLI(t)

	tmp, err := os.MkdirTemp("", "repo-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	nested := filepath.Join(tmp, "nested", "dir")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	// init git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmp
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, string(out))
	}

	if _, stderr, err := runBin(t, bin, nested, nil, "add", "*.local"); err != nil {
		t.Fatalf("add failed: %v, %s", err, stderr)
	}

	out, _, err := runBin(t, bin, nested, nil, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "*.local") {
		t.Fatalf("expected pattern in list: %s", out)
	}

	if _, stderr, err := runBin(t, bin, nested, nil, "remove", "*.local"); err != nil {
		t.Fatalf("remove failed: %v, %s", err, stderr)
	}

	out, _, err = runBin(t, bin, nested, nil, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if strings.Contains(out, "*.local") {
		t.Fatalf("pattern still present after remove")
	}

	excludeFile := filepath.Join(tmp, ".git", "info", "exclude")
	if _, err := os.Stat(excludeFile); err != nil {
		t.Fatalf("exclude file missing: %v", err)
	}
}

func TestAddPatternStartingWithDash(t *testing.T) {
	bin := buildCLI(t)

	tmpRepo, err := os.MkdirTemp("", "repo-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	init := exec.Command("git", "init")
	init.Dir = tmpRepo
	if out, err := init.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v, %s", err, string(out))
	}

	if _, stderr, err := runBin(t, bin, tmpRepo, nil, "add", "--cache"); err != nil {
		t.Fatalf("add --cache failed: %v, %s", err, stderr)
	}

	out, _, err := runBin(t, bin, tmpRepo, nil, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "--cache") {
		t.Fatalf("expected --cache in list, got: %s", out)
	}
}

func TestListFiltersPatternsByGlob(t *testing.T) {
	bin := buildCLI(t)

	tmpRepo, err := os.MkdirTemp("", "repo-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	init := exec.Command("git", "init")
	init.Dir = tmpRepo
	if out, err := init.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v, %s", err, string(out))
	}

	patterns := []string{
		"src/pkg/main.go",
		"docs/pkg/README.md",
		"other.txt",
	}
	for _, pattern := range patterns {
		if _, stderr, err := runBin(t, bin, tmpRepo, nil, "add", pattern); err != nil {
			t.Fatalf("add %q failed: %v, %s", pattern, err, stderr)
		}
	}

	out, stderr, err := runBin(t, bin, tmpRepo, nil, "list", "*pkg*")
	if err != nil {
		t.Fatalf("list glob failed: %v, %s", err, stderr)
	}

	trimmed := strings.TrimSpace(out)
	var lines []string
	if trimmed != "" {
		lines = strings.Split(trimmed, "\n")
	}
	want := []string{"src/pkg/main.go", "docs/pkg/README.md"}
	if len(lines) != len(want) {
		t.Fatalf("unexpected list output: got %v, want %v", lines, want)
	}
	for i, line := range lines {
		if line != want[i] {
			t.Fatalf("unexpected list output: got %v, want %v", lines, want)
		}
	}
}

func TestListIgnoresCommentedLines(t *testing.T) {
	bin := buildCLI(t)

	tmpRepo, err := os.MkdirTemp("", "repo-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	init := exec.Command("git", "init")
	init.Dir = tmpRepo
	if out, err := init.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v, %s", err, string(out))
	}

	excludeFile := filepath.Join(tmpRepo, ".git", "info", "exclude")
	content := "# comment one\n*.log\n  # indented comment\nsrc/pkg/main.go\n"
	if err := os.WriteFile(excludeFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write exclude: %v", err)
	}

	out, stderr, err := runBin(t, bin, tmpRepo, nil, "list")
	if err != nil {
		t.Fatalf("list failed: %v, %s", err, stderr)
	}

	if strings.Contains(out, "# comment one") || strings.Contains(out, "# indented comment") {
		t.Fatalf("commented lines should not be listed: %s", out)
	}
	if !strings.Contains(out, "*.log") || !strings.Contains(out, "src/pkg/main.go") {
		t.Fatalf("expected patterns missing from list output: %s", out)
	}
}

func TestGlobalScopeSetsCoreExcludesfile(t *testing.T) {
	bin := buildCLI(t)

	tmpHome, err := os.MkdirTemp("", "home-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	config := filepath.Join(tmpHome, "global.gitconfig")

	env := map[string]string{
		"HOME":              tmpHome,
		"GIT_CONFIG_GLOBAL": config,
	}

	// run global add
	if _, stderr, err := runBin(t, bin, tmpHome, env, "--global", "add", "*.machine"); err != nil {
		t.Fatalf("global add failed: %v, %s", err, stderr)
	}

	expectedIgnore := filepath.Join(tmpHome, ".gitignore_global")
	data, err := os.ReadFile(expectedIgnore)
	if err != nil {
		t.Fatalf("read ignore: %v", err)
	}
	if !strings.Contains(string(data), "*.machine") {
		t.Fatalf("pattern not found in global ignore: %s", string(data))
	}

	// verify git config --global --get core.excludesfile
	cmd := exec.Command("git", "config", "--global", "--get", "core.excludesfile")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "HOME="+tmpHome, "GIT_CONFIG_GLOBAL="+config)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git config get failed: %v, %s", err, string(out))
	}
	if strings.TrimSpace(string(out)) != expectedIgnore {
		t.Fatalf("expected core.excludesfile %s, got %s", expectedIgnore, string(out))
	}
}

func TestGlobalScopeRespectsHomeRelativeExcludesfile(t *testing.T) {
	bin := buildCLI(t)

	tmpHome, err := os.MkdirTemp("", "home-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	config := filepath.Join(tmpHome, "global.gitconfig")

	env := map[string]string{
		"HOME":              tmpHome,
		"GIT_CONFIG_GLOBAL": config,
	}

	// set core.excludesfile to relative path
	cmd := exec.Command("git", "config", "--global", "core.excludesfile", "relative.ignore")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "HOME="+tmpHome, "GIT_CONFIG_GLOBAL="+config)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git config set failed: %v, %s", err, string(out))
	}

	// init a repo and run --global add
	tmpRepo, err := os.MkdirTemp("", "repo-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	init := exec.Command("git", "init")
	init.Dir = tmpRepo
	if out, err := init.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v, %s", err, string(out))
	}

	if _, stderr, err := runBin(t, bin, tmpRepo, env, "--global", "add", "*.machine"); err != nil {
		t.Fatalf("global add failed: %v, %s", err, stderr)
	}

	expectedIgnore := filepath.Join(tmpHome, "relative.ignore")
	data, err := os.ReadFile(expectedIgnore)
	if err != nil {
		t.Fatalf("read ignore: %v", err)
	}
	if !strings.Contains(string(data), "*.machine") {
		t.Fatalf("pattern not found in relative ignore: %s", string(data))
	}
}

func TestEditModeModifiesSameUnderlyingFile(t *testing.T) {
	bin := buildCLI(t)

	tmpRepo, err := os.MkdirTemp("", "repo-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	init := exec.Command("git", "init")
	init.Dir = tmpRepo
	if out, err := init.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v, %s", err, string(out))
	}

	editorDir := filepath.Join(tmpRepo, "editor tools")
	if err := os.MkdirAll(editorDir, 0o755); err != nil {
		t.Fatalf("mkdir editor dir: %v", err)
	}
	editor := buildEditorHelper(t, editorDir)

	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("lookpath git: %v", err)
	}
	gitDir := filepath.Join(tmpRepo, "git-bin")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("mkdir git dir: %v", err)
	}
	if err := os.Symlink(realGit, filepath.Join(gitDir, "git")); err != nil {
		t.Fatalf("symlink git: %v", err)
	}

	env := map[string]string{
		"PATH":   gitDir,
		"EDITOR": fmt.Sprintf("\"%s\" --write", editor),
	}
	if _, stderr, err := runBin(t, bin, tmpRepo, env, "edit"); err != nil {
		t.Fatalf("edit failed: %v, %s", err, stderr)
	}

	out, _, err := runBin(t, bin, tmpRepo, map[string]string{"PATH": gitDir}, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "from-editor") {
		t.Fatalf("expected from-editor in list, got: %s", out)
	}
}

func TestBinDirFlagIsRejected(t *testing.T) {
	bin := buildCLI(t)
	tmpRepo, err := os.MkdirTemp("", "repo-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	init := exec.Command("git", "init")
	init.Dir = tmpRepo
	if out, err := init.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v, %s", err, string(out))
	}

	_, stderr, err := runBin(t, bin, tmpRepo, nil, "--bin-dir", tmpRepo, "list")
	if err == nil {
		t.Fatalf("expected error for unknown arg, got none; stderr=%s", stderr)
	}
	if !strings.Contains(stderr, "unknown argument: --bin-dir") {
		t.Fatalf("unexpected stderr: %s", stderr)
	}
}

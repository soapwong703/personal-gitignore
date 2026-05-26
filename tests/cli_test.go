package tests

import (
	"bytes"
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
	// also provide compatibility alias
	alias := filepath.Join(tmp, "personal-gitignore")
	data, err := os.ReadFile(bin)
	if err != nil {
		t.Fatalf("read built binary: %v", err)
	}
	if err := os.WriteFile(alias, data, 0o755); err != nil {
		t.Fatalf("write alias: %v", err)
	}
	return bin
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

func TestLocalInlineCRUDUsesGitInfoExclude(t *testing.T) {
	bin := buildCLI(t)

	tmp, err := os.MkdirTemp("", "repo-")
	if err != nil {
		t.Fatalf("mktemp: %v", err)
	}
	// init git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmp
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, string(out))
	}

	// run setup
	if _, stderr, err := runBin(t, bin, tmp, nil, "setup"); err != nil {
		t.Fatalf("setup failed: %v, %s", err, stderr)
	}

	if _, stderr, err := runBin(t, bin, tmp, nil, "add", "*.local"); err != nil {
		t.Fatalf("add failed: %v, %s", err, stderr)
	}

	out, _, err := runBin(t, bin, tmp, nil, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "*.local") {
		t.Fatalf("expected pattern in list: %s", out)
	}

	if _, stderr, err := runBin(t, bin, tmp, nil, "remove", "*.local"); err != nil {
		t.Fatalf("remove failed: %v, %s", err, stderr)
	}

	out, _, err = runBin(t, bin, tmp, nil, "list")
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

	// create an editor helper script that writes a line
	editor := filepath.Join(tmpRepo, "editor.sh")
	if err := os.WriteFile(editor, []byte("#!/usr/bin/env sh\necho from-editor > \"$1\"\n"), 0o755); err != nil {
		t.Fatalf("write editor: %v", err)
	}

	env := map[string]string{"EDITOR": editor}
	if _, stderr, err := runBin(t, bin, tmpRepo, env, "edit"); err != nil {
		t.Fatalf("edit failed: %v, %s", err, stderr)
	}

	out, _, err := runBin(t, bin, tmpRepo, nil, "list")
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

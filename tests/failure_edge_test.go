package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestAddAcceptsMultipleArgs verifies add can ingest multiple positional patterns in one command.
func TestAddAcceptsMultipleArgs(t *testing.T) {
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

	if _, stderr, err := runBin(t, bin, tmpRepo, nil, "add", "one", "two", "three"); err != nil {
		t.Fatalf("add with multiple args failed: %v, %s", err, stderr)
	}

	out, _, err := runBin(t, bin, tmpRepo, nil, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	for _, want := range []string{"one", "two", "three"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in list output: %s", want, out)
		}
	}
}

// TestAddAndRemoveSplitWhitespace verifies whitespace inside args is split into separate pattern entries.
func TestAddAndRemoveSplitWhitespace(t *testing.T) {
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

	if _, stderr, err := runBin(t, bin, tmpRepo, nil, "add", "alpha beta", "gamma\tdelta"); err != nil {
		t.Fatalf("add split args failed: %v, %s", err, stderr)
	}

	out, _, err := runBin(t, bin, tmpRepo, nil, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	for _, want := range []string{"alpha", "beta", "gamma", "delta"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in list output: %s", want, out)
		}
	}

	if _, stderr, err := runBin(t, bin, tmpRepo, nil, "remove", "alpha beta", "gamma", "delta"); err != nil {
		t.Fatalf("remove split args failed: %v, %s", err, stderr)
	}

	out, _, err = runBin(t, bin, tmpRepo, nil, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	for _, removed := range []string{"alpha", "beta", "gamma", "delta"} {
		if strings.Contains(out, removed) {
			t.Fatalf("did not expect %q after remove, got: %s", removed, out)
		}
	}
}

// TestListInvalidGlob ensures an invalid glob returns an error instead of panicking.
func TestListInvalidGlob(t *testing.T) {
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

	// Add a sample pattern
	if _, stderr, err := runBin(t, bin, tmpRepo, nil, "add", "foo"); err != nil {
		t.Fatalf("add failed: %v, %s", err, stderr)
	}

	_, stderr, err := runBin(t, bin, tmpRepo, nil, "list", "[")
	if err == nil {
		t.Fatalf("expected error for invalid glob; stderr=%s", stderr)
	}
	if !strings.Contains(stderr, "syntax error") {
		t.Fatalf("expected syntax error in stderr, got: %s", stderr)
	}
}

// TestAddVeryLongPattern verifies the CLI handles very long patterns without failing.
func TestAddVeryLongPattern(t *testing.T) {
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

	long := strings.Repeat("a", 10000)
	if _, stderr, err := runBin(t, bin, tmpRepo, nil, "add", long); err != nil {
		t.Fatalf("add long pattern failed: %v, %s", err, stderr)
	}

	out, _, err := runBin(t, bin, tmpRepo, nil, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, long[:100]) {
		t.Fatalf("long pattern missing from list output: %s", out)
	}
}

// TestManyPatterns ensures reasonable performance and correctness when adding many patterns.
func TestManyPatterns(t *testing.T) {
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

	const n = 260
	for i := 0; i < n; i++ {
		pattern := filepath.Join("dir", "file") + string('a'+rune(i%26)) + "-" + strconv.Itoa((i/26)%10)
		if _, stderr, err := runBin(t, bin, tmpRepo, nil, "add", pattern); err != nil {
			t.Fatalf("add %d failed: %v, %s", i, err, stderr)
		}
	}

	out, _, err := runBin(t, bin, tmpRepo, nil, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	trimmed := strings.TrimSpace(out)
	var lines []string
	if trimmed != "" {
		lines = strings.Split(trimmed, "\n")
	}
	if len(lines) != n {
		t.Fatalf("unexpected number of patterns: got %d, want %d", len(lines), n)
	}
}

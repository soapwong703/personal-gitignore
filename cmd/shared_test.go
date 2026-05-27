package cmd

import (
	"errors"
	"testing"
)

func TestChooseEditorCommandPrefersEnvOrder(t *testing.T) {
	env := []string{
		"VISUAL=visual-editor",
		"EDITOR=primary-editor",
	}

	editor, source, err := chooseEditorCommand(env, func() (string, error) {
		return "git-editor", nil
	})
	if err != nil {
		t.Fatalf("chooseEditorCommand returned error: %v", err)
	}
	if editor != "primary-editor" {
		t.Fatalf("expected EDITOR value, got %q", editor)
	}
	if source != "EDITOR" {
		t.Fatalf("expected source EDITOR, got %q", source)
	}
}

func TestChooseEditorCommandFallsBackToGitEditor(t *testing.T) {
	env := []string{}
	editor, source, err := chooseEditorCommand(env, func() (string, error) {
		return "vi", nil
	})
	if err != nil {
		t.Fatalf("chooseEditorCommand returned error: %v", err)
	}
	if editor != "vi" {
		t.Fatalf("expected git editor vi, got %q", editor)
	}
	if source != "GIT_EDITOR" {
		t.Fatalf("expected source GIT_EDITOR, got %q", source)
	}
}

func TestResolveEditorPartsWindowsFallbackToNotepad(t *testing.T) {
	parts, err := resolveEditorParts("vi", "GIT_EDITOR", "windows", func(name string) (string, error) {
		switch name {
		case "vi":
			return "", errors.New("not found")
		case "notepad":
			return "C:/Windows/System32/notepad.exe", nil
		default:
			return "", errors.New("unexpected")
		}
	})
	if err != nil {
		t.Fatalf("resolveEditorParts returned error: %v", err)
	}
	if len(parts) != 1 || parts[0] != "notepad" {
		t.Fatalf("expected notepad fallback, got %v", parts)
	}
}

func TestResolveEditorPartsNoFallbackForEditorEnv(t *testing.T) {
	_, err := resolveEditorParts("vi", "EDITOR", "windows", func(name string) (string, error) {
		return "", errors.New("not found")
	})
	if err == nil {
		t.Fatalf("expected error when EDITOR command is missing")
	}
}

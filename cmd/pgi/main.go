package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

type options struct {
	local       bool
	globalScope bool
	help        bool
	command     string
	pattern     string
}

var validCommands = map[string]struct{}{
	"list":   {},
	"add":    {},
	"remove": {},
	"clear":  {},
	"edit":   {},
}

const helpText = `Usage:
	pgi [--local|--global] [--help] <command> [pattern]

Commands:
	list [glob]       Show the current ignore patterns, filtered by glob when provided
	add <pattern>     Add a pattern if it is not already present
	remove <pattern>  Remove a pattern if it exists
	clear             Remove all patterns
	edit              Open the ignore file in your editor

Examples:
	pgi --help
	pgi list
	pgi list "*.log"
	pgi add "*.log"
	pgi --global add "*.env"
	pgi edit

The default scope is local. Use --global to manage the global ignore file.
`

func runGit(args []string, cwd string, env []string) (string, error) {
	cmd := exec.Command("git", args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	if env != nil {
		cmd.Env = env
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = "git command failed"
		}
		return "", errors.New(msg)
	}
	return strings.TrimSpace(string(out)), nil
}

func resolveLocalIgnoreFile(cwd string, env []string) (string, error) {
	gitPath, err := runGit([]string{"rev-parse", "--git-path", "info/exclude"}, cwd, env)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(gitPath) {
		return gitPath, nil
	}
	return filepath.Join(cwd, gitPath), nil
}

func resolveGlobalIgnoreFile(env []string) (string, error) {
	cmd := exec.Command("git", "config", "--global", "--get", "core.excludesfile")
	cmd.Env = env
	out, err := cmd.Output()
	if err == nil {
		configured := strings.TrimSpace(string(out))
		if configured != "" {
			return expandHome(configured)
		}
	}

	defaultPath, err := expandHome("~/.gitignore_global")
	if err != nil {
		return "", err
	}
	if _, err := runGit([]string{"config", "--global", "core.excludesfile", defaultPath}, "", env); err != nil {
		return "", err
	}
	return defaultPath, nil
}

func ensureFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	return f.Close()
}

func readPatterns(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	patterns := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		patterns = append(patterns, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return patterns, nil
}

func writePatterns(path string, patterns []string) error {
	content := strings.Join(patterns, "\n")
	if content != "" {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

const globSeparatorPlaceholder = "\x00"

func isCommentLine(line string) bool {
	trimmed := strings.TrimLeft(line, " \t")
	return strings.HasPrefix(trimmed, "#")
}

func filterPatternsByGlob(patterns []string, glob string) ([]string, error) {
	compiledGlob := strings.ReplaceAll(glob, "/", globSeparatorPlaceholder)
	if _, err := path.Match(compiledGlob, ""); err != nil {
		return nil, err
	}

	filtered := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		row := strings.ReplaceAll(pattern, "/", globSeparatorPlaceholder)
		matched, err := path.Match(compiledGlob, row)
		if err != nil {
			return nil, err
		}
		if matched {
			filtered = append(filtered, pattern)
		}
	}
	return filtered, nil
}

func splitCommand(command string) ([]string, error) {
	parts := []string{}
	var current strings.Builder
	inSingle := false
	inDouble := false
	runes := []rune(command)

	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		if inSingle {
			if ch == '\'' {
				inSingle = false
				continue
			}
			current.WriteRune(ch)
			continue
		}

		if inDouble {
			switch ch {
			case '"':
				inDouble = false
			case '\\':
				if i+1 < len(runes) {
					next := runes[i+1]
					if next == '"' || next == '\\' || next == '$' || next == '`' || next == '\n' {
						current.WriteRune(next)
						i++
						continue
					}
				}
				current.WriteRune(ch)
			default:
				current.WriteRune(ch)
			}
			continue
		}

		if unicode.IsSpace(ch) {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			continue
		}

		switch ch {
		case '\'':
			inSingle = true
		case '"':
			inDouble = true
		case '\\':
			if i+1 < len(runes) {
				next := runes[i+1]
				if unicode.IsSpace(next) || next == '\\' || next == '\'' || next == '"' {
					current.WriteRune(next)
					i++
					continue
				}
			}
			current.WriteRune(ch)
		default:
			current.WriteRune(ch)
		}
	}

	if inSingle || inDouble {
		return nil, errors.New("unterminated quoted editor command")
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	if len(parts) == 0 {
		return nil, errors.New("empty editor command")
	}

	return parts, nil
}

func openEditor(path string, env []string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		var err error
		editor, err = runGit([]string{"var", "GIT_EDITOR"}, "", env)
		if err != nil {
			editor = ""
		}
	}
	if editor == "" {
		return errors.New("No editor found. Set EDITOR or VISUAL.")
	}
	editor = strings.TrimSpace(editor)
	parts, err := splitCommand(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(parts[0], append(parts[1:], path)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("Editor exited with status %d", exitErr.ExitCode())
		}
		return err
	}
	return nil
}

func parseArgs(args []string) (options, error) {
	opts := options{}
	positionals := []string{}
	commandSeen := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !commandSeen {
			switch arg {
			case "--local":
				opts.local = true
			case "--global":
				opts.globalScope = true
			case "--help":
				opts.help = true
			default:
				if strings.HasPrefix(arg, "--") {
					return options{}, fmt.Errorf("unknown argument: %s", arg)
				}
				commandSeen = true
				positionals = append(positionals, arg)
			}
			continue
		}
		positionals = append(positionals, arg)
	}

	if opts.help {
		return opts, nil
	}
	if opts.local && opts.globalScope {
		return options{}, errors.New("--local and --global cannot be used together")
	}
	if len(positionals) == 0 {
		return options{}, errors.New("missing command")
	}
	opts.command = positionals[0]
	if _, ok := validCommands[opts.command]; !ok {
		return options{}, fmt.Errorf("invalid command: %s", opts.command)
	}
	if len(positionals) > 1 {
		opts.pattern = positionals[1]
	}
	if len(positionals) > 2 {
		return options{}, errors.New("too many arguments")
	}
	return opts, nil
}

func expandHome(path string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			return home, nil
		}
		return filepath.Join(home, path[2:]), nil
	}
	expanded := os.ExpandEnv(path)
	if filepath.IsAbs(expanded) {
		return expanded, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, expanded), nil
}

func main() {
	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	if opts.help {
		fmt.Print(helpText)
		return
	}

	env := os.Environ()
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	var ignoreFile string
	if opts.globalScope {
		ignoreFile, err = resolveGlobalIgnoreFile(env)
	} else {
		ignoreFile, err = resolveLocalIgnoreFile(cwd, env)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if err := ensureFile(ignoreFile); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	switch opts.command {
	case "edit":
		if err := openEditor(ignoreFile, env); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		return
	case "list":
		patterns, err := readPatterns(ignoreFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		visible := make([]string, 0, len(patterns))
		for _, pattern := range patterns {
			if isCommentLine(pattern) {
				continue
			}
			visible = append(visible, pattern)
		}
		if opts.pattern != "" {
			visible, err = filterPatternsByGlob(visible, opts.pattern)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
		}
		for _, p := range visible {
			fmt.Println(p)
		}
		return
	case "add", "remove":
		if opts.pattern == "" {
			fmt.Fprintf(os.Stderr, "Error: '%s' requires a pattern\n", opts.command)
			os.Exit(1)
		}
	}

	patterns, err := readPatterns(ignoreFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	switch opts.command {
	case "add":
		for _, p := range patterns {
			if p == opts.pattern {
				return
			}
		}
		patterns = append(patterns, opts.pattern)
		if err := writePatterns(ignoreFile, patterns); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "remove":
		updated := make([]string, 0, len(patterns))
		for _, p := range patterns {
			if p != opts.pattern {
				updated = append(updated, p)
			}
		}
		if err := writePatterns(ignoreFile, updated); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "clear":
		if err := writePatterns(ignoreFile, []string{}); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	}
}

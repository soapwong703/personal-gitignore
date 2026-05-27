# Agent Instructions

- Use Go 1.24 from [go.mod](go.mod). The CLI behavior lives in [main.go](main.go).
- Treat [main.go](main.go) as the single CLI entrypoint. Run `go run .` from the repository root.
- Local scope writes to `.git/info/exclude`. Global scope uses `git config --global core.excludesfile` and falls back to `~/.gitignore_global`.
- Edit mode must honor `EDITOR`, then `VISUAL`, then `git var GIT_EDITOR`. Inline mode and editor mode should modify the same underlying ignore file.
- Keep README.md focused on project description, installation, and usage; do not add release-process or other development workflow docs there.
- Validate CLI changes with `go test ./...`. New behavior should be covered with tests in [tests/cli_test.go](tests/cli_test.go) or the closest equivalent integration test.
- For build or release changes, use `go build -trimpath -ldflags="-s -w" -o pgi .` as the baseline and keep [install.sh](install.sh), [README.md](README.md), and [.github/workflows/release.yml](.github/workflows/release.yml) aligned.
- Keep parser and command compatibility stable. Unknown flags should fail, commands should stay explicit, and existing stderr/stdout text may be asserted by tests.
- When testing global-scope behavior, isolate `HOME` and `GIT_CONFIG_GLOBAL` so the real user git config is not touched.
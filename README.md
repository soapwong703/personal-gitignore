# personal-gitignore

A command line tool to manage personal gitignore rules that never get committed.

## What it manages

- **Local (default):** `.git/info/exclude` in the current repository.
- **Global:** the file configured by `git config --global core.excludesfile` (defaults to `~/.gitignore_global`).

Both inline mode and editor mode operate on the same underlying file for the selected scope.

## Usage

```bash
./pgi [--local|--global] <command> [pattern]
```

`pgi` is the default command. `personal-gitignore` is still available for compatibility.
Running `./pgi` directly from this repository requires Go (uses `go run`).

Install directly from GitHub in one line:

```bash
curl -fsSL https://raw.githubusercontent.com/soapwong703/personal-gitignore/main/install.sh | sh
export PATH="$HOME/.local/bin:$PATH"
```

> The one-line installer builds a native binary, so Go is required during installation.

Custom install directory:

```bash
curl -fsSL https://raw.githubusercontent.com/soapwong703/personal-gitignore/main/install.sh | sh -s -- --bin-dir "$HOME/bin"
```

Commands:

- `setup` – ensure the file exists and print its path
- `install` – install `pgi` (default command) and `personal-gitignore` compatibility binary to `~/.local/bin` (or `--bin-dir`)
- `list` – show current rules
- `add <pattern>` – add a rule
- `remove <pattern>` – remove a rule
- `clear` – remove all rules
- `edit` – open the selected file in your editor (`$EDITOR`/`$VISUAL`)

Examples:

```bash
./pgi install
export PATH="$HOME/.local/bin:$PATH"

./pgi add .env.local
./pgi remove .env.local
./pgi --global add *.machine
./pgi edit
```

## GitHub Release Build

A GitHub Actions workflow builds and bundles CLI artifacts for Linux, macOS, and Windows, and publishes them to GitHub Releases when a tag matching `v*` is pushed.

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release contains bundled archives with both `pgi` and `personal-gitignore` binaries for each target platform.


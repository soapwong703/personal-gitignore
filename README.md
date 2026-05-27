# personal-gitignore

A command line tool to manage personal gitignore rules that never get committed.

## What it manages

- **Local (default):** `.git/info/exclude` in the current repository.
- **Global:** the file configured by `git config --global core.excludesfile` (defaults to `~/.gitignore_global`).

Both inline mode and editor mode operate on the same underlying file for the selected scope.

## Install

macOS and Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/soapwong703/personal-gitignore/main/install.sh | sh
```

Windows:

```powershell
irm https://raw.githubusercontent.com/soapwong703/personal-gitignore/main/install.ps1 | iex
```

## Update

Re-run the same installer for your platform. It downloads the latest release and replaces the existing `pgi` binary.

macOS and Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/soapwong703/personal-gitignore/main/install.sh | sh
```

Windows:

```powershell
irm https://raw.githubusercontent.com/soapwong703/personal-gitignore/main/install.ps1 | iex
```

## Usage

```bash
pgi [--local|--global] [--help] <command> [pattern]
```

Commands:

- `list [glob]` - show the current ignore patterns, skipping commented lines and filtered by glob when provided
- `add <pattern>` - add a pattern if it is not already present
- `remove <pattern>` - remove a pattern if it exists
- `clear` - remove all patterns
- `edit` - open the ignore file in your editor

The default scope is local. Use `--global` to manage the global ignore file.

The target file is created automatically the first time you run a command.

Examples:

```bash
pgi --help
pgi list
pgi list "*.log"
pgi add "*.log"
pgi --global add "*.env"
pgi edit
```

## Uninstall

This removes the installed command only. Your ignore files remain untouched.

macOS and Linux:

```sh
rm -f "$HOME/.local/bin/pgi"
```

Windows PowerShell:

```powershell
Remove-Item -Force "$HOME\.local\bin\pgi.exe"
```

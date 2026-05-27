# personal-gitignore

Manage personal gitignore rules from the command line without committing them to your repos.

## What it manages

- **Local, default:** `.git/info/exclude` in the current repository.
- **Global:** the file configured by `git config --global core.excludesfile`, which defaults to `~/.gitignore_global`.

All commands target the selected scope, whether you manage it from the command line or open it in your editor.

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

Re-run the installer for your platform. It downloads the latest release and replaces the installed `pgi` binary.

## Usage

```bash
pgi [--global] [--help] <command> [pattern ...]
```

Commands:

- `list [glob]` - show the current ignore patterns, filtered by glob when provided
- `add <pattern...>` - add one or more patterns if they are not already present
- `remove <pattern...>` - remove one or more patterns if they exist
- `clear` - remove all non-comment patterns
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
pgi remove -- --prefix
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

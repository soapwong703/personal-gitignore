# personal-gitignore

A command line tool to manage personal gitignore rules that never get committed.

## What it manages

- **Local (default):** `.git/info/exclude` in the current repository.
- **Global:** the file configured by `git config --global core.excludesfile` (defaults to `~/.gitignore_global`).

Both inline mode and editor mode operate on the same underlying file for the selected scope.

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/soapwong703/personal-gitignore/main/install.sh | sh
```

## Usage

```bash
pgi [--local|--global] <command> [pattern]
```

Commands:

- `list` – show current rules
- `add <pattern>` – add a rule
- `remove <pattern>` – remove a rule
- `clear` – remove all rules
- `edit` – open the selected file in your editor (`$EDITOR`/`$VISUAL`)

The target file is created automatically the first time you run a command.

Examples:

```bash
pgi add .env.local
pgi remove .env.local
pgi --global add '*.machine'
pgi edit
```

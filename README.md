# personal-gitignore

A command line tool to manage personal gitignore rules that never get committed.

## What it manages

- **Local (default):** `.git/info/exclude` in the current repository.
- **Global:** the file configured by `git config --global core.excludesfile` (defaults to `~/.gitignore_global`).

Both inline mode and editor mode operate on the same underlying file for the selected scope.

## Usage

```bash
./personal-gitignore [--local|--global] <command> [pattern]
```

Commands:

- `setup` – ensure the file exists and print its path
- `list` – show current rules
- `add <pattern>` – add a rule
- `remove <pattern>` – remove a rule
- `clear` – remove all rules
- `edit` – open the selected file in your editor (`$EDITOR`/`$VISUAL`)

Examples:

```bash
./personal-gitignore add .env.local
./personal-gitignore remove .env.local
./personal-gitignore --global add *.machine
./personal-gitignore edit
```

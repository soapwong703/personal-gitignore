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

`pgi` is the default command. `personal-gitignore` remains available as a compatibility alias.

Commands:

- `setup` – ensure the file exists and print its path
- `list` – show current rules
- `add <pattern>` – add a rule
- `remove <pattern>` – remove a rule
- `clear` – remove all rules
- `edit` – open the selected file in your editor (`$EDITOR`/`$VISUAL`)

Examples:

```bash
pgi add .env.local
pgi remove .env.local
pgi --global add '*.machine'
pgi edit
```

## GitHub Release Build

A GitHub Actions workflow builds the release bundles for Linux, macOS, and Windows, then publishes them to GitHub Releases when a tag matching `v*` is pushed.

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release contains bundled archives with both executables for each target platform.

# personal-gitignore

A command line tool to manage personal gitignore rules that never get committed.

## What it manages

- **Local (default):** `.git/info/exclude` in the current repository.
- **Global:** the file configured by `git config --global core.excludesfile` (defaults to `~/.gitignore_global`).

Both inline mode and editor mode operate on the same underlying file for the selected scope.

## Install (bundled release artifacts)

The intended distribution is via GitHub Release bundles produced by CI.

1. Open the latest release: https://github.com/soapwong703/personal-gitignore/releases/latest
2. Download the archive for your platform:
   - Linux/macOS: `personal-gitignore_<os>_<arch>.tar.gz`
   - Windows: `personal-gitignore_windows_amd64.zip`
3. Extract the archive and move `pgi` to a directory in your `PATH` (for example `~/.local/bin`).
4. (Optional) Keep `personal-gitignore` as a compatibility alias.

Example (Linux/macOS):

```bash
tar -xzf personal-gitignore_linux_amd64.tar.gz
install -m 0755 personal-gitignore_linux_amd64/pgi ~/.local/bin/pgi
install -m 0755 personal-gitignore_linux_amd64/personal-gitignore ~/.local/bin/personal-gitignore
export PATH="$HOME/.local/bin:$PATH"
```

## Usage

```bash
pgi [--local|--global] <command> [pattern]
```

`pgi` is the default command. `personal-gitignore` is available as a compatibility alias.

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

A GitHub Actions workflow builds and bundles CLI artifacts for Linux, macOS, and Windows, and publishes them to GitHub Releases when a tag matching `v*` is pushed.

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release contains bundled archives with both `pgi` and `personal-gitignore` binaries for each target platform.

## Alternative: build from source installer

If you want to build locally instead of using bundled release artifacts:

```bash
curl -fsSL https://raw.githubusercontent.com/soapwong703/personal-gitignore/main/install.sh | sh
```

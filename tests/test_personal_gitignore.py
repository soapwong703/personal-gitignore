import os
import subprocess
import tempfile
import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[1]
CLI = REPO_ROOT / "personal-gitignore"
ALIAS = REPO_ROOT / "pgi"
INSTALLER = REPO_ROOT / "install.sh"


class PersonalGitignoreCliTests(unittest.TestCase):
    def run_cli(self, args, cwd, env=None):
        process = subprocess.run(
            [str(CLI), *args],
            cwd=cwd,
            env=env,
            capture_output=True,
            text=True,
            check=False,
        )
        return process

    def test_local_inline_crud_uses_git_info_exclude(self):
        with tempfile.TemporaryDirectory() as tmp:
            repo = Path(tmp) / "repo"
            repo.mkdir()
            subprocess.run(["git", "init"], cwd=repo, check=True, capture_output=True)

            add = self.run_cli(["setup"], cwd=repo)
            self.assertEqual(add.returncode, 0, add.stderr)

            add = self.run_cli(["add", "*.local"], cwd=repo)
            self.assertEqual(add.returncode, 0, add.stderr)

            listed = self.run_cli(["list"], cwd=repo)
            self.assertIn("*.local", listed.stdout)

            remove = self.run_cli(["remove", "*.local"], cwd=repo)
            self.assertEqual(remove.returncode, 0, remove.stderr)

            listed = self.run_cli(["list"], cwd=repo)
            self.assertNotIn("*.local", listed.stdout)

            exclude_file = repo / ".git" / "info" / "exclude"
            self.assertTrue(exclude_file.exists())

    def test_global_scope_sets_core_excludesfile(self):
        with tempfile.TemporaryDirectory() as tmp:
            home = Path(tmp) / "home"
            home.mkdir()
            config = Path(tmp) / "global.gitconfig"

            env = os.environ.copy()
            env["HOME"] = str(home)
            env["GIT_CONFIG_GLOBAL"] = str(config)

            result = self.run_cli(["--global", "add", "*.machine"], cwd=home, env=env)
            self.assertEqual(result.returncode, 0, result.stderr)

            expected_ignore = home / ".gitignore_global"
            self.assertTrue(expected_ignore.exists())
            self.assertIn("*.machine", expected_ignore.read_text(encoding="utf-8"))

            configured = subprocess.run(
                ["git", "config", "--global", "--get", "core.excludesfile"],
                env=env,
                check=True,
                capture_output=True,
                text=True,
            )
            self.assertEqual(configured.stdout.strip(), str(expected_ignore))

    def test_edit_mode_modifies_same_underlying_file(self):
        with tempfile.TemporaryDirectory() as tmp:
            repo = Path(tmp) / "repo"
            repo.mkdir()
            subprocess.run(["git", "init"], cwd=repo, check=True, capture_output=True)

            editor_helper = Path(tmp) / "append_editor.py"
            editor_helper.write_text(
                "import pathlib, sys\npathlib.Path(sys.argv[1]).write_text('from-editor\\n', encoding='utf-8')\n",
                encoding="utf-8",
            )

            env = os.environ.copy()
            env["EDITOR"] = f"python {editor_helper}"

            result = self.run_cli(["edit"], cwd=repo, env=env)
            self.assertEqual(result.returncode, 0, result.stderr)

            listed = self.run_cli(["list"], cwd=repo)
            self.assertIn("from-editor", listed.stdout)

    def test_pgi_alias_invokes_same_cli(self):
        with tempfile.TemporaryDirectory() as tmp:
            repo = Path(tmp) / "repo"
            repo.mkdir()
            subprocess.run(["git", "init"], cwd=repo, check=True, capture_output=True)

            add = subprocess.run(
                [str(ALIAS), "add", "*.alias"],
                cwd=repo,
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(add.returncode, 0, add.stderr)

            listed = self.run_cli(["list"], cwd=repo)
            self.assertIn("*.alias", listed.stdout)

    def test_install_places_binaries_in_user_bin(self):
        with tempfile.TemporaryDirectory() as tmp:
            home = Path(tmp) / "home"
            home.mkdir()
            repo = Path(tmp) / "repo"
            repo.mkdir()
            subprocess.run(["git", "init"], cwd=repo, check=True, capture_output=True)

            env = os.environ.copy()
            env["HOME"] = str(home)

            install = self.run_cli(["install"], cwd=home, env=env)
            self.assertEqual(install.returncode, 0, install.stderr)

            installed_cli = home / ".local" / "bin" / "personal-gitignore"
            installed_alias = home / ".local" / "bin" / "pgi"
            self.assertTrue(installed_cli.exists())
            self.assertTrue(installed_alias.exists())

            add = subprocess.run(
                [str(installed_alias), "add", "*.from-install"],
                cwd=repo,
                env=env,
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(add.returncode, 0, add.stderr)

            listed = subprocess.run(
                [str(installed_cli), "list"],
                cwd=repo,
                env=env,
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(listed.returncode, 0, listed.stderr)
            self.assertIn("*.from-install", listed.stdout)

    def test_bin_dir_flag_rejected_for_non_install_commands(self):
        with tempfile.TemporaryDirectory() as tmp:
            repo = Path(tmp) / "repo"
            repo.mkdir()
            subprocess.run(["git", "init"], cwd=repo, check=True, capture_output=True)

            result = self.run_cli(["--bin-dir", str(repo), "list"], cwd=repo)
            self.assertNotEqual(result.returncode, 0)
            self.assertIn("--bin-dir can only be used with the install command", result.stderr)

    def test_install_script_supports_one_line_url_style_install(self):
        with tempfile.TemporaryDirectory() as tmp:
            source = Path(tmp) / "source"
            source.mkdir()
            (source / "personal-gitignore").write_text(CLI.read_text(encoding="utf-8"), encoding="utf-8")
            (source / "pgi").write_text(ALIAS.read_text(encoding="utf-8"), encoding="utf-8")

            home = Path(tmp) / "home"
            home.mkdir()
            bin_dir = home / "bin"

            env = os.environ.copy()
            env["HOME"] = str(home)
            env["PGI_INSTALL_BASE_URL"] = f"file://{source}"

            install = subprocess.run(
                ["sh", str(INSTALLER), "--bin-dir", str(bin_dir)],
                env=env,
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(install.returncode, 0, install.stderr)

            installed_cli = bin_dir / "personal-gitignore"
            installed_alias = bin_dir / "pgi"
            self.assertTrue(installed_cli.exists())
            self.assertTrue(installed_alias.exists())

            repo = Path(tmp) / "repo"
            repo.mkdir()
            subprocess.run(["git", "init"], cwd=repo, check=True, capture_output=True)

            add = subprocess.run(
                [str(installed_alias), "add", "*.from-url-install"],
                cwd=repo,
                env=env,
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(add.returncode, 0, add.stderr)

            listed = subprocess.run(
                [str(installed_cli), "list"],
                cwd=repo,
                env=env,
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertEqual(listed.returncode, 0, listed.stderr)
            self.assertIn("*.from-url-install", listed.stdout)


if __name__ == "__main__":
    unittest.main()

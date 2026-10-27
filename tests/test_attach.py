import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import scripts.attach as attach_module
from scripts.attach import attach


def run(*args: str, cwd: Path) -> str:
    return subprocess.run(
        ["git", *args], cwd=cwd, check=True, capture_output=True, text=True
    ).stdout.strip()


def make_records(path: Path) -> None:
    path.mkdir()
    run("init", "-b", "main", cwd=path)
    run("config", "user.email", "test@example.invalid", cwd=path)
    run("config", "user.name", "Test", cwd=path)
    (path / "PROJECT.md").write_text("# Test project\n")
    run("add", ".", cwd=path)
    run("commit", "-m", "Add records", cwd=path)


class AttachTests(unittest.TestCase):
    def test_records_are_separate_and_invisible_to_application_git(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            app, records = root / "app", root / "records"
            app.mkdir()
            run("init", "-b", "main", cwd=app)
            make_records(records)

            self.assertEqual(attach(app, str(records)), "attached")
            self.assertEqual(attach(app, str(records)), "already attached")
            self.assertEqual(run("status", "--short", cwd=app), "")
            self.assertEqual(
                Path(run("rev-parse", "--show-toplevel", cwd=app / ".ignition")).resolve(),
                (app / ".ignition").resolve(),
            )
            (app / ".ignition" / "PROJECT.md").write_text("# Changed\n")
            self.assertEqual(run("status", "--short", cwd=app), "")

    def test_refuses_existing_non_git_folder(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            app = root / "app"
            app.mkdir()
            run("init", "-b", "main", cwd=app)
            (app / ".ignition").mkdir()
            with self.assertRaisesRegex(ValueError, "separate Git checkout"):
                attach(app, str(root / "records"))

    def test_attaches_inside_linked_worktree(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            app, worktree, records = root / "app", root / "linked", root / "records"
            app.mkdir()
            run("init", "-b", "main", cwd=app)
            run("config", "user.email", "test@example.invalid", cwd=app)
            run("config", "user.name", "Test", cwd=app)
            run("commit", "--allow-empty", "-m", "Initial", cwd=app)
            run("worktree", "add", "-b", "linked", str(worktree), cwd=app)
            make_records(records)

            self.assertEqual(attach(worktree, str(records)), "attached")
            self.assertEqual(run("status", "--short", cwd=worktree), "")
            exclude = Path(run("rev-parse", "--git-path", "info/exclude", cwd=worktree))
            if not exclude.is_absolute():
                exclude = worktree / exclude
            self.assertIn(".ignition/", exclude.read_text().splitlines())

    def test_failed_postcheck_removes_new_checkout(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            app, records = root / "app", root / "records"
            app.mkdir()
            run("init", "-b", "main", cwd=app)
            make_records(records)
            original_git = attach_module.git
            exclude = Path(run("rev-parse", "--git-path", "info/exclude", cwd=app))
            if not exclude.is_absolute():
                exclude = app / exclude
            original_exclude = exclude.read_bytes() if exclude.exists() else None

            def failing_status(*args: str, cwd: Path) -> str:
                if args and args[0] == "status":
                    return "?? .ignition/"
                return original_git(*args, cwd=cwd)

            with patch.object(attach_module, "git", side_effect=failing_status):
                with self.assertRaisesRegex(RuntimeError, "still appears"):
                    attach(app, str(records))
            self.assertFalse((app / ".ignition").exists())
            self.assertEqual(run("status", "--short", cwd=app), "")
            self.assertEqual(exclude.read_bytes() if exclude.exists() else None, original_exclude)
            self.assertEqual(attach(app, str(records)), "attached")

    def test_rejects_existing_checkout_from_other_records_repository(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            app, records_a, records_b = root / "app", root / "records-a", root / "records-b"
            app.mkdir()
            run("init", "-b", "main", cwd=app)
            make_records(records_a)
            make_records(records_b)

            self.assertEqual(attach(app, str(records_a)), "attached")
            with self.assertRaisesRegex(ValueError, "different records repository"):
                attach(app, str(records_b))
            self.assertEqual(run("remote", "get-url", "origin", cwd=app / ".ignition"), str(records_a))

    def test_existing_checkout_accepts_same_local_repo_via_relative_path(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            app, records = root / "app", root / "records"
            app.mkdir()
            run("init", "-b", "main", cwd=app)
            make_records(records)

            self.assertEqual(attach(app, "../records"), "attached")
            self.assertEqual(attach(app, str(records)), "already attached")

    def test_rejects_credentials_in_records_url_without_printing_them(self):
        with tempfile.TemporaryDirectory() as temp:
            app = Path(temp) / "app"
            app.mkdir()
            run("init", "-b", "main", cwd=app)
            with self.assertRaisesRegex(ValueError, "credential helper") as error:
                attach(app, "https://private-token-do-not-print@example.invalid/records.git")
            self.assertNotIn("private-token-do-not-print", str(error.exception))
            self.assertFalse((app / ".ignition").exists())

    def test_failed_clone_does_not_print_url(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            app = root / "app"
            app.mkdir()
            run("init", "-b", "main", cwd=app)
            secret = "private-token-do-not-print"
            invalid_url = root / secret
            script = Path(attach_module.__file__)

            result = subprocess.run(
                [sys.executable, str(script), "--app", str(app), "--records-url", str(invalid_url)],
                capture_output=True,
                text=True,
                check=False,
            )
            self.assertNotEqual(result.returncode, 0)
            self.assertNotIn(secret, result.stderr)
            self.assertFalse((app / ".ignition").exists())


if __name__ == "__main__":
    unittest.main()

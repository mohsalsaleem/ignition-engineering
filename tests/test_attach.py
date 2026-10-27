import subprocess
import tempfile
import unittest
from pathlib import Path

from scripts.attach import attach


def run(*args: str, cwd: Path) -> str:
    return subprocess.run(
        ["git", *args], cwd=cwd, check=True, capture_output=True, text=True
    ).stdout.strip()


class AttachTests(unittest.TestCase):
    def test_records_are_separate_and_invisible_to_application_git(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            app, records = root / "app", root / "records"
            app.mkdir()
            records.mkdir()
            run("init", "-b", "main", cwd=app)
            run("init", "-b", "main", cwd=records)
            run("config", "user.email", "test@example.invalid", cwd=records)
            run("config", "user.name", "Test", cwd=records)
            (records / "PROJECT.md").write_text("# Test project\n")
            run("add", ".", cwd=records)
            run("commit", "-m", "Add records", cwd=records)

            self.assertEqual(attach(app, str(records)), "attached")
            self.assertEqual(attach(app, str(records)), "already attached")
            self.assertEqual(run("status", "--short", cwd=app), "")
            self.assertEqual(
                Path(run("rev-parse", "--show-toplevel", cwd=app / ".ignite")).resolve(),
                (app / ".ignite").resolve(),
            )
            (app / ".ignite" / "PROJECT.md").write_text("# Changed\n")
            self.assertEqual(run("status", "--short", cwd=app), "")

    def test_refuses_existing_non_git_folder(self):
        with tempfile.TemporaryDirectory() as temp:
            root = Path(temp)
            app = root / "app"
            app.mkdir()
            run("init", "-b", "main", cwd=app)
            (app / ".ignite").mkdir()
            with self.assertRaisesRegex(ValueError, "separate Git checkout"):
                attach(app, str(root / "records"))


if __name__ == "__main__":
    unittest.main()

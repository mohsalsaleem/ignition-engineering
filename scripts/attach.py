#!/usr/bin/env python3
"""Attach a product's separate Ignite records repository at .ignite/."""

from __future__ import annotations

import argparse
import subprocess
import sys
from pathlib import Path


def git(*args: str, cwd: Path) -> str:
    result = subprocess.run(
        ["git", *args], cwd=cwd, text=True, capture_output=True, check=True
    )
    return result.stdout.strip()


def attach(app: Path, records_url: str) -> str:
    app = app.resolve()
    if not app.is_dir():
        raise ValueError(f"Application directory does not exist: {app}")
    root = Path(git("rev-parse", "--show-toplevel", cwd=app)).resolve()
    if root != app:
        raise ValueError(f"--app must be the application repository root: {root}")

    target = app / ".ignite"
    if git("ls-files", "--", ".ignite", cwd=app):
        raise ValueError("The application repository already tracks .ignite; stop and review it.")
    if target.exists():
        try:
            nested_root = Path(git("rev-parse", "--show-toplevel", cwd=target)).resolve()
        except subprocess.CalledProcessError as exc:
            raise ValueError(".ignite exists but is not an Ignite Git checkout.") from exc
        if nested_root != target:
            raise ValueError(".ignite exists but is not a separate Git checkout.")
        outcome = "already attached"
    else:
        subprocess.run(
            ["git", "clone", "--", records_url, str(target)],
            cwd=app,
            check=True,
            capture_output=True,
            text=True,
        )
        outcome = "attached"

    git_dir = Path(git("rev-parse", "--absolute-git-dir", cwd=app))
    exclude = git_dir / "info" / "exclude"
    exclude.parent.mkdir(parents=True, exist_ok=True)
    lines = exclude.read_text().splitlines() if exclude.exists() else []
    if ".ignite/" not in lines:
        with exclude.open("a") as output:
            if lines and exclude.stat().st_size and not exclude.read_bytes().endswith(b"\n"):
                output.write("\n")
            output.write(".ignite/\n")

    if git("status", "--short", "--untracked-files=all", "--", ".ignite", cwd=app):
        raise RuntimeError(".ignite still appears in application Git status.")
    return outcome


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--app", required=True, type=Path, help="Application Git root")
    parser.add_argument("--records-url", required=True, help="Private Ignite records Git URL")
    args = parser.parse_args()
    try:
        outcome = attach(args.app, args.records_url)
    except (ValueError, RuntimeError, subprocess.CalledProcessError) as exc:
        print(f"Ignite attach failed: {exc}", file=sys.stderr)
        return 1
    print(f"Ignite {outcome} at {args.app.resolve() / '.ignite'}")
    print("Application Git status is unchanged by Ignite records.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

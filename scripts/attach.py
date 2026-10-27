#!/usr/bin/env python3
"""Attach a product's separate Ignition records repository at .ignition/."""

from __future__ import annotations

import argparse
import re
import shutil
import subprocess
import sys
from pathlib import Path
from urllib.parse import urlsplit


def git(*args: str, cwd: Path) -> str:
    result = subprocess.run(
        ["git", *args], cwd=cwd, text=True, capture_output=True, check=True
    )
    return result.stdout.strip()


def records_identity(url: str, base: Path) -> str:
    """Compare local paths by location and remote URLs by their configured spelling."""
    if "://" in url or re.match(r"^[^/]+:", url):
        return url
    return str((base / url).resolve())


def attach(app: Path, records_url: str) -> str:
    app = app.resolve()
    if not app.is_dir():
        raise ValueError(f"Application directory does not exist: {app}")
    root = Path(git("rev-parse", "--show-toplevel", cwd=app)).resolve()
    if root != app:
        raise ValueError(f"--app must be the application repository root: {root}")
    parsed_url = urlsplit(records_url)
    if parsed_url.scheme in {"http", "https"} and "@" in parsed_url.netloc:
        raise ValueError("Use an SSH URL or Git credential helper, not credentials in --records-url.")

    target = app / ".ignition"
    if git("ls-files", "--", ".ignition", cwd=app):
        raise ValueError("The application repository already tracks .ignition; stop and review it.")
    target_preexisted = target.exists() or target.is_symlink()
    exclude: Path | None = None
    old_exclude: bytes | None = None
    restore_exclude = False
    try:
        if target_preexisted:
            if target.is_symlink():
                raise ValueError(".ignition exists but is not a separate Git checkout.")
            try:
                nested_root = Path(git("rev-parse", "--show-toplevel", cwd=target)).resolve()
            except subprocess.CalledProcessError as exc:
                raise ValueError(".ignition exists but is not an Ignition Git checkout.") from exc
            if nested_root != target:
                raise ValueError(".ignition exists but is not a separate Git checkout.")
            try:
                existing_url = git("config", "--get", "remote.origin.url", cwd=target)
            except subprocess.CalledProcessError as exc:
                raise ValueError(".ignition has no origin records repository.") from exc
            if records_identity(existing_url, target) != records_identity(records_url, app):
                raise ValueError(".ignition belongs to a different records repository.")
            outcome = "already attached"
        else:
            try:
                subprocess.run(
                    ["git", "clone", "--", records_url, str(target)],
                    cwd=app,
                    check=True,
                    capture_output=True,
                    text=True,
                )
            except subprocess.CalledProcessError:
                # CalledProcessError includes the command, which may contain credentials.
                raise RuntimeError("Could not clone records repository; check URL and access.") from None
            outcome = "attached"

        git_path = Path(git("rev-parse", "--git-path", "info/exclude", cwd=app))
        exclude = git_path if git_path.is_absolute() else app / git_path
        exclude.parent.mkdir(parents=True, exist_ok=True)
        old_exclude = exclude.read_bytes() if exclude.exists() else None
        if b".ignition/" not in (old_exclude or b"").splitlines():
            content = old_exclude or b""
            if content and not content.endswith(b"\n"):
                content += b"\n"
            restore_exclude = True
            exclude.write_bytes(content + b".ignition/\n")

        if git("status", "--short", "--untracked-files=all", "--", ".ignition", cwd=app):
            raise RuntimeError(".ignition still appears in application Git status.")
        return outcome
    except Exception:
        if restore_exclude and exclude is not None:
            if old_exclude is None:
                exclude.unlink(missing_ok=True)
            else:
                exclude.write_bytes(old_exclude)
        if not target_preexisted:
            if target.is_symlink() or target.is_file():
                target.unlink()
            elif target.exists():
                shutil.rmtree(target)
        raise


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--app", required=True, type=Path, help="Application Git root")
    parser.add_argument("--records-url", required=True, help="Private Ignition records Git URL")
    args = parser.parse_args()
    try:
        outcome = attach(args.app, args.records_url)
    except (ValueError, RuntimeError, subprocess.CalledProcessError) as exc:
        print(f"Ignition attach failed: {exc}", file=sys.stderr)
        return 1
    print(f"Ignition {outcome} at {args.app.resolve() / '.ignition'}")
    print("Application Git status is unchanged by Ignition records.")
    print("Local build contexts need a separate .ignition exclusion; see LAYOUT.md.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

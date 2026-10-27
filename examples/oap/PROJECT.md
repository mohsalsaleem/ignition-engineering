# Project: oap

## Owner

Mohamed Saleem.

## Repositories

Application: github.com/mohsalsaleem/OpenAppPlatform. This folder is a worked example kept in the Ignition repository; a real product keeps its records in its own repository, at `.ignition/`.

## Environments

Local: Docker Desktop on the owner's Mac (arm64). Test server: vultr_oap (amd64), which also hosts a live site.

## Commands

Before every push: `REWRITE_BASE=origin/main make rewrite-check`. Release: `make release VERSION=...`.

## Risks and boundaries

On a server, touch only OAP-labelled resources and ask before anything destructive. Passwords are typed only by the owner. Publishing needs the owner's approval.

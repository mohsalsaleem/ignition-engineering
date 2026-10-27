---
id: W-20261027-6yedp2
kind: bug
title: Install ships a local copy of a released oapd image
stage: closed
risk: normal
size: small
parent:
created: "2026-10-27T18:03:03Z"
---
# Install ships a local copy of a released oapd image

## Request

During R5.3 (onboarding the vultr_oap server with v0.2.0-rc.1), `oapctl host install` refused to install the released oapd image. Raised by Claude while running R5.3 for the owner. Recorded on 2026-10-27 from the R5.3 logs and the OAP repository, as a worked example.

## Expected and actual

Expected: the server pulls `ghcr.io/mohsalsaleem/oapd:v0.2.0-rc.1@sha256:3f8b3b2d…`, the multi-arch release pinned by digest, from its registry.

Actual: `oapctl: image ghcr.io/mohsalsaleem/oapd:v0.2.0-rc.1@sha256:3f8b3b2d… is arm64 but the server is amd64; build it for linux/amd64`

## Reproduction

On an arm64 Mac that has a local copy of the released image, run `oapctl host install` against an amd64 server.

Regression test: `TestReleasedImageIsPulledByTheServer` in `cmd/oapctl/shipimage_test.go` fakes a docker that reports a local arm64 copy, and requires that a released digest reference is not shipped with `docker save`, while `oapd:dev` still is.

## Root cause

`shipImage` in `cmd/oapctl/install.go` looked for a local copy of the image first. Finding one, it shipped it with `docker save`, after checking its architecture against the server's. A released image is a multi-arch registry reference pinned by digest, but the local copy is only this machine's platform variant: arm64 on the owner's Mac. So any machine that had pulled the release couldn't install it on a server of another architecture.

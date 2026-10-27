# Release: Install ships a local copy of a released oapd image

Work item: W-20261027-6yedp2

## Revision

Tag v0.2.0-rc.2 on e3b3f17, published as a GitHub pre-release on 2026-10-23T10:47:56Z. oapd image `ghcr.io/mohsalsaleem/oapd:v0.2.0-rc.2@sha256:27c9f8737d38ba7ea5a09aada7e8ee14cc33eec259f2bc82a0844275d01cbf68`.

## Checks

The regression test passes at c18a6ac (evidence recorded under verify). The release's `SHA256SUMS` were checked against the published binaries.

## Rollback

The change is only in oapctl's install path, with no server-side data. Rolling back means using rc.1's oapctl, which brings the bug back, so the route forward is a fix release. A server installed with rc.2 can be reverted with `oapctl host uninstall`, using its host manifest.

## Production verification

vultr_oap installed from v0.2.0-rc.2: the server pulled the released image itself (evidence under verify). The R5.3 release journey then passed on that installation (receipt r2, 587 s).

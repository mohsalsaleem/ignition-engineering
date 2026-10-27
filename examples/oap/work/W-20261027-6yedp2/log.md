# Log: Install ships a local copy of a released oapd image

Append-only. Add entries with `ignition log`; don't edit earlier ones.

## 2026-10-27T18:03:03Z · stage · Claude (Opus 5.5)

Created as a bug in intake.

## 2026-10-27T18:03:03Z · stage · Claude (Opus 5.5)

intake → reproduce

## 2026-10-27T18:03:03Z · note · Claude (Opus 5.5)

Evidence E-20261027-1ryen0 (fail): v0.2.0-rc.1 host install is refused on an amd64 server

## 2026-10-27T18:03:03Z · note · Claude (Opus 5.5)

Evidence E-20261027-29be3d (fail): Install ships a local copy of a released, digest-pinned image instead of letting the server pull it

## 2026-10-27T18:03:03Z · stage · Claude (Opus 5.5)

reproduce → diagnose

## 2026-10-27T18:03:03Z · stage · Claude (Opus 5.5)

diagnose → build

## 2026-10-27T18:03:03Z · change · Claude (Opus 5.5)

c18a6ac "Let the server pull a released oapd even when it's local": a registry reference pinned by digest is now pulled by the server, and `docker save` is used only for images that aren't in a registry, such as oapd:dev. The test came first.

## 2026-10-27T18:03:03Z · stage · Claude (Opus 5.5)

build → verify

## 2026-10-27T18:03:03Z · note · Claude (Opus 5.5)

Evidence E-20261027-w0gb75 (pass): A released, digest-pinned image is pulled by the server even when a local copy exists

## 2026-10-27T18:03:03Z · note · Claude (Opus 5.5)

Evidence E-20261027-hygyew (pass): v0.2.0-rc.2 host install pulls the released image on the server and installs OAP

## 2026-10-27T18:03:03Z · finding · Claude (Opus 5.5)

The first rc.2 install failed with "isn't on the server and couldn't be pulled": the GHCR package was private. The owner made it public. This is a separate bug, not this item. It also showed that R5.6's earlier evidence of an anonymous pull was wrong (Docker Desktop had used stored credentials); corrected in OAP 9385169.

## 2026-10-27T18:03:03Z · finding · Claude (Opus 5.5)

In hindsight this item touched release tooling, which the work model rates high risk. It was handled as normal risk. As high risk, it would have stopped for the owner after diagnosis as well as before release.

## 2026-10-27T18:03:03Z · checkpoint · Claude (Opus 5.5)

Leaving verify: Is the evidence enough to release to users?
Owner: Mohamed Saleem
Decision: "Cool", in reply to the proposal to publish v0.2.0-rc.2 with this fix. Taken as approval to publish.

## 2026-10-27T18:03:03Z · stage · Claude (Opus 5.5)

verify → release

## 2026-10-27T18:03:03Z · stage · Claude (Opus 5.5)

release → closed

## 2026-10-27T18:03:03Z · handoff · Claude (Opus 5.5)

Closed. Changed: c18a6ac, released in v0.2.0-rc.2. Verified: regression test failing before and passing after the fix; rc.2 install on vultr_oap. Open: the private-package bug and the hindsight risk finding above, both owned by R5.3's evidence pack.

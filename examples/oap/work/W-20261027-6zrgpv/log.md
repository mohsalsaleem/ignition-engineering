# Log: Install ships a local copy of a released oapd image

Append-only. Add entries with `ignition log`; don't edit earlier ones.

## 2026-10-27T18:20:51Z · stage · claude-code · claude-code · claude-opus-5-5 · -

Created as a bug in intake.

## 2026-10-27T18:20:51Z · stage · claude-code · claude-code · claude-opus-5-5 · -

intake → reproduce

## 2026-10-27T18:20:51Z · evidence · claude-code · claude-code · claude-opus-5-5 · -

E-20261027-yx3xky (fail): v0.2.0-rc.1 host install is refused on an amd64 server

## 2026-10-27T18:20:51Z · evidence · claude-code · claude-code · claude-opus-5-5 · -

E-20261027-j2gy0t (fail): Install ships a local copy of a released, digest-pinned image instead of letting the server pull it

## 2026-10-27T18:20:51Z · stage · claude-code · claude-code · claude-opus-5-5 · -

reproduce → diagnose

## 2026-10-27T18:20:51Z · stage · claude-code · claude-code · claude-opus-5-5 · -

diagnose → build

## 2026-10-27T18:20:51Z · change · claude-code · claude-code · claude-opus-5-5 · -

c18a6ac "Let the server pull a released oapd even when it's local": a registry reference pinned by digest is now pulled by the server, and `docker save` is used only for images that aren't in a registry, such as oapd:dev. The test came first.

## 2026-10-27T18:20:51Z · stage · claude-code · claude-code · claude-opus-5-5 · -

build → verify

## 2026-10-27T18:20:51Z · evidence · claude-code · claude-code · claude-opus-5-5 · -

E-20261027-0gfzwd (pass): A released, digest-pinned image is pulled by the server even when a local copy exists

## 2026-10-27T18:20:51Z · evidence · claude-code · claude-code · claude-opus-5-5 · -

E-20261027-hdrvrq (pass): v0.2.0-rc.2 host install pulls the released image on the server and installs OAP

## 2026-10-27T18:20:51Z · finding · claude-code · claude-code · claude-opus-5-5 · -

The first rc.2 install failed with "isn't on the server and couldn't be pulled": the GHCR package was private. The owner made it public. This is a separate bug, not this item. It also showed that R5.6's earlier evidence of an anonymous pull was wrong (Docker Desktop had used stored credentials); corrected in OAP 9385169.

## 2026-10-27T18:20:51Z · finding · claude-code · claude-code · claude-opus-5-5 · -

In hindsight this item touched release tooling, which the work model rates high risk. It was handled as normal risk. As high risk, it would have stopped for the owner after diagnosis as well as before release.

## 2026-10-27T18:20:51Z · checkpoint · claude-code · claude-code · claude-opus-5-5 · -

Leaving verify: Is the evidence enough to release to users?
Owner: Mohamed Saleem
Decision: "Cool", in reply to the proposal to publish v0.2.0-rc.2 with this fix. Taken as approval to publish.

## 2026-10-27T18:20:51Z · stage · claude-code · claude-code · claude-opus-5-5 · -

verify → release

## 2026-10-27T18:20:51Z · stage · claude-code · claude-code · claude-opus-5-5 · -

release → closed

## 2026-10-27T18:20:51Z · handoff · claude-code · claude-code · claude-opus-5-5 · -

Closed. Changed: c18a6ac, released in v0.2.0-rc.2. Verified: regression test failing before and passing after the fix; rc.2 install on vultr_oap. Open: the private-package bug and the hindsight risk finding above, both owned by R5.3's evidence pack.

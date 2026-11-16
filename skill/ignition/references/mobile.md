# Mobile apps

Read this when the product ships an app people install: iOS, Android,
desktop. Everything in [stages.md](stages.md) still applies. This adds what's
different about apps, stage by stage. Items that ship in an app set
`release: app` in the brief (see "An app" under release in stages.md).

An app differs from a service in four ways that drive everything below:
- A build can't be recalled once it's on people's devices.
- Old builds stay in use for months.
- The release build isn't the build you tested in development.
- A screen can pass every test and still be wrong.

## Set up the project once

Write these into `PROJECT.md`, so every agent and QA handoff gets them. Where
one is missing, that's a task worth doing before the next feature.

**Commands**
- Fast checks per platform: unit tests, and the shared core's tests if logic
  is shared.
- A **release smoke**: build the release configuration exactly as the release
  route does (optimized, shrunk, signed, the same compile flags), install it
  fresh, and run the critical journey. Debug builds hide release-only
  failures: R8 or ProGuard removing classes, a compile flag that drops a
  package's conditions, resources read from the wrong bundle.
- **Screenshot tests** of each screen's states, rendered headless. Pin the
  time zone, locale and clock in the test runner itself, not in a shell
  variable that a reused build daemon ignores, so the references match on
  any machine on any day.
- A **screen walk**: a UI test that opens every tab or main screen on
  made-up demo data, with no server or account, and saves a screenshot of
  each page, light and dark. A person or agent compares these with the
  design.

**Behaviour checks**
- Put each screen's rules in a plain model (state, actions, results) that
  tests drive without a UI, device or network. That's the fast, headless
  check QA needs; views then only draw the model.
- Content rules per screen, as tests: what each screen must never show
  because it belongs on another one. They catch a screen that quietly grows
  other screens' content.
- Name the **spare simulator or emulator** for tests, separate from any
  device someone is signed in on. UI tests uninstall and reset the app.

**Environments**
- Which server each build kind talks to (debug, staging or testing track,
  production), and how a build is switched. Never hard-code a server address
  in a view.
- The stores and tracks per platform, and how builds are numbered. Read the
  last number from the store, not from memory.
- The **release route**: one script or pipeline that checks, uploads and
  publishes. Put the look at the screens before upload as a gate it enforces,
  not a step people remember.

**Risks and boundaries**
- Signing keys, keystores and store API keys live outside the repository,
  and only the owner creates or rotates them. Agents use them through the
  release route and never copy them into records, logs or chat.
- Permissions the app asks for (health data, camera, location,
  notifications), and the store rules that apply to them: privacy labels,
  data-safety forms, review guidelines.
- In-app purchases and subscriptions: which keys exist, and which
  environment is sandbox.

## Shape and design

- **One design per screen**, as the reference QA compares against: what the
  screen holds, in what order, and what it never holds. Say which parts are
  deliberately not built yet.
- **Platform conventions:** follow the platform's own (Human Interface
  Guidelines, Material) unless the design says otherwise and why.
- **Both platforms:** say whether the change ships on each platform in the
  same release. A gap is a decision, logged with `--type scope`, not an
  accident found later.
- **Accessibility from the start:** the largest text sizes, screen readers
  (VoiceOver, TalkBack), contrast in light and dark, and touch target size.
  It's part of the acceptance criteria, not polish.
- **Offline and sync:** what the screen shows without a network, what's
  queued, and what wins on conflict.
- **Compatibility:** if the change touches the API, data or sync format,
  the design says how builds already installed keep working, or what minimum
  version is enforced and how people are told. The server can't assume
  every client updated.
- **Permissions:** when the app asks, what it says, and what works when
  people say no.

## Build

- Logic in models, not views, so it's testable without a device.
- Change both platforms in the same item when both ship it, and keep their
  wording and behaviour the same. Shared copy and shared test transcripts
  help.
- Every user-facing string goes through the platform's localization system,
  even with one language.
- Anything you might need to turn off after release goes behind a switch the
  server controls, because a shipped build can't be changed.
- Keep the build machine healthy. Archives and derived data fill disks, and
  a full disk fails builds with unrelated-looking errors. A locked keychain
  fails code signing. Clean up after each archive.

## Verify

QA tests what will ship, on the build revision. In addition to stages.md:

1. **Tests:** unit and core tests, and screenshot tests against their
   references, on both platforms.
2. **Release smoke** on a simulator or emulator, from a fresh install.
3. **Look at the screens.** Run the screen walk and compare every changed
   screen with its design, light and dark, at the largest text size, on the
   smallest supported screen. Check what's there, what's missing, what
   belongs on another screen, and what's clipped, overlapping or cramped.
   Record one evidence entry per screen, naming the screenshot. Tests
   written from what was built confirm what was built, design mistakes
   included. Only comparing with the design catches those.
4. **Device-only behaviour** on a real device: health data, camera,
   notifications, background work, purchases. Record it as `inconclusive`
   with `--limits` if no device was available, never as a pass.

## Release

Follow "An app" in stages.md. In practice:
- Start with testers (TestFlight, Play internal testing, a pre-release),
  then a staged rollout. Know how to halt it before you start.
- Expire or replace the builds this one supersedes, so testers only see the
  new one.
- Number from the store, not from memory, and never reuse a number.
- Production verification means installing from the store or track on a real
  device and running the journey, after store processing finishes.

## Observe

Use real signals:
- the crash-free rate per build;
- crash reports from the stores (TestFlight feedback, Play vitals) or the
  app's crash reporting;
- adoption of the new version;
- store reviews and feedback.

A crash that appears only in the release build becomes a bug, reproduced on a
release build.

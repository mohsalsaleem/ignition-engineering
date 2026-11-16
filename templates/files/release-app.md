# Release: {{.Title}}

Work item: {{.ID}}

## Revision

<!-- The exact commit every build below was made from. -->

## Builds

<!-- One line per platform: the app version and build number (iOS CFBundleVersion, Android versionCode), where it went (TestFlight, Play internal testing, a GitHub pre-release, App Store, Play production), and which server environment it talks to. Each store needs a number higher than its last upload. -->

## Checks

<!-- What passed on these builds, not only on debug builds: tests, a smoke run of the release build (release builds are optimized and signed differently, and fail in ways debug builds don't), screens compared with their designs, and store processing or review. Link the evidence IDs. -->

## Rollout

<!-- Who gets it: testers, a staged percentage, everyone. How to halt it (pause a phased release, halt a Play rollout, expire a TestFlight build). Which older builds were expired or replaced. -->

## Compatibility

<!-- People keep old builds installed for months. Which app versions the server must still serve, whether this release changes the API, data or sync format they rely on, and the minimum version, if one is enforced. -->

## Rollback

<!-- An installed app can't be rolled back. Say what you'd do instead: halt the rollout, expire the build, turn the feature off from the server, or ship a fixed build (and how long store review takes). And whether the server can go back without breaking the builds already out. -->

## Production verification

<!-- The real journey checked after release on a real device, installed from the store or testing track the way people get it, not a simulator or a local build. With evidence. -->

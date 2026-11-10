---
targets:
---
# Project: {{.Product}}

<!-- `targets` above: if the product ships on several platforms that must each be verified (for example `targets: ios, android`), list them. Verify then needs passing evidence on every one. Leave it empty otherwise. -->

## Owner

<!-- Who decides at checkpoints. -->

## Repositories

<!-- The repositories this product spans, and where its Ignition records live. -->

## Environments

<!-- Local, staging, production: how each is deployed and checked. -->

## Commands

<!-- Fast check, full test, build, and the critical journey check. -->

## Behaviour checks

<!-- How an agent proves behaviour without a person watching: the command that runs behaviour tests against the real logic (no UI, device or browser where possible), how long it takes, how to run one test, how to add test data, and how to drive the running app from a terminal if there is a way. The QA agent is given this. If the product has nothing like it, say so: building one is usually worth a task. -->

## Risks and boundaries

<!-- Sensitive data, external providers, spending limits, things agents must never do. -->

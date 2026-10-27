# Contract for AI-assisted engineering

Give the agent a work item and the product's `engineering/PROJECT.md` before asking it to code. The agent should:

1. Restate the intended outcome and inspect the affected code, tests, and operating constraints.
2. Work in a small, reviewable branch. Avoid unrelated rewrites and keep existing user data intact.
3. Add or update tests that prove changed behavior, especially failure, permission, and data lifecycle paths. Run the product's configured checks and report actual results.
4. Show the diff, assumptions, tradeoffs, deployment effects, and evidence in a pull request or equivalent review artifact.
5. Follow the product's staged release path and obtain the specified decision for high-risk production changes.
6. Verify the deployed journey and report uncertainty plainly. Do not claim a feature is live because code merged or an endpoint returned 200.

AI can help implement and review, but cannot be the sole source of truth for requirements, security authorization, real-world facts, or production acceptance. Record provenance for generated content and model decisions where the product needs it. Never put credentials, private user files, or production data in prompts, issues, or release notes unless a specific secure path is defined.

Suggested task prompt:

> Work item: [link]. Product profile: `engineering/PROJECT.md`. Implement only the accepted slice. Run the configured checks, show the diff and evidence, and stop at the product's release approval gate.

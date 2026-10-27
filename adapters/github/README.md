# GitHub adapter

The core system works without GitHub. For a GitHub-hosted product:

1. Enable Issues or name another single backlog in `.ignite/PROJECT.md`.
2. Add a pull request template based on `templates/PULL_REQUEST.md`.
3. Run the product's fast and relevant full checks on every pull request. Confirm real PRs receive results; a workflow file alone proves nothing.
4. Protect release branches with pull requests, required checks, and review where the repository plan allows it. If unavailable, document a manual gate and consider the plan change before scaling contributors.
5. Keep deployment and incident evidence linked from work items and release records.

GitHub's [branch protection documentation](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches) describes what can be enforced. Availability depends on repository visibility and plan. Do not claim enforcement until the repository settings and a test PR prove it.

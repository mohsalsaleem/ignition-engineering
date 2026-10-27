# Adopt this system in an existing product

Adopt it one real change at a time. Do not pause product delivery to produce a large documentation set. Follow [LAYOUT.md](LAYOUT.md): the product's `.ignition/` space is a separate Git checkout, not a tracked subdirectory of its application repository.

## Step 1 — establish the truth

Inspect the current repository, backlog, checks, branch settings, environments, deployments, backups, monitoring, secrets, and incident route. Run a recent test command and inspect a recent pull request and deployment. Create a private product Ignition records repository and run `python3 scripts/attach.py --app /path/to/product --records-url <records-git-url>`. The script checks it out at `.ignition/`, excludes it from application Git, and verifies that app Git remains clean. Record what actually works in `.ignition/PROJECT.md`; mark a missing control as missing.

## Step 2 — choose one work tracker

Use the tool already used by the product if it can link work items to code and releases. Enable Issues, Linear, or another tracker only when it will be the shared backlog. Import open work in priority order; do not import every old note. Use `templates/WORK_ITEM.md` for the next change and close or archive duplicate lists.

## Step 3 — make the checks real

Put the fast check, full check, build, and migration/config check in the product profile. Make CI run the appropriate checks on pull requests. Open a small test pull request, observe the check result, and confirm a failing check prevents release by policy or an enforced branch rule. If the hosting plan cannot enforce it, name the person who performs the manual check.

## Step 4 — pilot one change

Take a normal-risk deployed change through Ready → In progress → In review → On staging → Released → Observed → Done. Fill the work item, review, and release templates with actual evidence. Investigations, incidents, and work without a runtime change use their routes in [SYSTEM.md](SYSTEM.md). Check that a different person, or a separate deliberate review pass for a solo operator, can understand the decision and reproduce the verification.

## Step 5 — close the loop

Run a rollback or restore rehearsal appropriate to the product. Record one realistic failure drill or incident using `templates/INCIDENT.md`. At the next weekly review, remove steps that added no useful information and strengthen any control that failed. Update the shared system only for lessons that apply to multiple products.

## A 30-day starting sequence

| Period | Outcome |
| --- | --- |
| Days 1–3 | Product profile accurate; one backlog; current CI and release gaps visible |
| Week 1 | Pull request checks run and a normal change passes all states |
| Week 2 | Staging and production evidence become routine; rollback/restore path rehearsed |
| Weeks 3–4 | Review incidents and delivery friction; version the first improvement to this system |

Success means a new collaborator can answer: what should I work on, how do I prove it works, who may release it, what is running, and how do we recover?

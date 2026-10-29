---
id: {{.ID}}
work: {{.Work}}
claim: {{.Claim}}
criteria:
result: {{.Result}}
revision: {{.Revision}}
environment: {{.Environment}}
method: {{.Method}}
by: {{.By}}
agent:
model:
session:
at: {{.At}}
---
# Evidence: {{.Claim}}

## Observed

<!-- What was actually seen: output, numbers, links to CI runs or logs. No secrets. -->
{{.Observed}}

## Limits

<!-- What this does not show. -->
{{.Limits}}

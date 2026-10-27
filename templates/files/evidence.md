---
id: {{.ID}}
work: {{.Work}}
claim: {{.Claim}}
result: {{.Result}}
revision: {{.Revision}}
environment: {{.Environment}}
method: {{.Method}}
by: {{.By}}
at: {{.At}}
---
# Evidence: {{.Claim}}

## Observed

<!-- What was actually seen: output, numbers, links to CI runs or logs. No secrets. -->
{{.Observed}}

## Limits

<!-- What this does not show. -->
{{.Limits}}

---
type: feature
status: active
owner: "@armosec/backend"
scope: repo
---

# Paired MITRE mappings (NAUT-1586)

RuntimeRule.mitreMappings is an optional list of {tactic, technique} objects.
A non-empty list is authoritative. GetMitreMappings falls back to the existing
mitreTactic / mitreTechnique scalar fields for older stored rules.
Writers project the first pair into those scalar fields for older readers.
JSON, BSON and YAML use the same pair structure. No existing scalar changes type.
Deploy readers before publishing multi-mapping rules. R1189 uses TA0008/T1021
and TA0007/T1613; preserve that association rather than building parallel lists.

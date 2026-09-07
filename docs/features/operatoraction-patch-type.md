---
type: feature
status: active
owner: "@armosec/backend"
scope: repo
related_code:
  - apis/operatoraction.go
  - apis/operatoraction_test.go
---

# OperatorAction "patch" type

`OperatorActionArgs` gained a fifth remediation shape: `OperatorActionPatch` ("patch"), plus two
new args fields, `Patch` and `PatchType`.

## Why it exists

`kubescape/operator` PR #410 (https://github.com/kubescape/operator/pull/410) added a generic
patch remediation action — apply a Strategic Merge Patch or JSON Merge Patch to a
Deployment/StatefulSet/DaemonSet/Pod without sending the full workload object. This repo didn't
define the action yet, so the operator worked around it: a *local* `OperatorActionPatch` constant
in `mainhandler/remediators/patch.go`, and reading `"patch"`/`"patchType"` directly off the raw
`Command.Args` map (`extractPatchArgs` in `mainhandler/actionhandler.go`) instead of through the
typed `OperatorActionArgs` struct every other action field uses. This change adds the type support
here so callers can build patch commands through the typed struct like every other action.

## Shape

```go
const OperatorActionPatch OperatorActionType = "patch"

type OperatorActionArgs struct {
    // ... existing fields (Action, Target, Selector, FindingRef, DryRun, TTL, Reason) ...

    // Patch is the raw patch body for the "patch" action (required when
    // Action == OperatorActionPatch). It is a JSON or YAML object, encoded
    // as a string. Must be object-shaped: RFC 6902 JSON Patch arrays are
    // not supported.
    Patch string `json:"patch,omitempty"`
    // PatchType selects the patch action's patch type: "strategic" (the
    // default when empty) or "merge". Ignored for every other action.
    PatchType string `json:"patchType,omitempty"`
}
```

## Scope of this change

Purely the typed contract addition: the new constant, the two new fields, and round-trip test
coverage (`TestOperatorActionArgsRoundTripPatch` in `apis/operatoraction_test.go`). `IsDryRun`,
`ToArgs`, and `OperatorActionArgsFromMap` are unchanged — they're already generic (reflect-free
JSON marshal/unmarshal) and needed no changes to carry the new fields.

Patch semantics, validation, and safety rails (canonicalizing the patch body, rejecting JSON Patch
arrays, applying strategic-merge vs. JSON-merge) live entirely in the operator repo
(`mainhandler/remediators/patch.go`'s `Plan` / `canonicalizePatch`), not here. This schema only
defines the wire shape the operator already validates and tests against.

## Compatibility

Additive for serialization and for every non-patch action: both new fields are `omitempty`, no
existing field renamed/retyped/retagged, and commands that don't use `OperatorActionPatch`
serialize exactly as before.

**Not yet additive for deserializing an existing patch command.** The operator's
`extractPatchArgs` (`mainhandler/actionhandler.go`) currently accepts `patch` as *either* a JSON
string *or* a raw JSON object (it `json.Marshal`s the object case back to a string). But
`handleOperatorAction` parses the whole `Command.Args` map through
`apis.OperatorActionArgsFromMap` first and hard-errors on failure — so once the operator bumps its
`armoapi-go` pin to include this `Patch string` field, any producer still sending an
object-shaped `patch` value will fail that parse (`json: cannot unmarshal object into Go struct
field ... of type string`) before `extractPatchArgs` ever runs. This wasn't a problem before this
change because the object-shaped case was simply ignored by the typed parse and only
`extractPatchArgs` ever looked at it.

**Rollout order:** the `armoapi-go` bump that picks up this schema must land together with (or
after) switching the operator's `handleOperatorAction`/`extractPatchArgs` to read `args.Patch` /
`args.PatchType` directly and to require `patch` be sent as a string (dropping the raw-object
`default:` branch) — or with every `patch`-command producer confirmed to already send `patch` as
a string. Until that lands, do not bump the operator's `armoapi-go` pin past this schema on a
branch that also accepts patch commands with an object-shaped `patch` payload.

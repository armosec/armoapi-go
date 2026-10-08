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

// OperatorActionPatchBody is a JSON object, encoded as a string. Its
// UnmarshalJSON also accepts a raw JSON object directly, storing its JSON
// text verbatim; any other raw shape (array, number, bool) is rejected.
type OperatorActionPatchBody string

type OperatorActionArgs struct {
    // ... existing fields (Action, Target, Selector, FindingRef, DryRun, TTL, Reason) ...

    // Patch is the raw patch body for the "patch" action (required when
    // Action == OperatorActionPatch). See OperatorActionPatchBody for its accepted wire
    // shapes. Must be object-shaped: RFC 6902 JSON Patch arrays are not
    // supported.
    Patch OperatorActionPatchBody `json:"patch,omitempty"`
    // PatchType selects the patch action's patch type: "strategic" (the
    // default when empty) or "merge". Ignored for every other action.
    PatchType string `json:"patchType,omitempty"`
}
```

`OperatorActionPatch`'s doc comment deliberately doesn't call the patch "arbitrary" — the type only
carries the patch body. Validating its shape and safety (size limits, supported target kinds, and
an escalation-field denylist covering things like `hostNetwork`, `serviceAccountName`, and
container `image`/`privileged`/`capabilities` changes) is enforced entirely by the operator's
`PatchRemediator`, not by this package.

## Scope of this change

Purely the typed contract addition: the new constant, the two new fields, `OperatorActionPatchBody`'s
tolerant-but-shape-checked decoding (see Compatibility below), and test coverage in
`apis/operatoraction_test.go` (`TestOperatorActionArgsRoundTripPatch`,
`TestOperatorActionArgsFromMapPatchAcceptsObjectShape`,
`TestOperatorActionArgsFromMapPatchRejectsNonObjectShapes`). `IsDryRun`, `ToArgs`, and
`OperatorActionArgsFromMap` are unchanged — they're already generic (reflect-free JSON
marshal/unmarshal) and needed no changes to carry the new fields.

Patch semantics, validation, and safety rails (canonicalizing the patch body, rejecting JSON Patch
arrays, applying strategic-merge vs. JSON-merge) live entirely in the operator repo
(`mainhandler/remediators/patch.go`'s `Plan` / `canonicalizePatch`), not here. This schema only
defines the wire shape the operator already validates and tests against.

## Compatibility

Additive for serialization and for every non-patch action: both new fields are `omitempty`, no
existing field renamed/retyped/retagged, and commands that don't use `OperatorActionPatch`
serialize exactly as before.

Additive for deserialization too, including an existing patch command. The operator's
`extractPatchArgs` (`mainhandler/actionhandler.go`) has always accepted `patch` as *either* a JSON
string *or* a raw JSON object (it `json.Marshal`s the object case back to a string). A plain `Patch
string` field would have broken that: `handleOperatorAction` parses the whole `Command.Args` map
through `apis.OperatorActionArgsFromMap` as a single unmarshal *before* `extractPatchArgs` runs, so
an object-shaped `patch` value would have failed the whole parse — not just the patch field, but
`target`, `dryRun`, and `reason` along with it — the moment the operator bumped its `armoapi-go`
pin. `OperatorActionPatchBody.UnmarshalJSON` accepts both shapes for exactly this reason, storing the object
case's JSON text verbatim, so the operator's dependency bump is a no-op for any existing producer
regardless of which shape it sends.

That said, the operator's `extractPatchArgs` still duplicates this parsing (it reads `patch`/
`patchType` off the raw map, not through the typed struct) — collapsing it to read `args.Patch` /
`args.PatchType` directly is a follow-up in the operator repo, not part of this change.

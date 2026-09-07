package apis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func TestOperatorActionArgsRoundTrip(t *testing.T) {
	in := OperatorActionArgs{
		Action: OperatorActionQuarantine,
		Target: &OperatorActionTarget{
			Kind:      "Deployment",
			Namespace: "payments",
			Name:      "api",
		},
		Selector: &OperatorActionSelector{
			Control:     "C-0016",
			MinSeverity: "High",
		},
		FindingRef: "workloadconfigurationscansummaries/payments/api",
		DryRun:     boolPtr(true),
		TTL:        "24h",
		Reason:     "C-0016 allowPrivilegeEscalation",
	}

	m, err := in.ToArgs()
	require.NoError(t, err)

	out, err := OperatorActionArgsFromMap(m)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}

func TestOperatorActionArgsRoundTripPatch(t *testing.T) {
	in := OperatorActionArgs{
		Action:    OperatorActionPatch,
		Target:    &OperatorActionTarget{Kind: "Deployment", Namespace: "payments", Name: "api"},
		Patch:     `{"spec":{"template":{"spec":{"containers":[{"name":"api","securityContext":{"seccompProfile":{"type":"RuntimeDefault"}}}]}}}}`,
		PatchType: "merge",
		DryRun:    boolPtr(true),
		Reason:    "enforce RuntimeDefault seccompProfile",
	}

	m, err := in.ToArgs()
	require.NoError(t, err)
	assert.Equal(t, string(in.Patch), m["patch"], "wire key must be 'patch' — the operator reads it off the raw map")
	assert.Equal(t, in.PatchType, m["patchType"], "wire key must be 'patchType'")

	out, err := OperatorActionArgsFromMap(m)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}

// The operator's own extractPatchArgs has always accepted "patch" as either a
// JSON string or a raw JSON object (re-marshaling the object case back to a
// string). OperatorActionArgsFromMap must tolerate the same, or an
// object-shaped "patch" value fails the whole args parse — not just the
// patch field — since it's a single json.Unmarshal for the entire struct.
func TestOperatorActionArgsFromMapPatchAcceptsObjectShape(t *testing.T) {
	m := map[string]interface{}{
		"action": string(OperatorActionPatch),
		"target": map[string]interface{}{"kind": "Deployment", "namespace": "payments", "name": "api"},
		"patch": map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name": "api",
								"securityContext": map[string]interface{}{
									"seccompProfile": map[string]interface{}{"type": "RuntimeDefault"},
								},
							},
						},
					},
				},
			},
		},
		"patchType": "strategic",
	}

	out, err := OperatorActionArgsFromMap(m)
	require.NoError(t, err, "object-shaped patch must not fail the whole args parse")
	assert.Equal(t, OperatorActionPatch, out.Action)
	assert.Equal(t, "Deployment", out.Target.Kind, "other fields must still parse alongside an object-shaped patch")
	assert.JSONEq(t, `{"spec":{"template":{"spec":{"containers":[{"name":"api","securityContext":{"seccompProfile":{"type":"RuntimeDefault"}}}]}}}}`, string(out.Patch))
	assert.Equal(t, "strategic", out.PatchType)
}

// A command carrying the typed args through the generic map should be
// recoverable on the receiving (operator) side.
func TestOperatorActionArgsViaCommand(t *testing.T) {
	args, err := OperatorActionArgs{
		Action: OperatorActionAnnotate,
		Target: &OperatorActionTarget{Kind: "Deployment", Namespace: "default", Name: "web"},
		DryRun: boolPtr(true),
	}.ToArgs()
	require.NoError(t, err)

	c := Command{CommandName: TypeOperatorAction, Args: args}

	got, err := OperatorActionArgsFromMap(c.Args)
	require.NoError(t, err)
	assert.Equal(t, OperatorActionAnnotate, got.Action)
	assert.True(t, got.IsDryRun())
	assert.Equal(t, "web", got.Target.Name)
}

// Safe-by-default: a producer that omits DryRun must be treated as a dry-run by
// the operator, and must never round-trip into an explicit "apply".
func TestOperatorActionArgsDryRunDefaultsSafe(t *testing.T) {
	in := OperatorActionArgs{
		Action: OperatorActionQuarantine,
		Target: &OperatorActionTarget{Kind: "Deployment", Namespace: "payments", Name: "api"},
	}
	assert.True(t, in.IsDryRun(), "omitted DryRun must be treated as dry-run")

	m, err := in.ToArgs()
	require.NoError(t, err)
	_, present := m["dryRun"]
	assert.False(t, present, "omitted DryRun must not be serialized as an explicit value")

	out, err := OperatorActionArgsFromMap(m)
	require.NoError(t, err)
	assert.True(t, out.IsDryRun(), "operator must read an omitted DryRun as a dry-run")

	// Only an explicit false applies, and that must survive the wire round-trip:
	// DryRun=false serializes as "dryRun": false and reads back as apply.
	apply := OperatorActionArgs{Action: OperatorActionCordon, DryRun: boolPtr(false)}
	assert.False(t, apply.IsDryRun())

	am, err := apply.ToArgs()
	require.NoError(t, err)
	assert.Equal(t, false, am["dryRun"], "explicit DryRun=false must serialize on the wire")

	applyOut, err := OperatorActionArgsFromMap(am)
	require.NoError(t, err)
	assert.False(t, applyOut.IsDryRun(), "operator must read an explicit DryRun=false as apply")
}

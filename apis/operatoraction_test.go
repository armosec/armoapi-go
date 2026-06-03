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

	// Only an explicit false applies.
	apply := OperatorActionArgs{Action: OperatorActionCordon, DryRun: boolPtr(false)}
	assert.False(t, apply.IsDryRun())
}

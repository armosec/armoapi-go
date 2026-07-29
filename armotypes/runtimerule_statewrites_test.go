package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// nodeAgentEventTypes mirrors node-agent's utils.EventType value set
// (node-agent/pkg/utils/events.go:206-234). armotypes must define a constant
// for every one of these, because a stateWrites clause may name any of them.
var nodeAgentEventTypes = []string{
	"exec", "open", "dns", "network", "http", "syscall", "capabilities",
	"ptrace", "bpf", "kmod", "symlink", "hardlink", "unshare", "iouring",
	"ssh", "randomx", "procfs", "fork", "exit", "all",
}

func TestEventType_CoversNodeAgentEventSet(t *testing.T) {
	for _, want := range nodeAgentEventTypes {
		assert.True(t, IsKnownEventType(EventType(want)),
			"armotypes has no EventType constant for node-agent event type %q", want)
	}
}

func TestEventType_K8sAdmissionStillKnown(t *testing.T) {
	assert.True(t, IsKnownEventType(EventTypeK8sAdmission))
}

func TestEventType_UnknownIsRejected(t *testing.T) {
	assert.False(t, IsKnownEventType(EventType("not-a-real-event-type")))
	assert.False(t, IsKnownEventType(EventType("")))
}

func TestRuntimeRule_StateWrites_YAMLRoundTrip(t *testing.T) {
	const in = `
id: R1089
name: mount exec then external connection
severity: 8
profileDependency: 2
stateWrites:
  - eventType: exec
    when: 'true'
    scope: container
    name: mount_exec
    key: string(event.pid)
    value:
      argv: placeholder
    ttl: 10m
expressions:
  message: "'msg'"
  uniqueId: "'uid'"
  ruleExpression:
    - eventType: network
      expression: 'state.has("mount_exec", string(event.pid))'
`
	var rule RuntimeRule
	require.NoError(t, yaml.Unmarshal([]byte(in), &rule))

	require.Len(t, rule.StateWrites, 1)
	w := rule.StateWrites[0]
	assert.Equal(t, EventTypeExec, w.EventType)
	assert.Equal(t, StateScopeContainer, w.Scope)
	assert.Equal(t, "mount_exec", w.Name)
	assert.Equal(t, "string(event.pid)", w.Key)
	assert.Equal(t, "10m", w.TTL)
	assert.Equal(t, "true", w.When)
	assert.Equal(t, "placeholder", w.Value["argv"])
}

func TestRuntimeRule_StateWrites_JSONRoundTrip(t *testing.T) {
	original := RuntimeRule{
		ID: "R1089",
		StateWrites: []StateWrite{{
			EventType: EventTypeK8sAdmission,
			When:      `event.Operation == "CREATE"`,
			Scope:     StateScopeIdentity,
			Name:      "pod_created",
			TTL:       "30m",
		}},
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"stateWrites"`)

	var decoded RuntimeRule
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, original.StateWrites, decoded.StateWrites)
	// Key omitted: identity-scoped markers need no subject discriminator.
	assert.Empty(t, decoded.StateWrites[0].Key)
}

func TestRuntimeRule_NoStateWrites_OmittedFromJSON(t *testing.T) {
	data, err := json.Marshal(RuntimeRule{ID: "R0001"})
	require.NoError(t, err)
	assert.NotContains(t, string(data), "stateWrites",
		"rules without state must serialize exactly as before")
}

func TestStateWrite_EventTypeIsValidated(t *testing.T) {
	w := StateWrite{EventType: EventType("nonsense")}
	assert.False(t, IsKnownEventType(w.EventType))
	assert.True(t, IsKnownEventType(StateWrite{EventType: EventTypePtrace}.EventType))
}

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
	// Rule-selection names derived from a tracer stream, not streams themselves
	// (node-agent/pkg/utils/events.go: RuleEventType).
	"dnsQuery",
}

func TestEventType_CoversNodeAgentEventSet(t *testing.T) {
	for _, want := range nodeAgentEventTypes {
		assert.True(t, IsKnownEventType(EventType(want)),
			"armotypes has no EventType constant for node-agent event type %q", want)
	}
}

// "dnsQuery" is how a rule opts in to forwarded DNS queries and error
// responses. It was accepted in ruleExpression (never validated there) but
// rejected in stateWrites, so no write leg could be driven by a query under any
// spelling. "dns" keeps meaning answered responses; nothing changes for existing
// rules.
func TestEventType_DNSQueryIsKnownAndMayDriveAWrite(t *testing.T) {
	assert.True(t, IsKnownEventType(EventTypeDNSQuery))
	assert.True(t, IsValidStateWriteEventType(EventTypeDNSQuery))
	assert.True(t, IsKnownEventType(EventTypeDNS), "dns stays valid alongside dnsQuery")
	assert.NotEqual(t, EventTypeDNS, EventTypeDNSQuery)
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

// EventTypeAll is a rule-binding wildcard, not an event stream. A state write
// must name a concrete stream to drive it: `stateWrites.eventType: all` would
// otherwise pass validation and then never match anything, which is precisely
// the silently-dead-rule failure this contract exists to prevent.
func TestIsValidStateWriteEventType_RejectsTheAllWildcard(t *testing.T) {
	assert.True(t, IsKnownEventType(EventTypeAll),
		"all remains valid for rule bindings")
	assert.False(t, IsValidStateWriteEventType(EventTypeAll),
		"but it cannot drive a state write")
}

func TestIsValidStateWriteEventType_AcceptsConcreteStreams(t *testing.T) {
	for _, et := range []EventType{
		EventTypeExec, EventTypeOpen, EventTypeNetwork, EventTypePtrace,
		EventTypeKmod, EventTypeBPF, EventTypeK8sAdmission,
	} {
		assert.True(t, IsValidStateWriteEventType(et), "%s should drive a state write", et)
	}
}

func TestIsValidStateWriteEventType_RejectsUnknown(t *testing.T) {
	assert.False(t, IsValidStateWriteEventType(EventType("nonsense")))
	assert.False(t, IsValidStateWriteEventType(EventType("")))
}

func TestStateWrite_CounterFields_YAMLRoundTrip(t *testing.T) {
	in := "eventType: network\nscope: container\nname: ssh_out\nincrement: 1\nttl: 1m\n"
	var w StateWrite
	require.NoError(t, yaml.Unmarshal([]byte(in), &w))
	require.NotNil(t, w.Increment)
	assert.Equal(t, 1, *w.Increment)
	assert.Nil(t, w.Clear)

	clear := true
	out, err := yaml.Marshal(StateWrite{EventType: EventTypeNetwork, Scope: StateScopeContainer, Name: "ssh_out", Clear: &clear, TTL: "1m"})
	require.NoError(t, err)
	assert.Contains(t, string(out), "clear: true")
	assert.NotContains(t, string(out), "increment")
}

func TestStateWrite_CounterFields_OmittedWhenZero(t *testing.T) {
	b, err := json.Marshal(StateWrite{EventType: EventTypeExec, Scope: StateScopeContainer, Name: "x", TTL: "1m"})
	require.NoError(t, err)
	assert.NotContains(t, string(b), `"increment"`)
	assert.NotContains(t, string(b), `"clear"`)
}

// An explicit zero is a value, not an absence: the consuming engine must be
// able to reject `increment: 0` and `clear: false` instead of silently
// treating them as a set.
func TestStateWrite_CounterFields_ExplicitZeroIsPreserved(t *testing.T) {
	var w StateWrite
	require.NoError(t, yaml.Unmarshal([]byte("eventType: exec\nscope: container\nname: x\nincrement: 0\nclear: false\nttl: 1m\n"), &w))
	require.NotNil(t, w.Increment)
	require.NotNil(t, w.Clear)
	assert.Equal(t, 0, *w.Increment)
	assert.False(t, *w.Clear)

	var j StateWrite
	require.NoError(t, json.Unmarshal([]byte(`{"eventType":"exec","scope":"container","name":"x","increment":0,"clear":false,"ttl":"1m"}`), &j))
	require.NotNil(t, j.Increment)
	require.NotNil(t, j.Clear)

	b, err := json.Marshal(j)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"increment":0`)
	assert.Contains(t, string(b), `"clear":false`)
}

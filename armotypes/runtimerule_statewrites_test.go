package armotypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

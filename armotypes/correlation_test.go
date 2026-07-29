package armotypes

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestCorrelationEvidence_NodeAgentShape(t *testing.T) {
	ts := time.Date(2026, 7, 28, 12, 0, 3, 100000000, time.UTC)
	ev := CorrelationEvidence{
		Name:      "mount_exec",
		EventType: string(EventTypeExec),
		Timestamp: ts,
		Scope:     string(StateScopeContainer),
		Key:       "4471",
		Process: &Process{
			PID: 4471, PPID: 1201,
			Comm: "xmrig", Pcomm: "bash",
			Path: "/mnt/data/xmrig", Cwd: "/",
			Uid: ptr.To(uint32(0)), Gid: ptr.To(uint32(0)),
			UpperLayer: ptr.To(false),
		},
		Values: map[string]interface{}{"argv": []interface{}{"-o", "pool.example:4444"}},
	}

	data, err := json.Marshal(ev)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"process"`)
	assert.NotContains(t, string(data), `"admission"`,
		"admission must be omitted on a node-agent entry")

	var decoded CorrelationEvidence
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.NotNil(t, decoded.Process)
	assert.Equal(t, uint32(4471), decoded.Process.PID)
	assert.Nil(t, decoded.Admission)
	assert.Equal(t, ts.UTC(), decoded.Timestamp.UTC())
}

func TestCorrelationEvidence_AdmissionShape(t *testing.T) {
	ev := CorrelationEvidence{
		Name:      "pod_created",
		EventType: string(EventTypeK8sAdmission),
		Timestamp: time.Date(2026, 7, 28, 11, 42, 10, 0, time.UTC),
		Scope:     string(StateScopeIdentity),
		Admission: &AdmissionEvidence{
			Username:        "system:serviceaccount:kube-system:deployer",
			Operation:       "CREATE",
			Resource:        "pods",
			Kind:            "Pod",
			ObjectNamespace: "production",
		},
	}

	data, err := json.Marshal(ev)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"admission"`)
	assert.NotContains(t, string(data), `"process"`,
		"process must be omitted on an admission entry")
	// Key is empty for identity scope: the identity IS the subject.
	assert.NotContains(t, string(data), `"key"`)
	// objectName is absent on CREATE because generated names are assigned
	// after admission.
	assert.NotContains(t, string(data), `"objectName"`)

	var decoded CorrelationEvidence
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.NotNil(t, decoded.Admission)
	assert.Equal(t, "system:serviceaccount:kube-system:deployer", decoded.Admission.Username)
	assert.Nil(t, decoded.Process)
}

func TestCorrelationAlert_EmptyOmitsField(t *testing.T) {
	data, err := json.Marshal(CorrelationAlert{})
	require.NoError(t, err)
	assert.Equal(t, "{}", string(data),
		"an alert with no correlations must add nothing to the payload")
}

func TestRuntimeAlert_CorrelationsAreTopLevel(t *testing.T) {
	alert := RuntimeAlert{
		BaseRuntimeAlert: BaseRuntimeAlert{
			AlertName:   "Process executed from mount connected to external endpoint",
			InfectedPID: 4471,
			Severity:    8,
			Timestamp:   time.Date(2026, 7, 28, 12, 0, 3, 480000000, time.UTC),
		},
		RuleID: "R1089",
		CorrelationAlert: CorrelationAlert{
			Correlations: []CorrelationEvidence{{
				Name:      "mount_exec",
				EventType: string(EventTypeExec),
				Scope:     string(StateScopeContainer),
				Key:       "4471",
				Process:   &Process{PID: 4471, Comm: "xmrig", Path: "/mnt/data/xmrig"},
			}},
		},
	}

	data, err := json.Marshal(alert)
	require.NoError(t, err)

	// Inlined, so "correlations" sits beside "alertName" -- not nested under a
	// "correlationAlert" object.
	var flat map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &flat))
	require.Contains(t, flat, "correlations")
	assert.NotContains(t, flat, "correlationAlert")
	assert.Contains(t, flat, "alertName")

	var decoded RuntimeAlert
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Len(t, decoded.Correlations, 1)
	assert.Equal(t, "mount_exec", decoded.Correlations[0].Name)
	assert.Equal(t, uint32(4471), decoded.Correlations[0].Process.PID)
	// The triggering event's identity is untouched by correlation.
	assert.Equal(t, uint32(4471), decoded.InfectedPID)
	assert.Equal(t, "R1089", decoded.RuleID)
}

func TestRuntimeAlert_WithoutCorrelations_IsUnchanged(t *testing.T) {
	alert := RuntimeAlert{
		BaseRuntimeAlert: BaseRuntimeAlert{AlertName: "Unexpected process launched"},
		RuleID:           "R0001",
	}

	data, err := json.Marshal(alert)
	require.NoError(t, err)
	assert.NotContains(t, string(data), "correlations",
		"non-correlation alerts must serialize exactly as before")
}

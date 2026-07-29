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

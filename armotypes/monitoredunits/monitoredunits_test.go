package monitoredunits

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readGolden(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return b
}

// The golden files are the examples of the wire contract. Decoding and encoding
// one must give the same JSON: no field is lost, renamed or added.
func TestGoldenRoundTrip(t *testing.T) {
	for _, name := range []string{"start_entry.json", "periodic_list_with_stop.json"} {
		t.Run(name, func(t *testing.T) {
			golden := readGolden(t, name)
			var msg Message
			require.NoError(t, json.Unmarshal(golden, &msg))
			out, err := json.Marshal(msg)
			require.NoError(t, err)
			assert.JSONEq(t, string(golden), string(out))
		})
	}
}

func TestDecodeStartEntry(t *testing.T) {
	var msg Message
	require.NoError(t, json.Unmarshal(readGolden(t, "start_entry.json"), &msg))

	assert.Equal(t, SchemaVersion, msg.SchemaVersion)
	assert.Equal(t, "ci-cluster", msg.ClusterName)
	assert.Equal(t, time.Date(2026, 11, 18, 8, 0, 0, 0, time.UTC), msg.AgentStartedAt)
	require.NotNil(t, msg.CloudMetadata)
	assert.Equal(t, armotypes.ProviderAws, msg.CloudMetadata.Provider)
	assert.Equal(t, "i-0abc123def4567890", msg.CloudMetadata.InstanceID)
	assert.False(t, msg.IsPeriodicList())

	require.Len(t, msg.Entries, 1)
	e := msg.Entries[0]
	assert.Equal(t, EntryTypeStart, e.Type)
	assert.Equal(t, ReasonUnitStarted, e.Reason)
	assert.Equal(t, UnitKindK8sPod, e.Unit.Kind)
	require.NotNil(t, e.MonitoredFrom)
	assert.Equal(t, time.Date(2026, 11, 18, 10, 0, 3, 0, time.UTC), *e.MonitoredFrom)
	assert.Nil(t, e.MonitoredUntil)
	assert.Nil(t, e.CI)

	attrs, err := e.Unit.K8sPodAttributes()
	require.NoError(t, err)
	assert.Equal(t, "arc-runners", attrs.Namespace)
	assert.Equal(t, "arc-k8s-test-runner-2kr9d", attrs.Labels["runner-pod"])
	assert.Equal(t, "Pod", attrs.WorkloadKind)
	assert.Nil(t, attrs.WorkloadLabels, "the workload is the pod itself")
}

func TestStopEntryCarriesOnlyKindAndID(t *testing.T) {
	until := time.Date(2026, 11, 18, 10, 7, 12, 0, time.UTC)
	e := Entry{
		Type:           EntryTypeStop,
		Reason:         ReasonUnitEnded,
		Unit:           Unit{Kind: UnitKindK8sPod, ID: "uid-1"},
		MonitoredUntil: &until,
	}
	out, err := json.Marshal(e)
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"stop","reason":"unit-ended","unit":{"kind":"k8s-pod","id":"uid-1"},"monitoredUntil":"2026-11-18T10:07:12Z"}`, string(out))

	_, err = e.Unit.K8sPodAttributes()
	assert.Error(t, err, "a stop entry has no attributes")
}

// An absent list and an empty list mean different things, and both must
// survive the wire.
func TestPeriodicListAbsentVersusEmpty(t *testing.T) {
	var notList Message
	out, err := json.Marshal(notList)
	require.NoError(t, err)
	assert.NotContains(t, string(out), "monitoredUnitIds")

	var emptyList Message
	emptyList.SetMonitoredUnitIDs(nil)
	out, err = json.Marshal(emptyList)
	require.NoError(t, err)
	assert.Contains(t, string(out), `"monitoredUnitIds":[]`)

	var decoded Message
	require.NoError(t, json.Unmarshal(out, &decoded))
	require.True(t, decoded.IsPeriodicList())
	assert.Empty(t, *decoded.MonitoredUnitIDs)

	var absent Message
	require.NoError(t, json.Unmarshal([]byte(`{"schemaVersion":1}`), &absent))
	assert.False(t, absent.IsPeriodicList())
}

// A step pod whose owner has no labels still sends workloadLabels, as {}.
func TestWorkloadLabelsEmptyIsKeptNilIsDropped(t *testing.T) {
	started := time.Date(2026, 11, 18, 10, 0, 1, 0, time.UTC)

	u, err := NewK8sPodUnit("uid-1", started, K8sPodAttributes{
		PodName: "p", Namespace: "ns", WorkloadKind: "Job", WorkloadName: "j",
		WorkloadLabels: map[string]string{},
	})
	require.NoError(t, err)
	assert.Contains(t, string(u.Attributes), `"workloadLabels":{}`)
	attrs, err := u.K8sPodAttributes()
	require.NoError(t, err)
	assert.NotNil(t, attrs.WorkloadLabels)

	u, err = NewK8sPodUnit("uid-2", started, K8sPodAttributes{
		PodName: "p", Namespace: "ns", WorkloadKind: "Pod", WorkloadName: "p",
	})
	require.NoError(t, err)
	assert.NotContains(t, string(u.Attributes), "workloadLabels")
	attrs, err = u.K8sPodAttributes()
	require.NoError(t, err)
	assert.Nil(t, attrs.WorkloadLabels)
}

// A consumer keeps a kind, a reason and attributes that it does not know.
func TestUnknownValuesSurvive(t *testing.T) {
	in := `{"type":"start","reason":"some-new-reason","unit":{"kind":"vm","id":"i-1","attributes":{"hostname":"h","extra":{"a":1}}},"monitoredFrom":"2026-11-18T10:00:03Z"}`
	var e Entry
	require.NoError(t, json.Unmarshal([]byte(in), &e))
	assert.Equal(t, Reason("some-new-reason"), e.Reason)
	assert.Equal(t, UnitKind("vm"), e.Unit.Kind)

	_, err := e.Unit.K8sPodAttributes()
	assert.Error(t, err, "not a k8s-pod")

	out, err := json.Marshal(e)
	require.NoError(t, err)
	assert.JSONEq(t, in, string(out))
}

func TestEntryCIIdentity(t *testing.T) {
	in := `{"type":"start","reason":"unit-started","unit":{"kind":"vm","id":"i-1"},"ci":{"provider":"github","github":{"runnerName":"vm-runner-7","runId":9912,"attempt":2,"job":"build"}}}`
	var e Entry
	require.NoError(t, json.Unmarshal([]byte(in), &e))
	require.NotNil(t, e.CI)
	assert.Equal(t, armotypes.CIProviderGitHub, e.CI.Provider)
	require.NotNil(t, e.CI.GitHub)
	assert.Equal(t, int64(9912), e.CI.GitHub.RunID)

	out, err := json.Marshal(e)
	require.NoError(t, err)
	assert.JSONEq(t, in, string(out))
}

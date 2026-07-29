package armotypes

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

// TestProcessRefMarshalText pins the exact textual form of a ProcessRef. The
// format is part of the wire contract (it is the JSON object key of
// NetworkStream.Processes), so a change here is a breaking change for every
// consumer that decodes the map.
func TestProcessRefMarshalText(t *testing.T) {
	tests := []struct {
		name string
		ref  ProcessRef
		want string
	}{
		{name: "typical", ref: ProcessRef{PID: 12345, StartTime: 91827364}, want: "12345/91827364"},
		{name: "zero pid", ref: ProcessRef{PID: 0, StartTime: 91827364}, want: "0/91827364"},
		{name: "zero start time", ref: ProcessRef{PID: 12345, StartTime: 0}, want: "12345/0"},
		{name: "zero value", ref: ProcessRef{}, want: "0/0"},
		{name: "max pid", ref: ProcessRef{PID: math.MaxUint32, StartTime: 1}, want: "4294967295/1"},
		{name: "max start time", ref: ProcessRef{PID: 1, StartTime: math.MaxUint64}, want: "1/18446744073709551615"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ref.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

// TestProcessRefRoundTrip asserts MarshalText/UnmarshalText are exact inverses,
// including the edge cases that matter for identity: pid 0 and a zero start
// time must survive the round trip rather than being silently dropped, because
// the pair is what disambiguates a reused pid.
func TestProcessRefRoundTrip(t *testing.T) {
	tests := []ProcessRef{
		{PID: 12345, StartTime: 91827364},
		{PID: 0, StartTime: 91827364},
		{PID: 12345, StartTime: 0},
		{},
		{PID: math.MaxUint32, StartTime: math.MaxUint64},
		{PID: 1, StartTime: 1},
	}
	for _, want := range tests {
		text, err := want.MarshalText()
		require.NoError(t, err)

		t.Run(string(text), func(t *testing.T) {
			var got ProcessRef
			require.NoError(t, got.UnmarshalText(text))
			assert.Equal(t, want, got)
		})
	}
}

func TestProcessRefUnmarshalTextErrors(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "empty", text: ""},
		{name: "no separator", text: "12345"},
		{name: "extra separator", text: "1/2/3"},
		{name: "non numeric pid", text: "abc/1"},
		{name: "non numeric start time", text: "1/abc"},
		{name: "empty pid", text: "/1"},
		{name: "empty start time", text: "1/"},
		{name: "negative pid", text: "-1/1"},
		{name: "negative start time", text: "1/-1"},
		{name: "pid overflows uint32", text: "4294967296/1"},
		{name: "start time overflows uint64", text: "1/18446744073709551616"},
		{name: "separator only", text: "/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ProcessRef
			assert.Error(t, got.UnmarshalText([]byte(tt.text)))
		})
	}
}

// TestProcessRefUnmarshalTextDoesNotPartiallyMutate asserts a rejected input
// leaves the receiver untouched. A consumer that ignores the error must not end
// up attributing a connection to a half-parsed process.
func TestProcessRefUnmarshalTextDoesNotPartiallyMutate(t *testing.T) {
	got := ProcessRef{PID: 7, StartTime: 8}
	require.Error(t, got.UnmarshalText([]byte("99/bad")))
	assert.Equal(t, ProcessRef{PID: 7, StartTime: 8}, got)
}

// TestProcessRefAsJSONMapKey is the load-bearing test for the design: the same
// ProcessRef value is both the key of NetworkStream.Processes and the value of
// NetworkStreamEvent.Process, so a consumer can look up the tree with
// stream.Processes[*event.Process] and no string formatting at the call site.
func TestProcessRefAsJSONMapKey(t *testing.T) {
	ref := ProcessRef{PID: 4242, StartTime: 314159}
	in := NetworkStream{
		ProcessAttributionVersion: NetworkStreamProcessAttributionV1,
		Processes: map[ProcessRef]*ProcessTree{
			ref: {ProcessTree: Process{PID: 4242, Comm: "curl", Cmdline: "curl https://example.com"}, ContainerID: "abc123"},
		},
		Entities: map[string]NetworkStreamEntity{
			"entity-1": {
				Kind: NetworkStreamEntityKindContainer,
				Outbound: map[string]NetworkStreamEvent{
					"1.2.3.4/443/TCP": {
						IPAddress: "1.2.3.4",
						Port:      443,
						Protocol:  NetworkStreamEventProtocolTCP,
						Process:   &ref,
					},
				},
			},
		},
	}

	data, err := json.Marshal(in)
	require.NoError(t, err)

	// The map key must serialise as the composite string, and the per-connection
	// reference as the same string, so the two are joinable on the wire.
	assert.Contains(t, string(data), `"processes":{"4242/314159":`)
	assert.Contains(t, string(data), `"process":"4242/314159"`)

	var out NetworkStream
	require.NoError(t, json.Unmarshal(data, &out))
	assert.Equal(t, in, out)

	// The join a consumer actually performs.
	event := out.Entities["entity-1"].Outbound["1.2.3.4/443/TCP"]
	require.NotNil(t, event.Process)
	tree, ok := out.Processes[*event.Process]
	require.True(t, ok, "per-connection ProcessRef must key into Processes")
	assert.Equal(t, "curl", tree.ProcessTree.Comm)
}

// TestNetworkStreamProcessAttributionOmitted asserts the addition is invisible
// on the wire for an un-upgraded sensor. Both consumers decode with plain
// json.Unmarshal and neither sets DisallowUnknownFields, so absent fields are
// the entire backward-compatibility story.
func TestNetworkStreamProcessAttributionOmitted(t *testing.T) {
	in := NetworkStream{
		Entities: map[string]NetworkStreamEntity{
			"entity-1": {
				Kind:     NetworkStreamEntityKindContainer,
				Outbound: map[string]NetworkStreamEvent{"1.2.3.4/443/TCP": {IPAddress: "1.2.3.4", Port: 443}},
			},
		},
	}

	data, err := json.Marshal(in)
	require.NoError(t, err)

	assert.NotContains(t, string(data), "processAttributionVersion")
	assert.NotContains(t, string(data), "processes")
	assert.NotContains(t, string(data), `"process"`)
}

// TestNetworkStreamPreProcessAttributionPayloadDecodes asserts a payload
// produced by a sensor that predates process attribution still decodes, and
// that the absent version marker reads as "no attribution available" rather
// than "attribution ran and found nothing".
func TestNetworkStreamPreProcessAttributionPayloadDecodes(t *testing.T) {
	const legacy = `{"entities":{"e1":{"kind":"container","containerName":"nginx","outbound":{"1.2.3.4/443/TCP":{"ipAddress":"1.2.3.4","port":443,"protocol":"TCP"}}}}}`

	var out NetworkStream
	require.NoError(t, json.Unmarshal([]byte(legacy), &out))

	assert.Zero(t, out.ProcessAttributionVersion)
	assert.Nil(t, out.Processes)
	assert.Nil(t, out.Entities["e1"].Outbound["1.2.3.4/443/TCP"].Process)
}

// TestProcessRefBSONRoundTrip covers the bson tags. mongo-driver honours
// encoding.TextMarshaler/TextUnmarshaler for struct map keys, so the Processes
// map survives a bson round trip with the same composite key as JSON.
func TestProcessRefBSONRoundTrip(t *testing.T) {
	ref := ProcessRef{PID: 4242, StartTime: 314159}
	in := NetworkStream{
		ProcessAttributionVersion: NetworkStreamProcessAttributionV1,
		Processes:                 map[ProcessRef]*ProcessTree{ref: {ProcessTree: Process{PID: 4242, Comm: "curl"}}},
		Entities: map[string]NetworkStreamEntity{
			"entity-1": {Outbound: map[string]NetworkStreamEvent{"1.2.3.4/443/TCP": {IPAddress: "1.2.3.4", Process: &ref}}},
		},
	}

	data, err := bson.Marshal(in)
	require.NoError(t, err)

	var raw bson.M
	require.NoError(t, bson.Unmarshal(data, &raw))
	assert.Equal(t, int32(NetworkStreamProcessAttributionV1), raw["processAttributionVersion"])
	assert.Contains(t, raw["processes"], "4242/314159")

	var out NetworkStream
	require.NoError(t, bson.Unmarshal(data, &out))
	require.NotNil(t, out.Processes[ref])
	assert.Equal(t, "curl", out.Processes[ref].ProcessTree.Comm)
}

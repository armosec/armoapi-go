package armotypes

import (
	"encoding/json"
	"math"
	"sort"
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
	t.Parallel()
	tests := []struct {
		name string
		ref  ProcessRef
		want string
	}{
		{name: "typical", ref: ProcessRef{PID: 12345, StartTimeNs: 91827364}, want: "12345/91827364"},
		{name: "zero pid", ref: ProcessRef{PID: 0, StartTimeNs: 91827364}, want: "0/91827364"},
		{name: "unknown start time", ref: ProcessRef{PID: 12345, StartTimeNs: 0}, want: "12345/0"},
		{name: "zero value", ref: ProcessRef{}, want: "0/0"},
		{name: "max pid", ref: ProcessRef{PID: math.MaxUint32, StartTimeNs: 1}, want: "4294967295/1"},
		{name: "max start time", ref: ProcessRef{PID: 1, StartTimeNs: math.MaxUint64}, want: "1/18446744073709551615"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.ref.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
			// String must agree with the wire form so a logged ref is greppable
			// against the payload.
			assert.Equal(t, tt.want, tt.ref.String())
		})
	}
}

// TestProcessRefRoundTrip asserts ref -> text -> ref is lossless, including the
// edge cases that matter for identity: pid 0 and a zero start time must survive
// rather than being silently dropped, because the pair is what disambiguates a
// reused pid.
//
// Note this is one direction only. text -> ref -> text is NOT injective; see
// TestProcessRefUnmarshalTextAcceptsNonCanonical.
func TestProcessRefRoundTrip(t *testing.T) {
	t.Parallel()
	tests := []ProcessRef{
		{PID: 12345, StartTimeNs: 91827364},
		{PID: 0, StartTimeNs: 91827364},
		{PID: 12345, StartTimeNs: 0},
		{},
		{PID: math.MaxUint32, StartTimeNs: math.MaxUint64},
		{PID: 1, StartTimeNs: 1},
	}
	for _, want := range tests {
		text, err := want.MarshalText()
		require.NoError(t, err)

		t.Run(string(text), func(t *testing.T) {
			t.Parallel()
			var got ProcessRef
			require.NoError(t, got.UnmarshalText(text))
			assert.Equal(t, want, got)
		})
	}
}

// TestProcessRefUnmarshalTextIgnoresExtraComponents is the forward-compatibility
// guarantee. encoding/json aborts the ENTIRE decode when a TextUnmarshaler map
// key fails (unlike int keys, which only skip the entry), and both stream
// consumers nack or drop a message whose decode errors. So a future revision
// that appends a component to the key format must not be able to make older
// consumers reject whole messages.
func TestProcessRefUnmarshalTextIgnoresExtraComponents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		text string
		want ProcessRef
	}{
		{name: "one extra component", text: "12345/678/90", want: ProcessRef{PID: 12345, StartTimeNs: 678}},
		{name: "several extra components", text: "1/2/3/4/5", want: ProcessRef{PID: 1, StartTimeNs: 2}},
		{name: "extra component is non numeric", text: "1/2/abc", want: ProcessRef{PID: 1, StartTimeNs: 2}},
		{name: "trailing separator", text: "1/2/", want: ProcessRef{PID: 1, StartTimeNs: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got ProcessRef
			require.NoError(t, got.UnmarshalText([]byte(tt.text)))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProcessRefUnmarshalTextErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		text string
	}{
		{name: "empty", text: ""},
		{name: "no separator", text: "12345"},
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
			t.Parallel()
			var got ProcessRef
			assert.Error(t, got.UnmarshalText([]byte(tt.text)))
		})
	}
}

// TestProcessRefUnmarshalTextDoesNotPartiallyMutate asserts a rejected input
// leaves the receiver untouched for a DIRECT caller. Note this guarantee does
// not survive encoding/json — see TestNetworkStreamEventMalformedRefInstallsZeroPointer.
func TestProcessRefUnmarshalTextDoesNotPartiallyMutate(t *testing.T) {
	t.Parallel()
	got := ProcessRef{PID: 7, StartTimeNs: 8}
	require.Error(t, got.UnmarshalText([]byte("99/bad")))
	assert.Equal(t, ProcessRef{PID: 7, StartTimeNs: 8}, got)
}

// TestProcessRefUnmarshalTextAcceptsNonCanonical documents that leading zeros
// are accepted, so distinct wire keys for the same process collapse to one map
// entry (last wins). Tightening this would turn a cosmetically odd producer
// into a whole-message decode failure, which is the worse trade — see
// TestProcessRefUnmarshalTextIgnoresExtraComponents.
func TestProcessRefUnmarshalTextAcceptsNonCanonical(t *testing.T) {
	t.Parallel()
	const payload = `{"processes":{
		"12345/1":  {"processTree":{"comm":"first"}},
		"0012345/1":{"processTree":{"comm":"second"}},
		"12345/01": {"processTree":{"comm":"third"}}
	}}`

	var out NetworkStream
	require.NoError(t, json.Unmarshal([]byte(payload), &out))
	require.Len(t, out.Processes, 1, "non-canonical keys for one process collapse to one entry")

	tree := out.Processes[ProcessRef{PID: 12345, StartTimeNs: 1}]
	require.NotNil(t, tree)
	assert.Equal(t, "third", tree.ProcessTree.Comm, "last key wins")
}

// TestProcessRefAsJSONMapKey is the load-bearing test for the design: the same
// ProcessRef value is both the key of NetworkStream.Processes and the value of
// NetworkStreamEvent.ProcessRef, so a consumer resolves the tree without any
// string formatting at the call site.
func TestProcessRefAsJSONMapKey(t *testing.T) {
	t.Parallel()
	ref := ProcessRef{PID: 4242, StartTimeNs: 314159}
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
						IPAddress:  "1.2.3.4",
						Port:       443,
						Protocol:   NetworkStreamEventProtocolTCP,
						ProcessRef: &ref,
					},
				},
			},
		},
	}

	data, err := json.Marshal(in)
	require.NoError(t, err)

	// Both sides must render the ref as the same composite string, so the two
	// are joinable on the wire. Asserted structurally rather than as a
	// substring, so the test does not depend on key ordering.
	var top struct {
		Processes map[string]json.RawMessage `json:"processes"`
		Entities  map[string]struct {
			Outbound map[string]struct {
				ProcessRef string `json:"processRef"`
			} `json:"outbound"`
		} `json:"entities"`
	}
	require.NoError(t, json.Unmarshal(data, &top))
	assert.Equal(t, []string{"4242/314159"}, sortedKeys(top.Processes))
	assert.Equal(t, "4242/314159", top.Entities["entity-1"].Outbound["1.2.3.4/443/TCP"].ProcessRef)

	var out NetworkStream
	require.NoError(t, json.Unmarshal(data, &out))
	assert.Equal(t, in, out)

	// The join a consumer actually performs.
	event := out.Entities["entity-1"].Outbound["1.2.3.4/443/TCP"]
	tree := out.ProcessTreeFor(&event)
	require.NotNil(t, tree, "per-connection ProcessRef must resolve against Processes")
	assert.Equal(t, "curl", tree.ProcessTree.Comm)
}

// TestProcessTreeForDegradesToNil covers every way attribution can be absent or
// broken. ProcessTreeFor exists so consumers do not each reimplement these
// guards — a bare Processes[*event.ProcessRef] panics on a null map value,
// which is reachable from a malformed payload.
func TestProcessTreeForDegradesToNil(t *testing.T) {
	t.Parallel()
	ref := ProcessRef{PID: 1, StartTimeNs: 2}
	tests := []struct {
		name   string
		stream *NetworkStream
		event  *NetworkStreamEvent
	}{
		{name: "nil stream", stream: nil, event: &NetworkStreamEvent{ProcessRef: &ref}},
		{name: "nil event", stream: &NetworkStream{}, event: nil},
		{name: "unattributed event", stream: &NetworkStream{Processes: map[ProcessRef]*ProcessTree{ref: {}}}, event: &NetworkStreamEvent{}},
		{name: "no processes map", stream: &NetworkStream{}, event: &NetworkStreamEvent{ProcessRef: &ref}},
		{
			name:   "dangling ref",
			stream: &NetworkStream{Processes: map[ProcessRef]*ProcessTree{{PID: 9, StartTimeNs: 9}: {}}},
			event:  &NetworkStreamEvent{ProcessRef: &ref},
		},
		{
			name:   "nil map value",
			stream: &NetworkStream{Processes: map[ProcessRef]*ProcessTree{ref: nil}},
			event:  &NetworkStreamEvent{ProcessRef: &ref},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Nil(t, tt.stream.ProcessTreeFor(tt.event))
		})
	}
}

// TestProcessTreeForNilMapValueFromPayload proves the "nil map value" case above
// is reachable from the wire, not just constructible in Go.
func TestProcessTreeForNilMapValueFromPayload(t *testing.T) {
	t.Parallel()
	var out NetworkStream
	require.NoError(t, json.Unmarshal([]byte(`{"processes":{"1/2":null}}`), &out))

	tree, ok := out.Processes[ProcessRef{PID: 1, StartTimeNs: 2}]
	require.True(t, ok, "key present")
	require.Nil(t, tree, "value nil — a bare map index would return a nil tree that reads as found")

	assert.Nil(t, out.ProcessTreeFor(&NetworkStreamEvent{ProcessRef: &ProcessRef{PID: 1, StartTimeNs: 2}}))
}

// TestNetworkStreamEventMalformedRefInstallsZeroPointer documents that
// encoding/json allocates the pointer target before calling UnmarshalText and
// leaves it installed on failure. A caller that ignores the Unmarshal error
// therefore sees a NON-nil ref reading as pid 0. Both stream consumers do check
// the error, so this is a documented sharp edge rather than a live hazard.
func TestNetworkStreamEventMalformedRefInstallsZeroPointer(t *testing.T) {
	t.Parallel()
	for _, payload := range []string{`{"processRef":"abc/1"}`, `{"processRef":"1"}`, `{"processRef":123}`} {
		var ev NetworkStreamEvent
		require.Error(t, json.Unmarshal([]byte(payload), &ev), payload)
		if assert.NotNil(t, ev.ProcessRef, payload) {
			assert.Equal(t, ProcessRef{}, *ev.ProcessRef, payload)
		}
	}
}

// TestNetworkStreamMalformedKeyDiscardsWholeMap records the cost of a typed map
// key: unlike an int-keyed map (which skips the bad entry and continues),
// encoding/json returns immediately, so the whole Processes map is lost and the
// caller gets an error. This is why UnmarshalText tolerates extra components.
func TestNetworkStreamMalformedKeyDiscardsWholeMap(t *testing.T) {
	t.Parallel()
	const payload = `{"entities":{"e1":{"containerName":"nginx"}},"processes":{"BAD":{"processTree":{"comm":"x"}},"1/2":{"processTree":{"comm":"good"}}}}`

	var out NetworkStream
	err := json.Unmarshal([]byte(payload), &out)
	require.Error(t, err)
	assert.Empty(t, out.Processes, "one unparseable key discards every entry, including valid ones")
}

// TestNetworkStreamProcessAttributionOmitted asserts the addition is invisible
// on the wire for an un-upgraded sensor. Both consumers decode with plain
// json.Unmarshal and neither sets DisallowUnknownFields, so absent fields are
// the entire backward-compatibility story.
func TestNetworkStreamProcessAttributionOmitted(t *testing.T) {
	t.Parallel()
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

	var top map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &top))
	assert.Equal(t, []string{"entities"}, sortedKeys(top),
		"neither ProcessAttributionVersion nor Processes may appear when unset")

	var event map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(mustNestedEvent(t, data), &event))
	assert.NotContains(t, sortedKeys(event), "processRef")
}

// TestNetworkStreamPreProcessAttributionPayloadDecodes asserts a payload
// produced by a sensor that predates process attribution still decodes, and
// that the absent version marker reads as "no attribution available" rather
// than "attribution ran and found nothing".
func TestNetworkStreamPreProcessAttributionPayloadDecodes(t *testing.T) {
	t.Parallel()
	const legacy = `{"entities":{"e1":{"kind":"container","containerName":"nginx","outbound":{"1.2.3.4/443/TCP":{"ipAddress":"1.2.3.4","port":443,"protocol":"TCP"}}}}}`

	var out NetworkStream
	require.NoError(t, json.Unmarshal([]byte(legacy), &out))

	assert.Zero(t, out.ProcessAttributionVersion)
	assert.Nil(t, out.Processes)
	assert.Nil(t, out.Entities["e1"].Outbound["1.2.3.4/443/TCP"].ProcessRef)
}

// TestProcessRefBSONRoundTrip covers the bson tags, and pins the fact that bson
// and JSON do NOT agree on the per-connection representation: mongo-driver
// honours encoding.TextMarshaler for map KEYS only, so Processes is keyed by the
// same composite string in both, while the event-level ref is a string in JSON
// and a subdocument in BSON. It round-trips either way, but a Mongo query needs
// processRef.pid where JSONB needs the string.
func TestProcessRefBSONRoundTrip(t *testing.T) {
	t.Parallel()
	ref := ProcessRef{PID: 4242, StartTimeNs: 314159}
	in := NetworkStream{
		ProcessAttributionVersion: NetworkStreamProcessAttributionV1,
		Processes:                 map[ProcessRef]*ProcessTree{ref: {ProcessTree: Process{PID: 4242, Comm: "curl"}}},
		Entities: map[string]NetworkStreamEntity{
			"entity-1": {Outbound: map[string]NetworkStreamEvent{"1.2.3.4/443/TCP": {IPAddress: "1.2.3.4", ProcessRef: &ref}}},
		},
	}

	data, err := bson.Marshal(in)
	require.NoError(t, err)

	var raw bson.M
	require.NoError(t, bson.Unmarshal(data, &raw))
	assert.Equal(t, int32(NetworkStreamProcessAttributionV1), raw["processAttributionVersion"])
	assert.Contains(t, raw["processes"], ref.String(), "map key is the composite string in bson too")

	event := raw["entities"].(bson.M)["entity-1"].(bson.M)["outbound"].(bson.M)["1.2.3.4/443/TCP"].(bson.M)
	// Both components widen to int64: mongo-driver encodes Go unsigned integers
	// as int64 so the full unsigned range survives.
	assert.Equal(t, bson.M{"pid": int64(4242), "startTimeNs": int64(314159)}, event["processRef"],
		"bson emits the event-level ref as a subdocument, NOT the text form")

	var out NetworkStream
	require.NoError(t, bson.Unmarshal(data, &out))
	require.NotNil(t, out.Processes[ref])
	assert.Equal(t, "curl", out.Processes[ref].ProcessTree.Comm)
	assert.Equal(t, &ref, out.Entities["entity-1"].Outbound["1.2.3.4/443/TCP"].ProcessRef)
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// mustNestedEvent extracts the single outbound event object from a marshalled
// NetworkStream.
func mustNestedEvent(t *testing.T, data []byte) []byte {
	t.Helper()
	var top struct {
		Entities map[string]struct {
			Outbound map[string]json.RawMessage `json:"outbound"`
		} `json:"entities"`
	}
	require.NoError(t, json.Unmarshal(data, &top))
	for _, entity := range top.Entities {
		for _, event := range entity.Outbound {
			return event
		}
	}
	t.Fatal("no outbound event found")
	return nil
}

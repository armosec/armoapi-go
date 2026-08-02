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

// Pins the textual form: it is the JSON object key of NetworkStream.Processes,
// so changing it breaks every consumer that decodes the map.
func TestProcessRefMarshalText(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ref  ProcessRef
		want string
	}{
		{name: "typical", ref: ProcessRef{PID: 12345, StartTimeNs: 918273640000000}, want: "12345/918273640000000"},
		{name: "zero pid", ref: ProcessRef{PID: 0, StartTimeNs: 918273640000000}, want: "0/918273640000000"},
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

// ref -> text -> ref is lossless, including pid 0 and a zero start time. The
// reverse direction is not injective; see TestProcessRefUnmarshalTextAcceptsNonCanonical.
func TestProcessRefRoundTrip(t *testing.T) {
	t.Parallel()
	tests := []ProcessRef{
		{PID: 12345, StartTimeNs: 918273640000000},
		{PID: 0, StartTimeNs: 918273640000000},
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

// Forward compatibility: appending a key component must not make older consumers
// reject whole messages. See TestNetworkStreamMalformedKeyDiscardsWholeMap for why.
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

// Rejected input leaves the receiver untouched — for direct callers only; see
// TestNetworkStreamEventMalformedRefInstallsZeroPointer.
func TestProcessRefUnmarshalTextDoesNotPartiallyMutate(t *testing.T) {
	t.Parallel()
	got := ProcessRef{PID: 7, StartTimeNs: 8}
	require.Error(t, got.UnmarshalText([]byte("99/bad")))
	assert.Equal(t, ProcessRef{PID: 7, StartTimeNs: 8}, got)
}

// Leading zeros are accepted, so non-canonical keys for one process collapse to
// a single entry (last wins). Tightening this would fail the whole message.
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

// The load-bearing test: one ProcessRef value serves as both map key and event
// reference, so the join needs no string formatting at the call site.
func TestProcessRefAsJSONMapKey(t *testing.T) {
	t.Parallel()
	ref := ProcessRef{PID: 4242, StartTimeNs: 3141590000000}
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

	// Asserted structurally, not as a substring, so key ordering is irrelevant.
	var top struct {
		Processes map[string]json.RawMessage `json:"processes"`
		Entities  map[string]struct {
			Outbound map[string]struct {
				ProcessRef string `json:"processRef"`
			} `json:"outbound"`
		} `json:"entities"`
	}
	require.NoError(t, json.Unmarshal(data, &top))
	assert.Equal(t, []string{"4242/3141590000000"}, sortedKeys(top.Processes))
	assert.Equal(t, "4242/3141590000000", top.Entities["entity-1"].Outbound["1.2.3.4/443/TCP"].ProcessRef)

	var out NetworkStream
	require.NoError(t, json.Unmarshal(data, &out))
	assert.Equal(t, in, out)

	// The join a consumer actually performs.
	event := out.Entities["entity-1"].Outbound["1.2.3.4/443/TCP"]
	tree := out.ProcessTreeFor(&event)
	require.NotNil(t, tree, "per-connection ProcessRef must resolve against Processes")
	assert.Equal(t, "curl", tree.ProcessTree.Comm)
}

// Every way attribution can be absent or broken. ProcessTreeFor exists so each
// consumer need not reimplement these guards.
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

// encoding/json leaves the allocated pointer installed when UnmarshalText fails,
// so ignoring its error yields a non-nil ref reading as pid 0. Both consumers do
// check the error, so this is a sharp edge rather than a live hazard.
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

// The cost of a typed map key: encoding/json returns immediately rather than
// skipping the bad entry, so one unparseable key discards the whole map. This is
// why UnmarshalText is permissive.
func TestNetworkStreamMalformedKeyDiscardsWholeMap(t *testing.T) {
	t.Parallel()
	const payload = `{"entities":{"e1":{"containerName":"nginx"}},"processes":{"BAD":{"processTree":{"comm":"x"}},"1/2":{"processTree":{"comm":"good"}}}}`

	var out NetworkStream
	err := json.Unmarshal([]byte(payload), &out)
	require.Error(t, err)
	assert.Empty(t, out.Processes, "one unparseable key discards every entry, including valid ones")
}

// The addition must be invisible on the wire for an un-upgraded sensor: absent
// fields are the entire backward-compatibility story.
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

// A pre-attribution payload still decodes, and the absent version marker reads as
// "sensor cannot attribute" rather than "attributed nothing".
func TestNetworkStreamPreProcessAttributionPayloadDecodes(t *testing.T) {
	t.Parallel()
	const legacy = `{"entities":{"e1":{"kind":"container","containerName":"nginx","outbound":{"1.2.3.4/443/TCP":{"ipAddress":"1.2.3.4","port":443,"protocol":"TCP"}}}}}`

	var out NetworkStream
	require.NoError(t, json.Unmarshal([]byte(legacy), &out))

	assert.Zero(t, out.ProcessAttributionVersion)
	assert.Nil(t, out.Processes)
	assert.Nil(t, out.Entities["e1"].Outbound["1.2.3.4/443/TCP"].ProcessRef)
}

// Covers the bson tags, and pins that BSON and JSON disagree on the event-level
// ref: mongo-driver honours TextMarshaler for map keys only, so a Mongo query
// needs processRef.pid where JSONB needs the string.
func TestProcessRefBSONRoundTrip(t *testing.T) {
	t.Parallel()
	ref := ProcessRef{PID: 4242, StartTimeNs: 3141590000000}
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
	// Both widen to int64 — mongo-driver encodes Go unsigned ints that way.
	assert.Equal(t, bson.M{"pid": int64(4242), "startTimeNs": int64(3141590000000)}, event["processRef"],
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

// mustNestedEvent extracts the single outbound event from a marshalled NetworkStream.
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

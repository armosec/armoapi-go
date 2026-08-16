package armotypes

import (
	"encoding/json"
	"math"
	"reflect"
	"sort"
	"strings"
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

// Host and ECS workload identity. Design and the bridge to the Kubernetes-shaped
// fields: docs/features/network-stream-workload-identity.md

// `json:",inline"` is a Kubernetes convention that encoding/json does not
// implement — what actually flattens these structs is anonymous-field promotion.
// Pinned per platform because a nested identity object would be a silent break:
// consumers read these keys at the entity's top level.
//
// The fixtures are the agreed host/ECS mapping, so they also record that the
// typed fields ship *alongside* the bridged Kubernetes-shaped ones.
func TestNetworkStreamEntityIdentityIsFlattened(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		entity NetworkStreamEntity
		want   []string
	}{
		{
			name: "kubernetes container is unchanged by this addition",
			entity: NetworkStreamEntity{
				Kind: NetworkStreamEntityKindContainer,
				NetworkStreamEntityContainer: NetworkStreamEntityContainer{
					ContainerName: "nginx", ContainerID: "cid", PodNamespace: "default",
					PodName: "nginx-abc", WorkloadName: "nginx", WorkloadKind: "Deployment",
				},
			},
			want: []string{"containerID", "containerName", "kind", "podName", "podNamespace", "workloadKind", "workloadName"},
		},
		{
			name: "host keeps kind host and carries the host struct plus the bridge",
			entity: NetworkStreamEntity{
				Kind:                    NetworkStreamEntityKindHost,
				NetworkStreamEntityHost: NetworkStreamEntityHost{HostID: "machine-id-1", HostName: "ip-10-0-0-1"},
				NetworkStreamEntityContainer: NetworkStreamEntityContainer{
					PodNamespace: "host", WorkloadName: "machine-id-1", WorkloadKind: "Host",
				},
			},
			want: []string{"hostID", "hostName", "kind", "podNamespace", "workloadKind", "workloadName"},
		},
		{
			name: "ecs service keeps kind container so the traffic view still writes it",
			entity: NetworkStreamEntity{
				Kind: NetworkStreamEntityKindContainer,
				NetworkStreamEntityECS: NetworkStreamEntityECS{
					ECSClusterName: "prod", ClusterARN: "arn:aws:ecs:us-east-1:1:cluster/prod",
					ServiceName: "checkout", TaskFamily: "checkout-task", TaskRevision: "7",
					TaskARN: "arn:aws:ecs:us-east-1:1:task/prod/abc", ECSTaskID: "abc", LaunchType: "EC2",
				},
				NetworkStreamEntityContainer: NetworkStreamEntityContainer{
					ContainerName: "app", ContainerID: "cid",
					PodNamespace: "checkout", WorkloadName: "checkout-checkout-task-7", WorkloadKind: "ECSService",
				},
			},
			want: []string{
				"clusterArn", "containerID", "containerName", "ecsClusterName", "ecsTaskID", "kind", "launchType",
				"podNamespace", "serviceName", "taskArn", "taskFamily", "taskRevision", "workloadKind", "workloadName",
			},
		},
		{
			name: "standalone ecs task omits serviceName and bridges as ECSTask",
			entity: NetworkStreamEntity{
				Kind: NetworkStreamEntityKindContainer,
				NetworkStreamEntityECS: NetworkStreamEntityECS{
					ECSClusterName: "prod", TaskFamily: "batch", TaskRevision: "3", LaunchType: "FARGATE",
				},
				NetworkStreamEntityContainer: NetworkStreamEntityContainer{
					PodNamespace: "prod", WorkloadName: "batch-3", WorkloadKind: "ECSTask",
				},
			},
			want: []string{"ecsClusterName", "kind", "launchType", "podNamespace", "taskFamily", "taskRevision", "workloadKind", "workloadName"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := json.Marshal(tt.entity)
			require.NoError(t, err)

			var got map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(data, &got))
			assert.Equal(t, tt.want, sortedKeys(got), "identity must be flat, and unset identity must cost nothing")

			var out NetworkStreamEntity
			require.NoError(t, json.Unmarshal(data, &out))
			assert.Equal(t, tt.entity, out)
		})
	}
}

// encoding/json drops **both** fields when two promoted fields at the same depth
// share a JSON key — silently, with no error. Reflecting over every field of the
// three inlined structs is what makes a future colliding field fail here instead
// of erasing identity in production.
func TestNetworkStreamEntityIdentityKeysAreUnique(t *testing.T) {
	t.Parallel()
	var (
		container NetworkStreamEntityContainer
		ecs       NetworkStreamEntityECS
		host      NetworkStreamEntityHost
	)
	// Seeded with NetworkStreamEntity's own keys: a promoted field colliding with
	// one of those is *shadowed* by the outer field rather than dropped, which the
	// key-set assertion below would report only as an opaque diff.
	want := map[string]string{"kind": "", "inbound": "", "outbound": ""}
	fillStringFields(t, &container, want)
	fillStringFields(t, &ecs, want)
	fillStringFields(t, &host, want)
	for _, outer := range []string{"kind", "inbound", "outbound"} {
		delete(want, outer)
	}

	data, err := json.Marshal(NetworkStreamEntity{
		Kind:                         NetworkStreamEntityKindContainer,
		NetworkStreamEntityContainer: container,
		NetworkStreamEntityECS:       ecs,
		NetworkStreamEntityHost:      host,
	})
	require.NoError(t, err)

	var got map[string]string
	require.NoError(t, json.Unmarshal(data, &got))
	delete(got, "kind")
	assert.Equal(t, want, got, "every inlined identity field must reach the wire under its own key")
}

// A payload from a sensor that predates host/ECS identity decodes unchanged, and
// the absent structs read as zero rather than failing the message.
func TestNetworkStreamPreIdentityPayloadDecodes(t *testing.T) {
	t.Parallel()
	const legacy = `{"entities":{"e1":{"kind":"host","workloadKind":"Host","outbound":{"1.2.3.4/443/TCP":{"ipAddress":"1.2.3.4","port":443,"protocol":"TCP"}}}}}`

	var out NetworkStream
	require.NoError(t, json.Unmarshal([]byte(legacy), &out))

	entity := out.Entities["e1"]
	assert.Equal(t, NetworkStreamEntityKindHost, entity.Kind)
	assert.Equal(t, "Host", entity.WorkloadKind, "the bridge field still decodes")
	assert.Equal(t, NetworkStreamEntityECS{}, entity.NetworkStreamEntityECS)
	assert.Equal(t, NetworkStreamEntityHost{}, entity.NetworkStreamEntityHost)
}

// One message per node, but a node can be a host reporting itself alongside the
// ECS containers it runs. Pins that identity stays per entity and that process
// attribution still resolves across the addition.
func TestNetworkStreamMixedPlatformEntitiesRoundTrip(t *testing.T) {
	t.Parallel()
	ref := ProcessRef{PID: 99, StartTimeNs: 10000000}
	in := NetworkStream{
		ProcessAttributionVersion: NetworkStreamProcessAttributionV1,
		Processes:                 map[ProcessRef]*ProcessTree{ref: {ProcessTree: Process{PID: 99, Comm: "curl"}}},
		Entities: map[string]NetworkStreamEntity{
			"ip-10-0-0-1": {
				Kind:                    NetworkStreamEntityKindHost,
				NetworkStreamEntityHost: NetworkStreamEntityHost{HostID: "machine-id-1", HostName: "ip-10-0-0-1"},
				Outbound: map[string]NetworkStreamEvent{
					"1.2.3.4/443/TCP": {IPAddress: "1.2.3.4", Port: 443, Protocol: NetworkStreamEventProtocolTCP, ProcessRef: &ref},
				},
			},
			"container-abc": {
				Kind:                   NetworkStreamEntityKindContainer,
				NetworkStreamEntityECS: NetworkStreamEntityECS{ECSClusterName: "prod", ServiceName: "checkout", LaunchType: "EC2"},
			},
		},
	}

	data, err := json.Marshal(in)
	require.NoError(t, err)
	var out NetworkStream
	require.NoError(t, json.Unmarshal(data, &out))
	assert.Equal(t, in, out)

	host := out.Entities["ip-10-0-0-1"]
	assert.Equal(t, "machine-id-1", host.HostID)
	assert.Empty(t, host.ECSClusterName, "identity does not bleed between entities")
	assert.Empty(t, out.Entities["container-abc"].HostID)

	event := host.Outbound["1.2.3.4/443/TCP"]
	require.NotNil(t, out.ProcessTreeFor(&event), "attribution still resolves alongside identity")
}

// The new structs carry bson:",inline" so BSON agrees with JSON on their keys.
// The pre-existing container embed has no bson tags, so it lands as a nested,
// default-lowercased subdocument instead. Both halves are pinned so neither
// changes by accident. No BSON path consumes the network stream today — see the
// bson section of docs/features/network-stream-workload-identity.md, which carries
// the four-repo survey forward from the process-attribution doc.
func TestNetworkStreamEntityIdentityBSONRoundTrip(t *testing.T) {
	t.Parallel()
	in := NetworkStreamEntity{
		Kind: NetworkStreamEntityKindContainer,
		NetworkStreamEntityECS: NetworkStreamEntityECS{
			ECSClusterName: "prod", ServiceName: "checkout", TaskRevision: "7", LaunchType: "EC2",
		},
		NetworkStreamEntityContainer: NetworkStreamEntityContainer{
			WorkloadName: "checkout-checkout-task-7", WorkloadKind: "ECSService",
		},
	}

	data, err := bson.Marshal(in)
	require.NoError(t, err)

	var raw bson.M
	require.NoError(t, bson.Unmarshal(data, &raw))
	assert.Equal(t, "prod", raw["ecsClusterName"], "inlined, and under the same key as json")
	assert.Equal(t, "EC2", raw["launchType"])
	assert.NotContains(t, raw, "networkstreamentityecs", "bson:\",inline\" must keep the ecs struct flat")
	assert.NotContains(t, raw, "hostID", "an empty identity struct contributes no keys in bson either")

	nested, ok := raw["networkstreamentitycontainer"].(bson.M)
	require.True(t, ok, "the container embed predates this change and has no bson tags, so it nests")
	assert.Equal(t, "ECSService", nested["workloadkind"], "and its keys are default-lowercased, unlike json")

	var out NetworkStreamEntity
	require.NoError(t, bson.Unmarshal(data, &out))
	assert.Equal(t, in, out)
}

// fillStringFields sets every string field of the struct behind ptr to a unique
// marker and records it under the field's JSON key, failing if that key was
// already claimed — which is the collision encoding/json would swallow.
func fillStringFields(t *testing.T, ptr any, keys map[string]string) {
	t.Helper()
	value := reflect.ValueOf(ptr).Elem()
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		key, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		require.NotEmpty(t, key, "%s.%s declares no json name", typ.Name(), field.Name)
		if key == "-" {
			continue // deliberately off the wire, so it cannot collide
		}
		require.Equal(t, reflect.String, field.Type.Kind(),
			"%s.%s is not a string — extend this helper so the field stays covered", typ.Name(), field.Name)
		_, duplicate := keys[key]
		require.False(t, duplicate,
			"json key %q is already claimed on NetworkStreamEntity; encoding/json silently drops both "+
				"fields when two promoted ones collide, or lets an outer field shadow a promoted one", key)

		marker := typ.Name() + "." + field.Name
		value.Field(i).SetString(marker)
		keys[key] = marker
	}
}

// Cloud metadata of the sending host. Top level rather than per entity because it
// describes the machine that sent the message, not any one workload on it.

// The consumer reconstructs cloud metadata from designator attributes, and treats
// an absent cloud provider as "no cloud metadata at all", so the provider has to
// survive the trip.
//
// The nested keys are pinned literally, not just round-tripped through Go: a
// Go-to-Go round trip passes whatever the tags say, and producer and consumer run
// different library versions during a rollout, so a renamed tag is a real
// cross-version break. Note they are snake_case, unlike the camelCase keys of the
// enclosing message — CloudMetadata predates this file and keeps its own spelling.
func TestNetworkStreamCloudMetadataRoundTrip(t *testing.T) {
	t.Parallel()
	in := NetworkStream{
		CloudMetadata: &CloudMetadata{
			Provider:   ProviderAws,
			Region:     "us-east-1",
			AccountID:  "123456789012",
			InstanceID: "i-0abc123",
			Hostname:   "ip-192-0-2-10",
			HostType:   HostTypeEc2,
		},
		Entities: map[string]NetworkStreamEntity{
			"entity-1": {Outbound: map[string]NetworkStreamEvent{"1.2.3.4/443/TCP": {IPAddress: "1.2.3.4"}}},
		},
	}

	data, err := json.Marshal(in)
	require.NoError(t, err)

	var top map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &top))
	require.Contains(t, top, "cloudMetadata", "the key is read at the top level of the message")

	var nested map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(top["cloudMetadata"], &nested))
	assert.Equal(t,
		[]string{"account_id", "host_type", "hostname", "instance_id", "provider", "region"},
		sortedKeys(nested),
		"these are the wire keys the consumer maps to designator attributes")

	var out NetworkStream
	require.NoError(t, json.Unmarshal(data, &out))
	assert.Equal(t, in.CloudMetadata, out.CloudMetadata)
}

// A pointer, not a value: encoding/json does not honour omitempty on a struct, so
// a value field would add "cloudMetadata":{} to every Kubernetes payload that has
// never carried one. Absent has to stay absent.
//
// TestNetworkStreamProcessAttributionOmitted already covers this incidentally, by
// pinning the same top-level key set. This is not independent coverage; it is here
// so the failure names the field whose shape caused it.
func TestNetworkStreamCloudMetadataOmittedWhenAbsent(t *testing.T) {
	t.Parallel()
	in := NetworkStream{
		Entities: map[string]NetworkStreamEntity{
			"entity-1": {Outbound: map[string]NetworkStreamEvent{"1.2.3.4/443/TCP": {IPAddress: "1.2.3.4"}}},
		},
	}

	data, err := json.Marshal(in)
	require.NoError(t, err)

	var top map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &top))
	assert.Equal(t, []string{"entities"}, sortedKeys(top), "cloudMetadata may not appear when unset")
}

// A payload from a sensor that predates the field decodes, and reads as "the
// sender does not report cloud metadata" rather than as empty metadata.
func TestNetworkStreamPreCloudMetadataPayloadDecodes(t *testing.T) {
	t.Parallel()
	const legacy = `{"entities":{"e1":{"kind":"container","outbound":{"1.2.3.4/443/TCP":{"ipAddress":"1.2.3.4","port":443}}}}}`

	var out NetworkStream
	require.NoError(t, json.Unmarshal([]byte(legacy), &out))

	assert.Nil(t, out.CloudMetadata)
}

// Covers the bson tags. The stream is persisted, so a missing tag would store the
// field under a Go-cased key that no query looks for.
func TestNetworkStreamCloudMetadataBSONRoundTrip(t *testing.T) {
	t.Parallel()
	in := NetworkStream{
		CloudMetadata: &CloudMetadata{Provider: ProviderAws, Region: "us-east-1", InstanceID: "i-0abc123"},
	}

	data, err := bson.Marshal(in)
	require.NoError(t, err)

	var raw bson.M
	require.NoError(t, bson.Unmarshal(data, &raw))
	require.Contains(t, raw, "cloudMetadata")
	assert.Equal(t, "us-east-1", raw["cloudMetadata"].(bson.M)["region"])

	var out NetworkStream
	require.NoError(t, bson.Unmarshal(data, &out))
	assert.Equal(t, in.CloudMetadata, out.CloudMetadata)
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

package armotypes

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"k8s.io/utils/ptr"
)

// newKeys are the JSON keys this change adds to Trace, StackFrame and
// TraceModule. An alert produced by a consumer that predates the change must
// not carry any of them.
var newKeys = []string{
	`"status"`, `"runtime"`, `"moduleIdx"`, `"fileOffset"`, `"source"`,
	`"statusReason"`, `"truncatedAt"`, `"modules"`, `"hook"`, `"joinKey"`,
	`"tid"`, `"processStartBootNs"`, `"recordBootNs"`, `"eventBootNs"`,
	`"mechanismVersion"`, `"buildId"`, `"inode"`, `"device"`,
}

// fullTrace is a trace that sets every new field: two frames with different
// statuses and runtimes, and two modules they index into.
func fullTrace() Trace {
	return Trace{
		TraceID:  "9f2c1f5a-0f3d-4c6e-9a44-2d0f5b8c71aa",
		Package:  "openssl",
		Language: "Native",
		Stack: []StackFrame{
			{
				FrameID:    "0",
				Function:   "SSL_write",
				File:       "/usr/lib/x86_64-linux-gnu/libssl.so.3",
				Line:       ptr.To(412),
				Address:    "0x7f3c2a1b4e20",
				Arguments:  []string{"ssl", "buf", "num"},
				UserSpace:  true,
				NativeCode: ptr.To(true),
				Anomaly:    true,
				Status:     FrameStatusResolved,
				Runtime:    FrameRuntimeNative,
				ModuleIdx:  ptr.To(0),
				FileOffset: "0x1b4e20",
				Source:     "symtab",
			},
			{
				FrameID:    "1",
				Function:   "main.handleRequest",
				Address:    "0x55a1c0de1234",
				UserSpace:  true,
				NativeCode: ptr.To(false),
				Status:     FrameStatusHeuristic,
				Runtime:    FrameRuntimeGo,
				ModuleIdx:  ptr.To(1),
				FileOffset: "0xde1234",
				Source:     "deferred",
			},
		},
		Status:       StackStatusTruncated,
		StatusReason: "rate-limited",
		TruncatedAt:  ptr.To(2),
		Modules: []TraceModule{
			{
				Path:    "/usr/lib/x86_64-linux-gnu/libssl.so.3",
				BuildID: "3f8a1c2d4e5b6079",
				Inode:   "0x120001",
				Device:  "0x10302",
			},
			{
				Path:    "/usr/local/bin/server",
				BuildID: "aa11bb22cc33dd44",
				Inode:   "0x220002",
				Device:  "0x10302",
			},
		},
		Hook:               "execve",
		JoinKey:            "4471:139872345",
		Tid:                4473,
		ProcessStartBootNs: 90123456789,
		RecordBootNs:       91123456789,
		EventBootNs:        91123400000,
		MechanismVersion:   "wire/13+vocab/1",
	}
}

// fullAlert is a realistic alert carrying fullTrace.
func fullAlert() RuntimeAlert {
	ra := RuntimeAlert{}
	ra.AlertName = "Unexpected process launched"
	ra.InfectedPID = 4471
	ra.Severity = 8
	ra.Timestamp = time.Date(2026, 9, 16, 10, 30, 0, 0, time.UTC)
	ra.Nanoseconds = 1789561800000000000
	ra.BaseRuntimeAlert.UniqueID = "8c1f0d22-5b3a-4f9e-8c2d-1a7b6e4f0d93"
	ra.AgentVersion = "v0.2.41"
	ra.Trace = fullTrace()
	ra.AlertType = AlertTypeRule
	ra.AlertSourcePlatform = AlertSourcePlatformK8sAgent
	ra.RuleID = "R0001"
	ra.HostName = "node-agent-abcde"
	ra.Message = "Unexpected process launched in container"
	ra.ClusterName = "cluster-a"
	ra.PodName = "web-7d9f8c5b6-xk2lm"
	ra.ContainerName = "web"
	return ra
}

// TestRuntimeAlertTrace_JSONRoundTrip pins the JSON wire form of every new
// field. A wrong or missing json tag drops the value here.
func TestRuntimeAlertTrace_JSONRoundTrip(t *testing.T) {
	original := fullAlert()

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded RuntimeAlert
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, original.Trace, decoded.Trace)
	assert.Equal(t, original, decoded)
}

// TestRuntimeAlertTrace_BSONRoundTrip pins the MongoDB write path, which the
// JSON round trip cannot reach: config-service persists runtime incidents.
// It also catches duplicated inline keys — the driver reports those as an
// error, where encoding/json silently drops both fields.
func TestRuntimeAlertTrace_BSONRoundTrip(t *testing.T) {
	original := fullAlert()

	raw, err := bson.Marshal(original)
	require.NoError(t, err)

	var decoded RuntimeAlert
	require.NoError(t, bson.Unmarshal(raw, &decoded))

	assert.Equal(t, original.Trace, decoded.Trace)
	assert.Equal(t, original, decoded)
}

// TestTraceBSONKeys pins the stored key of each new field. A renamed bson tag
// orphans values already written to MongoDB.
func TestTraceBSONKeys(t *testing.T) {
	raw, err := bson.Marshal(fullTrace())
	require.NoError(t, err)

	for _, key := range []string{
		"status", "statusReason", "truncatedAt", "modules", "hook", "joinKey",
		"tid", "processStartBootNs", "recordBootNs", "eventBootNs", "mechanismVersion",
	} {
		_, err := bson.Raw(raw).LookupErr(key)
		assert.NoError(t, err, "trace must be stored under key %q", key)
	}

	frameRaw, err := bson.Marshal(fullTrace().Stack[0])
	require.NoError(t, err)
	for _, key := range []string{"status", "runtime", "moduleIdx", "fileOffset", "source"} {
		_, err := bson.Raw(frameRaw).LookupErr(key)
		assert.NoError(t, err, "stack frame must be stored under key %q", key)
	}

	moduleRaw, err := bson.Marshal(fullTrace().Modules[0])
	require.NoError(t, err)
	for _, key := range []string{"path", "buildId", "inode", "device"} {
		_, err := bson.Raw(moduleRaw).LookupErr(key)
		assert.NoError(t, err, "trace module must be stored under key %q", key)
	}
}

// TestTraceJSONKeys pins the wire key of each new field, the JSON analogue of
// TestTraceBSONKeys. A round trip through the same updated Go type still
// succeeds after a json tag is renamed or removed, because encoding/json
// unmarshals unknown fields as zero values; asserting the literal key is the
// only way to pin the documented wire contract.
func TestTraceJSONKeys(t *testing.T) {
	raw, err := json.Marshal(fullTrace())
	require.NoError(t, err)
	var traceFields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &traceFields))

	for _, key := range []string{
		"status", "statusReason", "truncatedAt", "modules", "hook", "joinKey",
		"tid", "processStartBootNs", "recordBootNs", "eventBootNs", "mechanismVersion",
	} {
		_, ok := traceFields[key]
		assert.True(t, ok, "trace must be encoded under key %q", key)
	}

	frameRaw, err := json.Marshal(fullTrace().Stack[0])
	require.NoError(t, err)
	var frameFields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(frameRaw, &frameFields))
	for _, key := range []string{"status", "runtime", "moduleIdx", "fileOffset", "source"} {
		_, ok := frameFields[key]
		assert.True(t, ok, "stack frame must be encoded under key %q", key)
	}

	moduleRaw, err := json.Marshal(fullTrace().Modules[0])
	require.NoError(t, err)
	var moduleFields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(moduleRaw, &moduleFields))
	for _, key := range []string{"path", "buildId", "inode", "device"} {
		_, ok := moduleFields[key]
		assert.True(t, ok, "trace module must be encoded under key %q", key)
	}
}

// oldShapeGoldenJSON is the exact JSON that the code at origin/main produced
// for the alert built in TestRuntimeAlertTrace_OldShapeIsByteIdentical. It was
// captured before this change and must never be regenerated to make a test
// pass: a diff against it means an existing consumer's payload changed.
const oldShapeGoldenJSON = `{"alertName":"Unexpected process launched","timestamp":"2026-09-16T10:30:00Z","trace":{"traceId":"9f2c1f5a-0f3d-4c6e-9a44-2d0f5b8c71aa","stack":[{"frameId":"0","function":"SSL_write","file":"/usr/lib/x86_64-linux-gnu/libssl.so.3","line":412,"address":"0x7f3c2a1b4e20","arguments":["ssl","buf","num"],"userSpace":true,"nativeCode":true,"anomaly":true}],"package":"openssl","language":"Native"},"malwareFile":{"hashes":{},"timestamps":{"creationTime":"0001-01-01T00:00:00Z","modificationTime":"0001-01-01T00:00:00Z","accessTime":"0001-01-01T00:00:00Z"},"ownership":{},"attributes":{}},"processTree":{"processTree":{"startTime":"0001-01-01T00:00:00Z"}},"signature":{"first_seen":"0001-01-01T00:00:00Z"},"kind":{"Group":"","Version":"","Kind":""},"resource":{"Group":"","Version":"","Resource":""},"cdrevent":{"cloudMetadata":{},"eventData":{}},"request":{},"response":{},"sourcePodInfo":{},"networkscan":{},"alertType":0,"alertSourcePlatform":0,"ruleID":"R0001","hostName":"node-agent-abcde","message":"Unexpected process launched in container"}`

// TestRuntimeAlertTrace_OldShapeIsByteIdentical asserts backward compatibility:
// an alert that sets none of the new fields marshals exactly as it did before
// this change, because every new field carries omitempty.
func TestRuntimeAlertTrace_OldShapeIsByteIdentical(t *testing.T) {
	ra := RuntimeAlert{}
	ra.AlertName = "Unexpected process launched"
	ra.Timestamp = time.Date(2026, 9, 16, 10, 30, 0, 0, time.UTC)
	ra.Trace = Trace{
		TraceID:  "9f2c1f5a-0f3d-4c6e-9a44-2d0f5b8c71aa",
		Package:  "openssl",
		Language: "Native",
		Stack: []StackFrame{{
			FrameID:    "0",
			Function:   "SSL_write",
			File:       "/usr/lib/x86_64-linux-gnu/libssl.so.3",
			Line:       ptr.To(412),
			Address:    "0x7f3c2a1b4e20",
			Arguments:  []string{"ssl", "buf", "num"},
			UserSpace:  true,
			NativeCode: ptr.To(true),
			Anomaly:    true,
		}},
	}
	ra.AlertType = AlertTypeRule
	ra.RuleID = "R0001"
	ra.HostName = "node-agent-abcde"
	ra.Message = "Unexpected process launched in container"

	data, err := json.Marshal(ra)
	require.NoError(t, err)

	assert.Equal(t, oldShapeGoldenJSON, string(data),
		"an alert that sets no new field must marshal exactly as it did before this change")

	for _, key := range newKeys {
		assert.NotContains(t, string(data), key,
			"an alert that sets no new field must not emit %s", key)
	}

	// The empty forms are the tightest statement of the same rule.
	emptyTrace, err := json.Marshal(Trace{})
	require.NoError(t, err)
	assert.Equal(t, "{}", string(emptyTrace))

	emptyFrame, err := json.Marshal(StackFrame{})
	require.NoError(t, err)
	assert.Equal(t, "{}", string(emptyFrame))

	emptyModule, err := json.Marshal(TraceModule{})
	require.NoError(t, err)
	assert.Equal(t, "{}", string(emptyModule))
}

// oldTrace copies the Trace shape as it was before this change. It stands in
// for a consumer that has not bumped armoapi-go yet.
type oldTrace struct {
	TraceID  string     `json:"traceId,omitempty" bson:"traceId,omitempty"`
	Stack    []oldFrame `json:"stack,omitempty" bson:"stack,omitempty"`
	Package  string     `json:"package,omitempty" bson:"package,omitempty"`
	Language string     `json:"language,omitempty" bson:"language,omitempty"`
}

type oldFrame struct {
	FrameID    string   `json:"frameId,omitempty" bson:"frameId,omitempty"`
	Function   string   `json:"function,omitempty" bson:"function,omitempty"`
	File       string   `json:"file,omitempty" bson:"file,omitempty"`
	Line       *int     `json:"line,omitempty" bson:"line,omitempty"`
	Address    string   `json:"address,omitempty" bson:"address,omitempty"`
	Arguments  []string `json:"arguments,omitempty" bson:"arguments,omitempty"`
	UserSpace  bool     `json:"userSpace,omitempty" bson:"userSpace,omitempty"`
	NativeCode *bool    `json:"nativeCode,omitempty" bson:"nativeCode,omitempty"`
	Anomaly    bool     `json:"anomaly,omitempty" bson:"anomaly,omitempty"`
}

// assertOldShapeIntact checks that the fields an old consumer already knew
// survived, whichever codec it decoded with.
func assertOldShapeIntact(t *testing.T, old oldTrace) {
	t.Helper()
	assert.Equal(t, "9f2c1f5a-0f3d-4c6e-9a44-2d0f5b8c71aa", old.TraceID)
	assert.Equal(t, "openssl", old.Package)
	assert.Equal(t, "Native", old.Language)
	require.Len(t, old.Stack, 2)
	assert.Equal(t, "SSL_write", old.Stack[0].Function)
	require.NotNil(t, old.Stack[0].Line)
	assert.Equal(t, 412, *old.Stack[0].Line)
	assert.Equal(t, "0x7f3c2a1b4e20", old.Stack[0].Address)
	assert.True(t, old.Stack[0].Anomaly)
}

// TestTrace_ForwardCompatibleWithOldConsumer asserts that a consumer still on
// the previous armoapi-go decodes a new-shape trace without erroring, and keeps
// every field it already knew.
func TestTrace_ForwardCompatibleWithOldConsumer(t *testing.T) {
	data, err := json.Marshal(fullTrace())
	require.NoError(t, err)

	var old oldTrace
	require.NoError(t, json.Unmarshal(data, &old),
		"an old consumer must ignore the new fields, not fail on them")
	assertOldShapeIntact(t, old)

	// The storage path matters more than the wire here: config-service reads
	// stored incidents back into Go structs, so a stored document carrying the
	// new keys must decode into a consumer that has not bumped yet.
	stored, err := bson.Marshal(fullTrace())
	require.NoError(t, err)

	var oldFromBSON oldTrace
	require.NoError(t, bson.Unmarshal(stored, &oldFromBSON),
		"a stored document with the new keys must decode into the old shape")
	assertOldShapeIntact(t, oldFromBSON)
}

// TestTraceModule_FullRangeIdentityRoundTrips is the reason Inode and Device,
// like Address and FileOffset, are hex strings rather than uint64. BSON has no
// unsigned integer type, so the driver encodes uint64 as a signed 64-bit
// integer and returns "value out of range" for anything above math.MaxInt64 --
// and that error fails the whole incident document, not just the field. Some
// filesystems set the top bits of an inode number (overlayfs with xino), so no
// bound can be assumed. As strings, the full 64-bit range stores and reads back.
func TestTraceModule_FullRangeIdentityRoundTrips(t *testing.T) {
	in := TraceModule{
		Path:   "/usr/local/bin/server",
		Inode:  fmt.Sprintf("0x%x", uint64(math.MaxUint64)),
		Device: fmt.Sprintf("0x%x", uint64(math.MaxUint64)),
	}
	b, err := bson.Marshal(in)
	require.NoError(t, err, "an inode above MaxInt64 must not fail the document")
	var out TraceModule
	require.NoError(t, bson.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}

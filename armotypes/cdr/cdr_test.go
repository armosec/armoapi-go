package cdr

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCdrAlertBatchConnectionLevel locks the wire contract for ConnectionLevel: it
// round-trips under its json tag, is omitted when empty (so existing AWS producers that
// never set it stay byte-for-byte back-compatible), and the constants hold their string
// values (consumers compare against them).
func TestCdrAlertBatchConnectionLevel(t *testing.T) {
	assert.Equal(t, ConnectionLevel("account"), ConnectionLevelAccount)
	assert.Equal(t, ConnectionLevel("organization"), ConnectionLevelOrganization)

	t.Run("account level round-trips", func(t *testing.T) {
		b, err := json.Marshal(CdrAlertBatch{
			CustomerGUID:    "cg",
			CloudAccountID:  "sub-1",
			OrgID:           "tenant-1",
			Provider:        Azure,
			ConnectionLevel: ConnectionLevelAccount,
		})
		require.NoError(t, err)
		assert.Contains(t, string(b), `"connectionLevel":"account"`)

		var got CdrAlertBatch
		require.NoError(t, json.Unmarshal(b, &got))
		assert.Equal(t, ConnectionLevelAccount, got.ConnectionLevel)
	})

	t.Run("omitted when empty (AWS back-compat)", func(t *testing.T) {
		b, err := json.Marshal(CdrAlertBatch{CustomerGUID: "cg", Provider: AWS})
		require.NoError(t, err)
		assert.NotContains(t, jsonKeys(t, b), "connectionLevel",
			"an unset ConnectionLevel must not appear on the wire")
	})
}

// jsonKeys returns the top-level key set of a marshalled batch. Absence is
// asserted on the key set rather than with strings.Contains on the raw JSON: a
// substring check can pass or fail on an unrelated field's *value* that happens to
// contain the name, so it does not actually test what it claims to.
func jsonKeys(t *testing.T, b []byte) map[string]struct{} {
	t.Helper()
	var m map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(b, &m))
	keys := make(map[string]struct{}, len(m))
	for k := range m {
		keys[k] = struct{}{}
	}
	return keys
}

// TestCdrAlertBatchLogsSeen locks the wire contract for LogsSeen. The load-bearing
// property is that "not reported" and "reported as zero" stay distinguishable on
// the wire: conflating them breaks the Pending -> Connected gate in both
// directions — an AWS/Azure producer that reports nothing would be held Pending
// forever, and a GCP collector whose log pipe is silently dead would be reported
// Connected.
func TestCdrAlertBatchLogsSeen(t *testing.T) {
	t.Run("omitted when not reported (AWS/Azure back-compat)", func(t *testing.T) {
		b, err := json.Marshal(CdrAlertBatch{CustomerGUID: "cg", Provider: AWS, IsHeartbeat: true})
		require.NoError(t, err)
		assert.NotContains(t, jsonKeys(t, b), "logsSeen",
			"a producer that does not report logsSeen must not put the key on the wire")
	})

	t.Run("an explicit zero IS serialized", func(t *testing.T) {
		// omitempty on a pointer tests nil, not the pointee. This is the whole
		// reason the field is a pointer: a collector saying "I report this signal
		// and no log has traversed the pipe yet" must be distinguishable from a
		// collector that does not report it at all.
		b, err := json.Marshal(CdrAlertBatch{
			CustomerGUID: "cg", Provider: GCP, IsHeartbeat: true,
			ConnectionLevel: ConnectionLevelAccount,
			LogsSeen:        NewLogsSeen(0),
		})
		require.NoError(t, err)
		assert.Contains(t, jsonKeys(t, b), "logsSeen",
			"an explicit zero must reach the wire, or it is indistinguishable from not reporting")

		var got CdrAlertBatch
		require.NoError(t, json.Unmarshal(b, &got))
		count, reported := got.LogsSeenValue()
		assert.True(t, reported)
		assert.Equal(t, uint64(0), count)
	})

	t.Run("non-zero round-trips", func(t *testing.T) {
		b, err := json.Marshal(CdrAlertBatch{
			CustomerGUID: "cg", CloudAccountID: "my-project", Provider: GCP,
			IsHeartbeat: true, ConnectionLevel: ConnectionLevelAccount,
			LogsSeen: NewLogsSeen(42),
		})
		require.NoError(t, err)
		assert.Contains(t, jsonKeys(t, b), "logsSeen")

		var got CdrAlertBatch
		require.NoError(t, json.Unmarshal(b, &got))
		count, reported := got.LogsSeenValue()
		assert.True(t, reported)
		assert.Equal(t, uint64(42), count)
	})

	t.Run("LogsSeenValue distinguishes absent from zero", func(t *testing.T) {
		var absent CdrAlertBatch
		require.NoError(t, json.Unmarshal([]byte(`{"customerGUID":"cg","isHeartbeat":true}`), &absent))
		count, reported := absent.LogsSeenValue()
		assert.False(t, reported, "an absent count must not read as reported")
		assert.Equal(t, uint64(0), count)

		var zero CdrAlertBatch
		require.NoError(t, json.Unmarshal([]byte(`{"customerGUID":"cg","isHeartbeat":true,"logsSeen":0}`), &zero))
		count, reported = zero.LogsSeenValue()
		assert.True(t, reported, "an explicit zero must read as reported")
		assert.Equal(t, uint64(0), count)
	})

	t.Run("LogsSeenValue is nil-receiver safe", func(t *testing.T) {
		var b *CdrAlertBatch
		count, reported := b.LogsSeenValue()
		assert.False(t, reported)
		assert.Equal(t, uint64(0), count)
	})

	t.Run("carried on a real alert batch too, not only heartbeats", func(t *testing.T) {
		// The count is a property of the collector, not of the message kind, so a
		// normal alert batch may carry it and the consumer must read it the same way.
		b, err := json.Marshal(CdrAlertBatch{
			CustomerGUID: "cg", CloudAccountID: "my-project", Provider: GCP,
			ConnectionLevel: ConnectionLevelOrganization,
			OrgID:           "123456789012",
			RuleFailures:    []CdrAlert{{RuleID: "G0001"}},
			LogsSeen:        NewLogsSeen(7),
		})
		require.NoError(t, err)

		var got CdrAlertBatch
		require.NoError(t, json.Unmarshal(b, &got))
		count, reported := got.LogsSeenValue()
		assert.True(t, reported)
		assert.Equal(t, uint64(7), count)
		assert.False(t, got.IsHeartbeat)
	})
}

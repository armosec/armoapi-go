package cdr

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// GcpAuditLog is persisted, not just sent on the wire: config-service stores
// runtime incidents with it embedded. These tests pin the BSON write path, which
// a JSON round-trip test cannot reach.

// Most audit events omit numResponseItems — only list/query methods return items —
// leaving the zero json.Number, the empty string. The driver encodes json.Number by
// parsing it, "" parses as neither int nor float, and because BSON encoding is not
// per-field that error fails the WHOLE document. Every GCP CDR incident failed to
// persist with:
//
//	cannot marshal type types.Document[*config-service/types.RuntimeAlert] to a
//	BSON Document: strconv.ParseFloat: parsing "": invalid syntax
func TestGcpAuditLog_MarshalsWithoutNumResponseItems(t *testing.T) {
	const create = `{
	  "methodName": "google.pubsub.v1.Publisher.CreateTopic",
	  "serviceName": "pubsub.googleapis.com"
	}`

	var pl GcpAuditLogPayload
	require.NoError(t, json.Unmarshal([]byte(create), &pl))
	require.Equal(t, json.Number(""), pl.NumResponseItems, "an absent field must leave the zero value")

	raw, err := bson.Marshal(pl)
	require.NoError(t, err, "an event without numResponseItems must still persist")

	_, err = bson.Raw(raw).LookupErr("numResponseItems")
	assert.ErrorIs(t, err, bsoncore.ErrElementNotFound, "an empty count must not be stored at all")
}

// When the field IS present it must survive as a real number, in both the forms
// the protobuf-JSON spec allows: protojson emits int64 as a quoted string, but a
// bare number is also valid on parse.
func TestGcpAuditLog_NumResponseItemsRoundTrips(t *testing.T) {
	for name, body := range map[string]string{
		"protojson quoted string": `{"numResponseItems":"3"}`,
		"bare number":             `{"numResponseItems":3}`,
	} {
		t.Run(name, func(t *testing.T) {
			var pl GcpAuditLogPayload
			require.NoError(t, json.Unmarshal([]byte(body), &pl))
			require.Equal(t, json.Number("3"), pl.NumResponseItems)

			raw, err := bson.Marshal(pl)
			require.NoError(t, err)
			assert.Equal(t, bson.TypeInt64, bson.Raw(raw).Lookup("numResponseItems").Type,
				"a count must be stored as a number, so it stays queryable")

			var back GcpAuditLogPayload
			require.NoError(t, bson.Unmarshal(raw, &back))
			assert.Equal(t, json.Number("3"), back.NumResponseItems)
		})
	}
}

// The whole audit-log record must survive the JSON -> BSON -> JSON path an
// incident actually takes, for the common no-count case.
func TestGcpAuditLog_FullRecordSurvivesBSON(t *testing.T) {
	const record = `{
	  "protoPayload": {
	    "methodName": "google.pubsub.v1.Publisher.CreateTopic",
	    "serviceName": "pubsub.googleapis.com",
	    "resourceName": "projects/cdr-project/topics/t1"
	  }
	}`

	var ev GcpAuditLogEvent
	require.NoError(t, json.Unmarshal([]byte(record), &ev))

	raw, err := bson.Marshal(ev)
	require.NoError(t, err)

	var back GcpAuditLogEvent
	require.NoError(t, bson.Unmarshal(raw, &back))
	require.NotNil(t, back.ProtoPayload)
	assert.Equal(t, "google.pubsub.v1.Publisher.CreateTopic", back.ProtoPayload.MethodName)
}

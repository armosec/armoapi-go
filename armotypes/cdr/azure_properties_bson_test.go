package cdr

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// The Properties bag crosses BSON as well as JSON: config-service persists
// runtime incidents with AzureActivityLogEvent embedded. These tests are the
// guard the JSON round-trip test could not be — they pin every shape the
// incidents collection actually holds, including the two that a bare
// json.RawMessage field got wrong in opposite directions.

// storedDoc is what every Azure incident written before AzureProperties existed
// holds: the bag as a plain embedded document, from the map[string]any field.
// Decoding one into a json.RawMessage field failed the whole find with
// "cannot decode document into json.RawMessage", 500ing every Azure incident.
func TestAzureProperties_DecodesStoredDocument(t *testing.T) {
	stored := bson.M{"properties": bson.M{
		"eventCategory": "Administrative",
		"statusCode":    "OK",
		"entity":        "/subscriptions/aaa/resourceGroups/rg1",
	}}
	raw, err := bson.Marshal(stored)
	require.NoError(t, err)

	var ev AzureActivityLogEvent
	require.NoError(t, bson.Unmarshal(raw, &ev))

	var got map[string]any
	require.NoError(t, json.Unmarshal(ev.Properties, &got))
	assert.Equal(t, map[string]any{
		"eventCategory": "Administrative",
		"statusCode":    "OK",
		"entity":        "/subscriptions/aaa/resourceGroups/rg1",
	}, got)
}

// storedBinary is what incidents written while the field was a bare
// json.RawMessage hold: []byte marshals to BSON binary, so the bag became an
// opaque blob no Mongo filter or index could reach. Those documents still have
// to read back.
func TestAzureProperties_DecodesStoredBinary(t *testing.T) {
	body := []byte(`{"eventCategory":"Administrative","statusCode":"OK"}`)
	raw, err := bson.Marshal(bson.M{"properties": body})
	require.NoError(t, err)

	var ev AzureActivityLogEvent
	require.NoError(t, bson.Unmarshal(raw, &ev))
	assert.JSONEq(t, string(body), string(ev.Properties))
}

// Some operations send the bag as a string holding JSON. It is stored as a
// string and must come back as the JSON it is, not as a quoted string.
func TestAzureProperties_DecodesStoredJSONString(t *testing.T) {
	raw, err := bson.Marshal(bson.M{"properties": `{"message":"denied"}`})
	require.NoError(t, err)

	var ev AzureActivityLogEvent
	require.NoError(t, bson.Unmarshal(raw, &ev))
	assert.JSONEq(t, `{"message":"denied"}`, string(ev.Properties))
}

// A bag that is a bare string is kept, quoted, rather than dropped.
func TestAzureProperties_DecodesStoredPlainString(t *testing.T) {
	raw, err := bson.Marshal(bson.M{"properties": "not json at all"})
	require.NoError(t, err)

	var ev AzureActivityLogEvent
	require.NoError(t, bson.Unmarshal(raw, &ev))
	assert.Equal(t, `"not json at all"`, string(ev.Properties))
}

// An absent or null bag decodes to nil, and never fails the incident it belongs to.
func TestAzureProperties_DecodesMissingAndNull(t *testing.T) {
	for name, doc := range map[string]bson.M{
		"absent": {"operationName": "X/Y"},
		"null":   {"properties": nil},
	} {
		t.Run(name, func(t *testing.T) {
			raw, err := bson.Marshal(doc)
			require.NoError(t, err)

			var ev AzureActivityLogEvent
			require.NoError(t, bson.Unmarshal(raw, &ev))
			assert.Nil(t, ev.Properties)
		})
	}
}

// A bag of a BSON type this field never writes — nothing stores a date or an
// ObjectID here — yields nil rather than failing the decode, the same bargain the
// JSON side makes. Note this is about genuinely foreign types: every type the
// writer CAN emit is covered by TestAzureProperties_WriterAndReaderAgreeOnEveryShape.
func TestAzureProperties_ForeignTypeDoesNotFailDecode(t *testing.T) {
	for name, bag := range map[string]interface{}{
		"datetime": primitive.NewDateTimeFromTime(time.Unix(0, 0)),
		"objectID": primitive.NewObjectID(),
	} {
		t.Run(name, func(t *testing.T) {
			raw, err := bson.Marshal(bson.M{"properties": bag})
			require.NoError(t, err)

			var ev AzureActivityLogEvent
			require.NoError(t, bson.Unmarshal(raw, &ev))
			assert.Nil(t, ev.Properties)
		})
	}
}

// Writing stores an object bag as a document, so it stays queryable and
// indexable — the property the binary form silently lost.
func TestAzureProperties_StoresObjectAsDocument(t *testing.T) {
	ev := AzureActivityLogEvent{
		Properties: AzureProperties(`{"eventCategory":"Administrative","statusCode":"OK"}`),
	}
	raw, err := bson.Marshal(ev)
	require.NoError(t, err)

	assert.Equal(t, bson.TypeEmbeddedDocument, bson.Raw(raw).Lookup("properties").Type,
		"an object bag must be stored as a document, not binary")

	var back AzureActivityLogEvent
	require.NoError(t, bson.Unmarshal(raw, &back))
	assert.JSONEq(t, string(ev.Properties), string(back.Properties))
}

// A bag that is not valid JSON is stored as a string rather than dropped.
func TestAzureProperties_StoresInvalidJSONAsString(t *testing.T) {
	ev := AzureActivityLogEvent{Properties: AzureProperties(`{not json`)}
	raw, err := bson.Marshal(ev)
	require.NoError(t, err)
	assert.Equal(t, bson.TypeString, bson.Raw(raw).Lookup("properties").Type)
}

// An empty bag is omitted from both encodings — omitempty still holds with a
// custom marshaller on the field.
func TestAzureProperties_EmptyBag(t *testing.T) {
	ev := AzureActivityLogEvent{}

	raw, err := bson.Marshal(ev)
	require.NoError(t, err)
	_, err = bson.Raw(raw).LookupErr("properties")
	assert.ErrorIs(t, err, bsoncore.ErrElementNotFound, "an empty bag must not be stored at all")

	b, err := json.Marshal(ev)
	require.NoError(t, err)
	assert.NotContains(t, string(b), `"properties"`)
}

// The bag survives the full JSON -> BSON -> JSON path an incident actually
// takes: decoded from the collector's alert, stored by config-service, read
// back and served to the UI, which renders the whole azureActivityLog object.
func TestAzureProperties_JSONToBSONToJSON(t *testing.T) {
	// "time" has no omitempty and is always served; it is part of the shape.
	const record = `{
	  "time": "2026-08-12T15:11:47Z",
	  "operationName": "MICROSOFT.STORAGE/STORAGEACCOUNTS/LISTKEYS/ACTION",
	  "properties": {"eventCategory":"Administrative","statusCode":"OK","entity":"/subscriptions/aaa"}
	}`

	var ev AzureActivityLogEvent
	require.NoError(t, json.Unmarshal([]byte(record), &ev))

	raw, err := bson.Marshal(ev)
	require.NoError(t, err)

	var back AzureActivityLogEvent
	require.NoError(t, bson.Unmarshal(raw, &back))

	served, err := json.Marshal(back)
	require.NoError(t, err)
	assert.JSONEq(t, record, string(served))
}

// Every shape MarshalBSONValue can emit must be read back by
// UnmarshalBSONValue. An asymmetry here is silent data loss on the round trip:
// a top-level scalar bag was stored as a BSON double/boolean and read back as
// nil, with no error, because the reader only recognised documents and strings.
func TestAzureProperties_WriterAndReaderAgreeOnEveryShape(t *testing.T) {
	for _, tc := range []struct {
		bag    string
		stored bsontype.Type
	}{
		{`{"eventCategory":"Administrative"}`, bson.TypeEmbeddedDocument},
		{`[1,2,3]`, bson.TypeArray},
		{`"a string bag"`, bson.TypeString},
		{`42`, bson.TypeDouble},
		{`3.5`, bson.TypeDouble},
		{`true`, bson.TypeBoolean},
		{`null`, bson.TypeNull},
	} {
		t.Run(tc.bag, func(t *testing.T) {
			ev := AzureActivityLogEvent{Properties: AzureProperties(tc.bag)}

			raw, err := bson.Marshal(ev)
			require.NoError(t, err, "a valid JSON bag must never fail to encode")
			assert.Equal(t, tc.stored, bson.Raw(raw).Lookup("properties").Type)

			var back AzureActivityLogEvent
			require.NoError(t, bson.Unmarshal(raw, &back))
			if tc.bag == `null` {
				assert.Nil(t, back.Properties)
				return
			}
			assert.JSONEq(t, tc.bag, string(back.Properties),
				"the writer emitted a shape the reader dropped")
		})
	}
}

// A bag that is literally JSON null must not fail the encode. bson.MarshalValue(nil)
// errors, which would have failed the whole incident insert, not just the bag.
func TestAzureProperties_NullBagEncodes(t *testing.T) {
	_, err := bson.Marshal(AzureActivityLogEvent{Properties: AzureProperties(`null`)})
	require.NoError(t, err)
}

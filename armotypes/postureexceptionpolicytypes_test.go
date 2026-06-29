package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func objectSelectorFixture() *LabelSelector {
	return &LabelSelector{
		MatchLabels: map[string]string{"app": "nginx"},
		MatchExpressions: []LabelSelectorRequirement{
			{
				Key:      "tier",
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{"frontend", "backend"},
			},
			{
				Key:      "track",
				Operator: metav1.LabelSelectorOpNotIn,
				Values:   []string{"canary"},
			},
		},
	}
}

// TestPostureExceptionPolicyObjectSelectorRoundTrip verifies that a full label
// selector (matchLabels + matchExpressions with multiple operators) survives a
// JSON marshal/unmarshal round-trip on the ObjectSelector field.
func TestPostureExceptionPolicyObjectSelectorRoundTrip(t *testing.T) {
	policy := PostureExceptionPolicy{ObjectSelector: objectSelectorFixture()}

	raw, err := json.Marshal(policy)
	require.NoError(t, err)

	var decoded PostureExceptionPolicy
	require.NoError(t, json.Unmarshal(raw, &decoded))

	require.NotNil(t, decoded.ObjectSelector)
	assert.Equal(t, policy.ObjectSelector, decoded.ObjectSelector)
}

// TestPostureExceptionPolicyObjectSelectorBSONRoundTrip pins the selector contract for
// the bson tag, since the field is persisted via mongo-driver too. It also asserts the
// nested keys are stored as camelCase (matchLabels / matchExpressions) rather than the
// lowercased fallback the metav1 type would produce, so persisted documents match the
// JSON/API/CRD shape.
func TestPostureExceptionPolicyObjectSelectorBSONRoundTrip(t *testing.T) {
	policy := PostureExceptionPolicy{ObjectSelector: objectSelectorFixture()}

	raw, err := bson.Marshal(policy)
	require.NoError(t, err)

	// Inspect the raw document keys, not just the symmetric round-trip.
	var doc bson.M
	require.NoError(t, bson.Unmarshal(raw, &doc))
	selector, ok := doc["objectSelector"].(bson.M)
	require.True(t, ok, "objectSelector must persist as a nested document")
	assert.Contains(t, selector, "matchLabels")
	assert.Contains(t, selector, "matchExpressions")
	assert.NotContains(t, selector, "matchlabels")
	assert.NotContains(t, selector, "matchexpressions")

	var decoded PostureExceptionPolicy
	require.NoError(t, bson.Unmarshal(raw, &decoded))
	require.NotNil(t, decoded.ObjectSelector)
	assert.Equal(t, policy.ObjectSelector, decoded.ObjectSelector)
}

// TestLabelSelectorToMetaV1 verifies the conversion to metav1.LabelSelector used by the
// downstream exception comparator, including the nil-receiver guard.
func TestLabelSelectorToMetaV1(t *testing.T) {
	assert.Nil(t, (*LabelSelector)(nil).ToMetaV1())

	got := objectSelectorFixture().ToMetaV1()
	want := &metav1.LabelSelector{
		MatchLabels: map[string]string{"app": "nginx"},
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{Key: "tier", Operator: metav1.LabelSelectorOpIn, Values: []string{"frontend", "backend"}},
			{Key: "track", Operator: metav1.LabelSelectorOpNotIn, Values: []string{"canary"}},
		},
	}
	assert.Equal(t, want, got)
}

// TestPostureExceptionPolicyObjectSelectorOmitEmpty pins the contract that a nil
// ObjectSelector is absent from both serialized forms (no label constraint), so a
// selector-less exception is never persisted as a match-all selector.
func TestPostureExceptionPolicyObjectSelectorOmitEmpty(t *testing.T) {
	jsonRaw, err := json.Marshal(PostureExceptionPolicy{})
	require.NoError(t, err)
	var jsonFields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(jsonRaw, &jsonFields))
	assert.NotContains(t, jsonFields, "objectSelector")

	bsonRaw, err := bson.Marshal(PostureExceptionPolicy{})
	require.NoError(t, err)
	var bsonFields map[string]bson.RawValue
	require.NoError(t, bson.Unmarshal(bsonRaw, &bsonFields))
	assert.NotContains(t, bsonFields, "objectSelector")
}

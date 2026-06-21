package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func objectSelectorFixture() *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: map[string]string{"app": "nginx"},
		MatchExpressions: []metav1.LabelSelectorRequirement{
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

// TestPostureExceptionPolicyObjectSelectorBSONRoundTrip pins the same selector
// contract for the bson tag, since the field is persisted via mongo-driver too.
func TestPostureExceptionPolicyObjectSelectorBSONRoundTrip(t *testing.T) {
	policy := PostureExceptionPolicy{ObjectSelector: objectSelectorFixture()}

	raw, err := bson.Marshal(policy)
	require.NoError(t, err)

	var decoded PostureExceptionPolicy
	require.NoError(t, bson.Unmarshal(raw, &decoded))

	require.NotNil(t, decoded.ObjectSelector)
	assert.Equal(t, policy.ObjectSelector, decoded.ObjectSelector)
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

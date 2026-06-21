package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestPostureExceptionPolicyObjectSelectorRoundTrip verifies that a full label
// selector (matchLabels + matchExpressions) survives a JSON marshal/unmarshal
// round-trip on the ObjectSelector field.
func TestPostureExceptionPolicyObjectSelectorRoundTrip(t *testing.T) {
	policy := PostureExceptionPolicy{
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"app": "nginx"},
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "tier",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"frontend", "backend"},
				},
			},
		},
	}

	raw, err := json.Marshal(policy)
	require.NoError(t, err)

	var decoded PostureExceptionPolicy
	require.NoError(t, json.Unmarshal(raw, &decoded))

	require.NotNil(t, decoded.ObjectSelector)
	assert.Equal(t, policy.ObjectSelector, decoded.ObjectSelector)
}

// TestPostureExceptionPolicyObjectSelectorOmitEmpty pins the contract that a nil
// ObjectSelector is omitted from the serialized form (no label constraint), so a
// selector-less exception is never persisted as a match-all selector.
func TestPostureExceptionPolicyObjectSelectorOmitEmpty(t *testing.T) {
	raw, err := json.Marshal(PostureExceptionPolicy{})
	require.NoError(t, err)
	assert.NotContains(t, string(raw), "objectSelector")
}

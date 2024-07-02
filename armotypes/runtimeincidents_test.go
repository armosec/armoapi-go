package armotypes

import (
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/admission"
)

func TestAdmissionAlert(t *testing.T) {
	alert := AdmissionAlert{
		Kind:             schema.GroupVersionKind{},
		RequestNamespace: "test-namespace",
		ObjectName:       "test-object",
		Resource:         schema.GroupVersionResource{},
		Subresource:      "test-subresource",
		Operation:        admission.Create,
		Options:          nil,
		DryRun:           false,
		Object:           nil,
		OldObject:        nil,
		UserInfo:         nil,
	}

	assert.Equal(t, "test-namespace", alert.RequestNamespace)
	assert.Equal(t, "test-object", alert.ObjectName)
	assert.Equal(t, "test-subresource", alert.Subresource)
	assert.Equal(t, admission.Create, alert.Operation)
	assert.False(t, alert.DryRun)
}

// Test json marshalling and unmarshalling of AdmissionAlert
func TestAdmissionAlertJSON(t *testing.T) {
	alert := AdmissionAlert{
		Kind:             schema.GroupVersionKind{},
		RequestNamespace: "test-namespace",
		ObjectName:       "test-object",
		Resource:         schema.GroupVersionResource{},
		Subresource:      "test-subresource",
		Operation:        admission.Create,
		Options:          nil,
		DryRun:           false,
		Object:           nil,
		OldObject:        nil,
		UserInfo:         nil,
	}

	jsonData, err := json.Marshal(alert)
	assert.Nil(t, err)

	var alert2 AdmissionAlert
	err = json.Unmarshal(jsonData, &alert2)
	assert.Nil(t, err)

	assert.Equal(t, alert, alert2)
}

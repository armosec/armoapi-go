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
		DryRun:           false,
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

func TestFindProcessRecursive(t *testing.T) {
	tree := Process{
		PID:     1,
		Comm:    "p1",
		Cmdline: "init",
		ChildrenMap: map[CommPID]*Process{
			{Comm: "p2", PID: 2}: {PID: 2, Comm: "p2", Cmdline: "bash"},
			{Comm: "p3", PID: 3}: {PID: 3, Comm: "p3", Cmdline: "nginx",
				ChildrenMap: map[CommPID]*Process{{Comm: "p4", PID: 4}: {PID: 4, Comm: "p4", Cmdline: "worker"}}},
		},
	}

	tests := []struct {
		pid            uint32
		expectedName   string
		expectedToFind bool
	}{
		{pid: 1, expectedName: "p1", expectedToFind: true},
		{pid: 2, expectedName: "p2", expectedToFind: true},
		{pid: 3, expectedName: "p3", expectedToFind: true},
		{pid: 4, expectedName: "p4", expectedToFind: true},
		{pid: 5, expectedName: "", expectedToFind: false},
	}

	for _, tt := range tests {
		result := findProcessRecursive(&tree, tt.pid)

		if tt.expectedToFind && result == nil {
			t.Errorf("Expected to find process with PID %d, but it was not found", tt.pid)
		}
		if result != nil && result.Comm != tt.expectedName {
			t.Errorf("Expected process with PID %d, got %v", tt.pid, result)
		}
	}
}

package armotypes

import (
	"testing"

	"encoding/json"

	"github.com/goradd/maps"
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
	// First initialize the process struct
	p := Process{
		PID:     1,
		Comm:    "p1",
		Cmdline: "cmd",
	}

	// Create the child processes
	p2 := &Process{PID: 2, Comm: "p2", Cmdline: "bash"}
	p4 := &Process{PID: 4, Comm: "p4", Cmdline: "worker"}
	p3 := &Process{PID: 3, Comm: "p3", Cmdline: "nginx"}

	// Create SafeMap instances and populate them
	p3ChildrenMap := maps.NewSafeMap[CommPID, *Process]()
	p3ChildrenMap.Set(CommPID{Comm: "p4", PID: 4}, p4)
	p3.ChildrenMap = *p3ChildrenMap

	pChildrenMap := maps.NewSafeMap[CommPID, *Process]()
	pChildrenMap.Set(CommPID{Comm: "p2", PID: 2}, p2)
	pChildrenMap.Set(CommPID{Comm: "p3", PID: 3}, p3)
	p.ChildrenMap = *pChildrenMap

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
		result := findProcessRecursive(&p, tt.pid)

		if tt.expectedToFind && result == nil {
			t.Errorf("Expected to find process with PID %d, but it was not found", tt.pid)
		}
		if result != nil && result.Comm != tt.expectedName {
			t.Errorf("Expected process with PID %d, got %v", tt.pid, result)
		}
	}
}

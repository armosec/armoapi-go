package armotypes

import (
	"reflect"
	"testing"
)

func TestIsKDRMonitored(t *testing.T) {
	tests := []struct {
		name      string
		nodeAgent bool
		runtime   bool
		expected  bool
	}{
		{"Both true", true, true, true},
		{"NodeAgent false", false, true, false},
		{"Runtime false", true, false, false},
		{"Both false", false, false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			np := NodeProfile{
				NodeAgentRunning:        test.nodeAgent,
				RuntimeDetectionEnabled: test.runtime,
			}
			if result := np.IsKDRMonitored(); result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestGetMonitoredNamespaces(t *testing.T) {
	np := NodeProfile{
		NodeAgentRunning:        true,
		RuntimeDetectionEnabled: true,
		PodStatuses: []PodStatus{
			{Namespace: "default", IsKDRMonitored: true, Phase: "Running"},
			{Namespace: "kube-system", IsKDRMonitored: true},
			{Namespace: "default", IsKDRMonitored: true}, // Duplicate
			{Namespace: "test", IsKDRMonitored: false},
		},
	}

	expected := []string{"default"}
	namespaces := np.GetMonitoredNamespaces()
	if !reflect.DeepEqual(namespaces, expected) {
		t.Errorf("Expected %v, got %v", expected, namespaces)
	}
}

func TestGetMonitoredPods(t *testing.T) {
	np := NodeProfile{
		NodeAgentRunning:        true,
		RuntimeDetectionEnabled: true,
		PodStatuses: []PodStatus{
			{Name: "pod1", IsKDRMonitored: true, Phase: "Running"},
			{Name: "pod2", IsKDRMonitored: false, Phase: "Running"},
		},
	}

	expected := []PodStatus{{Name: "pod1", IsKDRMonitored: true, Phase: "Running"}}
	monitoredPods := np.GetMonitoredPods()
	if !reflect.DeepEqual(monitoredPods, expected) {
		t.Errorf("Expected %v, got %v", expected, monitoredPods)
	}
}

func TestCountMonitoredNamespaces(t *testing.T) {
	np := NodeProfile{
		NodeAgentRunning:        true,
		RuntimeDetectionEnabled: true,
		PodStatuses: []PodStatus{
			{Namespace: "default", IsKDRMonitored: true, Phase: "Running"},
			{Namespace: "kube-system", IsKDRMonitored: true, Phase: "Running"},
			{Namespace: "default", IsKDRMonitored: true, Phase: "Running"},
			{Namespace: "test", IsKDRMonitored: false, Phase: "Pending"},
		},
	}

	expectedCount := 2
	count := np.CountMonitoredNamespaces()
	if count != expectedCount {
		t.Errorf("Expected %d, got %d", expectedCount, count)
	}
}

func TestCountMonitoredPods(t *testing.T) {
	np := NodeProfile{
		NodeAgentRunning:        true,
		RuntimeDetectionEnabled: true,
		PodStatuses: []PodStatus{
			{IsKDRMonitored: true, Phase: "Running"},
			{IsKDRMonitored: true, Phase: "Pending"},
			{IsKDRMonitored: false, Phase: "Running"},
		},
	}

	expectedCount := 1
	count := np.CountMonitoredPods()
	if count != expectedCount {
		t.Errorf("Expected %d, got %d", expectedCount, count)
	}
}

func TestCountMonitoredContainers(t *testing.T) {
	np := NodeProfile{
		NodeAgentRunning:        true,
		RuntimeDetectionEnabled: true,
		PodStatuses: []PodStatus{
			{
				Name:           "pod1",
				IsKDRMonitored: true,
				Phase:          "Running",
				Containers: []PodContainer{
					{Name: "container1", IsKDRMonitored: true},
					{Name: "container2", IsKDRMonitored: false},
				},
			},
			{
				Name:           "pod2",
				IsKDRMonitored: true,
				Containers: []PodContainer{
					{Name: "container3", IsKDRMonitored: true},
				},
			},
			{
				Name:           "pod3",
				IsKDRMonitored: false, // This pod is not monitored
				Containers: []PodContainer{
					{Name: "container4", IsKDRMonitored: true},
				},
			},
			{
				Name:           "pod4",
				IsKDRMonitored: true,
				Containers:     []PodContainer{}, // Monitored pod but no containers
			},
		},
	}

	expectedCount := 1
	count := np.CountMonitoredContainers()

	if count != expectedCount {
		t.Errorf("Expected %d, got %d", expectedCount, count)
	}
}

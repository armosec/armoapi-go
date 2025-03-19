package identifiers

import (
	"testing"
)

func TestCalcHashFNV(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestCalcHashFNV",
			args: args{
				id: "123",
			},
			want: "5003431119771845851",
		},
		{
			name: "TestCalcHashFNV-1",
			args: args{
				id: "1234",
			},
			want: "2282126479029740061",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcHashFNV(tt.args.id); got != tt.want {
				t.Errorf("CalcHashFNV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertResourceIDToResourceHashFNV(t *testing.T) {
	tests := []struct {
		customerGUID string
		clusterName  string
		resourceID   string
		expectedHash string // The expected string should be a valid hash; you would need to calculate this in advance
	}{
		{
			customerGUID: "6972d19c-16a1-41f9-bcda-ec1638f885d2",
			clusterName:  "kind-test-ac8",
			resourceID:   "apps/v1/default/Deployment/alpine-deployment",
			expectedHash: "7033989581143626640",
		},
		{
			customerGUID: "6972d19c-16a1-41f9-bcda-ec1638f885d2",
			clusterName:  "kind-test-ac8",
			resourceID:   "/v1/default/Service/httpbin",
			expectedHash: "1233188593520133341",
		},
		{
			customerGUID: "1e3a88bf-92ce-44f8-914e-cbe71830d566",
			clusterName:  "kind-test-ac8",
			resourceID:   "rbac.authorization.k8s.io/v1//ClusterRoleBinding/kubescape",
			expectedHash: "7404923029724202573",
		},
	}

	for _, test := range tests {
		actual := ConvertResourceIDToResourceHashFNV(test.customerGUID, test.clusterName, test.resourceID)
		if actual != test.expectedHash {
			t.Errorf("ConvertResourceIDToResourceHashFNV(%s) = %s; want %s", test.resourceID, actual, test.expectedHash)
		}
	}
}

func TestCalcContainerHashFNV(t *testing.T) {
	tests := []struct {
		name          string
		customerGUID  string
		cluster       string
		podName       string
		containerName string
		namespace     string
		wantHash      string // Expected hash value
	}{
		{
			name:          "Test case 1",
			customerGUID:  "customer123",
			cluster:       "clusterA",
			podName:       "podX",
			containerName: "containerY",
			namespace:     "namespaceZ",
			wantHash:      "8445573203380894384",
		},
		{
			name:          "Test case 2 - Different input",
			customerGUID:  "cust456",
			cluster:       "clusB",
			podName:       "podY",
			containerName: "contZ",
			namespace:     "nsX",
			wantHash:      "5926078173930249943",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcContainerHashFNV(tt.customerGUID, tt.cluster, tt.podName, tt.containerName, tt.namespace)
			if got != tt.wantHash {
				t.Errorf("CalcContainerHashFNV() = %v, want %v", got, tt.wantHash)
			}
		})
	}
}

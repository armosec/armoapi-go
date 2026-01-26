package armotypes

import "testing"

func TestRuntimeAlertValidateRequiredFieldsByPlatform(t *testing.T) {
	tests := []struct {
		name    string
		alert   RuntimeAlert
		wantErr bool
	}{
		{
			name: "ruleID required for ec2",
			alert: RuntimeAlert{
				AlertSourcePlatform: AlertSourcePlatformHost,
			},
			wantErr: true,
		},
		{
			name: "k8s missing workload fields",
			alert: RuntimeAlert{
				AlertSourcePlatform: AlertSourcePlatformK8s,
				RuleID:              "R0001",
				RuntimeAlertK8sDetails: RuntimeAlertK8sDetails{
					PodNamespace:  "ns",
					PodName:       "pod",
					ContainerName: "container",
				},
			},
			wantErr: true,
		},
		{
			name: "k8s ok",
			alert: RuntimeAlert{
				AlertSourcePlatform: AlertSourcePlatformK8s,
				RuleID:              "R0001",
				RuntimeAlertK8sDetails: RuntimeAlertK8sDetails{
					WorkloadNamespace: "wns",
					WorkloadKind:      "Deployment",
					WorkloadName:      "app",
					PodNamespace:      "ns",
					PodName:           "pod",
					ContainerName:     "container",
				},
			},
			wantErr: false,
		},
		{
			name: "cloud skips platform validation",
			alert: RuntimeAlert{
				AlertSourcePlatform: AlertSourcePlatformCloud,
				RuleID:              "R0001",
			},
			wantErr: false,
		},
		{
			name: "unknown skips platform validation",
			alert: RuntimeAlert{
				AlertSourcePlatform: AlertSourcePlatformUnknown,
				RuleID:              "R0001",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.alert.Validate()
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

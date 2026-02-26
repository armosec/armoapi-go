package armotypes

import (
	"testing"
	"time"
)

func TestGetAlertSourcePlatform(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		alert    RuntimeAlert
		expected AlertSourcePlatform
	}{
		{
			name: "Fargate task returns PtraceAgent",
			alert: RuntimeAlert{
				BaseRuntimeAlert: BaseRuntimeAlert{
					Timestamp: now,
				},
				RuntimeAlertECSDetails: RuntimeAlertECSDetails{
					TaskARN:    "arn:aws:ecs:us-east-1:123456789:task/my-cluster/12345678",
					ClusterARN: "arn:aws:ecs:us-east-1:123456789:cluster/my-cluster",
					LaunchType: "FARGATE",
				},
			},
			expected: AlertSourcePlatformPtraceAgent,
		},
		{
			name: "EC2 task returns ECSAgent",
			alert: RuntimeAlert{
				BaseRuntimeAlert: BaseRuntimeAlert{
					Timestamp: now,
				},
				RuntimeAlertECSDetails: RuntimeAlertECSDetails{
					TaskARN:    "arn:aws:ecs:us-east-1:123456789:task/my-cluster/12345678",
					ClusterARN: "arn:aws:ecs:us-east-1:123456789:cluster/my-cluster",
					LaunchType: "EC2",
				},
			},
			expected: AlertSourcePlatformECSAgent,
		},
		{
			name: "TaskARN without LaunchType returns ECSAgent (backward compatibility)",
			alert: RuntimeAlert{
				BaseRuntimeAlert: BaseRuntimeAlert{
					Timestamp: now,
				},
				RuntimeAlertECSDetails: RuntimeAlertECSDetails{
					TaskARN:    "arn:aws:ecs:us-east-1:123456789:task/my-cluster/12345678",
					ClusterARN: "arn:aws:ecs:us-east-1:123456789:cluster/my-cluster",
				},
			},
			expected: AlertSourcePlatformECSAgent,
		},
		{
			name: "ClusterARN without LaunchType returns ECSAgent (backward compatibility)",
			alert: RuntimeAlert{
				BaseRuntimeAlert: BaseRuntimeAlert{
					Timestamp: now,
				},
				RuntimeAlertECSDetails: RuntimeAlertECSDetails{
					ClusterARN: "arn:aws:ecs:us-east-1:123456789:cluster/my-cluster",
				},
			},
			expected: AlertSourcePlatformECSAgent,
		},
		{
			name: "Empty ECS data returns HostAgent",
			alert: RuntimeAlert{
				BaseRuntimeAlert: BaseRuntimeAlert{
					Timestamp: now,
				},
			},
			expected: AlertSourcePlatformHostAgent,
		},
		{
			name: "CDR alert returns Cloud",
			alert: RuntimeAlert{
				BaseRuntimeAlert: BaseRuntimeAlert{
					Timestamp: now,
				},
				AlertType: AlertTypeCdr,
			},
			expected: AlertSourcePlatformCloud,
		},
		{
			name: "PodName present returns K8sAgent",
			alert: RuntimeAlert{
				BaseRuntimeAlert: BaseRuntimeAlert{
					Timestamp: now,
				},
				RuntimeAlertK8sDetails: RuntimeAlertK8sDetails{
					PodName: "my-pod",
				},
			},
			expected: AlertSourcePlatformK8sAgent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.alert.GetAlertSourcePlatform()
			if got != tt.expected {
				t.Errorf("GetAlertSourcePlatform() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRuntimeAlertValidateRequiredFieldsByPlatform(t *testing.T) {
	tests := []struct {
		name    string
		alert   RuntimeAlert
		wantErr bool
	}{
		{
			name: "ruleID required for host-agent",
			alert: RuntimeAlert{
				AlertSourcePlatform: AlertSourcePlatformHostAgent,
			},
			wantErr: true,
		},
		{
			name: "k8s missing workload fields",
			alert: RuntimeAlert{
				AlertSourcePlatform: AlertSourcePlatformK8sAgent,
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
				AlertSourcePlatform: AlertSourcePlatformK8sAgent,
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

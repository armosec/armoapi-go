package armotypes

import (
	"testing"

	"github.com/armosec/armoapi-go/identifiers"
	"github.com/stretchr/testify/assert"
)

func TestTrimSpacesAroundCommas(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "multi-value list with whitespace around separators",
			input: "val1, val2 , val3",
			want:  "val1,val2,val3",
		},
		{
			name:  "single value with escaped comma preserves surrounding whitespace",
			input: `(KHTML\, like Gecko)`,
			want:  `(KHTML\, like Gecko)`,
		},
		{
			name:  "real-world userAgent with escaped comma",
			input: `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/147.0.0.0 Safari/537.36`,
			want:  `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/147.0.0.0 Safari/537.36`,
		},
		{
			name:  "mixed: multi-value list with one item containing escaped comma",
			input: `val1, (KHTML\, like Gecko) , val3`,
			want:  `val1,(KHTML\, like Gecko),val3`,
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "single value no separators no whitespace",
			input: "val1",
			want:  "val1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trimSpacesAroundCommas(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEscapeV2ListOperatorSeparator(t *testing.T) {
	assert.Equal(t, `no separators`, escapeV2ListOperatorSeparator(`no separators`))
	assert.Equal(t, `a\|b`, escapeV2ListOperatorSeparator(`a|b`))
	assert.Equal(t, `a\|b\|c`, escapeV2ListOperatorSeparator(`a|b|c`))
	assert.Equal(t, ``, escapeV2ListOperatorSeparator(``))
}

func TestGetRuntimeIncidentsRequestFilterFromExceptionPolicy(t *testing.T) {
	tests := []struct {
		name   string
		policy BaseExceptionPolicy
		want   []map[string]string
	}{
		{
			name: "no PolicyIDs returns nil",
			policy: BaseExceptionPolicy{
				Resources: []identifiers.PortalDesignator{
					{Attributes: map[string]string{identifiers.AttributeCluster: "c1"}},
				},
			},
			want: nil,
		},
		{
			name: "no resources returns nil filters slice",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I001"},
			},
			want: nil,
		},
		{
			name: "resource with nil Attributes is skipped",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I001"},
				Resources: []identifiers.PortalDesignator{
					{Attributes: nil},
				},
			},
			want: nil,
		},
		{
			name: "K8s designator produces filter on designators.attributes.* only",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I001"},
				Resources: []identifiers.PortalDesignator{
					{
						Attributes: map[string]string{
							identifiers.AttributeCluster:   "my-cluster",
							identifiers.AttributeNamespace: "default",
							identifiers.AttributeKind:      "Deployment",
							identifiers.AttributeName:      "frontend",
						},
					},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                   "I001",
					"status":                           "Open",
					"designators.attributes.cluster":   "my-cluster",
					"designators.attributes.namespace": "default",
					"designators.attributes.kind":      "Deployment",
					"designators.attributes.name":      "frontend",
				},
			},
		},
		{
			name: "cloud designator filters cloudProvider/accountID via cloudMetadata.*",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I007"},
				Resources: []identifiers.PortalDesignator{
					{
						Attributes: map[string]string{
							identifiers.AttributeCloudProvider:  "aws",
							identifiers.AttributeCloudAccountID: "123456789012",
							identifiers.AttributeRegion:         "us-east-1",
							identifiers.AttributeInstanceId:     "i-abc",
							identifiers.AttributeHostType:       "ecs-ec2",
						},
					},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                    "I007",
					"status":                            "Open",
					"cloudMetadata.provider":            "aws",
					"cloudMetadata.account_id":          "123456789012",
					"designators.attributes.region":     "us-east-1",
					"designators.attributes.instanceId": "i-abc",
					"designators.attributes.hostType":   "ecs-ec2",
				},
			},
		},
		{
			name: "GlobalRegex (*/*) values are omitted",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I001"},
				Resources: []identifiers.PortalDesignator{
					{
						Attributes: map[string]string{
							identifiers.AttributeCloudProvider:  GlobalRegex,
							identifiers.AttributeCloudAccountID: GlobalRegex,
							identifiers.AttributeRegion:         "us-east-1",
						},
					},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                "I001",
					"status":                        "Open",
					"designators.attributes.region": "us-east-1",
				},
			},
		},
		{
			name: "in operator: multi-value list canonicalized; trailing/leading whitespace trimmed",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I001"},
				Resources: []identifiers.PortalDesignator{
					{
						Attributes: map[string]string{identifiers.AttributeCluster: "c1"},
					},
				},
				AdvancedScopes: []AdvancedScopeEntity{
					{Entity: "process.name", Operator: "in", Values: "chrome , python ,firefox , java "},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                 "I001",
					"status":                         "Open",
					"designators.attributes.cluster": "c1",
					"identifiers.process.name":       "chrome,python,firefox,java",
				},
			},
		},
		{
			name: "in operator: escaped comma keeps surrounding whitespace inside the item",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I001"},
				Resources: []identifiers.PortalDesignator{
					{Attributes: map[string]string{identifiers.AttributeCluster: "c1"}},
				},
				AdvancedScopes: []AdvancedScopeEntity{
					{
						Entity:   "sourceInformation.userAgent",
						Operator: "in",
						Values:   `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/147.0.0.0 Safari/537.36`,
					},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                          "I001",
					"status":                                  "Open",
					"designators.attributes.cluster":          "c1",
					"identifiers.sourceInformation.userAgent": `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/147.0.0.0 Safari/537.36`,
				},
			},
		},
		{
			name: "contains operator: value gets |like suffix",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I001"},
				Resources: []identifiers.PortalDesignator{
					{Attributes: map[string]string{identifiers.AttributeCluster: "c1"}},
				},
				AdvancedScopes: []AdvancedScopeEntity{
					{Entity: "file.name", Operator: "contains", Values: ".exe"},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                 "I001",
					"status":                         "Open",
					"designators.attributes.cluster": "c1",
					"identifiers.file.name":          ".exe|like",
				},
			},
		},
		{
			name: "raw value containing | gets escaped before processing",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I001"},
				Resources: []identifiers.PortalDesignator{
					{Attributes: map[string]string{identifiers.AttributeCluster: "c1"}},
				},
				AdvancedScopes: []AdvancedScopeEntity{
					{Entity: "process.commandLine", Operator: "in", Values: "a|b"},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                  "I001",
					"status":                          "Open",
					"designators.attributes.cluster":  "c1",
					"identifiers.process.commandLine": `a\|b`,
				},
			},
		},
		{
			name: "K8s and cloud designators combined on the same resource",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I003"},
				Resources: []identifiers.PortalDesignator{
					{
						Attributes: map[string]string{
							identifiers.AttributeCluster:        "my-cluster",
							identifiers.AttributeNamespace:      "default",
							identifiers.AttributeCloudProvider:  "aws",
							identifiers.AttributeCloudAccountID: "111111111111",
							identifiers.AttributeRegion:         "eu-west-1",
						},
					},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                   "I003",
					"status":                           "Open",
					"designators.attributes.cluster":   "my-cluster",
					"designators.attributes.namespace": "default",
					"designators.attributes.region":    "eu-west-1",
					"cloudMetadata.provider":           "aws",
					"cloudMetadata.account_id":         "111111111111",
				},
			},
		},
		{
			name: "cloud designators + advanced scopes merge into a single filter",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I007"},
				Resources: []identifiers.PortalDesignator{
					{
						Attributes: map[string]string{
							identifiers.AttributeCloudProvider:  "aws",
							identifiers.AttributeCloudAccountID: "015253967648",
							identifiers.AttributeRegion:         "us-east-1",
							identifiers.AttributeInstanceId:     "i-09c616188a05401a6",
							identifiers.AttributeHostType:       "ecs-ec2",
						},
					},
				},
				AdvancedScopes: []AdvancedScopeEntity{
					{Entity: "process.name", Operator: "in", Values: "cat"},
					{Entity: "file.name", Operator: "in", Values: "shadow"},
					{Entity: "file.directory", Operator: "in", Values: "/etc"},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                    "I007",
					"status":                            "Open",
					"cloudMetadata.provider":            "aws",
					"cloudMetadata.account_id":          "015253967648",
					"designators.attributes.region":     "us-east-1",
					"designators.attributes.instanceId": "i-09c616188a05401a6",
					"designators.attributes.hostType":   "ecs-ec2",
					"identifiers.process.name":          "cat",
					"identifiers.file.name":             "shadow",
					"identifiers.file.directory":        "/etc",
				},
			},
		},
		{
			name: "all scopes GlobalRegex -> only seed fields remain",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I003"},
				Resources: []identifiers.PortalDesignator{
					{
						Attributes: map[string]string{
							identifiers.AttributeCluster:   GlobalRegex,
							identifiers.AttributeNamespace: GlobalRegex,
							identifiers.AttributeName:      GlobalRegex,
						},
					},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID": "I003",
					"status":         "Open",
				},
			},
		},
		{
			name: "advanced scope with empty Operator passes through (no normalization, no like suffix)",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I007"},
				Resources: []identifiers.PortalDesignator{
					{Attributes: map[string]string{identifiers.AttributeNamespace: "armo-node-agent"}},
				},
				AdvancedScopes: []AdvancedScopeEntity{
					{Entity: "process.name", Values: "cat"},
					{Entity: "file.name", Values: "shadow"},
					{Entity: "file.directory", Values: "/etc"},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                   "I007",
					"status":                           "Open",
					"designators.attributes.namespace": "armo-node-agent",
					"identifiers.process.name":         "cat",
					"identifiers.file.name":            "shadow",
					"identifiers.file.directory":       "/etc",
				},
			},
		},
		{
			name: "multiple advanced scopes with mixed operators on the same policy (in/contains)",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I001"},
				Resources: []identifiers.PortalDesignator{
					{
						Attributes: map[string]string{
							identifiers.AttributeCluster:   "my-cluster",
							identifiers.AttributeNamespace: "default",
							identifiers.AttributeKind:      "Deployment",
						},
					},
				},
				AdvancedScopes: []AdvancedScopeEntity{
					{Entity: "process.name", Operator: "in", Values: "chrome ,python , firefox,java "},
					{Entity: "process.commandLine", Operator: "in", Values: "wait-shutdown "},
					{Entity: "network.dstPort", Operator: "in", Values: "80,81,443"},
					{Entity: "file.name", Operator: "contains", Values: ".exe"},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                   "I001",
					"status":                           "Open",
					"designators.attributes.cluster":   "my-cluster",
					"designators.attributes.namespace": "default",
					"designators.attributes.kind":      "Deployment",
					"identifiers.process.name":         "chrome,python,firefox,java",
					"identifiers.process.commandLine":  "wait-shutdown",
					"identifiers.network.dstPort":      "80,81,443",
					"identifiers.file.name":            ".exe|like",
				},
			},
		},
		{
			name: "multiple resources with mixed GlobalRegex per attribute",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I003"},
				Resources: []identifiers.PortalDesignator{
					{
						Attributes: map[string]string{
							identifiers.AttributeCluster:   GlobalRegex,
							identifiers.AttributeNamespace: "ns1",
							identifiers.AttributeName:      "name1",
						},
					},
					{
						Attributes: map[string]string{
							identifiers.AttributeCluster:   "cluster1",
							identifiers.AttributeNamespace: GlobalRegex,
							identifiers.AttributeName:      GlobalRegex,
						},
					},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                   "I003",
					"status":                           "Open",
					"designators.attributes.namespace": "ns1",
					"designators.attributes.name":      "name1",
				},
				{
					"incidentTypeID":                 "I003",
					"status":                         "Open",
					"designators.attributes.cluster": "cluster1",
				},
			},
		},
		{
			name: "multiple resources -> one filter per resource, advanced scopes merged into each",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I001"},
				Resources: []identifiers.PortalDesignator{
					{Attributes: map[string]string{identifiers.AttributeCluster: "c1"}},
					{Attributes: map[string]string{identifiers.AttributeCluster: "c2"}},
				},
				AdvancedScopes: []AdvancedScopeEntity{
					{Entity: "process.name", Operator: "in", Values: "python"},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                 "I001",
					"status":                         "Open",
					"designators.attributes.cluster": "c1",
					"identifiers.process.name":       "python",
				},
				{
					"incidentTypeID":                 "I001",
					"status":                         "Open",
					"designators.attributes.cluster": "c2",
					"identifiers.process.name":       "python",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRuntimeIncidentsRequestFilterFromExceptionPolicy(tt.policy)
			assert.Equal(t, tt.want, got)
		})
	}
}

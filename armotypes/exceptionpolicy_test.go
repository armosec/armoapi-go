package armotypes

import (
	"testing"

	"github.com/armosec/armoapi-go/armotypes/common"
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
			name: "azure account id is matched case-insensitively (anchored regex + ignorecase)",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I008"},
				Resources: []identifiers.PortalDesignator{
					{
						Attributes: map[string]string{
							identifiers.AttributeCloudProvider:  "azure",
							identifiers.AttributeCloudAccountID: "94F9BBD2-E5A9-4920-AB09-577FF2AFDA21",
							identifiers.AttributeRegion:         "eastus",
						},
					},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                "I008",
					"status":                        "Open",
					"cloudMetadata.provider":        "azure",
					"cloudMetadata.account_id":      "^94F9BBD2-E5A9-4920-AB09-577FF2AFDA21$|regex&ignorecase",
					"designators.attributes.region": "eastus",
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
					"incidentTypeID":                 "I001",
					"status":                         "Open",
					"designators.attributes.cluster": "c1",
					"identifiers.sourceInformation.userAgent": `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/147.0.0.0 Safari/537.36`,
				},
			},
		},
		{
			// The Mongo half of risk acceptance: this filter drives the
			// "how many incidents would this affect" count and the retroactive
			// resolve of open incidents when an exception is created. A hash scope
			// has to reach it, and the entity name has to become the identifier
			// path verbatim.
			name: "file hash scopes become identifier filters",
			policy: BaseExceptionPolicy{
				PolicyIDs: []string{"I002"},
				Resources: []identifiers.PortalDesignator{
					{Attributes: map[string]string{identifiers.AttributeCluster: "c1"}},
				},
				AdvancedScopes: []AdvancedScopeEntity{
					{Entity: "file.sha256", Operator: "in", Values: "9f2b1e1d3c4a5b6c7d8e9f0a1b2c3d4e5f60718293a4b5c6d7e8f9a0b1c2d3e4"},
					{Entity: "file.md5", Operator: "in", Values: "5d41402abc4b2a76b9719d911017c592"},
					{Entity: "file.sha1", Operator: "contains", Values: "aaf4c61d"},
				},
			},
			want: []map[string]string{
				{
					"incidentTypeID":                 "I002",
					"status":                         "Open",
					"designators.attributes.cluster": "c1",
					"identifiers.file.sha256":        "9f2b1e1d3c4a5b6c7d8e9f0a1b2c3d4e5f60718293a4b5c6d7e8f9a0b1c2d3e4",
					"identifiers.file.md5":           "5d41402abc4b2a76b9719d911017c592",
					"identifiers.file.sha1":          "aaf4c61d|like",
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

// The fold has to reach every hash entity and no other entity. File paths, process
// names and command lines are case-significant on Linux, so folding one would change
// what the analyst's rule matches.
func TestNormalizeHashScopeValues(t *testing.T) {
	scopes := []AdvancedScopeEntity{
		{Entity: "file.md5", Operator: "in", Values: "5D41402ABC4B2A76B9719D911017C592"},
		{Entity: "file.sha1", Operator: "in", Values: "AAF4C61DDCC5E8A2DABEDE0F3B482CD9AEA9434D"},
		{Entity: "file.sha256", Operator: "contains", Values: "9F2B1E1D"},
		{Entity: "file.name", Operator: "in", Values: "Mirai.ELF"},
		{Entity: "file.directory", Operator: "in", Values: "/Tmp/Payloads"},
		{Entity: "process.commandLine", Operator: "contains", Values: "--Flag=Value"},
	}

	NormalizeHashScopeValues(scopes)

	assert.Equal(t, "5d41402abc4b2a76b9719d911017c592", scopes[0].Values)
	assert.Equal(t, "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d", scopes[1].Values)
	assert.Equal(t, "9f2b1e1d", scopes[2].Values)
	assert.Equal(t, "Mirai.ELF", scopes[3].Values)
	assert.Equal(t, "/Tmp/Payloads", scopes[4].Values)
	assert.Equal(t, "--Flag=Value", scopes[5].Values)
}

// A hash entity that Flatten() can emit but the fold does not know about would be
// stored unfolded and never match. Derive the list from Flatten() so the two cannot
// drift.
func TestHashScopeEntitiesCoverEveryHashFlattenEmits(t *testing.T) {
	populated := &common.Identifiers{
		File: &common.FileEntity{
			MD5:    "5d41402abc4b2a76b9719d911017c592",
			SHA1:   "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
			SHA256: "9f2b1e1d3c4a5b6c7d8e9f0a1b2c3d4e5f60718293a4b5c6d7e8f9a0b1c2d3e4",
		},
	}

	for field := range populated.Flatten() {
		assert.True(t, IsHashScopeEntity(field), "%s is a hash identifier but the case fold does not cover it", field)
	}
	assert.Len(t, hashScopeEntities, 3)
}

func TestNormalizeHashScopeValuesHandlesNoScopes(t *testing.T) {
	assert.NotPanics(t, func() { NormalizeHashScopeValues(nil) })
	assert.NotPanics(t, func() { NormalizeHashScopeValues([]AdvancedScopeEntity{}) })
}

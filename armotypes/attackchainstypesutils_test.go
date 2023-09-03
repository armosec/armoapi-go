package armotypes

import (
	"encoding/json"
	"reflect"
	"testing"
)

var (
	mockNode1 = `{
		"name": "Workload Exposure",
		"description": "",
		"controlIDs": [
			"C-0256",
			"C-0044"
		],
		"nextNodes": [
			{
				"name": "Vulnerable Image",
				"description": "",
				"vulnerabilities": [
					{
						"imageScanID": "5837849405239139707",
						"names": [
							"CVE-2020-27846",
							"CVE-2022-26148",
							"CVE-2022-48174",
							"GHSA-mpv3-g8m3-3fjc"
						]
					}
				],
				"nextNodes": [
					{
						"name": "Credential access",
						"description": "",
						"controlIDs": [
							"C-0261"
						]
					},
					{
						"name": "Persistence",
						"description": "",
						"controlIDs": [
							"C-0017"
						]
					},
					{
						"name": "Network",
						"description": "",
						"controlIDs": [
							"C-0260"
						]
					}
				]
			}
		]
	}`
	mockNode2 = `{
		"name": "Workload Exposure",
		"description": "",
		"controlIDs": [
			"C-0256"
		],
		"nextNodes": [
			{
				"name": "Service Destruction",
				"description": "",
				"controlIDs": [
					"C-0009"
				]
			}
		]
	}`
)

func TestGetControlIDsFromAllNodes(t *testing.T) {
	testCases := []struct {
		name           string
		node           string
		expectedResult []string
	}{
		{
			name:           "workload-external-track",
			node:           mockNode1,
			expectedResult: []string{"C-0256", "C-0044", "C-0261", "C-0017", "C-0260"},
		},
		{
			name:           "service-destruction",
			node:           mockNode2,
			expectedResult: []string{"C-0256", "C-0009"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node := &AttackChainNode{}
			err := json.Unmarshal([]byte(tc.node), node)
			if err != nil {
				t.Fatalf("failed to unmarshal node: %v", err)
			}
			result := node.getControlIDsFromAllNodes([]string{})
			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Fatalf("expected: %v, got: %v", tc.expectedResult, result)
			}
		})
	}
}

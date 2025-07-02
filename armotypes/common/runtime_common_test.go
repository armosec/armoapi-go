package common

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

// Helper: recursively flatten JSON map to dot notation
func groupJSONKeysWithDot(prefix string, raw_json map[string]any, with_dot_json map[string]string) {
	for key, value := range raw_json {
		new_key := key
		if prefix != "" {
			new_key = prefix + "." + key
		}

		switch vv := value.(type) {
		case map[string]any:
			groupJSONKeysWithDot(new_key, vv, with_dot_json)
		case float64:
			with_dot_json[new_key] = strconv.Itoa(int(vv))
		case string:
			if vv != "" {
				with_dot_json[new_key] = vv
			}
		default:
		}
	}
}

func TestIdentifiersFlatten(t *testing.T) {
	tests := []struct {
		Name        string
		Identifiers *Identifiers
	}{
		{
			Name: "Check identifiers with all fields",
			Identifiers: &Identifiers{
				Process: &ProcessEntity{
					Name:        "python",
					CommandLine: "/usr/bin/python",
				},
				File: &FileEntity{
					Name:      "file.txt",
					Directory: "/tmp",
				},
				Dns: &DnsEntity{
					Domain: "example.com",
				},
				Network: &NetworkEntity{
					DstIP:    "1.1.1.1",
					DstPort:  8080,
					Protocol: "TCP",
				},
				Http: &HttpEntity{
					Method:    "POST",
					Domain:    "api.example.com",
					UserAgent: "curl/7.68.0",
					Endpoint:  "/api/v1/resource",
					Payload:   "data",
				},
				CloudAPI: &CloudAPIEntity{
					Service:  "AWS",
					APICall:  "ListBuckets",
					Resource: "bucket-name",
					User:     "admin",
				},
			},
		},
		{
			Name: "Check identifiers with missing fields",
			Identifiers: &Identifiers{
				Process: &ProcessEntity{
					Name: "python",
				},
				File: &FileEntity{
					Name:      "file.txt",
					Directory: "/tmp",
				},
				Dns: &DnsEntity{
					Domain: "example.com",
				},
				Http: &HttpEntity{
					Method:  "POST",
					Payload: "data",
				},
			},
		},
	}

	for _, test := range tests {
		flatten_identifiers := test.Identifiers.Flatten()

		// Marshal the entire struct to JSON
		raw, err := json.Marshal(test.Identifiers)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}

		// Unmarshal back to generic map
		var jsonMap map[string]any
		if err := json.Unmarshal(raw, &jsonMap); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		// Recursively group the JSON keys with dot notation
		want := map[string]string{}
		groupJSONKeysWithDot("", jsonMap, want)

		diff := cmp.Diff(flatten_identifiers, want)
		assert.Empty(t, diff, "expected to have no diff")

	}
}

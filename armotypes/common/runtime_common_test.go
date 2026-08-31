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
					MD5:       "d41d8cd98f00b204e9800998ecf8427e",
					SHA1:      "da39a3ee5e6b4b0d3255bfef95601890afd80709",
					SHA256:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				},
				Dns: &DnsEntity{
					Domain: "example.com",
				},
				Network: &NetworkEntity{
					DstIP:    "1.1.1.1",
					DstPort:  8080,
					Protocol: "TCP",
					SourceIP: "2.2.2.2",
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
		{
			// The shape a malware alert produces: a name and a hash, with no
			// directory, which is what a bare file path yields.
			Name: "File entity with a name and a hash, no directory",
			Identifiers: &Identifiers{
				File: &FileEntity{
					Name:   "mirai",
					SHA256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				},
			},
		},
		{
			Name: "CDR entities - all fields populated",
			Identifiers: &Identifiers{
				Event: &EventDetailsEntity{
					EventName:   "CreateUser",
					EventSource: "iam.amazonaws.com",
				},
				UserIdentity: &UserIdentityEntity{
					UserName:    "alice",
					Type:        "IAMUser",
					ARN:         "arn:aws:iam::123456789012:user/alice",
					PrincipalID: "AIDA1234EXAMPLE",
					AccessKeyID: "AKIA1234EXAMPLE",
				},
				SourceInformation: &SourceInformationEntity{
					SourceIPAddress: "203.0.113.42",
					UserAgent:       "aws-cli/2.13.0",
				},
			},
		},
		{
			Name: "CDR entities - partial fields",
			Identifiers: &Identifiers{
				Event: &EventDetailsEntity{
					EventName: "DeleteBucket",
				},
				UserIdentity: &UserIdentityEntity{
					ARN:  "arn:aws:iam::123456789012:role/admin",
					Type: "AssumedRole",
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

// The Flatten test above derives its expectation from the json tags, so it cannot
// fail if a key string changes on both sides at once. These keys are a cross-service
// contract: they are the entity names stored in exception rows and queried as Mongo
// paths, so a rename is a silent data break. Pin the literals.
func TestFileIdentifierKeysAreStable(t *testing.T) {
	flattened := (&Identifiers{
		File: &FileEntity{
			Name:      "mirai.elf",
			Directory: "/tmp",
			MD5:       "5d41402abc4b2a76b9719d911017c592",
			SHA1:      "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
			SHA256:    "9f2b1e1d3c4a5b6c7d8e9f0a1b2c3d4e5f60718293a4b5c6d7e8f9a0b1c2d3e4",
		},
	}).Flatten()

	assert.Equal(t, map[string]string{
		"file.name":      "mirai.elf",
		"file.directory": "/tmp",
		"file.md5":       "5d41402abc4b2a76b9719d911017c592",
		"file.sha1":      "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
		"file.sha256":    "9f2b1e1d3c4a5b6c7d8e9f0a1b2c3d4e5f60718293a4b5c6d7e8f9a0b1c2d3e4",
	}, flattened)
}

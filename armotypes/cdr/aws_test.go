package cdr

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getValueAtJsonPath(genericMap map[string]interface{}, path string) (interface{}, bool) {
	keys := strings.Split(path, ".")
	var current interface{} = genericMap

	for _, key := range keys {
		asMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current, ok = asMap[key]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

// Tests that the AccountID and OrgID JSON paths matches the json tags on CdrAlert
// Later we use these paths to access the values in the config service
func TestCdrAlertJsonPath(t *testing.T) {
	tests := []struct {
		Name          string
		Path          string
		CdrAlert      CdrAlert
		ExpectedValue string
	}{
		{
			Name: "Test CdrEventAccountIDJsonPath",
			Path: CdrEventAccountIDJsonPath,
			CdrAlert: CdrAlert{
				EventData: EventData{
					AWSCloudTrail: &CloudTrailEvent{
						UserIdentity: UserIdentity{
							AccountID: "1111111111",
						},
					},
				},
			},
			ExpectedValue: "1111111111",
		},
		{
			Name: "Test CdrEventOrgIDJsonPath",
			Path: CdrEventOrgIDJsonPath,
			CdrAlert: CdrAlert{
				EventData: EventData{
					AWSCloudTrail: &CloudTrailEvent{
						UserIdentity: UserIdentity{
							OrgID: "o-1111111111",
						},
					},
				},
			},
			ExpectedValue: "o-1111111111",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			// 1. Marshal to JSON
			data, err := json.Marshal(test.CdrAlert)
			require.NoError(t, err)

			// 2. Unmarshal into a generic map
			var genericCdrAlert map[string]interface{}
			require.NoError(t, json.Unmarshal(data, &genericCdrAlert))

			PathWithoutPrefix := strings.TrimPrefix(test.Path, "cdrevent.")

			// 3. Traverse the path
			val, ok := getValueAtJsonPath(genericCdrAlert, PathWithoutPrefix)
			require.True(t, ok)
			assert.Equal(t, test.ExpectedValue, val)
		})
	}
}

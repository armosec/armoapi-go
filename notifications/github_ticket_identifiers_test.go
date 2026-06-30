package notifications

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// GitHubTicketIdentifiers carries routing identifiers at the top level and all
// issue content under a generic `fields` map, mirroring JiraTicketIdentifiers
// and LinearTicketIdentifiers. Content keys match the field schema's fieldId
// (labels, assignees, milestone), and select values may be bare or the UI's
// option-object shape — coercion happens downstream in cadashboardbe, so the
// map carries whatever was submitted verbatim.

func TestGitHubTicketIdentifiers_Serialization(t *testing.T) {
	identifiers := GitHubTicketIdentifiers{
		CollaborationGUID: "collab-guid-123",
		OrganizationName:  "my-org",
		RepositoryName:    "my-repo",
		Fields: map[string]interface{}{
			"labels":    []interface{}{"bug", "priority:high"},
			"assignees": []interface{}{"alice"},
			"milestone": "7",
		},
	}

	data, err := json.Marshal(identifiers)
	require.NoError(t, err)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &raw))
	require.Contains(t, raw, "collaborationGUID")
	require.Contains(t, raw, "organizationName")
	require.Contains(t, raw, "repositoryName")
	require.Contains(t, raw, "fields")
	// Content lives under `fields`, not at the top level — matches Jira/Linear.
	require.NotContains(t, raw, "labels")
	require.NotContains(t, raw, "assignees")
	require.NotContains(t, raw, "milestoneId")
	require.NotContains(t, raw, "milestone")

	var decoded GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, "collab-guid-123", decoded.CollaborationGUID)
	require.Equal(t, "my-org", decoded.OrganizationName)
	require.Equal(t, "my-repo", decoded.RepositoryName)
	require.Equal(t, identifiers.Fields, decoded.Fields)
}

func TestGitHubTicketIdentifiers_OmitsEmptyOptionals(t *testing.T) {
	identifiers := GitHubTicketIdentifiers{}

	data, err := json.Marshal(identifiers)
	require.NoError(t, err)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &raw))

	require.Empty(t, raw, "empty GitHubTicketIdentifiers should serialize to {}")
}

// The map carries the UI's select-option shape verbatim — the workflow path no
// longer needs a special milestoneId key; everything rides `fields` keyed by
// the schema fieldId (labels, assignees, milestone), and cadashboardbe coerces
// option objects downstream.
func TestGitHubTicketIdentifiers_PreservesSelectOptionShape(t *testing.T) {
	in := `{
		"collaborationGUID": "c1",
		"organizationName": "org",
		"repositoryName": "repo",
		"fields": {
			"labels": [{"id": "bug"}, {"id": "priority:high"}],
			"assignees": [{"id": "alice"}],
			"milestone": {"id": "12"}
		}
	}`

	var got GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal([]byte(in), &got))

	require.Equal(t, "c1", got.CollaborationGUID)
	require.Equal(t, "org", got.OrganizationName)
	require.Equal(t, "repo", got.RepositoryName)
	require.NotNil(t, got.Fields)

	milestone, ok := got.Fields["milestone"].(map[string]interface{})
	require.True(t, ok, "milestone option object preserved in fields")
	require.Equal(t, "12", milestone["id"])

	labels, ok := got.Fields["labels"].([]interface{})
	require.True(t, ok)
	require.Len(t, labels, 2)
}

func TestGitHubTicketIdentifiers_RoutingOnly(t *testing.T) {
	in := `{"collaborationGUID": "c1", "organizationName": "org", "repositoryName": "repo"}`

	var got GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal([]byte(in), &got))

	require.Equal(t, "c1", got.CollaborationGUID)
	require.Equal(t, "org", got.OrganizationName)
	require.Equal(t, "repo", got.RepositoryName)
	require.Nil(t, got.Fields)
}

// Structural parity: GitHub now nests content under `fields` just like Jira.
func TestGitHubTicketIdentifiers_MirrorsJiraFieldsShape(t *testing.T) {
	gh, err := json.Marshal(GitHubTicketIdentifiers{
		CollaborationGUID: "c1",
		Fields:            map[string]interface{}{"labels": []interface{}{"bug"}},
	})
	require.NoError(t, err)
	jira, err := json.Marshal(JiraTicketIdentifiers{
		CollaborationGUID: "c1",
		Fields:            map[string]interface{}{"labels": []interface{}{"bug"}},
	})
	require.NoError(t, err)

	var ghRaw, jiraRaw map[string]interface{}
	require.NoError(t, json.Unmarshal(gh, &ghRaw))
	require.NoError(t, json.Unmarshal(jira, &jiraRaw))
	require.Contains(t, ghRaw, "fields")
	require.Contains(t, jiraRaw, "fields")
	require.Equal(t, ghRaw["fields"], jiraRaw["fields"])
}

func TestAlertChannel_GitHubTicketIdentifiers(t *testing.T) {
	channel := AlertChannel{
		GitHubTicketIdentifiers: []GitHubTicketIdentifiers{
			{
				CollaborationGUID: "collab-guid-123",
				OrganizationName:  "my-org",
				RepositoryName:    "my-repo",
			},
		},
	}

	data, err := json.Marshal(channel)
	require.NoError(t, err)

	var decoded AlertChannel
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Len(t, decoded.GitHubTicketIdentifiers, 1)
	require.Equal(t, "collab-guid-123", decoded.GitHubTicketIdentifiers[0].CollaborationGUID)
	require.Equal(t, "my-org", decoded.GitHubTicketIdentifiers[0].OrganizationName)
	require.Equal(t, "my-repo", decoded.GitHubTicketIdentifiers[0].RepositoryName)
}

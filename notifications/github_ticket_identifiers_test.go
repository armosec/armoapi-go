package notifications

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitHubTicketIdentifiers_Serialization(t *testing.T) {
	milestoneID := 7
	identifiers := GitHubTicketIdentifiers{
		CollaborationGUID: "collab-guid-123",
		OrganizationName:  "my-org",
		RepositoryName:    "my-repo",
		Labels:            []string{"bug", "priority:high"},
		Assignees:         []string{"alice"},
		MilestoneID:       &milestoneID,
	}

	data, err := json.Marshal(identifiers)
	require.NoError(t, err)

	var decoded GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.Equal(t, identifiers.CollaborationGUID, decoded.CollaborationGUID)
	require.Equal(t, identifiers.OrganizationName, decoded.OrganizationName)
	require.Equal(t, identifiers.RepositoryName, decoded.RepositoryName)
	require.Equal(t, identifiers.Labels, decoded.Labels)
	require.Equal(t, identifiers.Assignees, decoded.Assignees)
	require.Equal(t, *identifiers.MilestoneID, *decoded.MilestoneID)
}

func TestGitHubTicketIdentifiers_OmitsEmptyOptionals(t *testing.T) {
	identifiers := GitHubTicketIdentifiers{}

	data, err := json.Marshal(identifiers)
	require.NoError(t, err)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &raw))

	require.Empty(t, raw, "empty GitHubTicketIdentifiers should serialize to {}")
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

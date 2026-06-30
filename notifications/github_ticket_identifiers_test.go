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

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &raw))
	require.Contains(t, raw, "collaborationGUID")
	require.Contains(t, raw, "organizationName")
	require.Contains(t, raw, "repositoryName")
	require.Contains(t, raw, "labels")
	require.Contains(t, raw, "assignees")
	require.Contains(t, raw, "milestoneId")

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

// The workflow/runtime payload historically accepted only the bare forms
// (labels: ["bug"], milestoneId: 7). The UI submits the same select-option
// shape it uses for every provider — {"id": ...} / {"name": ...} objects —
// so decoding must tolerate both, normalizing to the typed fields.

func TestGitHubTicketIdentifiers_DecodesBareForm(t *testing.T) {
	in := `{
		"collaborationGUID": "c1",
		"organizationName": "org",
		"repositoryName": "repo",
		"labels": ["bug", "priority:high"],
		"assignees": ["alice"],
		"milestoneId": 7
	}`

	var got GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal([]byte(in), &got))

	require.Equal(t, "c1", got.CollaborationGUID)
	require.Equal(t, "org", got.OrganizationName)
	require.Equal(t, "repo", got.RepositoryName)
	require.Equal(t, []string{"bug", "priority:high"}, got.Labels)
	require.Equal(t, []string{"alice"}, got.Assignees)
	require.NotNil(t, got.MilestoneID)
	require.Equal(t, 7, *got.MilestoneID)
}

func TestGitHubTicketIdentifiers_DecodesSelectOptionShape(t *testing.T) {
	// The UI's shape: option objects for labels/assignees, a milestone object.
	in := `{
		"collaborationGUID": "c1",
		"organizationName": "org",
		"repositoryName": "repo",
		"labels": [{"id": "bug"}, {"id": "priority:high"}],
		"assignees": [{"id": "alice"}],
		"milestone": {"id": "12"}
	}`

	var got GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal([]byte(in), &got))

	require.Equal(t, []string{"bug", "priority:high"}, got.Labels)
	require.Equal(t, []string{"alice"}, got.Assignees)
	require.NotNil(t, got.MilestoneID)
	require.Equal(t, 12, *got.MilestoneID)
}

func TestGitHubTicketIdentifiers_NameFallbackAndNumericMilestone(t *testing.T) {
	in := `{
		"labels": [{"name": "bug"}],
		"milestone": {"id": 3}
	}`

	var got GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal([]byte(in), &got))

	require.Equal(t, []string{"bug"}, got.Labels)
	require.NotNil(t, got.MilestoneID)
	require.Equal(t, 3, *got.MilestoneID)
}

func TestGitHubTicketIdentifiers_MilestoneIdWinsOverMilestone(t *testing.T) {
	in := `{"milestoneId": 7, "milestone": {"id": "12"}}`

	var got GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal([]byte(in), &got))

	require.NotNil(t, got.MilestoneID)
	require.Equal(t, 7, *got.MilestoneID)
}

func TestGitHubTicketIdentifiers_MixedAndEmptyOptionsSkipped(t *testing.T) {
	// Mixed bare + object elements; empty/malformed options are dropped.
	in := `{
		"labels": ["bug", {"id": "enhancement"}, {"id": ""}, {}],
		"assignees": []
	}`

	var got GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal([]byte(in), &got))

	require.Equal(t, []string{"bug", "enhancement"}, got.Labels)
	require.Empty(t, got.Assignees)
	require.Nil(t, got.MilestoneID)
}

func TestGitHubTicketIdentifiers_NullsAndAbsentAreNil(t *testing.T) {
	in := `{"collaborationGUID": "c1", "labels": null}`

	var got GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal([]byte(in), &got))

	require.Equal(t, "c1", got.CollaborationGUID)
	require.Nil(t, got.Labels)
	require.Nil(t, got.Assignees)
	require.Nil(t, got.MilestoneID)
}

func TestGitHubTicketIdentifiers_RoundTripMarshalsBareForm(t *testing.T) {
	// After decoding the option shape, re-marshaling produces the bare,
	// typed form (what stored configs and UNS expect).
	in := `{"labels": [{"id": "bug"}], "milestone": {"id": "12"}}`

	var got GitHubTicketIdentifiers
	require.NoError(t, json.Unmarshal([]byte(in), &got))

	out, err := json.Marshal(got)
	require.NoError(t, err)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal(out, &raw))
	require.Equal(t, []interface{}{"bug"}, raw["labels"])
	require.Equal(t, float64(12), raw["milestoneId"])
	require.NotContains(t, raw, "milestone")
}

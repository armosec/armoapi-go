package armotypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgenticEntityType(t *testing.T) {
	tests := []struct {
		name            string
		clientProviders []string
		serverProviders []string
		want            string
	}{
		{
			name:            "client only -> AI Agent",
			clientProviders: []string{"AWS Bedrock"},
			serverProviders: nil,
			want:            EntityTypeAIAgent,
		},
		{
			name:            "server only -> MCP Server",
			clientProviders: nil,
			serverProviders: []string{"OpenAI"},
			want:            EntityTypeMCPServer,
		},
		{
			name:            "both -> MCP Server (server wins)",
			clientProviders: []string{"AWS Bedrock"},
			serverProviders: []string{"OpenAI"},
			want:            EntityTypeMCPServer,
		},
		{
			name:            "neither -> empty",
			clientProviders: nil,
			serverProviders: nil,
			want:            "",
		},
		{
			name:            "empty (non-nil) slices -> empty",
			clientProviders: []string{},
			serverProviders: []string{},
			want:            "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, AgenticEntityType(tt.clientProviders, tt.serverProviders))
		})
	}
}

func TestIsAgentic(t *testing.T) {
	tests := []struct {
		name            string
		clientProviders []string
		serverProviders []string
		want            bool
	}{
		{"client only", []string{"AWS Bedrock"}, nil, true},
		{"server only", nil, []string{"OpenAI"}, true},
		{"both", []string{"AWS Bedrock"}, []string{"OpenAI"}, true},
		{"neither (nil)", nil, nil, false},
		{"neither (empty)", []string{}, []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsAgentic(tt.clientProviders, tt.serverProviders))
		})
	}
}

// TestAgenticEntityType_matchesIsAgentic asserts the two helpers agree: a
// non-empty entity type iff IsAgentic is true (single classification rule).
func TestAgenticEntityType_matchesIsAgentic(t *testing.T) {
	cases := [][2][]string{
		{{"AWS Bedrock"}, nil},
		{nil, {"OpenAI"}},
		{{"AWS Bedrock"}, {"OpenAI"}},
		{nil, nil},
	}
	for _, c := range cases {
		entityType := AgenticEntityType(c[0], c[1])
		assert.Equal(t, entityType != "", IsAgentic(c[0], c[1]))
	}
}

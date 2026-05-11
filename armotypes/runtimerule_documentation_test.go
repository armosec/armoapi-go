package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v3"
)

func TestRuntimeRule_DocumentationJSONRoundTrip(t *testing.T) {
	rule := RuntimeRule{
		ID:            "R0001",
		Name:          "Unexpected Process Launched",
		Documentation: "# R0001 — Unexpected Process Launched\n\n## Description\nFoo.\n",
	}

	data, err := json.Marshal(rule)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"documentation":"# R0001`)

	var got RuntimeRule
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, rule.Documentation, got.Documentation)
}

func TestRuntimeRule_DocumentationOmitEmpty(t *testing.T) {
	rule := RuntimeRule{ID: "R0001", Name: "X"}

	data, err := json.Marshal(rule)
	require.NoError(t, err)
	assert.NotContains(t, string(data), `"documentation"`,
		"empty documentation should be omitted from JSON")
}

func TestRuntimeRule_DocumentationBSONRoundTrip(t *testing.T) {
	rule := RuntimeRule{
		ID:            "R0001",
		Documentation: "# Rule\n\n## Description\n",
	}

	data, err := bson.Marshal(rule)
	require.NoError(t, err)

	var got RuntimeRule
	require.NoError(t, bson.Unmarshal(data, &got))
	assert.Equal(t, rule.Documentation, got.Documentation)
}

func TestRuntimeRule_DocumentationYAMLRoundTrip(t *testing.T) {
	rule := RuntimeRule{
		ID:            "R0001",
		Documentation: "# Rule\n",
	}

	data, err := yaml.Marshal(rule)
	require.NoError(t, err)

	var got RuntimeRule
	require.NoError(t, yaml.Unmarshal(data, &got))
	assert.Equal(t, rule.Documentation, got.Documentation)
}

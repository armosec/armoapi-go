package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUnmarshalRegistryKeepsScanAllRepositories(t *testing.T) {
	t.Run("flag and repository list survive unmarshal", func(t *testing.T) {
		payload := []byte(`{"provider":"quay","containerRegistryName":"quay.io/org","repositories":["org/a","org/b"],"scanAllRepositories":true}`)

		registry, err := UnmarshalRegistry(payload)

		require.NoError(t, err)
		assert.True(t, registry.GetBase().ScanAllRepositories)
		assert.Equal(t, []string{"org/a", "org/b"}, registry.GetBase().Repositories)
	})

	t.Run("flag is omitted from JSON when false", func(t *testing.T) {
		data, err := json.Marshal(BaseContainerImageRegistry{Provider: Quay, Repositories: []string{"org/a"}})

		require.NoError(t, err)
		assert.NotContains(t, string(data), "scanAllRepositories")
	})

	t.Run("registry without the flag decodes as false", func(t *testing.T) {
		registry, err := UnmarshalRegistry([]byte(`{"provider":"quay","repositories":["org/a"]}`))

		require.NoError(t, err)
		assert.False(t, registry.GetBase().ScanAllRepositories)
	})

	t.Run("flag survives a BSON round trip", func(t *testing.T) {
		data, err := bson.Marshal(BaseContainerImageRegistry{Provider: Quay, ScanAllRepositories: true})
		require.NoError(t, err)

		var decoded BaseContainerImageRegistry
		require.NoError(t, bson.Unmarshal(data, &decoded))
		assert.True(t, decoded.ScanAllRepositories)
	})
}

func TestCleanRegistryURL(t *testing.T) {
	assert.Equal(t, "quay.io", cleanRegistryURL("quay.io"))
	assert.Equal(t, "quay.io", cleanRegistryURL("www.quay.io"))
	assert.Equal(t, "quay.io", cleanRegistryURL("https://quay.io"))
	assert.Equal(t, "quay.io", cleanRegistryURL("https://www.quay.io"))
	assert.Equal(t, "quay.io", cleanRegistryURL("http://www.quay.io"))
	assert.Equal(t, "quay.io:5000", cleanRegistryURL("https://www.quay.io:5000"))
}

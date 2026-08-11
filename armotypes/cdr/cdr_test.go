package cdr

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCdrAlertBatchConnectionLevel locks the wire contract for ConnectionLevel: it
// round-trips under its json tag, is omitted when empty (so existing AWS producers that
// never set it stay byte-for-byte back-compatible), and the constants hold their string
// values (consumers compare against them).
func TestCdrAlertBatchConnectionLevel(t *testing.T) {
	assert.Equal(t, ConnectionLevel("account"), ConnectionLevelAccount)
	assert.Equal(t, ConnectionLevel("organization"), ConnectionLevelOrganization)

	t.Run("account level round-trips", func(t *testing.T) {
		b, err := json.Marshal(CdrAlertBatch{
			CustomerGUID:    "cg",
			CloudAccountID:  "sub-1",
			OrgID:           "tenant-1",
			Provider:        Azure,
			ConnectionLevel: ConnectionLevelAccount,
		})
		require.NoError(t, err)
		assert.Contains(t, string(b), `"connectionLevel":"account"`)

		var got CdrAlertBatch
		require.NoError(t, json.Unmarshal(b, &got))
		assert.Equal(t, ConnectionLevelAccount, got.ConnectionLevel)
	})

	t.Run("omitted when empty (AWS back-compat)", func(t *testing.T) {
		b, err := json.Marshal(CdrAlertBatch{CustomerGUID: "cg", Provider: AWS})
		require.NoError(t, err)
		assert.False(t, strings.Contains(string(b), "connectionLevel"),
			"an unset ConnectionLevel must not appear on the wire")
	})
}

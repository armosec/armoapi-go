package apis

import (
	"encoding/json"
	"testing"

	"github.com/moby/moby/api/types/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests pin the wire format of registry credentials after the switch from
// github.com/docker/docker to github.com/moby/moby/api (see
// docs/features/registry-authconfig-moby.md). Producers and consumers deploy
// independently, so the JSON keys must not drift.

var fullAuthConfig = registry.AuthConfig{
	Username:      "user",
	Password:      "pass",
	Auth:          "dXNlcjpwYXNz",
	ServerAddress: "registry.example.com",
	IdentityToken: "id-token",
	RegistryToken: "reg-token",
}

var fullAuthConfigJSON = map[string]any{
	"username":      "user",
	"password":      "pass",
	"auth":          "dXNlcjpwYXNz",
	"serveraddress": "registry.example.com",
	"identitytoken": "id-token",
	"registrytoken": "reg-token",
}

func TestRegistryCredentialsMarshalKeys(t *testing.T) {
	cmd := WebsocketScanCommand{
		Credentials:     &fullAuthConfig,
		ImageScanParams: ImageScanParams{Credentialslist: []registry.AuthConfig{fullAuthConfig}},
	}
	b, err := json.Marshal(cmd)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, fullAuthConfigJSON, got["credentials"])
	assert.Equal(t, []any{fullAuthConfigJSON}, got["credentialsList"])
}

func TestRegistryCredentialsUnmarshalIgnoresEmail(t *testing.T) {
	// A producer still on github.com/docker/docker may send the deprecated
	// "email" key; it must decode cleanly with every other field intact.
	payload := `{
		"credentials": {"username":"user","password":"pass","auth":"dXNlcjpwYXNz","email":"a@b.c",
			"serveraddress":"registry.example.com","identitytoken":"id-token","registrytoken":"reg-token"},
		"credentialsList": [{"username":"user","password":"pass","auth":"dXNlcjpwYXNz","email":"a@b.c",
			"serveraddress":"registry.example.com","identitytoken":"id-token","registrytoken":"reg-token"}]
	}`
	var cmd WebsocketScanCommand
	require.NoError(t, json.Unmarshal([]byte(payload), &cmd))
	require.NotNil(t, cmd.Credentials)
	assert.Equal(t, fullAuthConfig, *cmd.Credentials)
	assert.Equal(t, []registry.AuthConfig{fullAuthConfig}, cmd.Credentialslist)
}

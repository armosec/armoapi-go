package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestProfileDataPattern_YAMLRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		yaml string
		want ProfileDataPattern
	}{
		{"exact", `exact: "/var/run/docker.sock"`, ProfileDataPattern{Exact: "/var/run/docker.sock"}},
		{"prefix", `prefix: "/etc/cron.d/"`, ProfileDataPattern{Prefix: "/etc/cron.d/"}},
		{"suffix", `suffix: "/authorized_keys"`, ProfileDataPattern{Suffix: "/authorized_keys"}},
		{"contains", `contains: "authorized_keys"`, ProfileDataPattern{Contains: "authorized_keys"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got ProfileDataPattern
			require.NoError(t, yaml.Unmarshal([]byte(c.yaml), &got))
			assert.Equal(t, c.want, got)

			out, err := yaml.Marshal(got)
			require.NoError(t, err)
			var roundTripped ProfileDataPattern
			require.NoError(t, yaml.Unmarshal(out, &roundTripped))
			assert.Equal(t, c.want, roundTripped)
		})
	}
}

func TestProfileDataPattern_JSONRoundTrip(t *testing.T) {
	p := ProfileDataPattern{Prefix: "/etc/cron.d/"}
	b, err := json.Marshal(p)
	require.NoError(t, err)
	assert.JSONEq(t, `{"prefix":"/etc/cron.d/"}`, string(b))

	var got ProfileDataPattern
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, p, got)
}

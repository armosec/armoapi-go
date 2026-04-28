package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
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

func TestProfileDataField_UnmarshalYAML_All(t *testing.T) {
	var f ProfileDataField
	require.NoError(t, yaml.Unmarshal([]byte(`all`), &f))
	assert.True(t, f.All)
	assert.Empty(t, f.Patterns)
}

func TestProfileDataField_UnmarshalYAML_Patterns(t *testing.T) {
	var f ProfileDataField
	require.NoError(t, yaml.Unmarshal([]byte(`
- exact: "/var/run/docker.sock"
- prefix: "/etc/cron.d/"
`), &f))
	assert.False(t, f.All)
	require.Len(t, f.Patterns, 2)
	assert.Equal(t, "/var/run/docker.sock", f.Patterns[0].Exact)
	assert.Equal(t, "/etc/cron.d/", f.Patterns[1].Prefix)
}

func TestProfileDataField_UnmarshalYAML_RejectsUnknown(t *testing.T) {
	var f ProfileDataField
	err := yaml.Unmarshal([]byte(`42`), &f)
	assert.Error(t, err)

	err = yaml.Unmarshal([]byte(`"some"`), &f)
	assert.Error(t, err)
}

func TestProfileDataField_MarshalYAML_All(t *testing.T) {
	f := ProfileDataField{All: true}
	out, err := yaml.Marshal(f)
	require.NoError(t, err)
	assert.Equal(t, "all\n", string(out))
}

func TestProfileDataField_MarshalYAML_Patterns(t *testing.T) {
	f := ProfileDataField{Patterns: []ProfileDataPattern{{Exact: "/a"}}}
	out, err := yaml.Marshal(f)
	require.NoError(t, err)
	assert.Contains(t, string(out), `exact: /a`)
}

func TestProfileDataField_JSONRoundTrip(t *testing.T) {
	for _, f := range []ProfileDataField{
		{All: true},
		{Patterns: []ProfileDataPattern{{Prefix: "/x"}, {Suffix: "/y"}}},
	} {
		b, err := json.Marshal(f)
		require.NoError(t, err)
		var got ProfileDataField
		require.NoError(t, json.Unmarshal(b, &got))
		assert.Equal(t, f, got)
	}
}

func TestProfileDataField_BSONRoundTrip(t *testing.T) {
	type wrapper struct {
		Field ProfileDataField `bson:"field"`
	}
	for _, f := range []ProfileDataField{
		{All: true},
		{Patterns: []ProfileDataPattern{{Prefix: "/x"}, {Suffix: "/y"}}},
	} {
		b, err := bson.Marshal(wrapper{Field: f})
		require.NoError(t, err)
		var got wrapper
		require.NoError(t, bson.Unmarshal(b, &got))
		assert.Equal(t, f, got.Field)
	}
}

func TestProfileDataRequired_Validate_Valid(t *testing.T) {
	p := &ProfileDataRequired{
		Opens: &ProfileDataField{Patterns: []ProfileDataPattern{{Exact: "/a"}}},
		Execs: &ProfileDataField{All: true},
	}
	assert.NoError(t, p.Validate())
}

func TestProfileDataRequired_Validate_FieldEmpty(t *testing.T) {
	p := &ProfileDataRequired{Opens: &ProfileDataField{}}
	err := p.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "opens")
}

func TestProfileDataRequired_Validate_PatternMultiKey(t *testing.T) {
	p := &ProfileDataRequired{
		Opens: &ProfileDataField{Patterns: []ProfileDataPattern{{Exact: "a", Prefix: "b"}}},
	}
	err := p.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exactly one")
}

func TestProfileDataRequired_Validate_PatternEmpty(t *testing.T) {
	p := &ProfileDataRequired{
		Opens: &ProfileDataField{Patterns: []ProfileDataPattern{{}}},
	}
	err := p.Validate()
	require.Error(t, err)
}

func TestProfileDataRequired_Validate_AllAndPatterns(t *testing.T) {
	p := &ProfileDataRequired{
		Opens: &ProfileDataField{All: true, Patterns: []ProfileDataPattern{{Exact: "a"}}},
	}
	err := p.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "all")
}

func TestProfileDataRequired_IsEmpty(t *testing.T) {
	var nilP *ProfileDataRequired
	assert.True(t, nilP.IsEmpty())
	assert.True(t, (&ProfileDataRequired{}).IsEmpty())
	assert.False(t, (&ProfileDataRequired{Opens: &ProfileDataField{All: true}}).IsEmpty())
}

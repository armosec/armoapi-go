package armotypes

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestRuntimeRuleMitreMappings(t *testing.T) {
	pairs := []MitreMapping{{"TA0008", "T1021"}, {"TA0007", "T1613"}}
	r := RuntimeRule{MitreMappings: pairs, MitreTactic: "TA0008", MitreTechnique: "T1021"}
	for name, codec := range map[string]struct {
		marshal   func(any) ([]byte, error)
		unmarshal func([]byte, any) error
	}{
		"json": {json.Marshal, json.Unmarshal}, "bson": {bson.Marshal, bson.Unmarshal}, "yaml": {yaml.Marshal, yaml.Unmarshal},
	} {
		t.Run(name, func(t *testing.T) {
			data, err := codec.marshal(r)
			require.NoError(t, err)
			var decoded RuntimeRule
			require.NoError(t, codec.unmarshal(data, &decoded))
			require.Equal(t, pairs, decoded.GetMitreMappings())
		})
	}
	require.Equal(t, pairs[:1], (RuntimeRule{MitreTactic: "TA0008", MitreTechnique: "T1021"}).GetMitreMappings())
	require.Nil(t, (RuntimeRule{}).GetMitreMappings())
	require.Equal(t, pairs, (RuntimeRule{MitreMappings: pairs, MitreTactic: "other"}).GetMitreMappings())
}

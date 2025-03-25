package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshallChildrenMap(t *testing.T) {
	childrenMap := map[CommPID]*Process{
		{Comm: "test", PID: 1}: {PID: 1},
	}
	data, err := json.Marshal(childrenMap)
	assert.NoError(t, err)
	var childrenMap2 map[CommPID]*Process
	err = json.Unmarshal(data, &childrenMap2)
	assert.NoError(t, err)
	assert.Equal(t, childrenMap, childrenMap2)
}

package armotypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInnerFilterToElementMatchString(t *testing.T) {
	innerFilters := map[string]string{
		"filter1": "value1",
		"filter2": "value2",
		"filter3": "|operator",
	}

	expectedResult := "filter1:value1;filter2:value2;filter3:|operator|elemMatch"
	result := Filter2ElementMatchString(innerFilters)

	assert.Equal(t, expectedResult, result)
}

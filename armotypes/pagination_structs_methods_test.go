package armotypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInnerFilterToElementMatchString(t *testing.T) {
	innerFilters := map[string]string{
		"address": "|missing",
		"age":     "30|greater",
		"name":    "john",
	}

	expectedResult := "address:|missing;age:30|greater;name:john|elemMatch"
	result := Filter2ElementMatchString(innerFilters)

	assert.Equal(t, expectedResult, result)
}

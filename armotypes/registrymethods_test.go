package armotypes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanRegistryURL(t *testing.T) {
	assert.Equal(t, "quay.io", cleanRegistryURL("quay.io"))
	assert.Equal(t, "quay.io", cleanRegistryURL("www.quay.io"))
	assert.Equal(t, "quay.io", cleanRegistryURL("https://quay.io"))
	assert.Equal(t, "quay.io", cleanRegistryURL("https://www.quay.io"))
	assert.Equal(t, "quay.io", cleanRegistryURL("http://www.quay.io"))
	assert.Equal(t, "quay.io:5000", cleanRegistryURL("https://www.quay.io:5000"))
	assert.Equal(t, "gitlab.com", cleanRegistryURL("gitlab.com"))
}

package notifications

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDomainScope(t *testing.T) {
	cluster1Scope := AlertScope{
		Cluster:    "cluster1FullName",
		Namespaces: []string{"s1", "s2"},
	}
	cluster2Scope := AlertScope{
		Cluster:    "cluster2FullName",
		Namespaces: []string{"s1", "s2"},
	}
	alertChannelAPI := AlertChannelAPI{
		Scope: []EnrichedScope{
			{
				AlertScope:       cluster1Scope,
				ClusterShortName: "cluster1ShortName",
			},
			{
				AlertScope:       cluster2Scope,
				ClusterShortName: "cluster2ShortName",
			},
		},
	}

	domainScope := alertChannelAPI.GetDomainScope()

	assert.Equal(t, len(domainScope), 2)
	assert.Equal(t, domainScope, []AlertScope{cluster1Scope, cluster2Scope})
}

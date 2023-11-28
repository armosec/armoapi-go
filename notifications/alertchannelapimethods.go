package notifications

func (ac *AlertChannelAPI) GetDomainScope() []AlertScope {
	scope := make([]AlertScope, 0, len(ac.Scope))
	for _, enrichedScope := range ac.Scope {
		scope = append(scope, AlertScope{
			Cluster:    enrichedScope.Cluster,
			Namespaces: enrichedScope.Namespaces,
		})
	}
	return scope
}

package armotypes

func MockCustomerConfig() *CustomerConfig {
	scope := *MockPortalDesignator()
	settings := *MockSettings()
	return &CustomerConfig{
		Name:     "test-cluster1",
		Scope:    scope,
		Settings: settings,
	}
}

func MockPortalDesignator() *PortalDesignator {
	return &PortalDesignator{
		DesignatorType: "Attributes",
		WLID:           "",
		WildWLID:       "",
		SID:            "",
		Attributes:     map[string]string{"cluster": "test-cluster1"},
	}
}

func MockSettings() *Settings {
	postureControlInputs := map[string][]string{
		"public_registries":     {"quay.io/kiali/", "quay.io/datawire/", "quay.io/keycloak/", "quay.io/bitnami/"},
		"dangerousCapabilities": {"ALL", "SYS_ADMIN", "NET_ADMIN"},
	}
	return &Settings{
		PostureControlInputs:    postureControlInputs,
		PostureScanConfig:       PostureScanConfig{},
		VulnerabilityScanConfig: VulnerabilityScanConfig{},
	}
}

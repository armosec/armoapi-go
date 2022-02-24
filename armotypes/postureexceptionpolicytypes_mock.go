package armotypes

func mockException() *PostureExceptionPolicy {
	return &PostureExceptionPolicy{
		PortalBase: PortalBase{
			Name: "",
		},
		Actions: []PostureExceptionPolicyActions{AlertOnly},
		Resources: []PortalDesignator{
			{
				DesignatorType: DesignatorAttributes,
				Attributes: map[string]string{
					AttributeKind:      "",
					AttributeNamespace: "",
				},
			},
		},
	}
}

package armotypes

import (
	"maps"
	"strings"
	"time"

	"github.com/armosec/armoapi-go/identifiers"
)

const (
	// SecurityRiskPolicy - policy for security risks
	SecurityRiskExceptionPolicyType PolicyType = "securityRiskExceptionPolicy"

	// RuntimeIncidentPolicy - policy for runtime incidents
	RuntimeIncidentExceptionPolicyType PolicyType = "runtimeIncidentExceptionPolicy"

	// CSPM - policy for CSPM
	CSPMExceptionPolicyType PolicyType = "cspmExceptionPolicy"
)

type AdvancedScopeEntity struct {
	Entity   string `json:"entity" bson:"entity"`
	Operator string `json:"condition" bson:"operator"`
	Values   string `json:"values" bson:"values"`
}

type BaseExceptionPolicy struct {
	PortalBase `json:",inline" bson:"inline"`
	PolicyType PolicyType `json:"policyType,omitempty" bson:"policyType,omitempty"`

	// IDs of the policies (SecurityRiskID, ControlID, etc.)
	PolicyIDs      []string                       `json:"policyIDs,omitempty" bson:"policyIDs,omitempty"`
	CreationTime   string                         `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	Reason         string                         `json:"reason,omitempty" bson:"reason,omitempty"`
	ExpirationDate *time.Time                     `json:"expirationDate,omitempty" bson:"expirationDate,omitempty"`
	CreatedBy      string                         `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	Resources      []identifiers.PortalDesignator `json:"resources,omitempty" bson:"resources,omitempty"`
	AdvancedScopes []AdvancedScopeEntity          `json:"advancedScopes,omitempty" bson:"advancedScopes,omitempty"`
}

// Used by cadashboardbe (countIncidents API) and event-ingester (retroactive resolve on exception creation).
func GetRuntimeIncidentsRequestFilterFromExceptionPolicy(exceptionPolicy BaseExceptionPolicy) []map[string]string {
	if len(exceptionPolicy.PolicyIDs) == 0 {
		return nil
	}

	advancedScopesFilters := make(map[string]string)
	for _, advancedScope := range exceptionPolicy.AdvancedScopes {
		value := escapeV2ListOperatorSeparator(advancedScope.Values)
		switch advancedScope.Operator {
		case "in":
			value = trimSpacesAroundCommas(value)
		case "contains":
			value = value + V2ListOperatorSeparator + V2ListLikeOperator
		}
		advancedScopesFilters["identifiers."+advancedScope.Entity] = value
	}

	var filters []map[string]string
	for _, designator := range exceptionPolicy.Resources {
		if designator.Attributes == nil {
			continue
		}
		filter := map[string]string{
			"incidentTypeID": exceptionPolicy.PolicyIDs[0], // only incidents of this type are marked as resolved
			"status":         "Open",                       // only unresolved incidents are marked as resolved
		}

		if cluster, ok := designator.Attributes[identifiers.AttributeCluster]; ok && cluster != GlobalRegex {
			filter["designators.attributes.cluster"] = cluster
		}
		if namespace, ok := designator.Attributes[identifiers.AttributeNamespace]; ok && namespace != GlobalRegex {
			filter["designators.attributes.namespace"] = namespace
		}
		if kind, ok := designator.Attributes[identifiers.AttributeKind]; ok && kind != GlobalRegex {
			filter["designators.attributes.kind"] = kind
		}
		if workloadName, ok := designator.Attributes[identifiers.AttributeName]; ok && workloadName != GlobalRegex {
			filter["designators.attributes.name"] = workloadName
		}
		if service, ok := designator.Attributes[identifiers.AttributeService]; ok && service != GlobalRegex {
			filter["designators.attributes.service"] = service
		}
		if region, ok := designator.Attributes[identifiers.AttributeRegion]; ok && region != GlobalRegex {
			filter["designators.attributes.region"] = region
		}
		// cloudProvider and accountID are filtered against cloudMetadata.* rather than
		// designators.attributes.* because CDR incidents historically stored these
		// only in cloudMetadata.
		if provider, ok := designator.Attributes[identifiers.AttributeCloudProvider]; ok && provider != GlobalRegex {
			filter["cloudMetadata.provider"] = provider
		}
		if accountID, ok := designator.Attributes[identifiers.AttributeCloudAccountID]; ok && accountID != GlobalRegex {
			filter["cloudMetadata.account_id"] = accountID
		}
		if instanceID, ok := designator.Attributes[identifiers.AttributeInstanceId]; ok && instanceID != GlobalRegex {
			filter["designators.attributes.instanceId"] = instanceID
		}
		if hostType, ok := designator.Attributes[identifiers.AttributeHostType]; ok && hostType != GlobalRegex {
			filter["designators.attributes.hostType"] = hostType
		}

		maps.Copy(filter, advancedScopesFilters)

		filters = append(filters, filter)
	}
	return filters
}

// escapeV2ListOperatorSeparator escapes the V2List operator separator ("|") in a
// raw user-provided filter value so that config-service's SplitIgnoreEscaped treats
// embedded "|" characters as literal rather than as the start of an operator suffix.
func escapeV2ListOperatorSeparator(value string) string {
	return strings.ReplaceAll(value, V2ListOperatorSeparator, V2ListEscapeChar+V2ListOperatorSeparator)
}

// trimSpacesAroundCommas canonicalizes a comma-separated value list by trimming
// whitespace around each separator, treating backslash-escaped commas as literal
// (i.e., not a separator).
func trimSpacesAroundCommas(value string) string {
	parts := splitIgnoreEscaped(value, V2ListValueSeparator, V2ListEscapeChar)
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return strings.Join(parts, V2ListValueSeparator)
}

// splitIgnoreEscaped splits s by sep, treating sep preceded by escape as a literal
// (the resulting parts retain the escape sequence). Mirrors config-service's
// utils.SplitIgnoreEscaped on the producer side.
func splitIgnoreEscaped(s, sep, escape string) []string {
	pieces := strings.Split(s, sep)
	var out []string
	buf := ""
	for _, p := range pieces {
		if buf != "" {
			buf += sep
		}
		buf += p
		if !strings.HasSuffix(p, escape) {
			out = append(out, buf)
			buf = ""
		}
	}
	if buf != "" {
		out = append(out, buf)
	}
	return out
}

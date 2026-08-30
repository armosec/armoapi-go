package armotypes

import (
	"maps"
	"regexp"
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

// hashScopeEntities are the advanced-scope entities whose value is a hex file hash.
//
// Unexported on purpose. An exported map is mutable by any consumer, and this one
// decides whether a value is folded before it is compared - so a stray write would
// change matching behaviour globally and silently. Read it through
// IsHashScopeEntity.
var hashScopeEntities = map[string]bool{
	"file.md5":    true,
	"file.sha1":   true,
	"file.sha256": true,
}

// IsHashScopeEntity reports whether an advanced-scope entity's value is a file hash,
// and therefore must be compared without regard to the case it was typed in.
func IsHashScopeEntity(entity string) bool {
	return hashScopeEntities[entity]
}

// NormalizeHashScopeValues lower-cases the values of hash-scoped entities in place,
// and leaves every other entity untouched.
//
// The exception is matched in postgres-connector by a plain string comparison:
// operator "in" is an equality, operator "contains" is a LIKE. PostgreSQL folds case
// for neither. The alert carries a lower-case hash, and a hash copied from VirusTotal
// or from a threat report is usually upper case, so an unfolded scope would be saved,
// listed and previewed exactly like a working rule while suppressing nothing.
//
// This deliberately differs from the rule in docs/features/cdr-azure-event-shape.md,
// which says not to normalize casing at any layer and to match case-insensitively
// instead. That rule exists because an Azure subscription id has a customer-chosen
// original casing that is displayed and re-used, so rewriting it loses information
// and a partial rewrite reintroduces mismatches. A hex digest has no original casing
// to preserve and no information in its case, and unlike the Mongo path there is no
// case-insensitive operator available on the Postgres comparison - so folding on
// write is the only option that works there, and it is lossless here.
//
// Kept here rather than in each caller so the write paths cannot drift apart: both
// cadashboardbe, which parses the analyst's request, and event-ingester-service,
// which persists it, fold through this one function.
func NormalizeHashScopeValues(scopes []AdvancedScopeEntity) {
	for i := range scopes {
		if IsHashScopeEntity(scopes[i].Entity) {
			scopes[i].Values = strings.ToLower(scopes[i].Values)
		}
	}
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

// AnchoredIgnoreCaseRegexFilter renders a value as a V2List case-insensitive exact
// match ("^<escaped>$|regex&ignorecase"): an anchored, regex-escaped value with the
// ignorecase option, which config-service resolves to $regex + $options "i".
// Exported as the single source of this encoding so consumers that today
// re-implement it (e.g. event-ingester's BuildGetAccountWithFeatureQuery) can route
// through it instead of drifting out of sync.
func AnchoredIgnoreCaseRegexFilter(value string) string {
	return "^" + regexp.QuoteMeta(value) + "$" +
		V2ListOperatorSeparator + V2ListRegexOperator +
		V2ListSubQuerySeparator + V2ListIgnoreCaseOption
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
		// Hoisted so the accountID branch below can reuse it. Intentionally skips a
		// present-but-empty provider (the old code would have set the filter to ""):
		// a no-op improvement, not an accidental behavior change.
		provider := designator.Attributes[identifiers.AttributeCloudProvider]
		if provider != "" && provider != GlobalRegex {
			filter["cloudMetadata.provider"] = provider
		}
		if accountID, ok := designator.Attributes[identifiers.AttributeCloudAccountID]; ok && accountID != GlobalRegex {
			// Azure subscription ids arrive with inconsistent casing: incidents derive
			// them from the uppercase Activity Log resourceId, while onboarding/storage
			// keep the customer's original casing. Match case-insensitively — the same
			// anchored-regex + ignorecase mechanism config-service's
			// BuildGetAccountWithFeatureQuery uses — so an exact-equality miss can't
			// leave an Azure risk-acceptance ineffective. AWS/GCP keep exact equality.
			if strings.EqualFold(provider, string(ProviderAzure)) {
				filter["cloudMetadata.account_id"] = AnchoredIgnoreCaseRegexFilter(accountID)
			} else {
				filter["cloudMetadata.account_id"] = accountID
			}
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

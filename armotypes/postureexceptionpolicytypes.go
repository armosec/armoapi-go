package armotypes

import (
	"encoding/json"
	"time"

	"github.com/armosec/armoapi-go/identifiers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PostureExceptionPolicyActions string

const AlertOnly PostureExceptionPolicyActions = "alertOnly"
const Disable PostureExceptionPolicyActions = "disable"

type PolicyType string

const GlobalRegex = "*/*"

const PostureExceptionPolicyType PolicyType = "postureExceptionPolicy"
const VulnerabilityExceptionPolicyType PolicyType = "vulnerabilityExceptionPolicy"

type IgnoreRuleUserInputMessage struct {
	PolicyType PolicyType      `json:"policyType"`
	NewData    json.RawMessage `json:"newData"`
	OldData    json.RawMessage `json:"oldData"`
}

type PostureExceptionPolicy struct {
	PortalBase      `json:",inline" bson:"inline"`
	PolicyType      string                          `json:"policyType,omitempty" bson:"policyType,omitempty"`
	CreationTime    string                          `json:"creationTime,omitempty" bson:"creationTime,omitempty"`
	Actions         []PostureExceptionPolicyActions `json:"actions,omitempty" bson:"actions,omitempty"`
	Resources       []identifiers.PortalDesignator  `json:"resources" bson:"resources,omitempty"`
	// ObjectSelector carries a full Kubernetes label selector (matchLabels + matchExpressions)
	// for the exception's workload-matching axis. Unlike Resources (which encodes per-key regex
	// designators), this is evaluated with labels.Selector semantics by the exception processor.
	//
	// Both a nil selector and a non-nil but empty selector ({} / no matchLabels and no
	// matchExpressions) mean "no label constraint": the consumer must skip the label axis
	// entirely in those cases. Do NOT pass either straight to metav1.LabelSelectorAsSelector,
	// whose conversion disagrees with that intent — nil yields labels.Nothing() (matches nothing)
	// and an empty selector yields labels.Everything() (matches every workload). Guard for
	// nil/empty before converting; ObjectSelector.ToMetaV1() returns nil for a nil receiver.
	ObjectSelector  *LabelSelector                  `json:"objectSelector,omitempty" bson:"objectSelector,omitempty"`
	PosturePolicies []PosturePolicy                 `json:"posturePolicies,omitempty" bson:"posturePolicies,omitempty"`
	Reason          *string                         `json:"reason,omitempty" bson:"reason,omitempty"`
	ExpirationDate  *time.Time                      `json:"expirationDate,omitempty" bson:"expirationDate"`
	CreatedBy       string                          `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
}

// LabelSelector mirrors metav1.LabelSelector but carries explicit bson tags so the
// persisted (Mongo) document keeps the same camelCase keys the JSON/API and CRD shapes
// expose (matchLabels / matchExpressions). The upstream metav1 type declares json tags
// only; the mongo-driver default encoder would otherwise lowercase the nested keys to
// matchlabels / matchexpressions and diverge from the API contract.
type LabelSelector struct {
	MatchLabels      map[string]string          `json:"matchLabels,omitempty" bson:"matchLabels,omitempty"`
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty" bson:"matchExpressions,omitempty"`
}

// LabelSelectorRequirement mirrors metav1.LabelSelectorRequirement with explicit bson tags.
type LabelSelectorRequirement struct {
	Key      string                       `json:"key" bson:"key"`
	Operator metav1.LabelSelectorOperator `json:"operator" bson:"operator"`
	Values   []string                     `json:"values,omitempty" bson:"values,omitempty"`
}

// ToMetaV1 converts the selector to a metav1.LabelSelector for evaluation via
// metav1.LabelSelectorAsSelector. It returns nil for a nil receiver; callers must still
// treat a nil or empty result as "no label constraint" (see ObjectSelector).
func (s *LabelSelector) ToMetaV1() *metav1.LabelSelector {
	if s == nil {
		return nil
	}
	out := &metav1.LabelSelector{MatchLabels: s.MatchLabels}
	if len(s.MatchExpressions) > 0 {
		out.MatchExpressions = make([]metav1.LabelSelectorRequirement, len(s.MatchExpressions))
		for i, r := range s.MatchExpressions {
			out.MatchExpressions[i] = metav1.LabelSelectorRequirement{
				Key:      r.Key,
				Operator: r.Operator,
				Values:   r.Values,
			}
		}
	}
	return out
}

type PosturePolicy struct {
	FrameworkName string `json:"frameworkName" bson:"frameworkName"`
	// deprecated - use ControlID instead
	ControlName   string `json:"controlName,omitempty" bson:"controlName,omitempty"`
	ControlID     string `json:"controlID,omitempty" bson:"controlID,omitempty"`
	RuleName      string `json:"ruleName,omitempty" bson:"ruleName,omitempty"`
	SeverityScore int    `json:"severityScore,omitempty" bson:"severityScore,omitempty"`
}

func (exceptionPolicy *PostureExceptionPolicy) IsAlertOnly() bool {
	if exceptionPolicy.IsDisable() {
		return false
	}

	for i := range exceptionPolicy.Actions {
		if exceptionPolicy.Actions[i] == AlertOnly {
			return true
		}
	}
	return false
}
func (exceptionPolicy *PostureExceptionPolicy) IsDisable() bool {
	for i := range exceptionPolicy.Actions {
		if exceptionPolicy.Actions[i] == Disable {
			return true
		}
	}
	return false
}

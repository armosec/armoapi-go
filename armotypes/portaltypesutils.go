package armotypes

import (
	"strings"
	"time"

	wlidpkg "github.com/armosec/utils-k8s-go/wlid"
	"github.com/francoispqt/gojay"
	"k8s.io/utils/strings/slices"
)

var IgnoreLabels = []string{AttributeCluster, AttributeNamespace}

// AttributeDesignators describe a kubernetes object, with its labels.
type AttributesDesignators struct {
	cluster   string
	namespace string
	kind      string
	name      string
	path      string
	labels    map[string]string
}

func (ad *AttributesDesignators) GetCluster() string {
	return ad.cluster
}

func (ad *AttributesDesignators) GetNamespace() string {
	return ad.namespace
}

func (ad *AttributesDesignators) GetKind() string {
	return ad.kind
}

func (ad *AttributesDesignators) GetName() string {
	return ad.name
}

func (ad *AttributesDesignators) GetPath() string {
	return ad.path
}

func (ad *AttributesDesignators) GetLabels() map[string]string {
	return ad.labels
}

func AttributesDesignatorsFromWLID(wlid string) *PortalDesignator {
	wlidSlices := wlidpkg.RestoreMicroserviceIDs(wlid)
	pd := &PortalDesignator{
		DesignatorType: DesignatorAttributes,
		Attributes:     make(map[string]string, 4),
	}
	if len(wlidSlices) > 0 {
		pd.Attributes[AttributeCluster] = wlidSlices[0]
	}
	if len(wlidSlices) > 1 {
		pd.Attributes[AttributeNamespace] = wlidSlices[1]
	}
	if len(wlidSlices) > 2 {
		pd.Attributes[AttributeKind] = wlidSlices[2]
	}
	if len(wlidSlices) > 3 {
		pd.Attributes[AttributeName] = wlidSlices[3]
	}
	return pd
}

func (designator *PortalDesignator) GetCluster() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.cluster
}

func (designator *PortalDesignator) GetNamespace() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.namespace
}

func (designator *PortalDesignator) GetKind() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.kind
}

func (designator *PortalDesignator) GetName() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.name
}
func (designator *PortalDesignator) GetPath() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.path
}
func (designator *PortalDesignator) GetLabels() map[string]string {
	attributes := designator.DigestPortalDesignator()
	return attributes.labels
}

// DigestPortalDesignator - get cluster namespace and labels from designator
func (designator *PortalDesignator) DigestPortalDesignator() AttributesDesignators {
	switch designator.DesignatorType {
	case DesignatorAttributes, DesignatorAttribute:
		return designator.DigestAttributesDesignator()
	case DesignatorWlid.ToLower(), DesignatorWildWlid.ToLower():
		return AttributesDesignators{
			cluster:   wlidpkg.GetClusterFromWlid(designator.WLID),
			namespace: wlidpkg.GetNamespaceFromWlid(designator.WLID),
			kind:      wlidpkg.GetKindFromWlid(designator.WLID),
			name:      wlidpkg.GetNameFromWlid(designator.WLID),
			path:      "",
			labels:    map[string]string{},
		}
	// case DesignatorSid: // TODO
	default:
		// TODO - Do not print from here!
		// glog.Warningf("in 'digestPortalDesignator' designator type: '%v' not yet supported. please contact Armo team", designator.DesignatorType)
	}
	return AttributesDesignators{}
}

func (designator *PortalDesignator) DigestAttributesDesignator() AttributesDesignators {
	var attributes AttributesDesignators
	attr := designator.Attributes
	attributes.labels = make(map[string]string, len(attr))

	if attr == nil {
		return attributes
	}

	for k, v := range attr {
		switch k {
		case AttributeNamespace:
			attributes.namespace = v
		case AttributeCluster:
			attributes.cluster = v
		case AttributeKind:
			attributes.kind = v
		case AttributeName:
			attributes.name = v
		case AttributePath:
			attributes.path = v
		default:
			attributes.labels[k] = v
		}
	}

	return attributes
}

// DigestPortalDesignator DEPRECATED. use designator.DigestPortalDesignator() - get cluster namespace and labels from designator
func DigestPortalDesignator(designator *PortalDesignator) (string, string, map[string]string) {
	switch designator.DesignatorType {
	case DesignatorAttributes, DesignatorAttribute:
		return DigestAttributesDesignator(designator.Attributes)
	case DesignatorWlid, DesignatorWildWlid:
		return wlidpkg.GetClusterFromWlid(designator.WLID), wlidpkg.GetNamespaceFromWlid(designator.WLID), map[string]string{}
	// case DesignatorSid: // TODO
	default:
		// TODO - Do not print from here!
		// glog.Warningf("in 'digestPortalDesignator' designator type: '%v' not yet supported. please contact Armo team", designator.DesignatorType)
	}
	return "", "", nil
}
func DigestAttributesDesignator(attributes map[string]string) (string, string, map[string]string) {
	cluster := ""
	namespace := ""
	labels := map[string]string{}
	if attributes == nil {
		return cluster, namespace, labels
	}
	for k, v := range attributes {
		labels[k] = v
	}
	if v, ok := attributes[AttributeNamespace]; ok {
		namespace = v
		delete(labels, AttributeNamespace)
	}
	if v, ok := attributes[AttributeCluster]; ok {
		cluster = v
		delete(labels, AttributeCluster)
	}

	return cluster, namespace, labels
}

type mapString2String map[string]string

func (designatorMap mapString2String) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	str := ""
	err = dec.AddString(&str)
	if err != nil {
		return err
	}
	designatorMap[key] = str
	return nil
}

func (designatorMap mapString2String) NKeys() int {
	return 0
}

func (designator *PortalDesignator) UnmarshalJSONObject(dec *gojay.Decoder, key string) (err error) {
	switch key {
	case "designatorType":
		err = dec.String((*string)(&designator.DesignatorType))
	case "attributes":
		designatorAttributes := mapString2String{}
		if err = dec.Object(designatorAttributes); err == nil {
			designator.Attributes = designatorAttributes
		}
	}
	return err
}
func (designator *PortalDesignator) NKeys() int {
	return 2
}

func AttributesDesignatorsFromImageTag(imageTag string) *PortalDesignator {
	repoNameStart := strings.LastIndex(imageTag, "/")
	if repoNameStart < 0 {
		repoNameStart = len(imageTag)
	}
	tagNameStart := strings.LastIndex(imageTag, ":")
	if tagNameStart < 0 || tagNameStart < repoNameStart {
		tagNameStart = len(imageTag)
	}
	pd := &PortalDesignator{
		DesignatorType: DesignatorAttributes,
		Attributes:     make(map[string]string, 3),
	}
	pd.Attributes[AttributeRegistryName] = imageTag[:repoNameStart]

	if repoNameStart < len(imageTag)-1 {
		pd.Attributes[AttributeRepository] = imageTag[repoNameStart+1 : tagNameStart]
		if tagNameStart < len(imageTag)-1 {
			pd.Attributes[AttributeTag] = imageTag[tagNameStart+1:]
		}
	}
	return pd
}

// Getters & Setter used by derived types for interfaces implementation
func (p *PortalBase) GetGUID() string {
	return p.GUID
}
func (p *PortalBase) SetGUID(guid string) {
	p.GUID = guid
}
func (p *PortalBase) GetName() string {
	return p.Name
}
func (p *PortalBase) SetName(name string) {
	p.Name = name
}
func (p *PortalBase) GetAttributes() map[string]interface{} {
	return p.Attributes
}
func (p *PortalBase) SetAttributes(attributes map[string]interface{}) {
	p.Attributes = attributes
}

func (p *PortalBase) SetUpdatedTime(updatedTime *time.Time) {
	if updatedTime == nil {
		p.UpdatedTime = time.Now().UTC().Format(time.RFC3339)
		return
	}
	p.UpdatedTime = updatedTime.UTC().Format(time.RFC3339)
}

func (p *PortalBase) GetUpdatedTime() *time.Time {
	if p.UpdatedTime == "" {
		return nil
	}
	updatedTime, err := time.Parse(time.RFC3339, p.UpdatedTime)
	if err != nil {
		return nil
	}
	return &updatedTime
}

// GetCustomerStatus returns one of the 4 customer statuses:
// "paying" - has active subscription "Team" or "Enterprise"
// and subscription status is one of the ActiveSubscriptionStatuses
// "free" - has no active subscription or active subscription is free
// "trial" - has a "Team" active subscription, has not paid yet, and trial end has not passed
// "blocked" - has a "Team" active subscription, has not paid yet, and trial end has passed
func (p *PortalCustomer) GetCustomerStatus(now int64) CustomerStatus {
	if p.ActiveSubscription == nil || p.ActiveSubscription.LicenseType == LicenseTypeFree {
		return FreeCustomer
	}

	if p.ActiveSubscription.LicenseType == LicenseTypeTeam {
		if !slices.Contains(ActiveSubscriptionStatuses, p.ActiveSubscription.SubscriptionStatus) {
			if p.ActiveSubscription.TrialEnd > now {
				return TrialCustomer
			} else {
				return BlockedCustomer
			}
		} else {
			return PayingCustomer
		}
	}

	if p.ActiveSubscription.LicenseType == LicenseTypeEnterprise {
		if slices.Contains(ActiveSubscriptionStatuses, p.ActiveSubscription.SubscriptionStatus) {
			return PayingCustomer
		}
	}

	return FreeCustomer
}

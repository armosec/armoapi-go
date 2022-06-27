package armotypes

import (
	wlidpkg "github.com/armosec/utils-k8s-go/wlid"
	"github.com/francoispqt/gojay"
)

var IgnoreLabels = []string{AttributeCluster, AttributeNamespace}

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
	cluster, _, _, _, _ := designator.DigestPortalDesignator()
	return cluster
}

func (designator *PortalDesignator) GetNamespace() string {
	_, namespace, _, _, _ := designator.DigestPortalDesignator()
	return namespace
}

func (designator *PortalDesignator) GetKind() string {
	_, _, kind, _, _ := designator.DigestPortalDesignator()
	return kind
}

func (designator *PortalDesignator) GetName() string {
	_, _, _, name, _ := designator.DigestPortalDesignator()
	return name
}
func (designator *PortalDesignator) GetLabels() map[string]string {
	_, _, _, _, labels := designator.DigestPortalDesignator()
	return labels
}

// DigestPortalDesignator - get cluster namespace and labels from designator
func (designator *PortalDesignator) DigestPortalDesignator() (string, string, string, string, map[string]string) {
	switch designator.DesignatorType {
	case DesignatorAttributes, DesignatorAttribute:
		return designator.DigestAttributesDesignator()
	case DesignatorWlid.ToLower(), DesignatorWildWlid.ToLower():
		return wlidpkg.GetClusterFromWlid(designator.WLID), wlidpkg.GetNamespaceFromWlid(designator.WLID), wlidpkg.GetKindFromWlid(designator.WLID), wlidpkg.GetNameFromWlid(designator.WLID), map[string]string{}
	// case DesignatorSid: // TODO
	default:
		// TODO - Do not print from here!
		// glog.Warningf("in 'digestPortalDesignator' designator type: '%v' not yet supported. please contact Armo team", designator.DesignatorType)
	}
	return "", "", "", "", nil
}

func (designator *PortalDesignator) DigestAttributesDesignator() (string, string, string, string, map[string]string) {
	cluster := ""
	namespace := ""
	kind := ""
	name := ""
	labels := map[string]string{}
	attributes := designator.Attributes
	if attributes == nil {
		return cluster, namespace, kind, name, labels
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
	if v, ok := attributes[AttributeKind]; ok {
		kind = v
		delete(labels, AttributeKind)
	}
	if v, ok := attributes[AttributeName]; ok {
		name = v
		delete(labels, AttributeName)
	}
	return cluster, namespace, kind, name, labels
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

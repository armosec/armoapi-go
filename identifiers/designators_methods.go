package identifiers

import (
	"strings"

	"github.com/armosec/gojay"
	wlidpkg "github.com/armosec/utils-k8s-go/wlid"
)

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

func (designator *PortalDesignator) GetResourceID() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.resourceID
}

func (designator *PortalDesignator) GetK8sResourceHash() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.k8sResourceHash

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
		case AttributeResourceID:
			attributes.resourceID = v
		case AttributeK8sResourceHash:
			attributes.k8sResourceHash = v
		default:
			attributes.labels[k] = v
		}
	}

	return attributes
}

func (designator *PortalDesignator) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	var err error
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

func (ad *AttributesDesignators) GetResourceID() string {
	return ad.resourceID
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

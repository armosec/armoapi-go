package armotypes

import (
	"strings"

	wlidpkg "github.com/armosec/utils-k8s-go/wlid"
	"github.com/francoispqt/gojay"
)

var IgnoreLabels = []string{AttributeCluster, AttributeNamespace}

type attributesDesignators struct {
	Cluster   string
	Namespace string
	Kind      string
	Name      string
	Path      string
	Labels    map[string]string
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
	return attributes.Cluster
}

func (designator *PortalDesignator) GetNamespace() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.Namespace
}

func (designator *PortalDesignator) GetKind() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.Kind
}

func (designator *PortalDesignator) GetName() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.Name
}
func (designator *PortalDesignator) GetPath() string {
	attributes := designator.DigestPortalDesignator()
	return attributes.Path
}
func (designator *PortalDesignator) GetLabels() map[string]string {
	attributes := designator.DigestPortalDesignator()
	return attributes.Labels
}

// DigestPortalDesignator - get cluster namespace and labels from designator
func (designator *PortalDesignator) DigestPortalDesignator() attributesDesignators {
	switch designator.DesignatorType {
	case DesignatorAttributes, DesignatorAttribute:
		return designator.DigestAttributesDesignator()
	case DesignatorWlid.ToLower(), DesignatorWildWlid.ToLower():
		return attributesDesignators{wlidpkg.GetClusterFromWlid(designator.WLID), wlidpkg.GetNamespaceFromWlid(designator.WLID), wlidpkg.GetKindFromWlid(designator.WLID), wlidpkg.GetNameFromWlid(designator.WLID), "", map[string]string{}}
	// case DesignatorSid: // TODO
	default:
		// TODO - Do not print from here!
		// glog.Warningf("in 'digestPortalDesignator' designator type: '%v' not yet supported. please contact Armo team", designator.DesignatorType)
	}
	return attributesDesignators{}
}

func (designator *PortalDesignator) DigestAttributesDesignator() attributesDesignators {
	var attributes attributesDesignators
	attributes.Labels = map[string]string{}
	attr := designator.Attributes
	if attr == nil {
		return attributes
	}
	for k, v := range attr {
		attributes.Labels[k] = v
	}
	if v, ok := attr[AttributeNamespace]; ok {
		attributes.Namespace = v
		delete(attributes.Labels, AttributeNamespace)
	}
	if v, ok := attr[AttributeCluster]; ok {
		attributes.Cluster = v
		delete(attributes.Labels, AttributeCluster)
	}
	if v, ok := attr[AttributeKind]; ok {
		attributes.Kind = v
		delete(attributes.Labels, AttributeKind)
	}
	if v, ok := attr[AttributeName]; ok {
		attributes.Name = v
		delete(attributes.Labels, AttributeName)
	}
	if v, ok := attr[AttributePath]; ok {
		attributes.Path = v
		delete(attributes.Labels, AttributePath)
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

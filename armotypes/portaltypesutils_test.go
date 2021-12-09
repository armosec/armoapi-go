package armotypes

import "testing"

func TestAttributesDesignatorsFromWLID(t *testing.T) {
	attDesig := AttributesDesignatorsFromWLID("wlid://cluster-liortest1/namespace-default/deployment-payment")
	if attDesig.Attributes[AttributeCluster] != "liortest1" ||
		attDesig.Attributes[AttributeNamespace] != "default" ||
		attDesig.Attributes[AttributeKind] != "deployment" ||
		attDesig.Attributes[AttributeName] != "payment" {
		t.Errorf("wrong attributes desigantors:%v", attDesig)
	}

	attDesig = AttributesDesignatorsFromWLID("wlid://cluster-liortest1/namespace-default/deployment")
	if attDesig.Attributes[AttributeCluster] != "liortest1" ||
		attDesig.Attributes[AttributeNamespace] != "default" ||
		attDesig.Attributes[AttributeKind] != "deployment" {
		t.Errorf("wrong attributes desigantors:%v", attDesig)
	}
	attDesig = AttributesDesignatorsFromWLID("wlid://cluster-liortest1/namespace-default/")
	if attDesig.Attributes[AttributeCluster] != "liortest1" ||
		attDesig.Attributes[AttributeNamespace] != "default" {
		t.Errorf("wrong attributes desigantors:%v", attDesig)
	}
	attDesig = AttributesDesignatorsFromWLID("wlid://cluster-liortest1")
	if attDesig.Attributes[AttributeCluster] != "liortest1" {
		t.Errorf("wrong attributes desigantors:%v", attDesig)
	}
}

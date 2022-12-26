package armotypes

import (
	_ "embed"
	"strings"
	"testing"
	"time"

	"github.com/francoispqt/gojay"
	"github.com/stretchr/testify/assert"
)

var attribute = map[string]string{AttributeCluster: "cluster1", AttributeNamespace: "namespace1", AttributeKind: "kind1", AttributeName: "name1", AttributePath: "path1"}
var portalDesignator = PortalDesignator{DesignatorType: DesignatorAttribute, Attributes: attribute}

func TestGetCluster(t *testing.T) {
	er := attribute[AttributeCluster]
	cluster := portalDesignator.GetCluster()
	assert.Equal(t, er, cluster)
}

func TestGetNamespace(t *testing.T) {
	er := attribute[AttributeNamespace]
	namespace := portalDesignator.GetNamespace()
	assert.Equal(t, er, namespace)
}

func TestGetKind(t *testing.T) {
	er := attribute[AttributeKind]
	kind := portalDesignator.GetKind()
	assert.Equal(t, er, kind)
}

func TestGetName(t *testing.T) {
	er := attribute[AttributeName]
	name := portalDesignator.GetName()
	assert.Equal(t, er, name)
}

func TestGetPath(t *testing.T) {
	er := attribute[AttributePath]
	path := portalDesignator.GetPath()
	assert.Equal(t, er, path)
}

func TestSetUpdatedTime(t *testing.T) {
	now := time.Now()
	nowString := now.UTC().Format(time.RFC3339)
	
	validDateString:= "2022-12-26T15:05:23Z"
	validDate, _ := time.Parse(time.RFC3339, validDateString)

	type testCase struct {
		name     string
		time     *time.Time
		expected PortalBase
	}
	testTable := []testCase{
		{
			name:     "valid time",
			time:     &validDate,
			expected: PortalBase{UpdatedTime: validDateString},
		},
		{
			name:     "default time",
			time:     nil,
			expected: PortalBase{UpdatedTime: nowString},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			p := PortalBase{}
			p.SetUpdatedTime(test.time)
			assert.Equal(t, test.expected, p)
		})
	}
}

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

//go:embed fixtures/designatorTestCase.json
var designatorTestCase string

func TestDesignatorDecoding(t *testing.T) {
	designator := &PortalDesignator{}
	er := gojay.NewDecoder(strings.NewReader(designatorTestCase)).DecodeObject(designator)
	if er != nil {
		t.Errorf("decode failed due to: %v", er.Error())
	}
	assert.Equal(t, DesignatorAttributes, designator.DesignatorType)
	assert.Equal(t, "myCluster", designator.Attributes[AttributeCluster])
	assert.Equal(t, "8190928904639901517", designator.Attributes[AttributeWorkloadHash])
	assert.Equal(t, "myName", designator.Attributes[AttributeName])
	assert.Equal(t, "myNS", designator.Attributes[AttributeNamespace])
	assert.Equal(t, "deployment", designator.Attributes[AttributeKind])
	assert.Equal(t, "e57ec5a0-695f-4777-8366-1c64fada00a0", designator.Attributes[AttributeCustomerGUID])
	assert.Equal(t, "myContainer", designator.Attributes[AttributeContainerName])
}

func TestAttributesDesignatorsFromImageTag(t *testing.T) {
	deisgs := AttributesDesignatorsFromImageTag("docker.elastic.co/elasticsearch/elasticsearch:7.9.2")

	assert.Equal(t, "docker.elastic.co/elasticsearch", deisgs.Attributes[AttributeRegistryName])
	assert.Equal(t, "elasticsearch", deisgs.Attributes[AttributeRepository])
	assert.Equal(t, "7.9.2", deisgs.Attributes[AttributeTag])
	assert.Equal(t, 3, len(deisgs.Attributes))

	deisgs = AttributesDesignatorsFromImageTag("docker.elastic.co/elasticsearch/elasticsearch")

	assert.Equal(t, "docker.elastic.co/elasticsearch", deisgs.Attributes[AttributeRegistryName])
	assert.Equal(t, "elasticsearch", deisgs.Attributes[AttributeRepository])
	assert.Equal(t, "", deisgs.Attributes[AttributeTag])
	assert.Equal(t, 2, len(deisgs.Attributes))

	deisgs = AttributesDesignatorsFromImageTag("docker.elastic.co/elasticsearch")

	assert.Equal(t, "docker.elastic.co", deisgs.Attributes[AttributeRegistryName])
	assert.Equal(t, "elasticsearch", deisgs.Attributes[AttributeRepository])
	assert.Equal(t, "", deisgs.Attributes[AttributeTag])
	assert.Equal(t, 2, len(deisgs.Attributes))

	deisgs = AttributesDesignatorsFromImageTag("docker.elastic.co/")

	assert.Equal(t, "docker.elastic.co", deisgs.Attributes[AttributeRegistryName])
	assert.Equal(t, "", deisgs.Attributes[AttributeRepository])
	assert.Equal(t, "", deisgs.Attributes[AttributeTag])
	assert.Equal(t, 1, len(deisgs.Attributes))

	deisgs = AttributesDesignatorsFromImageTag("docker.elastic.co")

	assert.Equal(t, "docker.elastic.co", deisgs.Attributes[AttributeRegistryName])
	assert.Equal(t, "", deisgs.Attributes[AttributeRepository])
	assert.Equal(t, "", deisgs.Attributes[AttributeTag])
	assert.Equal(t, 1, len(deisgs.Attributes))

	deisgs = AttributesDesignatorsFromImageTag("")

	assert.Equal(t, "", deisgs.Attributes[AttributeRegistryName])
	assert.Equal(t, "", deisgs.Attributes[AttributeRepository])
	assert.Equal(t, "", deisgs.Attributes[AttributeTag])
	assert.Equal(t, 1, len(deisgs.Attributes))
}

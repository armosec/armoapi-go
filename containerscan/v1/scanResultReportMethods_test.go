package v1

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/armosec/armoapi-go/apis"
	"github.com/armosec/armoapi-go/containerscan"
	"github.com/armosec/armoapi-go/identifiers"
	"github.com/armosec/gojay"
	"github.com/stretchr/testify/assert"
)

func TestIsLastReport(t *testing.T) {
	lastReport := ScanResultReport{PaginationInfo: apis.PaginationMarks{IsLastReport: true}}
	if !lastReport.IsLastReport() {
		t.Errorf("Expected IsLastReport() to return true for the last report, but it returned false")
	}

	notLastReport := ScanResultReport{PaginationInfo: apis.PaginationMarks{IsLastReport: false}}
	if notLastReport.IsLastReport() {
		t.Errorf("Expected IsLastReport() to return false for a report that is not the last, but it returned true")
	}
}

func TestGetContainerScanID(t *testing.T) {
	report := ScanResultReport{ContainerScanID: "12345"}
	if report.GetContainerScanID() != "12345" {
		t.Errorf("Expected GetContainerScanID() to return \"12345\", but it returned %q", report.GetContainerScanID())
	}
}

func TestGetTimestamp(t *testing.T) {
	report := ScanResultReport{Timestamp: 1234567890}
	if report.GetTimestamp() != 1234567890 {
		t.Errorf("Expected GetTimestamp() to return 1234567890, but it returned %d", report.GetTimestamp())
	}
}

func TestGetWorkloadHash(t *testing.T) {
	report := ScanResultReport{Designators: identifiers.PortalDesignator{Attributes: map[string]string{"workloadHash": "hash123"}}}
	if report.GetWorkloadHash() != "hash123" {
		t.Errorf("Expected GetWorkloadHash() to return \"hash123\", but it returned %q", report.GetWorkloadHash())
	}
}

func TestGetCustomerGUID(t *testing.T) {
	report := ScanResultReport{Summary: &containerscan.CommonContainerScanSummaryResult{CustomerGUID: "abc123"}}
	if report.GetCustomerGUID() != "abc123" {
		t.Errorf("Expected GetCustomerGUID() to return \"abc123\", but it returned %q", report.GetCustomerGUID())
	}
}

func TestSetWorkloadHash(t *testing.T) {
	designators := identifiers.PortalDesignator{Attributes: map[string]string{"workloadHash": ""}}
	report := ScanResultReport{Designators: designators}

	workloadHash := "abc123"
	report.SetWorkloadHash(workloadHash)

	if report.GetDesignators().Attributes["workloadHash"] != workloadHash {
		t.Errorf("SetWorkloadHash failed, expected %s but got %s", workloadHash, report.GetDesignators().Attributes["workloadHash"])
	}
}

func TestSetCustomerGUID(t *testing.T) {
	summary := containerscan.CommonContainerScanSummaryResult{CustomerGUID: ""}
	report := ScanResultReport{Summary: &summary}

	customerGUID := "123abc"
	report.SetCustomerGUID(customerGUID)

	if report.GetSummary().GetCustomerGUID() != customerGUID {
		t.Errorf("SetCustomerGUID failed, expected %s but got %s", customerGUID, report.GetSummary().GetCustomerGUID())
	}
}

func TestGetSummary(t *testing.T) {
	report := ScanResultReport{Summary: &containerscan.CommonContainerScanSummaryResult{}}

	assert.Equal(t, report.GetSummary(), report.Summary)

	report.Summary = nil
	assert.Nil(t, report.GetSummary())
}

func TestGetPaginationInfo(t *testing.T) {
	report := ScanResultReport{PaginationInfo: apis.PaginationMarks{IsLastReport: true}}
	assert.Equal(t, report.GetPaginationInfo(), report.PaginationInfo)
}

func TestSetPaginationInfo(t *testing.T) {
	report := ScanResultReport{}
	paginationInfo := apis.PaginationMarks{IsLastReport: true}
	report.SetPaginationInfo(paginationInfo)
	assert.Equal(t, report.GetPaginationInfo(), paginationInfo)
}

//go:embed testdata/scanReportV1TestCase.json
var scanReportTestCase string

func TestScanResultReportDecoding(t *testing.T) {
	scanReport := &ScanResultReport{}
	er := gojay.NewDecoder(strings.NewReader(scanReportTestCase)).DecodeObject(scanReport)
	if er != nil {
		t.Errorf("decode failed due to: %v", er.Error())
	}
	assert.Equal(t, "5969736482532194479", scanReport.ContainerScanID)
	assert.Equal(t, int64(1656250322), scanReport.Timestamp)
	assert.True(t, scanReport.Validate(), "cannot validate scan report after gojay decoding")
	assert.Equal(t, identifiers.DesignatorAttributes, scanReport.Designators.DesignatorType)
	assert.Equal(t, "myCluster", scanReport.Designators.Attributes[identifiers.AttributeCluster])
	assert.Equal(t, "8190928904639901517", scanReport.Designators.Attributes[identifiers.AttributeWorkloadHash])
	assert.Equal(t, "myName", scanReport.Designators.Attributes[identifiers.AttributeName])
	assert.Equal(t, "myNS", scanReport.Designators.Attributes[identifiers.AttributeNamespace])
	assert.Equal(t, "deployment", scanReport.Designators.Attributes[identifiers.AttributeKind])
	assert.Equal(t, "e57ec5a0-695f-4777-8366-1c64fada00a0", scanReport.Designators.Attributes[identifiers.AttributeCustomerGUID])
	assert.Equal(t, "myContainer", scanReport.Designators.Attributes[identifiers.AttributeContainerName])

}

func TestSetContainerScanID(t *testing.T) {
	report := ScanResultReport{}

	report.SetContainerScanID("12345")
	if report.GetContainerScanID() != "12345" {
		t.Error("Expected SetContainerScanID() to set the container scan ID of a ScanResultReportV1, but it was not set")
	}
}

func TestSetTimestamp(t *testing.T) {
	report := ScanResultReport{Timestamp: 0}

	report.SetTimestamp(12345)
	if report.GetTimestamp() != 12345 {
		t.Errorf("Expected SetTimestamp() to set the timestamp of a ScanResultReportV1 to 12345, but it was set to %d", report.GetTimestamp())
	}
}

func TestGetVulnerabilities(t *testing.T) {
	report := ScanResultReport{Vulnerabilities: []containerscan.CommonContainerVulnerabilityResult{
		{
			Vulnerability: containerscan.Vulnerability{
				Name:    "CVE-1",
				ImageID: "sha256:1",
			},
			ContainerScanID: "12345",
		},
		{
			Vulnerability: containerscan.Vulnerability{
				Name:    "CVE-2",
				ImageID: "sha256:2",
			},
			ContainerScanID: "12345",
		},
		{
			Vulnerability: containerscan.Vulnerability{
				Name:    "CVE-3",
				ImageID: "sha256:3",
			},
			ContainerScanID: "12345",
		},
		{
			Vulnerability: containerscan.Vulnerability{
				Name:    "CVE-4",
				ImageID: "sha256:4",
			},
			ContainerScanID: "3333",
		},
	}}

	vulns := report.GetVulnerabilities()

	assert.Equal(t, 4, len(vulns))

	assert.Equal(t, "CVE-1", vulns[0].GetVulnerability().GetName())
	assert.Equal(t, "CVE-2", vulns[1].GetVulnerability().GetName())
	assert.Equal(t, "CVE-3", vulns[2].GetVulnerability().GetName())
	assert.Equal(t, "CVE-4", vulns[3].GetVulnerability().GetName())

	assert.Equal(t, "sha256:1", vulns[0].GetVulnerability().GetImageID())
	assert.Equal(t, "sha256:2", vulns[1].GetVulnerability().GetImageID())
	assert.Equal(t, "sha256:3", vulns[2].GetVulnerability().GetImageID())
	assert.Equal(t, "sha256:4", vulns[3].GetVulnerability().GetImageID())

	assert.Equal(t, "12345", vulns[0].GetContainerScanID())
	assert.Equal(t, "12345", vulns[1].GetContainerScanID())
	assert.Equal(t, "12345", vulns[2].GetContainerScanID())
	assert.Equal(t, "3333", vulns[3].GetContainerScanID())

}

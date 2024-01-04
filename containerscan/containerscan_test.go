package containerscan

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/armoapi-go/identifiers"
	"github.com/armosec/gojay"
	"github.com/stretchr/testify/assert"
)

func TestImageRegisteryInformation(t *testing.T) {
	mock := ScanResultReport{Designators: identifiers.PortalDesignator{WLID: "wlid://cluster-testc,namespace-testns/deployment-testname",
		Attributes: map[string]string{"registryName": "gcr.io/elated-pottery-310110",
			"project":  "testProject",
			"chikmook": "chikmook",
		}}}

	summary := mock.Summarize()
	assert.Equal(t, summary.Designators.Attributes["registryName"], mock.Designators.Attributes["registryName"], "missing registryName")
	assert.Equal(t, summary.Designators.Attributes["project"], mock.Designators.Attributes["project"], "missing project")
	assert.Equal(t, summary.Designators.Attributes["chikmook"], mock.Designators.Attributes["chikmook"], "missing random property")

}

func TestDecodeScanWIthDangearousArtifacts(t *testing.T) {
	rhs := &ScanResultReport{}
	er := gojay.NewDecoder(strings.NewReader(nginxScanJSON)).DecodeObject(rhs)
	if er != nil {
		t.Errorf("decode failed due to: %v", er.Error())
	}
	sumObj := rhs.Summarize()
	if sumObj.Registry != "" {
		t.Errorf("sumObj.Registry = %v", sumObj.Registry)
	}
	if sumObj.ImageTagSuffix != "nginx:1.18.0" {
		t.Errorf("sumObj.VersionImage = %v", sumObj.Registry)
	}
	if sumObj.ImageTag != "nginx:1.18.0" {
		t.Errorf("sumObj.ImgTag = %v", sumObj.ImageTag)
	}
	if sumObj.Status != "Success" {
		t.Errorf("sumObj.Status = %v", sumObj.Status)
	}
}

//go:embed fixtures/scanReportV1TestCase.json
var scanReportV1TestCase string

func TestExceptions(t *testing.T) {
	rhs := &ScanResultReport{}
	er := gojay.NewDecoder(strings.NewReader(nginxScanJSON)).DecodeObject(rhs)
	if er != nil {
		t.Errorf("decode failed due to: %v", er.Error())
	}
	regularSum := rhs.Summarize()
	exception := armotypes.MockVulnerabilityException()
	rhs.Layers[0].Vulnerabilities[0].ExceptionApplied = append(rhs.Layers[0].Vulnerabilities[0].ExceptionApplied, *exception)
	rhs.Layers[0].Vulnerabilities[1].ExceptionApplied = append(rhs.Layers[0].Vulnerabilities[1].ExceptionApplied, *exception)

	sumObj := rhs.Summarize()
	if len(sumObj.ExcludedSeveritiesStats) != 1 {
		t.Errorf("len(sumObj.ExcludedSeveritiesStats) = %v", len(sumObj.ExcludedSeveritiesStats))
	}
	excludedStats := sumObj.ExcludedSeveritiesStats[0]

	// Alone
	if excludedStats.Severity != "Medium" {
		t.Errorf("excludedStats.Severity = %v", excludedStats.Severity)
	}
	if excludedStats.TotalCount != 2 {
		t.Errorf("excludedStats.TotalCount = %v", excludedStats.TotalCount)
	}
	if excludedStats.FixAvailableOfTotalCount != 2 {
		t.Errorf("excludedStats.FixAvailableOfTotalCount = %v", excludedStats.FixAvailableOfTotalCount)
	}
	if excludedStats.RelevantCount != 0 {
		t.Errorf("excludedStats.RelevantCount = %v", excludedStats.RelevantCount)
	}
	if excludedStats.RelevantFixCount != 0 {
		t.Errorf("excludedStats.FixAvailableForRelevantCount = %v", excludedStats.RelevantFixCount)
	}
	if excludedStats.RCECount != 1 {
		t.Errorf("excludedStats.RCECount = %v", excludedStats.RCECount)
	}
	if excludedStats.UrgentCount != 0 {
		t.Errorf("excludedStats.UrgentCount = %v", excludedStats.UrgentCount)
	}
	if excludedStats.NeglectedCount != 0 {
		t.Errorf("excludedStats.NeglectedCount = %v", excludedStats.NeglectedCount)
	}

	// With-exceptions-summary VS regular
	if regularSum.TotalCount != (sumObj.TotalCount + excludedStats.TotalCount) {
		t.Errorf("sumObj.TotalCount = %v", sumObj.TotalCount)
	}
	if regularSum.FixAvailableOfTotalCount != (sumObj.FixAvailableOfTotalCount + excludedStats.FixAvailableOfTotalCount) {
		t.Errorf("sumObj.FixAvailableOfTotalCount = %v", sumObj.FixAvailableOfTotalCount)
	}
	if regularSum.RelevantCount != (sumObj.RelevantCount + excludedStats.RelevantCount) {
		t.Errorf("sumObj.RelevantCount = %v", sumObj.RelevantCount)
	}
	if regularSum.RelevantFixCount != (sumObj.RelevantFixCount + excludedStats.RelevantFixCount) {
		t.Errorf("sumObj.FixAvailableForRelevantCoun = %v", sumObj.RelevantFixCount)
	}
	if regularSum.RCECount != (sumObj.RCECount + excludedStats.RCECount) {
		t.Errorf("sumObj.RCECount = %v", sumObj.RCECount)
	}
	if regularSum.UrgentCount != (sumObj.UrgentCount + excludedStats.UrgentCount) {
		t.Errorf("sumObj.UrgentCount = %v", sumObj.UrgentCount)
	}
	if regularSum.NeglectedCount != (sumObj.NeglectedCount + excludedStats.NeglectedCount) {
		t.Errorf("sumObj.NeglectedCount = %v", sumObj.NeglectedCount)
	}
	// if regularSum.HealthStatus != (sumObj.HealthStatus + excludedStats.HealthStatus) {
	// 	t.Errorf("sumObj.HealthStatus = %v", sumObj.HealthStatus)
	// }
	// if excludedStats.HealthStatus                 != val {
	// 	t.Errorf("excludedStats.HealthStatus = %v", excludedStats.HealthStatus)
	// }
}

func TestUnmarshalScanReport(t *testing.T) {
	ds := GenerateContainerScanReportMock(GenerateVulnerability)
	str1 := ds.AsFNVHash()

	str2 := identifiers.CalcHashFNV(fmt.Sprintf("%v", ds))

	if str1 != str2 {
		t.Errorf("hashes are different: %v != %v", str1, str2)
	}

	rhs := &ScanResultReport{}

	bolB, _ := json.Marshal(ds)
	r := bytes.NewReader(bolB)

	er := gojay.NewDecoder(r).DecodeObject(rhs)
	if er != nil {
		t.Errorf("marshalling failed due to: %v", er.Error())
	}

	if rhs.AsFNVHash() != str1 {
		t.Errorf("marshalling failed different values after marshal:\nOriginal:\n%v\nParsed:\n%v\n\n===\n", string(bolB), rhs)
	}
}

func TestRCEFixCount(t *testing.T) {
	// RCE and fixed
	ds := GenerateContainerScanReportMock(GenerateVulnerabilityRCEAndFixed)
	summary := ds.Summarize()
	assert.Equal(t, summary.RCECount, summary.FixAvailableOfTotalCount)
	assert.Equal(t, summary.SeveritiesStats[0].RCEFixCount, summary.RCEFixCount)
	assert.Equal(t, summary.RCEFixCount, summary.RCECount)

	// RCE not fixed
	ds = GenerateContainerScanReportMock(GenerateVulnerabilityRCENotFixed)
	summary = ds.Summarize()
	assert.NotEqual(t, summary.RCECount, int64(0))
	assert.Equal(t, summary.FixAvailableOfTotalCount, int64(0))
	assert.Equal(t, summary.SeveritiesStats[0].RCEFixCount, summary.RCEFixCount)
	assert.Equal(t, summary.RCEFixCount, int64(0))

	//No RCE and fixed
	ds = GenerateContainerScanReportMock(GenerateVulnerabilityNoRCEAndFixed)
	summary = ds.Summarize()
	assert.Equal(t, summary.RCECount, int64(0))
	assert.NotEqual(t, summary.FixAvailableOfTotalCount, int64(0))
	assert.Equal(t, summary.SeveritiesStats[0].RCEFixCount, summary.RCEFixCount)
	assert.Equal(t, summary.RCEFixCount, int64(0))

	//No RCE and no fix
	ds = GenerateContainerScanReportMock(GenerateVulnerabilityNoRCENoFixed)
	summary = ds.Summarize()
	assert.Equal(t, summary.FixAvailableOfTotalCount, int64(0))
	assert.Equal(t, summary.RCEFixCount, int64(0))
	assert.Equal(t, summary.SeveritiesStats[0].RCEFixCount, summary.RCEFixCount)
	assert.Equal(t, summary.RCECount, int64(0))
}
func TestUnmarshalScanReport1(t *testing.T) {
	ds := Vulnerability{}
	if err := GenerateVulnerability(&ds); err != nil {
		t.Errorf("%v\n%v\n", ds, err)
	}
}

func TestGetByPkgNameSuccess(t *testing.T) {
	ds := GenerateContainerScanReportMock(GenerateVulnerability)
	a := ds.Layers[0].GetFilesByPackage("coreutils")
	if a != nil {

		fmt.Printf("%+v\n", *a)
	}

}

func TestGetByPkgNameMissing(t *testing.T) {
	ds := GenerateContainerScanReportMock(GenerateVulnerability)
	a := ds.Layers[0].GetFilesByPackage("s")
	if a != nil && len(*a) > 0 {
		t.Errorf("expected - no such package should be in that layer %v\n\n; found - %v", ds, a)
	}

}

func TestCalculateFixed(t *testing.T) {
	res := CalculateFixed([]FixedIn{{
		Name:    "",
		ImgTag:  "",
		Version: "",
	}})
	if res != 0 {
		t.Errorf("wrong fix status: %v", res)
	}
}

func TestIsRCE(t *testing.T) {
	ds := Vulnerability{}

	ds.Description = "Online Railway Reservation System 1.0 - Remote Code Execution (RCE) (Unauthenticated)"
	if true != ds.IsRCE() {
		t.Errorf("IsRCE failed")
	}
	ds.Description = "Gerapy 0.9.7 - Remote Code Execution (RCE) (Authenticated)"
	if true != ds.IsRCE() {
		t.Errorf("IsRCE failed")
	}
	ds.Description = "FORCEHENEW"
	if false != ds.IsRCE() {
		t.Errorf("IsRCE failed")
	}
}

func TestReportValidate(t *testing.T) {
	scanresult := &ScanResultReport{}

	if scanresult.Validate() {
		t.Error("empty scan passed validation")
	}
	scanresult.Timestamp = time.Now().Unix()
	scanresult.ImgHash = "fsdfsdf"
	scanresult.ImgTag = "yuy43434"
	scanresult.CustomerGUID = "<MY_GUID>"
	if scanresult.Validate() {
		t.Error("invalid customer guid passed validation")
	}
	scanresult.CustomerGUID = ""
	if scanresult.Validate() {
		t.Error("empty CustomerGUID passed validation")
	}
	scanresult.ImgHash = ""
	scanresult.ImgTag = ""
	scanresult.CustomerGUID = "8c338c97-383e-4083-a42f-d9b4e0448b13"
	if scanresult.Validate() {
		t.Error("empty scan passed validation")
	}
	scanresult.Timestamp = 0
	scanresult.ImgHash = "fsdfsdf"
	scanresult.ImgTag = "yuy43434"
	if scanresult.Validate() {
		t.Error("empty timestamp passed validation")
	}
	scanresult.Timestamp = time.Now().Unix()
	scanresult.ImgHash = "fsdfsdf"
	scanresult.ImgTag = "yuy43434"
	if !scanresult.Validate() {
		t.Error("valid timestamp failed the validation")
	}
}

func TestVulnerabilityToShort(t *testing.T) {
	vul := Vulnerability{
		Name:               "name",
		ImageID:            "imageHash",
		ImageTag:           "imageTag",
		RelatedPackageName: "packageName",
		PackageType:        "packageType",
		PackageVersion:     "packageVersion",
		Link:               "link",
		Description:        "description",
		Severity:           "severity",
		SeverityScore:      5,
	}
	short := vul.ToShortVulnerabilityResult()
	if short.Name != vul.Name {
		t.Errorf("ToShortVulnerabilityResult failed")
	}
}

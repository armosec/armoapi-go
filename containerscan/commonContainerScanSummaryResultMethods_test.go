package containerscan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWLID(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		WLID: "123",
	}
	if summary.GetWLID() != "123" {
		t.Errorf("Expected WLID to be %s, but got %s", "123", summary.GetWLID())
	}
}

func TestGetClusterName(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		ClusterName: "cluster1",
	}
	if summary.GetClusterName() != "cluster1" {
		t.Errorf("Expected ClusterName to be %s, but got %s", "cluster1", summary.GetClusterName())
	}
}

func TestGetNamespace(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		Namespace: "namespace1",
	}
	if summary.GetNamespace() != "namespace1" {
		t.Errorf("Expected Namespace to be %s, but got %s", "namespace1", summary.GetNamespace())
	}
}

func TestGetContainerName(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		ContainerName: "container1",
	}
	if summary.GetContainerName() != "container1" {
		t.Errorf("Expected ContainerName to be %s, but got %s", "container1", summary.GetContainerName())
	}
}

func TestGetStatus(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		Status: "running",
	}
	if summary.GetStatus() != "running" {
		t.Errorf("Expected Status to be %s, but got %s", "running", summary.GetStatus())
	}
}

func TestGetRegistry(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		Registry: "docker.io",
	}
	if summary.GetRegistry() != "docker.io" {
		t.Errorf("Expected Registry to be %s, but got %s", "docker.io", summary.GetRegistry())
	}
}

func TestGetImageTageSuffix(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		ImageTagSuffix: "latest",
	}
	if summary.GetImageTageSuffix() != "latest" {
		t.Errorf("Expected ImageTagSuffix to be %s, but got %s", "latest", summary.GetImageTageSuffix())
	}
}

func TestGetVersion(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		Version: "1.0.0",
	}
	if summary.GetVersion() != "1.0.0" {
		t.Errorf("Expected Version to be %s, but got %s", "1.0.0", summary.GetVersion())
	}
}

func TestCommonContainerScanSummaryResult_GetCustomerGUID(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		CustomerGUID: "123456",
	}
	if summary.GetCustomerGUID() != "123456" {
		t.Errorf("GetCustomerGUID returned unexpected value")
	}
}

func TestCommonContainerScanSummaryResult_GetContainerScanID(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		ContainerScanID: "abcd1234",
	}
	if summary.GetContainerScanID() != "abcd1234" {
		t.Errorf("GetContainerScanID returned unexpected value")
	}
}

func TestCommonContainerScanSummaryResult_GetTimestamp(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		Timestamp: 1234567890,
	}
	if summary.GetTimestamp() != 1234567890 {
		t.Errorf("GetTimestamp returned unexpected value")
	}
}

func TestCommonContainerScanSummaryResult_GetJobIDs(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		JobIDs: []string{"job1", "job2", "job3"},
	}
	jobIDs := summary.GetJobIDs()
	if len(jobIDs) != 3 || jobIDs[0] != "job1" || jobIDs[1] != "job2" || jobIDs[2] != "job3" {
		t.Errorf("GetJobIDs returned unexpected value")
	}
}

func TestCommonContainerScanSummaryResult_Validate(t *testing.T) {
	// Test case 1: valid summary
	summary1 := &CommonContainerScanSummaryResult{
		CustomerGUID:    "123456",
		ContainerScanID: "abcd1234",
		ImageTag:        "latest",
		Timestamp:       1234567890,
	}
	if !summary1.Validate() {
		t.Errorf("Validate returned unexpected value for valid summary")
	}

	// Test case 2: missing CustomerGUID
	summary2 := &CommonContainerScanSummaryResult{
		ContainerScanID: "abcd1234",
		ImageTag:        "latest",
		Timestamp:       1234567890,
	}
	if summary2.Validate() {
		t.Errorf("Validate returned unexpected value for summary with missing CustomerGUID")
	}

	// Test case 3: missing ContainerScanID
	summary3 := &CommonContainerScanSummaryResult{
		CustomerGUID: "123456",
		ImageTag:     "latest",
		Timestamp:    1234567890,
	}
	if summary3.Validate() {
		t.Errorf("Validate returned unexpected value for summary with missing ContainerScanID")
	}

	// Test case 4: missing ImageTag and ImageID
	summary4 := &CommonContainerScanSummaryResult{
		CustomerGUID:    "123456",
		ContainerScanID: "abcd1234",
		Timestamp:       1234567890,
	}
	if summary4.Validate() {
		t.Errorf("Validate returned unexpected value for summary with missing ImageTag and ImageID")
	}

	// Test case 5: invalid Timestamp
	summary5 := &CommonContainerScanSummaryResult{
		CustomerGUID:    "123456",
		ContainerScanID: "abcd1234",
		ImageTag:        "latest",
		Timestamp:       -1234567890,
	}
	if summary5.Validate() {
		t.Errorf("Validate returned unexpected value for summary with invalid Timestamp")
	}
}
func TestSetClusterName(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{}
	summary.SetClusterName("test")
	if summary.ClusterName != "test" {
		t.Errorf("Expected ClusterName to be 'test', but got '%v'", summary.ClusterName)
	}
}

func TestSetNamespace(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{}
	summary.SetNamespace("test")
	if summary.Namespace != "test" {
		t.Errorf("Expected Namespace to be 'test', but got '%v'", summary.Namespace)
	}
}

func TestSetContainerName(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{}
	summary.SetContainerName("test")
	if summary.ContainerName != "test" {
		t.Errorf("Expected ContainerName to be 'test', but got '%v'", summary.ContainerName)
	}
}

func TestSetStatus(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{}
	summary.SetStatus("test")
	if summary.Status != "test" {
		t.Errorf("Expected Status to be 'test', but got '%v'", summary.Status)
	}
}

func TestSetRegistry(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{}
	summary.SetRegistry("test")
	if summary.Registry != "test" {
		t.Errorf("Expected Registry to be 'test', but got '%v'", summary.Registry)
	}
}

func TestSetImageTagSuffix(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{}
	summary.SetImageTageSuffix("test")
	if summary.ImageTagSuffix != "test" {
		t.Errorf("Expected ImageTagSuffix to be 'test', but got '%v'", summary.ImageTagSuffix)
	}
}

func TestSetVersion(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{}
	summary.SetVersion("test")
	if summary.Version != "test" {
		t.Errorf("Expected Version to be 'test', but got '%v'", summary.Version)
	}
}

func TestGetSeverityStats(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{
		SeverityStats: SeverityStats{
			RCEFixCount: 1,
			TotalCount:  15,
			Severity:    "critical",
		},
	}
	stats := summary.GetSeverityStats()
	assert.Equal(t, int64(1), stats.RCEFixCount)
	assert.Equal(t, int64(15), stats.TotalCount)
	assert.Equal(t, "critical", stats.Severity)
}

func Test_CommonContainerScanSummaryResult_GetRelevantLabel(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{}
	assert.Equal(t, RelevantLabelNotExists, summary.GetRelevantLabel())

	summary.RelevantLabel = RelevantLabelYes
	assert.Equal(t, RelevantLabelYes, summary.GetRelevantLabel())

}

func Test_CommonContainerScanSummaryResult_SetRelevantLabel(t *testing.T) {
	summary := &CommonContainerScanSummaryResult{}
	assert.Equal(t, RelevantLabelNotExists, summary.RelevantLabel)

	summary.SetRelevantLabel(RelevantLabelNo)
	assert.Equal(t, RelevantLabelNo, summary.RelevantLabel)
}

func TestGetHasRelevancyData(t *testing.T) {
	reportWithRelevancy := CommonContainerScanSummaryResult{HasRelevancyData: true}
	if !reportWithRelevancy.GetHasRelevancyData() {
		t.Error("Expected GetHasRelevancyData() to return true for a ScanResultReportV1 with relevancy data, but it returned false")
	}

	reportWithoutRelevancy := CommonContainerScanSummaryResult{HasRelevancyData: false}
	if reportWithoutRelevancy.GetHasRelevancyData() {
		t.Error("Expected GetHasRelevancyData() to return false for a ScanResultReportV1 without relevancy data, but it returned true")
	}
}

package containerscan

import (
	"strings"
	"testing"
)

func TestToFlatVulnerabilities_EmptyReport(t *testing.T) {
	report := &ScanResultReport{}
	vulns := report.ToFlatVulnerabilities()
	if len(vulns) != 0 {
		t.Errorf("Expected 0 vulnerabilities, got %d", len(vulns))
	}
}

func TestGetVulnLink_GitHub(t *testing.T) {
	link := GetVulnLink("GHSA-xxxx-xxxx-xxxx")
	if !strings.HasPrefix(link, "https://github.com/advisories/") {
		t.Errorf("Expected GitHub advisory link, got %s", link)
	}
}

func TestGetVulnLink_NVD(t *testing.T) {
	link := GetVulnLink("CVE-2023-0001")
	if !strings.HasPrefix(link, "https://nvd.nist.gov/vuln/detail/") {
		t.Errorf("Expected NVD link, got %s", link)
	}
}

package containerscan

import (
	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/armoapi-go/identifiers"
	cautils "github.com/armosec/utils-k8s-go/armometadata"
)

var SeverityStr2Score = map[string]int{
	"Unknown":    1,
	"Negligible": 100,
	"Low":        200,
	"Medium":     300,
	"High":       400,
	"Critical":   500,
}

func (longVul *Vulnerability) ToShortVulnerabilityResult() *ShortVulnerabilityResult {
	ret := &ShortVulnerabilityResult{
		Name: longVul.Name,
	}
	return ret
}

// ToFlatVulnerabilities - returnsgit p
func (scanresult *ScanResultReport) ToFlatVulnerabilities() []ContainerScanVulnerabilityResult {
	vuls := make([]CommonContainerVulnerabilityResult, 0)
	vul2indx := make(map[string]int)
	scanID := scanresult.AsFNVHash()
	designatorsObj, ctxList := scanresult.GetDesignatorsNContext()
	for _, layer := range scanresult.Layers {
		for _, vul := range layer.Vulnerabilities {
			esLayer := ESLayer{LayerHash: layer.LayerHash, ParentLayerHash: layer.ParentLayerHash}
			if indx, isOk := vul2indx[vul.Name]; isOk {
				vuls[indx].Layers = append(vuls[indx].Layers, esLayer)
				continue
			}
			result := CommonContainerVulnerabilityResult{WLID: scanresult.WLID,
				Timestamp:   scanresult.Timestamp,
				Designators: *designatorsObj,
				Context:     ctxList,
				IsLastScan:  1,
			}

			vul.SeverityScore = SeverityStr2Score[vul.Severity]
			result.Vulnerability = vul
			result.Layers = make([]ESLayer, 0)
			result.Layers = append(result.Layers, esLayer)
			result.ContainerScanID = scanID

			result.IsFixed = CalculateFixed(vul.Fixes)
			result.RelevantLinks = append(result.RelevantLinks, "https://nvd.nist.gov/vuln/detail/"+vul.Name)
			result.RelevantLinks = append(result.RelevantLinks, vul.Link)
			result.Vulnerability.SetLink("https://nvd.nist.gov/vuln/detail/" + vul.Name)
			result.GetVulnerability().SetCategories(VulnerabilityCategory{IsRCE: vul.IsRCE()})
			vuls = append(vuls, result)
			vul2indx[vul.Name] = len(vuls) - 1

		}
	}
	// find first introduced
	for i, v := range vuls {
		earlyLayer := ""
		for _, layer := range v.Layers {
			if layer.ParentLayerHash == earlyLayer {
				earlyLayer = layer.LayerHash
			}
		}
		vuls[i].IntroducedInLayer = earlyLayer

	}
	vulnsArr := make([]ContainerScanVulnerabilityResult, len(vuls))
	for i, v := range vuls {
		vulnsArr[i] = &v
	}
	return vulnsArr
}

func (scanresult *ScanResultReport) Summarize() *CommonContainerScanSummaryResult {
	designatorsObj, ctxList := scanresult.GetDesignatorsNContext()
	summary := &CommonContainerScanSummaryResult{
		Designators:     *designatorsObj,
		Context:         ctxList,
		CustomerGUID:    scanresult.CustomerGUID,
		ImageTag:        scanresult.ImgTag,
		ImageID:         scanresult.ImgHash,
		WLID:            scanresult.WLID,
		Timestamp:       scanresult.Timestamp,
		ContainerName:   scanresult.ContainerName,
		ContainerScanID: scanresult.AsFNVHash(),
		JobIDs:          scanresult.Session.JobIDs,

		ImageSignatureValid:           scanresult.ImageSignatureValid,
		ImageHasSignature:             scanresult.ImageHasSignature,
		ImageSignatureValidationError: scanresult.ImageSignatureValidationError,
	}

	summary.ClusterName = designatorsObj.Attributes[identifiers.AttributeCluster]
	summary.Namespace = designatorsObj.Attributes[identifiers.AttributeNamespace]

	imageInfo, e2 := cautils.ImageTagToImageInfo(scanresult.ImgTag)
	if e2 == nil {
		summary.Registry = imageInfo.Registry
		summary.ImageTagSuffix = imageInfo.VersionImage
	}

	summary.PackagesName = make([]string, 0)

	actualSeveritiesStats := map[string]SeverityStats{}
	exculdedSeveritiesStats := map[string]SeverityStats{}

	vulnsList := make([]ShortVulnerabilityResult, 0)
	uniqueVulsMap := make(map[string]bool)
	uniqueExceptionVulsMap := make(map[string]bool)
	for _, layer := range scanresult.Layers {
		summary.PackagesName = append(summary.PackagesName, (layer.GetPackagesNames())...)
		for _, vul := range layer.Vulnerabilities {
			if _, isOk := uniqueVulsMap[vul.Name]; isOk {
				continue
			}
			if _, isOk := uniqueExceptionVulsMap[vul.Name]; isOk {
				continue
			}
			isIgnored := (len(vul.ExceptionApplied) > 0 &&
				len(vul.ExceptionApplied[0].Actions) > 0 &&
				vul.ExceptionApplied[0].Actions[0] == armotypes.Ignore)

			severitiesStats := exculdedSeveritiesStats
			if !isIgnored {
				summary.TotalCount++
				uniqueVulsMap[vul.Name] = true
				vulnsList = append(vulnsList, *(vul.ToShortVulnerabilityResult()))
				severitiesStats = actualSeveritiesStats
			} else {
				uniqueExceptionVulsMap[vul.Name] = true
			}

			// TODO: maybe add all severities just to have a placeholders
			if !KnownSeverities[vul.Severity] {
				vul.Severity = UnknownSeverity
			}

			vulnSeverityStats, ok := severitiesStats[vul.Severity]
			if !ok {
				vulnSeverityStats = SeverityStats{Severity: vul.Severity}
			}

			vulnSeverityStats.TotalCount++
			isFixed := CalculateFixed(vul.Fixes) > 0
			if isFixed {
				vulnSeverityStats.FixAvailableOfTotalCount++
				incrementCounter(&summary.FixAvailableOfTotalCount, true, isIgnored)
			}
			isRCE := vul.IsRCE()
			if isRCE {
				vulnSeverityStats.RCECount++
				incrementCounter(&summary.RCECount, true, isIgnored)
				if isFixed {
					summary.RCEFixCount++
					vulnSeverityStats.RCEFixCount++
				}
			}
			severitiesStats[vul.Severity] = vulnSeverityStats
		}
	}
	summary.Status = "Success"
	summary.Vulnerabilities = vulnsList

	// if criticalStats, hasCritical := severitiesStats[CriticalSeverity]; hasCritical && criticalStats.TotalCount > 0 {
	// 	summary.Status = "Fail"
	// }
	// if highStats, hasHigh := severitiesStats[HighSeverity]; hasHigh && highStats.RelevantCount > 0 {
	// 	summary.Status = "Fail"
	// }

	for sever := range actualSeveritiesStats {
		summary.SeveritiesStats = append(summary.SeveritiesStats, actualSeveritiesStats[sever])
	}
	for sever := range exculdedSeveritiesStats {
		summary.ExcludedSeveritiesStats = append(summary.ExcludedSeveritiesStats, exculdedSeveritiesStats[sever])
	}
	return summary
}

func incrementCounter(counter *int64, isGlobal, isIgnored bool) {
	if isGlobal && isIgnored {
		return
	}
	(*counter)++
}

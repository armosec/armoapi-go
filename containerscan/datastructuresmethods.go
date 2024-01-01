package containerscan

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strings"

	"github.com/armosec/armoapi-go/identifiers"
	"github.com/google/uuid"
)

func (layer *ScanResultLayer) GetFilesByPackage(pkgname string) (files *PkgFiles) {
	for _, pkg := range layer.Packages {
		if pkg.PackageName == pkgname {
			return &pkg.Files
		}
	}

	return &PkgFiles{}
}

func (layer *ScanResultLayer) GetPackagesNames() []string {
	pkgsNames := []string{}
	for _, pkg := range layer.Packages {
		pkgsNames = append(pkgsNames, pkg.PackageName)
	}
	return pkgsNames
}

// GenerateBogusHash - generate the old (bogus) hash for the workload
func GenerateBogusHash(context map[string]string) string {
	context[identifiers.AttributeNamespace] = ""
	return generateWorkloadHash(context)
}

// GenerateWorkloadHash - generate a hash for the workload
func GenerateWorkloadHash(context map[string]string) string {
	return generateWorkloadHash(context)
}

func generateWorkloadHash(context map[string]string) string {
	strForHash := context[identifiers.AttributeCluster] + context[identifiers.AttributeNamespace] + context[identifiers.AttributeKind] + context[identifiers.AttributeName] + context[identifiers.AttributeContainerName]
	hasher := fnv.New64a()
	hasher.Write([]byte(strForHash))
	return fmt.Sprintf("%v", hasher.Sum64())
}

func (scanresult *ScanResultReport) GetDesignatorsNContext() (*identifiers.PortalDesignator, []identifiers.ArmoContext) {

	designatorsObj := identifiers.AttributesDesignatorsFromWLID(scanresult.WLID)
	designatorsObj.Attributes["containerName"] = scanresult.ContainerName
	designatorsObj.Attributes["workloadHash"] = generateWorkloadHash(designatorsObj.Attributes)
	designatorsObj.Attributes["customerGUID"] = scanresult.CustomerGUID

	//Copy all missing attributes
	for k := range scanresult.Designators.Attributes {
		if _, ok := designatorsObj.Attributes[k]; !ok {
			designatorsObj.Attributes[k] = scanresult.Designators.Attributes[k]
		}
	}

	contextObj := identifiers.DesignatorToArmoContext(designatorsObj, "designators")
	return designatorsObj, contextObj
}

func (scanresult *ScanResultReport) Validate() bool {
	if scanresult.CustomerGUID == "" || (scanresult.ImgHash == "" && scanresult.ImgTag == "") || scanresult.Timestamp <= 0 {
		return false
	}

	if _, err := uuid.Parse(scanresult.CustomerGUID); err != nil {
		return false
	}

	//TODO validate layers & vuls

	return true
}

func (v *Vulnerability) IsRCE() bool {
	return IsRceFromDescription(v.Description)
}

func IsRceFromDescription(description string) bool {
	desc := strings.ToLower(description)

	isRCE, _ := regexp.MatchString(`[^A-Za-z]rce[^A-Za-z]`, description)

	return isRCE || strings.Contains(desc, "remote code execution") || strings.Contains(desc, "remote command execution") || strings.Contains(desc, "arbitrary code") || strings.Contains(desc, "code execution") || strings.Contains(desc, "code injection") || strings.Contains(desc, "command injection") || strings.Contains(desc, "inject arbitrary commands")
}

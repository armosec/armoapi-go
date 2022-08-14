package armotypes

const (
	// In-cluster namespaces
	ArmoSystemNamespace = "armo-system" // deprecated, kept for backward compatibility
	KubescapeNamespace  = "kubescape"

	ArmoKollectorContainerName = "armo-collector"

	// registry scan
	LowestHelmVersionSupportedRegistryScan = "v1.7.14"
	RegistryInfoArgKey                     = "registryInfo-v1"
	RegistryScanSecretName                 = "kubescape-registry-scan"

	// vulnerability scan
	LowestHelmVersionSupportedVulnerabilityScan = "v1.7.17"

	// cronjob template annotation and labels
	CronJobTemplateAnnotationJobIDKey        = "armo.jobid"
	CronJobTemplateAnnotationNamespaceKey    = "armo.namespace"
	CronJobTemplateAnnotationRegistryNameKey = "armo.cloud/registryname"
	CronJobTemplateAnnotationHostScannerKey  = "armo.host-scanner"
	CronJobTemplateAnnotationFrameworkKey    = "armo.framework"

	CronJobTemplateLabelKey               = "armo.tier"
	CronJobTemplateLabelValueKubescape    = "kubescape-scan"
	CronJobTemplateLabelValueVulnScan     = "vuln-scan"
	CronJobTemplateLabelValueRegistryScan = "registry-scan"
)

func GetInClusterSupportedNamespaces() []string {
	return []string{KubescapeNamespace, ArmoSystemNamespace}
}

package armotypes

const (
	K8sKindCluster   = "Cluster"
	K8sKindNamespace = "Namespace"

	K8sApiVersionV1      = "v1"
	K8sApiVersionRBAC    = "rbac.authorization.k8s.io"
	K8sApiVersionRBACV1  = K8sApiVersionRBAC + "/" + K8sApiVersionV1
	K8SApiVersionAppsV1  = "apps/v1"
	K8SApiVersionBatchV1 = "batch/v1"
)

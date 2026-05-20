package armotypes

// ClusterSBOMScanMessage is the Pulsar payload published by the cluster
// vulnerability scan dispatcher (one per ContainerProfile per tick) and
// consumed by the backend Vulnerability Scanner. The scanner reads the
// SBOM blob at SBOMObjectRef from S3, fetches the ContainerProfile row
// by ContainerProfileName, computes the filtered SBOM, runs Grype, and
// emits a vuln-scan summary tagged kind=cluster.
//
// All identifying fields needed for filtered-SBOM computation and for
// per-workload attribution downstream travel on the message — the
// scanner does not need to call back into postgres for them.
//
// Field stability: any field added later MUST be `omitempty` so older
// publishers continue to produce valid payloads (consumer treats absent
// fields as zero values).
type ClusterSBOMScanMessage struct {
	CustomerGUID         string `json:"customerGUID"`
	Cluster              string `json:"cluster"`
	WorkloadKind         string `json:"workloadKind"`
	WorkloadName         string `json:"workloadName"`
	WorkloadNamespace    string `json:"workloadNamespace"`
	WorkloadResourceHash string `json:"workloadResourceHash"`
	ContainerProfileName string `json:"containerProfileName"`
	ImageDigest          string `json:"imageDigest"`
	SyftVersion          string `json:"syftVersion"`
	SBOMObjectRef        string `json:"sbomObjectRef"`
}

package armotypes

// ClusterSBOMScanMessage is the Pulsar payload published by the cluster
// vulnerability scan dispatcher (one per ContainerProfile per tick) and
// consumed by the backend Vulnerability Scanner.
//
// The scanner reads the SBOM blob from S3 at SBOMObjectRef and fetches
// the ContainerProfile row from postgres (by ContainerProfileName) so
// it can compute the filtered SBOM before running Grype. It then emits
// a vuln-scan summary tagged kind=cluster.
//
// The message carries every identity field needed for downstream
// per-workload attribution (cluster, workload kind/name/namespace,
// resource hash, image digest, syft version) — that path is
// postgres-free. Only the filtered-SBOM computation step needs the CP
// row from the DB.
//
// SBOMObjectRef is a JSON-encoded `{"bucket": "...", "key": "..."}`
// string — kept as a string to pass through unchanged from the
// `sbom_metadata.resource_object_ref` column rather than serialise +
// deserialise on the dispatcher's hot path.
//
// Field stability: any field added later MUST be `omitempty` so older
// publishers continue to produce valid payloads (consumer treats
// absent fields as zero values).
type ClusterSBOMScanMessage struct {
	PortalBase           `json:",inline" bson:",inline"`
	CustomerGUID         string `json:"customerGUID" bson:"customerGUID"`
	Cluster              string `json:"cluster" bson:"cluster"`
	WorkloadKind         string `json:"workloadKind" bson:"workloadKind"`
	WorkloadName         string `json:"workloadName" bson:"workloadName"`
	WorkloadNamespace    string `json:"workloadNamespace" bson:"workloadNamespace"`
	WorkloadResourceHash string `json:"workloadResourceHash" bson:"workloadResourceHash"`
	ContainerProfileName string `json:"containerProfileName" bson:"containerProfileName"`
	ImageDigest          string `json:"imageDigest" bson:"imageDigest"`
	SyftVersion          string `json:"syftVersion" bson:"syftVersion"`
	SBOMObjectRef        string `json:"sbomObjectRef" bson:"sbomObjectRef"`
}

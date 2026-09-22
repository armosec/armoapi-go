# Kubernetes host machine-generation identity

`CloudMetadata.KubernetesHostIdentity` is optional nested JSON/BSON metadata under
`kubernetes_host_identity`. Legacy reports omit it. Existing cloud instance,
machine, provider, account and host-type fields retain their meanings.

The v1 envelope contains, in canonical JSON order: `version` (1), `cluster_uid`,
`cluster_name`, `node_uid`, `node_name`, `machine_fingerprint`, `key`, and optional
full `provider_id`. Required fields are not omitted. `CanonicalJSON` validates
and encodes this struct without whitespace using Go's standard JSON escaping.
This value metadata is not a Portal resource and has no PortalBase fields.

`NormalizeKubernetesHostMachineID` trims ASCII whitespace, removes UUID hyphens,
lowercases, and requires 32 hexadecimal digits other than the all-zero ID.
`KubernetesHostMachineFingerprint` computes `mf1-` plus lowercase, unpadded base32
of SHA-256 over the framed strings `kubernetes-host-machine-v1` and the normalized
machine-id. The raw machine-id is not stored in the envelope.

`KubernetesHostKey` computes `kh1-` plus the same digest encoding over framed
`kubernetes-host-key-v1`, cluster UID, Node UID, and machine fingerprint. Each
string is UTF-8 preceded by its unsigned 32-bit big-endian byte length. UIDs are
exact, nonempty API values; names and cloud metadata do not affect the digest.
Both identifiers are 56 characters.

The key survives ordinary agent restart and host reboot. A changed machine-id,
Node UID, or cluster UID creates a new key; recreating a Node on the same hardware
also starts a new generation. Cloned machine-ids on distinct Node objects remain
separate. A cloned machine-id retaining the same Node object is indistinguishable;
this scheme is not hardware attestation.

Synthetic host profiles use label `kubescape.io/k8s-host-key` and annotation
`kubescape.io/k8s-host-identity` containing the canonical envelope. String-valued
backend metadata uses `kubernetes_host_identity` for the same canonical JSON.
Ordinary pod profiles must not carry these markers.

`Validate` checks version, required fields, UTF-8, canonical fingerprint syntax,
and the derived key. It does not authenticate a customer or cluster or verify a
Node. Consumers must independently validate authenticated tenant routing,
registered cluster UID, and synced Node UID, name, and machine fingerprint.
Machine-id source selection and host-file consistency checks belong to producers,
not this pure shared contract.

Golden-vector, restart/replacement, malformed-input, and JSON/BSON compatibility
coverage lives in `armotypes/kuberneteshostidentity_test.go`.

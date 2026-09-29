package aisandbox

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// identity_key — the canonical, stable join key that an activity row's
// identity_ref points at. Format is a discriminated union: <tag>:<canonical-body>.
//
//   - The tag (segment before the FIRST ':') names the identity KIND, from the
//     fixed vocabulary below. It prevents cross-kind collisions and lets
//     consumers route. A tag is [a-z0-9-] and never contains ':'.
//   - The body is the natural globally-(or per-customer-)unique identifier for
//     that kind, in its canonical form. It is OPAQUE for joins: it may itself
//     contain ':' and '/' (AWS ARNs do) — equality is whole-string, so no
//     downstream parsing beyond the tag split is needed.
//
// Every producer (K8s/RBAC derivation, cloud IAM enrichment, and later the agent
// per-call signal) MUST build keys through the constructors here, so the format
// lives in exactly one audited place and cannot drift. Adding a new kind = a new
// tag constant + a new constructor; the format itself never changes.
const (
	AiSandboxIdentityTagK8sServiceAccount = "k8s-sa"
	AiSandboxIdentityTagAWSRole           = "aws-role"
	AiSandboxIdentityTagAWSUser           = "aws-user"
	AiSandboxIdentityTagGCPServiceAccount = "gcp-sa"
	AiSandboxIdentityTagAzureIdentity     = "azure-id"
	// AiSandboxIdentityTagAWSSession is the canonical kind for a workload's AWS
	// TEMPORARY session credentials (SUB-7956): STS tokens rotate, so keying
	// them by credential fingerprint minted a fresh "identity" every rotation.
	// A TTL credential is not inventory — the workload's session-ness is. The
	// body is therefore the workload's resource_hash, giving ONE stable
	// identity per workload that the read side joins back to the workload's
	// DERIVED cloud identity (the IRSA role from the SA annotation, or the ECS
	// task role) — the attribution point where the observed session and the
	// real IAM identity meet.
	AiSandboxIdentityTagAWSSession = "aws-session"
	// AiSandboxIdentityTagBearerSession is the canonical kind for a workload's
	// ROTATING bearer credentials against one destination (Azure AD / OAuth access
	// tokens, SA JWTs): the token rotates by design, so a per-fingerprint key mints
	// a fresh identity every rotation (measured on dev: 36 rows for one workload's
	// hourly Azure AD tokens, all to management.azure.com; ~5k unknown:bearer rows
	// customer-wide). Keyed per workload × destination host — the same
	// rotation-collapse rationale as aws-session, but per destination because two
	// different bearer credentials against two services must never merge.
	AiSandboxIdentityTagBearerSession = "bearer-session"
	// AiSandboxIdentityTagUnknown is the TOTAL fallback: a principal we observed but can
	// not classify into a known kind. Captured opaquely (never dropped, never
	// guessed) so the security signal — "some identity we don't model made this
	// call" — survives. See OpaqueIdentityKey.
	AiSandboxIdentityTagUnknown = "unknown"
)

// identityKeyMaxBodyLen bounds the opaque body length; a pathologically long or
// adversarial principal is replaced by its sha256 so identity_key can never grow
// unbounded. The human-readable value is carried separately on the row's Name —
// the key only needs to be stable and bounded, not readable.
const identityKeyMaxBodyLen = 512

// K8sServiceAccountKey → k8s-sa:<cluster>/<namespace>/<serviceAccount>.
// All three coordinates are required (a ServiceAccount is only unique within a
// namespace within a cluster). ok=false if any is empty, so the caller leaves
// the row unattributed rather than forming a colliding partial key. Case is
// preserved (K8s names are DNS-lowercase already; the cluster label is kept
// verbatim so it round-trips whatever the platform stores).
func K8sServiceAccountKey(cluster, namespace, serviceAccount string) (string, bool) {
	cluster = strings.TrimSpace(cluster)
	namespace = strings.TrimSpace(namespace)
	serviceAccount = strings.TrimSpace(serviceAccount)
	if cluster == "" || namespace == "" || serviceAccount == "" {
		return "", false
	}
	return joinIdentityKey(AiSandboxIdentityTagK8sServiceAccount, cluster+"/"+namespace+"/"+serviceAccount), true
}

// AWSRoleKey → aws-role:<roleARN>. The ARN is globally unique and self-contained
// (account id + role path). Case is PRESERVED — IAM resource names are
// case-sensitive, so lower-casing would corrupt the identity.
func AWSRoleKey(roleARN string) (string, bool) {
	return awsARNKey(AiSandboxIdentityTagAWSRole, roleARN)
}

// AWSUserKey → aws-user:<userARN>.
func AWSUserKey(userARN string) (string, bool) {
	return awsARNKey(AiSandboxIdentityTagAWSUser, userARN)
}

func awsARNKey(tag, arn string) (string, bool) {
	arn = strings.TrimSpace(arn)
	if arn == "" {
		return "", false
	}
	return joinIdentityKey(tag, arn), true
}

// GCPServiceAccountKey → gcp-sa:<email>. The email embeds the project and is
// globally unique; lower-cased (email identifiers are case-insensitive).
func GCPServiceAccountKey(email string) (string, bool) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return "", false
	}
	return joinIdentityKey(AiSandboxIdentityTagGCPServiceAccount, email), true
}

// AWSSessionKey → aws-session:<resourceHash>. ok=false on an empty hash so the
// caller leaves the row unattributed rather than forming one shared
// "aws-session:" bucket across workloads.
func AWSSessionKey(resourceHash string) (string, bool) {
	resourceHash = strings.TrimSpace(resourceHash)
	if resourceHash == "" {
		return "", false
	}
	return joinIdentityKey(AiSandboxIdentityTagAWSSession, resourceHash), true
}

// BearerSessionKey → bearer-session:<resourceHash>/<host>. The canonical key
// for a workload's rotating bearer credentials against one destination host
// (see AiSandboxIdentityTagBearerSession). Both coordinates are required —
// ok=false on either missing, so the caller falls back to the per-fingerprint
// opaque key rather than forming a shared bucket. The host is lower-cased
// (hostnames are case-insensitive); the resourceHash is kept verbatim. Like
// aws-session this deliberately loses the per-fingerprint identity for TTL
// credentials — a rotating token is not inventory; "this workload holds a
// bearer credential for this service" is.
func BearerSessionKey(resourceHash, host string) (string, bool) {
	resourceHash = strings.TrimSpace(resourceHash)
	host = strings.ToLower(strings.TrimSpace(host))
	if resourceHash == "" || host == "" {
		return "", false
	}
	return joinIdentityKey(AiSandboxIdentityTagBearerSession, resourceHash+"/"+host), true
}

// AzureIdentityKey → azure-id:<tenant>/<clientID>, or azure-id:<clientID> when
// the tenant is unknown. GUIDs are case-insensitive → lower-cased. The clientID
// is required; the tenant qualifies it when available.
func AzureIdentityKey(tenantID, clientID string) (string, bool) {
	tenantID = strings.ToLower(strings.TrimSpace(tenantID))
	clientID = strings.ToLower(strings.TrimSpace(clientID))
	if clientID == "" {
		return "", false
	}
	body := clientID
	if tenantID != "" {
		body = tenantID + "/" + clientID
	}
	return joinIdentityKey(AiSandboxIdentityTagAzureIdentity, body), true
}

// OpaqueIdentityKey is the TOTAL fallback for a principal we observed but can not
// classify into a known kind. Any non-empty raw principal yields a stable
// unknown:<body> key (the body is sha256-hashed when it exceeds
// identityKeyMaxBodyLen) — so an unrecognized identity is captured, deduped and
// joinable, never dropped or guessed. ok=false ONLY for genuinely-empty input,
// which leaves the row unattributed (UI shows '—'), distinct from an opaque
// identity we DID observe.
//
// Cardinality of the unknown bucket is bounded at the WRITER (top-N per workload
// + an 'unclassified' overflow bucket), not here — this constructor is stateless.
func OpaqueIdentityKey(rawPrincipal string) (string, bool) {
	raw := strings.TrimSpace(rawPrincipal)
	if raw == "" {
		return "", false
	}
	body := raw
	if len(body) > identityKeyMaxBodyLen {
		sum := sha256.Sum256([]byte(raw))
		body = "sha256:" + hex.EncodeToString(sum[:])
	}
	return joinIdentityKey(AiSandboxIdentityTagUnknown, body), true
}

// IdentityKeyTag returns the kind tag of an identity_key (the segment before the
// FIRST ':'), for routing/classification by consumers. Returns "" for a
// malformed key (no ':' or a leading ':'). Safe on bodies that themselves
// contain ':' (e.g. aws-role:arn:aws:iam::… → "aws-role").
func IdentityKeyTag(identityKey string) string {
	if i := strings.IndexByte(identityKey, ':'); i > 0 {
		return identityKey[:i]
	}
	return ""
}

func joinIdentityKey(tag, body string) string { return tag + ":" + body }

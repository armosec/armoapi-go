package armotypes

// CIProvider names the CI system that a CIIdentity describes.
type CIProvider string

const (
	CIProviderGitHub CIProvider = "github"
)

// CIIdentity is the CI identity that a sensor reports about itself: which CI
// system runs the workload, and the keys of the run in that system. It is
// optional enrichment, next to the Kubernetes and ECS details of an alert, and
// on an entry of a monitored units message.
//
// A sensor sets it only when it knows its own CI identity (for example a sensor
// on a VM runner). A sensor that reports Kubernetes pods does not set it: the
// consumer derives the identity from the pod and workload labels.
//
// The values are keys for a lookup in the CI system, not facts. The reporting
// workload may run untrusted code, so a consumer must confirm them with the CI
// system before it shows them as verified.
//
// Each provider has its own block, because the keys of a run differ between CI
// systems. Exactly one block matches Provider. A consumer that does not know a
// provider still gets Provider, and ignores the block.
//
// Wire contract: docs/features/monitored-units.md
type CIIdentity struct {
	Provider CIProvider        `json:"provider" bson:"provider"`
	GitHub   *GitHubCIIdentity `json:"github,omitempty" bson:"github,omitempty"`
}

// GitHubCIIdentity holds the keys of a GitHub Actions job run. Every field is
// optional: a sensor sends what it knows.
type GitHubCIIdentity struct {
	// RunnerName is the name of the runner as registered in GitHub.
	RunnerName string `json:"runnerName,omitempty" bson:"runnerName,omitempty"`
	// RunID is the ID of the workflow run.
	RunID int64 `json:"runId,omitempty" bson:"runId,omitempty"`
	// Attempt is the run attempt, starting at 1.
	Attempt int `json:"attempt,omitempty" bson:"attempt,omitempty"`
	// Job is the job name.
	Job string `json:"job,omitempty" bson:"job,omitempty"`
}

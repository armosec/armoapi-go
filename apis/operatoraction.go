package apis

import "encoding/json"

// OperatorActionType enumerates the concrete cluster operations that a
// TypeOperatorAction command can request. It is carried in
// OperatorActionArgs.Action and dispatched by the operator to the matching
// Remediator implementation.
type OperatorActionType string

const (
	// OperatorActionAnnotate adds a label/annotation to a workload (the
	// lowest-blast-radius action; shipped first to prove the pipeline).
	OperatorActionAnnotate OperatorActionType = "annotate"
	// OperatorActionQuarantine isolates a workload (deny-all NetworkPolicy +
	// quarantine label).
	OperatorActionQuarantine OperatorActionType = "quarantine"
	// OperatorActionCordon marks a node unschedulable.
	OperatorActionCordon OperatorActionType = "cordon"
	// OperatorActionRevert undoes a previously applied action on a target.
	OperatorActionRevert OperatorActionType = "revert"
)

// OperatorActionTarget identifies a single concrete object an action operates
// on.
type OperatorActionTarget struct {
	Kind      string `json:"kind,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

// OperatorActionSelector resolves a findings-driven target set from the
// scan-result CRDs already stored in the cluster (e.g. "workloads failing
// control C-0016 with severity >= High"). No re-scan is performed.
type OperatorActionSelector struct {
	// Control is the Kubescape control ID to select failing workloads for
	// (e.g. "C-0016").
	Control string `json:"control,omitempty"`
	// MinSeverity is the minimum finding severity to act on (e.g. "High").
	MinSeverity string `json:"minSeverity,omitempty"`
	// Namespace optionally limits the selector to a single namespace.
	Namespace string `json:"namespace,omitempty"`
}

// OperatorActionArgs is the typed schema for Command.Args when
// Command.CommandName == TypeOperatorAction.
//
// Either a fully-specified Target or a findings-driven Selector is provided;
// when both are present the Selector resolves an additional target set.
//
// Safe-by-default: DryRun is a *bool so its zero value (nil / absent on the
// wire) means dry-run. The operator therefore treats any request that does not
// carry an explicit DryRun=false as a plan-only dry-run. Only an explicit
// DryRun=false (the CLI's --confirm) performs a real cluster write. Use
// IsDryRun to interpret the field rather than dereferencing it directly.
type OperatorActionArgs struct {
	// Action is the operation to perform.
	Action OperatorActionType `json:"action"`
	// Target is an explicit single object to act on (optional if Selector is set).
	Target *OperatorActionTarget `json:"target,omitempty"`
	// Selector resolves targets from stored findings (optional if Target is set).
	Selector *OperatorActionSelector `json:"selector,omitempty"`
	// FindingRef references the scan-result CRD that justifies the action, for
	// audit (e.g. "workloadconfigurationscansummaries/payments/api").
	FindingRef string `json:"findingRef,omitempty"`
	// DryRun gates real cluster writes. nil (absent) or true means plan-only;
	// only an explicit false applies the change. See IsDryRun.
	DryRun *bool `json:"dryRun,omitempty"`
	// TTL optionally schedules an automatic revert after the given duration
	// (e.g. "24h").
	TTL string `json:"ttl,omitempty"`
	// Reason is a human-readable justification recorded in the audit trail.
	Reason string `json:"reason,omitempty"`
}

// IsDryRun reports whether the action should be treated as a plan-only dry-run.
// It returns true unless DryRun was explicitly set to false, so a request that
// omits the field can never silently perform a real cluster write.
func (a OperatorActionArgs) IsDryRun() bool {
	return a.DryRun == nil || *a.DryRun
}

// ToArgs serializes the typed action args into the generic Command.Args map.
func (a OperatorActionArgs) ToArgs() (map[string]interface{}, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// OperatorActionArgsFromMap parses a generic Command.Args map back into the
// typed action args.
func OperatorActionArgsFromMap(m map[string]interface{}) (OperatorActionArgs, error) {
	var a OperatorActionArgs
	b, err := json.Marshal(m)
	if err != nil {
		return a, err
	}
	err = json.Unmarshal(b, &a)
	return a, err
}

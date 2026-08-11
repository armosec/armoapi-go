package cdr

import "time"

// Provider-agnostic accessors over EventData's per-provider event variants
// (AWSCloudTrail / AzureActivityLog / GcpAuditLog). Consumers (e.g. the CDR
// ingester) use these instead of reaching into a specific provider's event, so
// their logic stays cloud-neutral and never nil-dereferences the wrong variant.
// Each returns a zero value when no event — or no relevant field — is present.

// EventTime returns the timestamp of whichever provider event is set.
func (e EventData) EventTime() time.Time {
	switch {
	case e.AWSCloudTrail != nil:
		return e.AWSCloudTrail.EventTime
	case e.AzureActivityLog != nil:
		return e.AzureActivityLog.Time
	case e.GcpAuditLog != nil:
		return e.GcpAuditLog.Timestamp
	}
	return time.Time{}
}

// SourceIP returns the caller/source IP of whichever provider event is set, or
// "" when absent.
func (e EventData) SourceIP() string {
	switch {
	case e.AWSCloudTrail != nil:
		return e.AWSCloudTrail.SourceIPAddress
	case e.AzureActivityLog != nil:
		return e.AzureActivityLog.CallerIPAddress
	case e.GcpAuditLog != nil:
		if e.GcpAuditLog.ProtoPayload != nil && e.GcpAuditLog.ProtoPayload.RequestMetadata != nil {
			return e.GcpAuditLog.ProtoPayload.RequestMetadata.CallerIP
		}
	}
	return ""
}

// Region returns the cloud region of whichever provider event is set, or "".
// Azure uses the Activity Log's `location` (often "global" for control-plane
// operations); GCP uses the first of the resource's current locations.
func (e EventData) Region() string {
	switch {
	case e.AWSCloudTrail != nil:
		return e.AWSCloudTrail.AWSRegion
	case e.AzureActivityLog != nil:
		return e.AzureActivityLog.Location
	case e.GcpAuditLog != nil:
		if e.GcpAuditLog.ProtoPayload != nil && e.GcpAuditLog.ProtoPayload.ResourceLocation != nil {
			if locs := e.GcpAuditLog.ProtoPayload.ResourceLocation.CurrentLocations; len(locs) > 0 {
				return locs[0]
			}
		}
	}
	return ""
}

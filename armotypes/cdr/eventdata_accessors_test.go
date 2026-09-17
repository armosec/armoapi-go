package cdr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventDataEventTime(t *testing.T) {
	ts := time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		ed   EventData
		want time.Time
	}{
		{"aws", EventData{AWSCloudTrail: &CloudTrailEvent{EventTime: ts}}, ts},
		{"azure", EventData{AzureActivityLog: &AzureActivityLogEvent{Time: ts}}, ts},
		{"gcp", EventData{GcpAuditLog: &GcpAuditLogEvent{Timestamp: ts}}, ts},
		{"none", EventData{}, time.Time{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ed.EventTime())
		})
	}
}

func TestEventDataSourceIP(t *testing.T) {
	tests := []struct {
		name string
		ed   EventData
		want string
	}{
		{"aws", EventData{AWSCloudTrail: &CloudTrailEvent{SourceIPAddress: "1.1.1.1"}}, "1.1.1.1"},
		{"azure", EventData{AzureActivityLog: &AzureActivityLogEvent{CallerIPAddress: "2.2.2.2"}}, "2.2.2.2"},
		{"gcp", EventData{GcpAuditLog: &GcpAuditLogEvent{ProtoPayload: &GcpAuditLogPayload{RequestMetadata: &GcpRequestMetadata{CallerIP: "3.3.3.3"}}}}, "3.3.3.3"},
		{"gcp nil payload", EventData{GcpAuditLog: &GcpAuditLogEvent{}}, ""},
		{"none", EventData{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ed.SourceIP())
		})
	}
}

func TestEventDataRegion(t *testing.T) {
	tests := []struct {
		name string
		ed   EventData
		want string
	}{
		{"aws", EventData{AWSCloudTrail: &CloudTrailEvent{AWSRegion: "us-east-1"}}, "us-east-1"},
		{"azure location", EventData{AzureActivityLog: &AzureActivityLogEvent{Location: "global"}}, "global"},
		{"azure empty location", EventData{AzureActivityLog: &AzureActivityLogEvent{}}, ""},
		{"gcp first current location", EventData{GcpAuditLog: &GcpAuditLogEvent{ProtoPayload: &GcpAuditLogPayload{ResourceLocation: &GcpResourceLocation{CurrentLocations: []string{"us-central1", "us-east1"}}}}}, "us-central1"},
		{"gcp nil payload", EventData{GcpAuditLog: &GcpAuditLogEvent{}}, ""},
		{"none", EventData{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ed.Region())
		})
	}
}

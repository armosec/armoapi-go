package armotypes

import (
	"strings"
	"time"
)

// Getters & Setter used by derived types for interfaces implementation
func (p *PortalBase) GetGUID() string {
	return p.GUID
}
func (p *PortalBase) SetGUID(guid string) {
	p.GUID = guid
}
func (p *PortalBase) GetName() string {
	return p.Name
}
func (p *PortalBase) SetName(name string) {
	p.Name = name
}
func (p *PortalBase) GetAttributes() map[string]interface{} {
	return p.Attributes
}
func (p *PortalBase) SetAttributes(attributes map[string]interface{}) {
	p.Attributes = attributes
}

func (p *PortalBase) SetUpdatedTime(updatedTime *time.Time) {
	if updatedTime == nil {
		p.UpdatedTime = time.Now().UTC().Format(time.RFC3339)
		return
	}
	p.UpdatedTime = updatedTime.UTC().Format(time.RFC3339)
}

func (p *PortalBase) GetUpdatedTime() *time.Time {
	if p.UpdatedTime == "" {
		return nil
	}
	updatedTime, err := time.Parse(time.RFC3339, p.UpdatedTime)
	if err != nil {
		return nil
	}
	return &updatedTime
}

func ValidateContainerScanID(containerScanID string) bool {
	return !strings.Contains(containerScanID, "/")
}

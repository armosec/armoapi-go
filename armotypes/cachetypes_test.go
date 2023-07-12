package armotypes

import (
	"testing"
)

func TestGetCacheData(t *testing.T) {
	cache := &PortalCache{
		Data: "test",
	}

	data, err := GetCacheData[string](cache)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if data != "test" {
		t.Errorf("expected data to be 'test', got '%v'", data)
	}

	cache.Data = int(123)
	num, err := GetCacheData[int](cache)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if num != 123 {
		t.Errorf("expected data to be 'test', got '%v'", data)
	}

	type myStruct struct {
		Name string
	}
	cache.Data = myStruct{
		Name: "test",
	}
	s, err := GetCacheData[myStruct](cache)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if s.Name != "test" {
		t.Errorf("expected s.Name to be 'test', got '%v'", s.Name)
	}
}

func TestDataTypeConversion(t *testing.T) {
	dataType := MakeCacheDataTypeV1("service", "domain", "purpose", "version")
	service, domain, purpose, version, err := ParseCacheDataTypeV1(dataType)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if service != "service" {
		t.Errorf("expected service to be 'service', got '%s'", service)
	}
	

	if domain != "domain" {
		t.Errorf("expected domain to be 'domain', got '%s'", domain)
	}

	if purpose != "purpose" {
		t.Errorf("expected purpose to be 'purpose', got '%s'", purpose)
	}

	if version != "version" {
		t.Errorf("expected version to be 'version', got '%s'", version)
	}

	// Test bad data type format
	_, _, _,_,  err = ParseCacheDataTypeV1("datatypeV1:invalid")
	if err == nil {
		t.Error("expected error, got nil")
	}

	// Test incomplete data type
	_, _, _,_,  err = ParseCacheDataTypeV1("datatypeV1:service-customer-domain")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

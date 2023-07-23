package armotypes

import (
	"testing"
	"time"
)

func TestSetExpiryTime(t *testing.T) {
	cache := PortalCache[int]{}
	expiryTime := time.Now().UTC().Add(time.Hour * 1)

	cache.SetExpiryTime(expiryTime)

	if !cache.ExpiryTime.Equal(expiryTime) {
		t.Errorf("Expected expiry time to be %v, but got %v", expiryTime, cache.ExpiryTime)
	}
}

func TestSetTTL(t *testing.T) {
	cache := PortalCache[int]{}
	ttl := time.Hour * 1

	cache.SetTTL(ttl)

	expectedExpiryTime := time.Now().UTC().Add(ttl)
	if cache.ExpiryTime.Before(expectedExpiryTime.Add(-time.Minute)) || cache.ExpiryTime.After(expectedExpiryTime.Add(time.Minute)) {
		t.Errorf("Expected expiry time to be approximately %v, but got %v", expectedExpiryTime, cache.ExpiryTime)
	}
}

package notifications

import (
	"encoding/json"
	"reflect"
	"testing"
)

func headersPtr(m map[string]string) *map[string]string {
	return &m
}

func TestSetWebhookConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *WebhookConfig
		wantHeaders interface{}
		wantPresent bool
	}{
		{
			name: "persists headers",
			config: &WebhookConfig{
				WebhookURL: "https://ingest.example.com/http-endpoint/v1/abc",
				Headers: headersPtr(map[string]string{
					"Authorization": "Bearer token123",
					"Content-Type":  "application/json",
				}),
			},
			wantHeaders: map[string]string{
				"Authorization": "Bearer token123",
				"Content-Type":  "application/json",
			},
			wantPresent: true,
		},
		{
			name: "nil headers omits the key",
			config: &WebhookConfig{
				WebhookURL: "https://ingest.example.com/http-endpoint/v1/abc",
			},
			wantPresent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SIEMIntegration{}
			s.SetWebhookConfig(tt.config)

			if s.Provider != SIEMProviderWebhook {
				t.Errorf("Provider = %q, want %q", s.Provider, SIEMProviderWebhook)
			}
			if got := s.Configuration["webhookURL"]; got != tt.config.WebhookURL {
				t.Errorf("webhookURL = %v, want %v", got, tt.config.WebhookURL)
			}

			got, present := s.Configuration["headers"]
			if present != tt.wantPresent {
				t.Fatalf("headers present = %v, want %v", present, tt.wantPresent)
			}
			if tt.wantPresent && !reflect.DeepEqual(got, tt.wantHeaders) {
				t.Errorf("headers = %#v, want %#v", got, tt.wantHeaders)
			}
		})
	}
}

// TestSetWebhookConfigPersistsThroughJSON reproduces the customer bug: headers
// entered in the console must survive the JSON marshal the persistence layer
// performs (see repositories/siem_setters.go). Before the fix SetWebhookConfig
// dropped them and they never reached storage.
func TestSetWebhookConfigPersistsThroughJSON(t *testing.T) {
	s := &SIEMIntegration{}
	s.SetWebhookConfig(&WebhookConfig{
		WebhookURL: "https://ingest.echo.taegis.secureworks.com/http-endpoint/v1/3f97419f",
		Headers: headersPtr(map[string]string{
			"Authorization": "Bearer TDGbyZ0zJXKq4gcJ2hy6BvY1cLVF4gGT0ntEbUKj4bmoQ",
			"Content-Type":  "application/json",
		}),
	})

	payload, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var restored SIEMIntegration
	if err := json.Unmarshal(payload, &restored); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	headers, ok := restored.Configuration["headers"].(map[string]interface{})
	if !ok {
		t.Fatalf("headers missing after round-trip, got %#v", restored.Configuration["headers"])
	}
	if headers["Authorization"] != "Bearer TDGbyZ0zJXKq4gcJ2hy6BvY1cLVF4gGT0ntEbUKj4bmoQ" {
		t.Errorf("Authorization header = %v, want the Bearer token", headers["Authorization"])
	}
	if headers["Content-Type"] != "application/json" {
		t.Errorf("Content-Type header = %v, want application/json", headers["Content-Type"])
	}
}

package notifications

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestNewRuntimeIncidentActorJSON(t *testing.T) {
	// payload shape produced by kubescape/messaging RuntimeIncidentIngesterOnFinishedMessage
	payload := []byte(`{"incidentGUID":"i1","userName":"alice@example.com","userGroups":["system:authenticated","devs"]}`)

	var incident NewRuntimeIncident
	if err := json.Unmarshal(payload, &incident); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if incident.UserName != "alice@example.com" {
		t.Errorf("UserName = %q, want %q", incident.UserName, "alice@example.com")
	}
	if want := []string{"system:authenticated", "devs"}; !reflect.DeepEqual(incident.UserGroups, want) {
		t.Errorf("UserGroups = %v, want %v", incident.UserGroups, want)
	}
}

func TestNewRuntimeIncidentActorOmittedWhenEmpty(t *testing.T) {
	data, err := json.Marshal(NewRuntimeIncident{IncidentGUID: "i1"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, key := range []string{`"userName"`, `"userGroups"`} {
		if strings.Contains(string(data), key) {
			t.Errorf("expected %s to be omitted, got %s", key, data)
		}
	}
}

func TestNewRuntimeIncidentActorBSONRoundTrip(t *testing.T) {
	in := NewRuntimeIncident{UserName: "alice@example.com", UserGroups: []string{"devs"}}
	data, err := bson.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var raw bson.M
	if err := bson.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal raw: %v", err)
	}
	if raw["userName"] != "alice@example.com" {
		t.Errorf(`bson key "userName" = %v, want %q`, raw["userName"], "alice@example.com")
	}
	if _, ok := raw["userGroups"]; !ok {
		t.Errorf(`bson key "userGroups" missing in %v`, raw)
	}
	var out NewRuntimeIncident
	if err := bson.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.UserName != in.UserName || !reflect.DeepEqual(out.UserGroups, in.UserGroups) {
		t.Errorf("round trip = %+v, want %+v", out, in)
	}
}

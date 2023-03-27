package apis

import (
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestCronJobParams(t *testing.T) {
	c := Command{}
	cjp := CronJobParams{
		CronTabSchedule: "* * * * *",
	}
	c.SetCronJobParams(cjp)
	cjpr := c.GetCronJobParams()
	assert.Equal(t, "* * * * *", cjpr.CronTabSchedule)
}

func TestLabels(t *testing.T) {
	c := Command{}
	cjp := map[string]string{
		"app": "game",
	}
	c.SetLabels(cjp)
	cjpr := c.GetLabels()
	assert.Equal(t, "game", cjpr["app"])
}

func TestFieldSelector(t *testing.T) {
	c := Command{}
	cjp := map[string]string{
		"app": "game",
	}
	c.SetFieldSelector(cjp)
	cjpr := c.GetFieldSelector()
	assert.Equal(t, "game", cjpr["app"])
}

func TestRegistryScanCommandGetWlid(t *testing.T) {
	r := &RegistryScanCommand{}
	if r.GetWlid() != "" {
		t.Errorf("Expected empty string, but got %s", r.GetWlid())
	}
}

func TestRegistryScanCommandGetCredentialsList(t *testing.T) {
	r := &RegistryScanCommand{ImageScanParams: ImageScanParams{Credentialslist: []types.AuthConfig{{
		Username: "username",
		Password: "password",
	}}}}
	expected := []types.AuthConfig{{
		Username: "username",
		Password: "password",
	}}
	if !reflect.DeepEqual(r.GetCredentialsList(), expected) {
		t.Errorf("Expected %v, but got %v", expected, r.GetCredentialsList())
	}
}

func TestRegistryScanCommandSetCredentialsList(t *testing.T) {
	r := &RegistryScanCommand{}
	expected := []types.AuthConfig{{
		Username: "username",
		Password: "password",
	}}
	r.SetCredentialsList(expected)
	if !reflect.DeepEqual(r.GetCredentialsList(), expected) {
		t.Errorf("Expected %v, but got %v", expected, r.GetCredentialsList())
	}
}

func TestRegistryScanCommandGetArgs(t *testing.T) {
	r := &RegistryScanCommand{ImageScanParams: ImageScanParams{Args: map[string]interface{}{"skipTLSVerify": true}}}
	expected := map[string]interface{}{"skipTLSVerify": true}
	if !reflect.DeepEqual(r.GetArgs(), expected) {
		t.Errorf("Expected %v, but got %v", expected, r.GetArgs())
	}
}

func TestRegistryScanCommandSetArgs(t *testing.T) {
	r := &RegistryScanCommand{}
	expected := map[string]interface{}{"skipTLSVerify": true}
	r.SetArgs(expected)
	if !reflect.DeepEqual(r.GetArgs(), expected) {
		t.Errorf("Expected %v, but got %v", expected, r.GetArgs())
	}
}

func TestRegistryScanCommandGetSession(t *testing.T) {
	r := &RegistryScanCommand{ImageScanParams: ImageScanParams{Session: SessionChain{
		JobIDs: []string{"1234"},
	}}}
	expected := SessionChain{
		JobIDs: []string{"1234"},
	}
	if !reflect.DeepEqual(r.GetSession(), expected) {
		t.Errorf("Expected %v, but got %v", expected, r.GetSession())
	}
}

func TestRegistryScanCommandSetSession(t *testing.T) {
	r := &RegistryScanCommand{}
	expected := SessionChain{
		JobIDs: []string{"1234"},
	}
	r.SetSession(expected)
	if !reflect.DeepEqual(r.GetSession(), expected) {
		t.Errorf("Expected %v, but got %v", expected, r.GetSession())
	}
}

func TestRegistryScanCommandGetImageTag(t *testing.T) {
	r := &RegistryScanCommand{ImageScanParams: ImageScanParams{ImageTag: "nginx:latest"}}
	if r.GetImageTag() != "nginx:latest" {
		t.Errorf("Expected nginx:latest, but got %s", r.GetImageTag())
	}
}

func TestRegistryScanCommandSetImageTag(t *testing.T) {
	r := &RegistryScanCommand{}
	expected := "nginx:latest"
	r.SetImageTag(expected)
	if r.GetImageTag() != expected {
		t.Errorf("Expected %s, but got %s", expected, r.GetImageTag())
	}
}

func TestRegistryScanCommandGetJobID(t *testing.T) {
	r := &RegistryScanCommand{ImageScanParams: ImageScanParams{JobID: "1234"}}
	if r.GetJobID() != "1234" {
		t.Errorf("Expected 1234, but got %s", r.GetJobID())
	}
}

func TestRegistryScanCommandSetJobID(t *testing.T) {
	registry := &RegistryScanCommand{}
	registry.SetJobID("job1")
	assert.Equal(t, "job1", registry.GetJobID())
}

func TestRegistryScanCommandGetParentJobID(t *testing.T) {
	registry := &RegistryScanCommand{}
	registry.SetParentJobID("parentJob1")
	assert.Equal(t, "parentJob1", registry.GetParentJobID())
}

func TestRegistryScanCommandSetParentJobID(t *testing.T) {
	registry := &RegistryScanCommand{}
	registry.SetParentJobID("parentJob2")
	assert.Equal(t, "parentJob2", registry.GetParentJobID())
}

func TestRegistryScanCommandGetCreds(t *testing.T) {
	registry := &RegistryScanCommand{}
	creds := registry.GetCreds()
	assert.Nil(t, creds)
}

func TestRegistryScanCommandGetImageHash(t *testing.T) {
	registry := &RegistryScanCommand{}
	imageHash := registry.GetImageHash()
	assert.Equal(t, "", imageHash)
}

func TestWebsocketScanCommandGetWlid(t *testing.T) {
	w := WebsocketScanCommand{Wlid: "wlid://cluster-marina/namespace-default/deployment-nginx"}
	if w.GetWlid() != "wlid://cluster-marina/namespace-default/deployment-nginx" {
		t.Errorf("GetWlid() returned unexpected value, expected: %s, got: %s", "wlid://cluster-marina/namespace-default/deployment-nginx", w.GetWlid())
	}
}

func TestWebsocketScanCommandGetImageTag(t *testing.T) {
	w := WebsocketScanCommand{
		ImageScanParams: ImageScanParams{
			ImageTag: "nginx:latest",
		},
	}
	if w.GetImageTag() != "nginx:latest" {
		t.Errorf("GetImageTag() returned unexpected value, expected: %s, got: %s", "nginx:latest", w.GetImageTag())
	}
}

func TestWebsocketScanCommandSetImageTag(t *testing.T) {
	w := WebsocketScanCommand{
		ImageScanParams: ImageScanParams{
			ImageTag: "nginx:latest",
		},
	}
	w.SetImageTag("nginx:1.21.0")
	if w.ImageTag != "nginx:1.21.0" {
		t.Errorf("SetImageTag() did not set the ImageTag value correctly, expected: %s, got: %s", "nginx:1.21.0", w.ImageTag)
	}
}

func TestWebsocketScanCommandGetJobID(t *testing.T) {
	w := WebsocketScanCommand{ImageScanParams: ImageScanParams{JobID: "12345"}}
	if w.GetJobID() != "12345" {
		t.Errorf("GetJobID() returned unexpected value, expected: %s, got: %s", "12345", w.GetJobID())
	}
}

func TestWebsocketScanCommand_SetJobID(t *testing.T) {
	w := WebsocketScanCommand{ImageScanParams: ImageScanParams{Session: SessionChain{
		JobIDs: []string{"1234"},
	}}}
	w.SetJobID("54321")
	if w.JobID != "54321" {
		t.Errorf("SetJobID() did not set the JobID value correctly, expected: %s, got: %s", "54321", w.JobID)
	}
}

func TestWebsocketScanCommandGetParentJobID(t *testing.T) {
	w := WebsocketScanCommand{ImageScanParams: ImageScanParams{ParentJobID: "12345"}}
	if w.GetParentJobID() != "12345" {
		t.Errorf("GetParentJobID() returned unexpected value, expected: %s, got: %s", "12345", w.GetParentJobID())
	}
}

func TestWebsocketScanCommandSetParentJobID(t *testing.T) {
	w := WebsocketScanCommand{ImageScanParams: ImageScanParams{ParentJobID: "12345"}}
	w.SetParentJobID("54321")
	if w.ParentJobID != "54321" {
		t.Errorf("SetParentJobID() did not set the ParentJobID value correctly, expected: %s, got: %s", "54321", w.ParentJobID)
	}
}

func TestWebsocketScanCommandGetCreds(t *testing.T) {
	w := WebsocketScanCommand{}
	creds := w.GetCreds()
	assert.Nil(t, creds)

	w.Credentials = &types.AuthConfig{
		Username: "user",
	}
	assert.Equal(t, "user", w.GetCreds().Username)
}

func TestWebsocketScanCommandGetImageHash(t *testing.T) {
	w := WebsocketScanCommand{}
	w.ImageHash = "imageHash"
	assert.Equal(t, "imageHash", w.GetImageHash())
}

func TestWebsocketScanCommandGetCredentialsList(t *testing.T) {
	authConfig := types.AuthConfig{
		Username:      "user1",
		Password:      "password1",
		ServerAddress: "https://registry1.com",
	}
	expected := []types.AuthConfig{authConfig}

	cmd := &WebsocketScanCommand{
		ImageScanParams: ImageScanParams{
			Credentialslist: expected,
		},
	}

	got := cmd.GetCredentialsList()

	if len(got) != len(expected) {
		t.Fatalf("expected %d credentials, got %d", len(expected), len(got))
	}

	if got[0] != expected[0] {
		t.Fatalf("expected %+v, got %+v", expected[0], got[0])
	}
}

func TestWebsocketScanCommandSetCredentialsList(t *testing.T) {
	authConfig := types.AuthConfig{
		Username:      "user1",
		Password:      "password1",
		ServerAddress: "https://registry1.com",
	}
	expected := []types.AuthConfig{authConfig}

	cmd := &WebsocketScanCommand{}

	cmd.SetCredentialsList(expected)

	if len(cmd.Credentialslist) != len(expected) {
		t.Fatalf("expected %d credentials, got %d", len(expected), len(cmd.Credentialslist))
	}

	if cmd.Credentialslist[0] != expected[0] {
		t.Fatalf("expected %+v, got %+v", expected[0], cmd.Credentialslist[0])
	}
}

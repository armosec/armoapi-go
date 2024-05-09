package armotypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericCRDWithNodeProfile(t *testing.T) {
	// Setting up a sample NodeProfile
	nodeProfile := NodeProfile{
		CustomerGUID:            "1234-5678-9101",
		Cluster:                 "ClusterA",
		Name:                    "Node1",
		PodStatuses:             []PodStatus{{Name: "Pod1", HasFinalApplicationProfile: true}},
		NodeAgentRunning:        true,
		RuntimeDetectionEnabled: true,
	}

	// Creating a GenericCRD with NodeProfile as Spec
	crd := GenericCRD[NodeProfile]{
		Kind:       "NodeAgentProfile",
		ApiVersion: "kubescape.io/v1",
		Spec:       nodeProfile,
	}

	// Marshalling to JSON
	jsonData, err := json.Marshal(crd)
	assert.Nil(t, err, "Marshalling should not produce an error")
	assert.NotNil(t, jsonData, "Marshalled JSON should not be nil")

	// Unmarshalling back to struct to check integrity
	var decodedCRD GenericCRD[NodeProfile]
	err = json.Unmarshal(jsonData, &decodedCRD)
	assert.Nil(t, err, "Unmarshalling should not produce an error")
	assert.Equal(t, crd, decodedCRD, "Original and decoded CRD should be equal")
}

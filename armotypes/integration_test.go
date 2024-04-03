package armotypes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntity_Validate(t *testing.T) {
	// Base entity setup for valid cases
	validBaseEntity := EntityIdentifiers{
		Cluster:          "testCluster",
		Namespace:        "testNamespace",
		Name:             "testName",
		Kind:             "testKind",
		ResourceHash:     "testResourceHash",
		ResourceID:       "testResourceID",
		CVEName:          "testCVEName",
		CVEID:            "testCVEID",
		Severity:         "testSeverity",
		SeverityScore:    10,
		Component:        "testComponent",
		ComponentVersion: "testComponentVersion",
		ImageReposiotry:  "testImageRepository",
		LayerHash:        "testLayerHash",
		ControlID:        "testControlID",
		BaseScore:        5.0,
	}

	// Test cases
	tests := []struct {
		name    string
		entity  EntityIdentifiers
		wantErr bool
	}{
		{
			name: "Valid posture resource entity",
			entity: EntityIdentifiers{
				Type:         EntityTypePostureResource,
				Cluster:      validBaseEntity.Cluster,
				Namespace:    validBaseEntity.Namespace,
				Name:         validBaseEntity.Name,
				Kind:         validBaseEntity.Kind,
				ResourceHash: validBaseEntity.ResourceHash,
				ResourceID:   validBaseEntity.ResourceID,
			},
			wantErr: false,
		},
		{
			name: "Invalid posture resource entity",
			entity: EntityIdentifiers{
				Type: EntityTypePostureResource,
			},
			wantErr: true,
		},
		{
			name: "Valid repository resource entity",
			entity: EntityIdentifiers{
				Type:       EntityTypeRepositoryResource,
				RepoHash:   "testRepoHash",
				Namespace:  validBaseEntity.Namespace,
				Name:       validBaseEntity.Name,
				Kind:       validBaseEntity.Kind,
				ResourceID: validBaseEntity.ResourceID,
			},
			wantErr: false,
		},
		{
			name: "Invalid repository resource entity",
			entity: EntityIdentifiers{
				Type: EntityTypeRepositoryResource,
			},
			wantErr: true,
		},
		{
			name: "Valid container scan workload entity",
			entity: EntityIdentifiers{
				Type:         EntityTypeContainerScanWorkload,
				Cluster:      validBaseEntity.Cluster,
				Namespace:    validBaseEntity.Namespace,
				Name:         validBaseEntity.Name,
				Kind:         validBaseEntity.Kind,
				ResourceHash: validBaseEntity.ResourceHash,
			},
			wantErr: false,
		},
		{
			name: "Invalid container scan workload entity",
			entity: EntityIdentifiers{
				Type: EntityTypeContainerScanWorkload,
			},
			wantErr: true,
		},
		{
			name: "Valid image entity",
			entity: EntityIdentifiers{
				Type:            EntityTypeImage,
				ImageReposiotry: validBaseEntity.ImageReposiotry,
			},
			wantErr: false,
		},
		{
			name: "Invalid image entity",
			entity: EntityIdentifiers{
				Type: EntityTypeImage,
			},
			wantErr: true,
		},
		{
			name: "Valid image layer entity",
			entity: EntityIdentifiers{
				Type:            EntityTypeImageLayer,
				ImageReposiotry: validBaseEntity.ImageReposiotry,
				LayerHash:       validBaseEntity.LayerHash,
			},
			wantErr: false,
		},
		{
			name: "Invalid image layer entity",
			entity: EntityIdentifiers{
				Type: EntityTypeImageLayer,
			},
			wantErr: true,
		},
		{
			name: "Valid vulnerability entity",
			entity: EntityIdentifiers{
				Type:             EntityTypeVulanrability,
				CVEName:          validBaseEntity.CVEName,
				CVEID:            validBaseEntity.CVEID,
				Severity:         validBaseEntity.Severity,
				SeverityScore:    validBaseEntity.SeverityScore,
				Component:        validBaseEntity.Component,
				ComponentVersion: validBaseEntity.ComponentVersion,
			},
			wantErr: false,
		},
		{
			name: "Invalid vulnerability entity",
			entity: EntityIdentifiers{
				Type: EntityTypeVulanrability,
			},
			wantErr: true,
		},
		{
			name: "Valid control entity",
			entity: EntityIdentifiers{
				Type:      EntityTypeControl,
				ControlID: validBaseEntity.ControlID,
				Severity:  validBaseEntity.Severity,
				BaseScore: validBaseEntity.BaseScore,
			},
			wantErr: false,
		},
		{
			name: "Invalid control entity",
			entity: EntityIdentifiers{
				Type: EntityTypeControl,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entity.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

package armotypes

import (
	"testing"
)

var securityIssueControl = SecurityIssueControl{
	SecurityIssue: SecurityIssue{
		Cluster:          "cluster1",
		ClusterShortName: "clusterShortName1",
	},
}

var securityIssueAttackPath = SecurityIssueAttackPath{
	SecurityIssue: SecurityIssue{
		Cluster:          "cluster1",
		ClusterShortName: "clusterShortName1",
	},
}

func checkISecurityIssueControlGetClusterName(test ISecurityIssue) string {
	return test.GetClusterName()

}

func TestSecurityIssueControl(t *testing.T) {
	// Test for the SecurityIssueControl struct
	if securityIssueControl.Cluster != "cluster1" {
		t.Errorf("Expected %s, got %s", "cluster1", securityIssueControl.Cluster)
	}
	if securityIssueControl.ClusterShortName != "clusterShortName1" {
		t.Errorf("Expected %s, got %s", "clusterShortName1", securityIssueControl.ClusterShortName)
	}

	securityIssueControl.SetClusterName("cluster2")

	if securityIssueControl.GetClusterName() != "cluster2" {
		t.Errorf("Expected %s, got %s", "cluster2", securityIssueControl.Cluster)
	}

	securityIssueControl.SetShortClusterName("clusterShortName2")

	if securityIssueControl.GetShortClusterName() != "clusterShortName2" {
		t.Errorf("Expected %s, got %s", "clusterShortName2", securityIssueControl.ClusterShortName)
	}

	// Test for the ISecurityIssue interface
	if checkISecurityIssueControlGetClusterName(&securityIssueControl) != "cluster2" {
		t.Errorf("Expected %s, got %s", "cluster2", securityIssueControl.Cluster)
	}

	// Test for the SecurityIssueAttackPath struct
	if securityIssueAttackPath.Cluster != "cluster1" {
		t.Errorf("Expected %s, got %s", "cluster1", securityIssueAttackPath.Cluster)
	}

	if securityIssueAttackPath.ClusterShortName != "clusterShortName1" {
		t.Errorf("Expected %s, got %s", "clusterShortName1", securityIssueAttackPath.ClusterShortName)
	}

	securityIssueAttackPath.SetClusterName("cluster2")

	if securityIssueAttackPath.GetClusterName() != "cluster2" {
		t.Errorf("Expected %s, got %s", "cluster2", securityIssueAttackPath.Cluster)
	}

	securityIssueAttackPath.SetShortClusterName("clusterShortName2")

	if securityIssueAttackPath.GetShortClusterName() != "clusterShortName2" {
		t.Errorf("Expected %s, got %s", "clusterShortName2", securityIssueAttackPath.ClusterShortName)
	}

	// Test for the ISecurityIssue interface
	if checkISecurityIssueControlGetClusterName(&securityIssueAttackPath) != "cluster2" {
		t.Errorf("Expected %s, got %s", "cluster2", securityIssueAttackPath.Cluster)
	}

}

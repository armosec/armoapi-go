package armotypes

// GetControlIDsFromAllNodes is a recursive func that returns a list of controlIDs from all nodes in the attack chain
func (attackChainNode *AttackChainNode) GetControlIDsFromAllNodes(controlIDs []string) []string {
	controlIDs = append(controlIDs, attackChainNode.ControlIDs...)
	for i := range attackChainNode.NextNodes {
		controlIDs = attackChainNode.NextNodes[i].GetControlIDsFromAllNodes(controlIDs)
	}
	return controlIDs
}

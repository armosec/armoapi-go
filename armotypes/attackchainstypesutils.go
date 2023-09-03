package armotypes

// getControlIDsFromAllNodes is a recursive func that returns a list of controlIDs from all nodes in the attack chain
func (attackChainNode *AttackChainNode) getControlIDsFromAllNodes(controlIDs []string) []string {
	controlIDs = append(controlIDs, attackChainNode.ControlIDs...)
	for i := range attackChainNode.NextNodes {
		controlIDs = attackChainNode.NextNodes[i].getControlIDsFromAllNodes(controlIDs)
	}
	return controlIDs
}

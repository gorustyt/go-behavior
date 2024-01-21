package decorators

import "github.com/gorustyt/go-behavior/core"

type ForceFailureNode struct {
	*core.DecoratorNode
}

func NewForceFailureNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &ForceFailureNode{DecoratorNode: core.NewDecoratorNode(name, cfg)}
	n.SetRegistrationID("ForceFailure")
	return n
}

func (n *ForceFailureNode) Tick() core.NodeStatus {
	n.SetStatus(core.NodeStatus_RUNNING)
	child_status := n.Child().ExecuteTick()
	if core.IsStatusCompleted(child_status) {
		n.ResetChild()
		return core.NodeStatus_FAILURE
	}
	// RUNNING or skipping
	return child_status
}

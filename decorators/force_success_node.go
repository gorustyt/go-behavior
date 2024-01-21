package decorators

import "github.com/gorustyt/go-behavior/core"

type ForceSuccessNode struct {
	*core.DecoratorNode
}

func NewForceSuccessNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &ForceSuccessNode{DecoratorNode: core.NewDecoratorNode(name, cfg)}
	n.SetRegistrationID("ForceSuccess")
	return n
}

func (n *ForceSuccessNode) Tick() core.NodeStatus {
	n.SetStatus(core.NodeStatus_RUNNING)
	child_status := n.Child().ExecuteTick()
	if core.IsStatusCompleted(child_status) {
		n.ResetChild()
		return core.NodeStatus_SUCCESS
	}
	// RUNNING or skipping
	return child_status
}

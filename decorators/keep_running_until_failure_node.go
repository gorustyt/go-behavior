package decorators

import "github.com/gorustyt/go-behavior/core"

type KeepRunningUntilFailureNode struct {
	*core.DecoratorNode
}

func NewKeepRunningUntilFailureNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &KeepRunningUntilFailureNode{DecoratorNode: core.NewDecoratorNode(name, cfg)}
	n.SetRegistrationID("KeepRunningUntilFailure")
	return n
}
func (n *KeepRunningUntilFailureNode) Tick() core.NodeStatus {
	n.SetStatus(core.NodeStatus_RUNNING)
	childState := n.Child().ExecuteTick()
	switch childState {
	case core.NodeStatus_FAILURE:
		n.ResetChild()
		return core.NodeStatus_FAILURE
	case core.NodeStatus_SUCCESS:
		n.ResetChild()
		return core.NodeStatus_RUNNING
	case core.NodeStatus_RUNNING:
		return core.NodeStatus_RUNNING

	default:
		panic("invalid status")
	}
	return n.Status()
}

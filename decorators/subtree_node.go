package decorators

import "github.com/gorustyt/go-behavior/core"

func init() {
	core.SetPorts(&SubTreeNode{}, core.InputPortWithDefaultValue("_autoremap", false, "If true, all the ports with the same name will be remapped"))
}

type SubTreeNode struct {
	*core.DecoratorNode
	subtreeId string
}

func NewSubTreeNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &SubTreeNode{DecoratorNode: core.NewDecoratorNode(name, cfg)}
	n.SetRegistrationID("SubTree")
	return n
}

func (n *SubTreeNode) SetSubtreeID(ID string) {
	n.subtreeId = ID
}

func (n *SubTreeNode) Tick() core.NodeStatus {
	prevStatus := n.Status()
	if prevStatus == core.NodeStatus_IDLE {
		n.SetStatus(core.NodeStatus_RUNNING)
	}
	childStatus := n.Child().ExecuteTick()
	if core.IsStatusCompleted(childStatus) {
		n.ResetChild()
	}
	return childStatus
}

package decorators

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

type InverterNode struct {
	*core.DecoratorNode
}

func NewInverterNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &InverterNode{DecoratorNode: core.NewDecoratorNode(name, cfg)}
	n.SetRegistrationID("Inverter")
	return n
}

func (n *InverterNode) Tick() core.NodeStatus {
	n.SetStatus(core.NodeStatus_RUNNING)
	childStatus := n.Child().ExecuteTick()
	switch childStatus {
	case core.NodeStatus_SUCCESS:
		n.ResetChild()
		return core.NodeStatus_FAILURE
	case core.NodeStatus_FAILURE:
		n.ResetChild()
		return core.NodeStatus_SUCCESS
	case core.NodeStatus_RUNNING:
	case core.NodeStatus_SKIPPED:
		return childStatus
	case core.NodeStatus_IDLE:
		panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))

	}
	return n.Status()
}

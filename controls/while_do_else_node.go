package controls

import "github.com/gorustyt/go-behavior/core"

type WhileDoElseNode struct {
	*core.ControlNode
}

func NewWhileDoElseNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &WhileDoElseNode{
		ControlNode: core.NewControlNode(name, cfg),
	}
	n.SetRegistrationID("WhileDoElse")
	return n
}

func (n *WhileDoElseNode) Tick() core.NodeStatus {
	childrenCount := len(n.Children)

	if childrenCount != 2 && childrenCount != 3 {
		panic("WhileDoElseNode must have either 2 or 3 children")
	}

	n.SetStatus(core.NodeStatus_RUNNING)

	conditionStatus := n.Children[0].ExecuteTick()

	if conditionStatus == core.NodeStatus_RUNNING {
		return conditionStatus
	}

	status := core.NodeStatus_IDLE

	if conditionStatus == core.NodeStatus_SUCCESS {
		if childrenCount == 3 {
			n.HaltChild(2)
		}
		status = n.Children[1].ExecuteTick()
	} else if conditionStatus == core.NodeStatus_FAILURE {
		if childrenCount == 3 {
			n.HaltChild(1)
			status = n.Children[2].ExecuteTick()
		} else if childrenCount == 2 {
			status = core.NodeStatus_FAILURE
		}
	}

	if status == core.NodeStatus_RUNNING {
		return core.NodeStatus_RUNNING
	} else {
		n.ResetChildren()
		return status
	}

}
func (n *WhileDoElseNode) Halt() {
	n.ControlNode.Halt()
}

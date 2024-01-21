package controls

import "github.com/gorustyt/go-behavior/core"

type IfThenElseNode struct {
	*core.ControlNode
	index int
}

func NewIfThenElseNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &IfThenElseNode{ControlNode: core.NewControlNode(name, cfg)}
	n.SetRegistrationID("IfThenElse")
	return n
}

func (n *IfThenElseNode) halt() {
	n.index = 0
	n.ControlNode.Halt()
}

func (n *IfThenElseNode) Tick() core.NodeStatus {
	childrenCount := len(n.Children)

	if childrenCount != 2 && childrenCount != 3 {
		panic("IfThenElseNode must have either 2 or 3 children")
	}

	n.SetStatus(core.NodeStatus_RUNNING)

	if n.index == 0 {
		conditionStatus := n.Children[0].ExecuteTick()

		if conditionStatus == core.NodeStatus_RUNNING {
			return conditionStatus
		} else if conditionStatus == core.NodeStatus_SUCCESS {
			n.index = 1
		} else if conditionStatus == core.NodeStatus_FAILURE {
			if childrenCount == 3 {
				n.index = 2
			} else {
				return conditionStatus
			}
		}
	}
	// not an else
	if n.index > 0 {
		status := n.Children[n.index].ExecuteTick()
		if status == core.NodeStatus_RUNNING {
			return core.NodeStatus_RUNNING
		} else {
			n.ResetChildren()
			n.index = 0
			return status
		}
	}

	panic("Something unexpected happened in IfThenElseNode")
}

package controls

import "github.com/gorustyt/go-behavior/core"

type IfThenElseNode struct {
	*core.ControlNode
	index int
}

func NewIfThenElseNode() *IfThenElseNode {
	return &IfThenElseNode{}
}

func (n *IfThenElseNode) halt() {
	n.index = 0
	n.ControlNode.Halt()
}

func (n *IfThenElseNode) Tick() core.NodeStatus {
	children_count := len(n.Children)

	if children_count != 2 && children_count != 3 {
		panic("IfThenElseNode must have either 2 or 3 children")
	}

	n.SetStatus(core.NodeStatus_RUNNING)

	if n.index == 0 {
		condition_status := n.Children[0].ExecuteTick()

		if condition_status == core.NodeStatus_RUNNING {
			return condition_status
		} else if condition_status == core.NodeStatus_SUCCESS {
			n.index = 1
		} else if condition_status == core.NodeStatus_FAILURE {
			if children_count == 3 {
				n.index = 2
			} else {
				return condition_status
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

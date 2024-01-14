package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

// 按顺序执行孩子结点直到其中一个孩子结点返回失败状态或所有孩子结点返回成功状态。
type SequenceNode struct {
	*core.ControlNode
	allSkipped bool
	asynch     bool
	index      int
}

func NewSequenceNode() *SequenceNode {
	return &SequenceNode{
		allSkipped: true,
	}
}

func (n *SequenceNode) Halt() {
	n.index = 0
	n.ControlNode.Halt()
}

func (n *SequenceNode) Tick() core.NodeStatus {
	children_count := len(n.Children)

	if n.Status() == core.NodeStatus_IDLE {
		n.allSkipped = true
	}

	n.SetStatus(core.NodeStatus_RUNNING)

	for n.index < children_count {
		current_child_node := n.Children[n.index]

		prev_status := current_child_node.Status()
		child_status := current_child_node.ExecuteTick()

		// switch to RUNNING state as soon as you find an active child
		n.allSkipped = (child_status == core.NodeStatus_SKIPPED)

		switch child_status {
		case core.NodeStatus_RUNNING:
			{
				return core.NodeStatus_RUNNING
			}
		case core.NodeStatus_FAILURE:
			{
				// Reset on failure
				n.ResetChildren()
				n.index = 0
				return child_status
			}
		case core.NodeStatus_SUCCESS:
			{
				n.index++
				// Return the execution flow if the child is async,
				// to make this interruptable.
				if n.asynch && requiresWakeUp() &&
					prev_status == core.NodeStatus_IDLE &&
					n.index < children_count {
					emitWakeUpSignal()
					return core.NodeStatus_RUNNING
				}
			}
			break

		case core.NodeStatus_SKIPPED:
			{
				// It was requested to skip this node
				n.index++
			}
			break

		case core.NodeStatus_IDLE:
			{
				panic(fmt.Sprintf("[%v]: A children should not return IDL", n.Name()))
			}
		} // end switch
	} // end while loop

	// The entire while loop completed. This means that all the children returned SUCCESS.
	if n.index == children_count {
		n.ResetChildren()
		n.index = 0
	}
	// Skip if ALL the nodes have been skipped
	if n.allSkipped {
		return core.NodeStatus_SKIPPED
	}
	return core.NodeStatus_SUCCESS

}

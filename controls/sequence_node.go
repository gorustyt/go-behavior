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

func NewSequenceNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	makeAynch := false
	if len(args) > 0 {
		makeAynch = args[0].(bool)
	}
	n := &SequenceNode{
		ControlNode: core.NewControlNode(name, cfg),
		allSkipped:  true,
		asynch:      makeAynch,
	}
	if n.asynch {
		n.SetRegistrationID("AsyncSequence")
	} else {
		n.SetRegistrationID("Sequence")
	}

	return n
}

func (n *SequenceNode) Halt() {
	n.index = 0
	n.ControlNode.Halt()
}

func (n *SequenceNode) Tick() core.NodeStatus {
	childrenCount := len(n.Children)

	if n.Status() == core.NodeStatus_IDLE {
		n.allSkipped = true
	}

	n.SetStatus(core.NodeStatus_RUNNING)

	for n.index < childrenCount {
		currentChildNode := n.Children[n.index]

		prevStatus := currentChildNode.Status()
		childStatus := currentChildNode.ExecuteTick()

		// switch to RUNNING state as soon as you find an active child
		n.allSkipped = childStatus == core.NodeStatus_SKIPPED

		switch childStatus {
		case core.NodeStatus_RUNNING:
			{
				return core.NodeStatus_RUNNING
			}
		case core.NodeStatus_FAILURE:
			// Reset on failure
			n.ResetChildren()
			n.index = 0
			return childStatus
		case core.NodeStatus_SUCCESS:
			n.index++
			// Return the execution flow if the child is async,
			// to make this interruptable.
			if n.asynch && n.RequiresWakeUp() &&
				prevStatus == core.NodeStatus_IDLE &&
				n.index < childrenCount {
				n.EmitWakeUpSignal()
				return core.NodeStatus_RUNNING
			}
		case core.NodeStatus_SKIPPED:

			// It was requested to skip this node
			n.index++
		case core.NodeStatus_IDLE:
			panic(fmt.Sprintf("[%v]: A children should not return IDL", n.Name()))
		} // end switch
	} // end while loop

	// The entire while loop completed. This means that all the children returned SUCCESS.
	if n.index == childrenCount {
		n.ResetChildren()
		n.index = 0
	}
	// Skip if ALL the nodes have been skipped
	if n.allSkipped {
		return core.NodeStatus_SKIPPED
	}
	return core.NodeStatus_SUCCESS

}

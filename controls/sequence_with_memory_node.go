package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

type SequenceWithMemory struct {
	*core.ControlNode
	currentChildIdx int
	allSkipped      bool
}

func NewSequenceWithMemory(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &SequenceWithMemory{
		allSkipped:  true,
		ControlNode: core.NewControlNode(name, cfg),
	}
	n.SetRegistrationID("SequenceWithMemory")
	return n
}
func (n *SequenceWithMemory) Halt() {
	n.ControlNode.Halt()
}
func (n *SequenceWithMemory) Tick() core.NodeStatus {
	childrenCount := len(n.Children)

	if n.Status() == core.NodeStatus_IDLE {
		n.allSkipped = true
	}
	n.SetStatus(core.NodeStatus_RUNNING)

	for n.currentChildIdx < childrenCount {
		currentChildNode := n.Children[n.currentChildIdx]

		prevStatus := currentChildNode.Status()
		childStatus := currentChildNode.ExecuteTick()

		// switch to RUNNING state as soon as you find an active child
		n.allSkipped = childStatus == core.NodeStatus_SKIPPED

		switch childStatus {
		case core.NodeStatus_RUNNING:
			return childStatus
		case core.NodeStatus_FAILURE:
			// DO NOT reset current_child_idx_ on failure
			for i := n.currentChildIdx; i < len(n.Children); i++ {
				n.HaltChild(i)
			}
			return childStatus
		case core.NodeStatus_SUCCESS:
			n.currentChildIdx++
			// Return the execution flow if the child is async,
			// to make this interruptable.
			if n.RequiresWakeUp() && prevStatus == core.NodeStatus_IDLE &&
				n.currentChildIdx < childrenCount {
				n.EmitWakeUpSignal()
				return core.NodeStatus_RUNNING
			}
		case core.NodeStatus_SKIPPED:
			// It was requested to skip this node
			n.currentChildIdx++
		case core.NodeStatus_IDLE:
			panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))
		} // end switch
	} // end while loop

	// The entire while loop completed. This means that all the children returned SUCCESS.
	if n.currentChildIdx == childrenCount {
		n.ResetChildren()
		n.currentChildIdx = 0
	}
	if n.allSkipped {
		return core.NodeStatus_SKIPPED
	}
	// Skip if ALL the nodes have been skipped
	return core.NodeStatus_SUCCESS
}

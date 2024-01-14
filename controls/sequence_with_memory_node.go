package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

type SequenceWithMemory struct {
	*core.ControlNode
	current_child_idx int
	all_skipped_      bool
}

func NewSequenceWithMemory() *SequenceWithMemory {
	return &SequenceWithMemory{
		all_skipped_: true,
	}
}
func (n *SequenceWithMemory) Halt() {
	n.ControlNode.Halt()
}
func (n *SequenceWithMemory) Tick() core.NodeStatus {
	children_count := len(n.Children)

	if n.Status() == core.NodeStatus_IDLE {
		n.all_skipped_ = true
	}
	n.SetStatus(core.NodeStatus_RUNNING)

	for n.current_child_idx < children_count {
		current_child_node := n.Children[n.current_child_idx]

		prev_status := current_child_node.Status()
		child_status := current_child_node.ExecuteTick()

		// switch to RUNNING state as soon as you find an active child
		n.all_skipped_ = (child_status == core.NodeStatus_SKIPPED)

		switch child_status {
		case core.NodeStatus_RUNNING:
			{
				return child_status
			}
		case core.NodeStatus_FAILURE:
			{
				// DO NOT reset current_child_idx_ on failure
				for i := n.current_child_idx; i < len(n.Children); i++ {
					n.HaltChild(i)
				}

				return child_status
			}
		case core.NodeStatus_SUCCESS:
			{
				n.current_child_idx++
				// Return the execution flow if the child is async,
				// to make this interruptable.
				if requiresWakeUp() && prev_status == core.NodeStatus_IDLE &&
					n.current_child_idx < children_count {
					emitWakeUpSignal()
					return core.NodeStatus_RUNNING
				}
			}
			break

		case core.NodeStatus_SKIPPED:
			{
				// It was requested to skip this node
				n.current_child_idx++
			}
			break

		case core.NodeStatus_IDLE:
			{
				panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))
			}
		} // end switch
	} // end while loop

	// The entire while loop completed. This means that all the children returned SUCCESS.
	if n.current_child_idx == children_count {
		n.ResetChildren()
		n.current_child_idx = 0
	}
	if n.all_skipped_ {
		return core.NodeStatus_SKIPPED
	}
	// Skip if ALL the nodes have been skipped
	return core.NodeStatus_SUCCESS
}

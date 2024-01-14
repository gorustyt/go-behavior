package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

type ReactiveSequence struct {
	*core.ControlNode
	runningChild           int
	throwIfMultipleRunning bool
}

func NewReactiveSequence() *ReactiveSequence {
	return &ReactiveSequence{
		runningChild:           -1,
		throwIfMultipleRunning: true,
	}
}

func (n *ReactiveSequence) EnableException(enable bool) {
	n.throwIfMultipleRunning = enable
}

func (n *ReactiveSequence) Tick() core.NodeStatus {
	all_skipped := true
	if n.Status() == core.NodeStatus_IDLE {
		n.runningChild = -1
	}
	n.SetStatus(core.NodeStatus_RUNNING)

	for index := 0; index < len(n.Children); index++ {
		current_child_node := n.Children[index]
		child_status := current_child_node.ExecuteTick()

		// switch to RUNNING state as soon as you find an active child
		all_skipped = child_status == core.NodeStatus_SKIPPED

		switch child_status {
		case core.NodeStatus_RUNNING:
			// reset the previous children, to make sure that they are
			// in IDLE state the next time we tick them
			for i := 0; i < len(n.Children); i++ {
				if i != index {
					n.HaltChild(i)
				}
			}
			if n.runningChild == -1 {
				n.runningChild = index
			} else if n.throwIfMultipleRunning && n.runningChild != int(index) {
				panic("[ReactiveSequence]: only a single child can return RUNNING.\n,This throw can be disabled with ReactiveSequence::EnableException(false)")
			}
			return core.NodeStatus_RUNNING
		case core.NodeStatus_FAILURE:

			n.ResetChildren()
			return core.NodeStatus_FAILURE

			// do nothing if SUCCESS
		case core.NodeStatus_SUCCESS:

		case core.NodeStatus_SKIPPED:

			// to allow it to be skipped again, we must reset the node
			n.HaltChild(index)

		case core.NodeStatus_IDLE:

			panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))

		} // end switch
	} //end for

	n.ResetChildren()
	if all_skipped {
		return core.NodeStatus_SKIPPED
	}
	// Skip if ALL the nodes have been skipped
	return core.NodeStatus_SUCCESS
}

func (n *ReactiveSequence) Halt() {
	n.runningChild = -1
	n.Halt()
}

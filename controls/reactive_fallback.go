package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

type ReactiveFallback struct {
	*core.ControlNode
	runningChild           int
	throwIfMultipleRunning bool
}

func NewReactiveFallback(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &ReactiveFallback{
		ControlNode:            core.NewControlNode(name, cfg),
		runningChild:           -1,
		throwIfMultipleRunning: true,
	}
	return n
}
func (n *ReactiveFallback) Tick() core.NodeStatus {
	allSkipped := true
	if n.Status() == core.NodeStatus_IDLE {
		n.runningChild = -1
	}
	n.SetStatus(core.NodeStatus_RUNNING)
	for index := 0; index < len(n.Children); index++ {
		currentChildNode := n.Children[index]
		childStatus := currentChildNode.ExecuteTick()
		// switch to RUNNING state as soon as you find an active child
		allSkipped = childStatus == core.NodeStatus_SKIPPED
		switch childStatus {
		case core.NodeStatus_RUNNING:
			// reset the previous children, to make sure that they are
			// in IDLE state the next time we tick them
			for i := 0; i < len(n.Children); i++ {
				if i != index {
					n.HaltChild(i)
				}
			}
			if n.runningChild == -1 {
				n.runningChild = int(index)
			} else if n.throwIfMultipleRunning && n.runningChild != int(index) {
				panic("[ReactiveFallback]: only a single child can return RUNNING.\nThis throw can be disabled with ReactiveFallback::EnableException(false)")
			}
			return core.NodeStatus_RUNNING
		case core.NodeStatus_FAILURE:
		case core.NodeStatus_SUCCESS:
			n.ResetChildren()
			return core.NodeStatus_SUCCESS
		case core.NodeStatus_SKIPPED:
			// to allow it to be skipped again, we must reset the node
			n.HaltChild(index)
			break

		case core.NodeStatus_IDLE:
			panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))
		} // end switch
	} //end for

	n.ResetChildren()
	if allSkipped {
		return core.NodeStatus_SKIPPED
	}
	// Skip if ALL the nodes have been skipped
	return core.NodeStatus_FAILURE
}

func (n *ReactiveFallback) Halt() {
	n.runningChild = -1
	n.Halt()
}

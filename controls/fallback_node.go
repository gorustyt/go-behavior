package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

//按顺序执行孩子结点直到其中一个孩子结点返回成功状态或所有孩子结点返回失败状态。一般用来实现角色的备选行为。

type FallbackNode struct {
	*core.ControlNode
	async           bool
	allSkipped      bool
	currentChildIdx int
}

func NewFallbackNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	makeAynch := false
	if len(args) > 0 {
		makeAynch = args[0].(bool)
	}
	n := &FallbackNode{
		ControlNode: core.NewControlNode(name, cfg),
		allSkipped:  true,
		async:       makeAynch,
	}
	if n.async {
		n.SetRegistrationID("AsyncFallback")
	} else {
		n.SetRegistrationID("Fallback")
	}
	return n
}

func (n *FallbackNode) Tick() core.NodeStatus {
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
		n.allSkipped = n.allSkipped && (childStatus == core.NodeStatus_SKIPPED)

		switch childStatus {
		case core.NodeStatus_RUNNING:
			return childStatus
		case core.NodeStatus_SUCCESS:
			n.ResetChildren()
			n.currentChildIdx = 0
			return childStatus
		case core.NodeStatus_FAILURE:
			n.currentChildIdx++
			// Return the execution flow if the child is async,
			// to make this interruptable.
			if n.async && n.RequiresWakeUp() &&
				prevStatus == core.NodeStatus_IDLE &&
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

	// The entire while loop completed. This means that all the children returned FAILURE.
	if n.currentChildIdx == childrenCount {
		n.ResetChildren()
		n.currentChildIdx = 0
	}
	if n.allSkipped {
		return core.NodeStatus_SKIPPED
	}
	// Skip if ALL the nodes have been skipped
	return core.NodeStatus_FAILURE
}

func (n *FallbackNode) Halt() {
	n.currentChildIdx = 0
	n.ControlNode.Halt()
}

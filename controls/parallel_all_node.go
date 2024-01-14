package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

type ParallelAllNode struct {
	*core.ControlNode
	failure_threshold_ int
	completed_list_    map[int]struct{}
	failure_count_     int
}

func NewParallelAllNode(name string, cfg *core.NodeConfig) *ParallelAllNode {
	return &ParallelAllNode{
		ControlNode:        core.NewControlNode(name, cfg),
		completed_list_:    map[int]struct{}{},
		failure_threshold_: 1,
	}
}
func (n *ParallelAllNode) tick() core.NodeStatus {
	max_failures := 0
	v, err := n.GetInput("max_failures")
	max_failures = v.(int)
	if err != nil {
		panic(fmt.Sprintf("Missing parameter [max_failures] in ParallelNode%v", err))
	}
	children_count := len(n.Children)
	n.setFailureThreshold(max_failures)

	skipped_count := 0

	if children_count < n.failure_threshold_ {
		panic("Number of children is less than threshold. Can never fail.")
	}

	n.SetStatus(core.NodeStatus_RUNNING)

	// Routing the tree according to the sequence node's logic:
	for index := 0; index < children_count; index++ {
		child_node := n.Children[index]

		// already completed
		if _, ok := n.completed_list_[index]; ok {
			continue
		}

		child_status := child_node.ExecuteTick()
		switch child_status {
		case core.NodeStatus_SUCCESS:
			n.completed_list_[index] = struct{}{}
		case core.NodeStatus_FAILURE:
			n.completed_list_[index] = struct{}{}
			n.failure_count_++
		case core.NodeStatus_RUNNING:
			// Still working. Check the next
		case core.NodeStatus_SKIPPED:
			skipped_count++
		case core.NodeStatus_IDLE:
			panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))
		}
	}

	if skipped_count == children_count {
		return core.NodeStatus_SKIPPED
	}
	if skipped_count+len(n.completed_list_) >= children_count {
		// DONE
		n.HaltChildren()
		n.completed_list_ = map[int]struct{}{}
		status := core.NodeStatus_SUCCESS
		if n.failure_count_ >= n.failure_threshold_ {
			status = core.NodeStatus_FAILURE
		}
		n.failure_count_ = 0
		return status
	}

	// Some children haven't finished, yet.
	return core.NodeStatus_RUNNING
}

func (n *ParallelAllNode) Halt() {
	n.completed_list_ = map[int]struct{}{}
	n.failure_count_ = 0
	n.ControlNode.Halt()
}

func (n *ParallelAllNode) failureThreshold() int {
	return n.failure_threshold_
}

func (n *ParallelAllNode) setFailureThreshold(threshold int) {
	if threshold < 0 {
		n.failure_threshold_ = max(int(len(n.Children))+threshold+1, 0)
	} else {
		n.failure_threshold_ = threshold
	}
}

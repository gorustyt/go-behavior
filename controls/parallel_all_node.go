package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

func init() {
	core.SetPorts(&ParallelAllNode{},
		core.InputPortWithDefaultValue("max_failures", 1,
			"If the number of children returning FAILURE exceeds this value, ParallelAll returns FAILURE"))
}

type ParallelAllNode struct {
	*core.ControlNode
	failure_threshold_ int
	completed_list_    map[int]struct{}
	failure_count_     int
}

func NewParallelAllNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	return &ParallelAllNode{
		ControlNode:        core.NewControlNode(name, cfg),
		completed_list_:    map[int]struct{}{},
		failure_threshold_: 1,
	}
}
func (n *ParallelAllNode) tick() core.NodeStatus {
	maxFailures := 0
	v, err := n.GetInput("max_failures", &maxFailures)
	maxFailures = v.(int)
	if err != nil {
		panic(fmt.Sprintf("Missing parameter [max_failures] in ParallelNode%v", err))
	}
	childrenCount := len(n.Children)
	n.setFailureThreshold(maxFailures)

	skippedCount := 0

	if childrenCount < n.failure_threshold_ {
		panic("Number of children is less than threshold. Can never fail.")
	}

	n.SetStatus(core.NodeStatus_RUNNING)

	// Routing the tree according to the sequence node's logic:
	for index := 0; index < childrenCount; index++ {
		childNode := n.Children[index]

		// already completed
		if _, ok := n.completed_list_[index]; ok {
			continue
		}

		childStatus := childNode.ExecuteTick()
		switch childStatus {
		case core.NodeStatus_SUCCESS:
			n.completed_list_[index] = struct{}{}
		case core.NodeStatus_FAILURE:
			n.completed_list_[index] = struct{}{}
			n.failure_count_++
		case core.NodeStatus_RUNNING:
			// Still working. Check the next
		case core.NodeStatus_SKIPPED:
			skippedCount++
		case core.NodeStatus_IDLE:
			panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))
		}
	}

	if skippedCount == childrenCount {
		return core.NodeStatus_SKIPPED
	}
	if skippedCount+len(n.completed_list_) >= childrenCount {
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

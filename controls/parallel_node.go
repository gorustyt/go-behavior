package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

func init() {
	core.SetPorts(&ParallelNode{}, core.InputPortWithDefaultValue(THRESHOLD_SUCCESS, -1,
		"number of children that need to succeed to trigger a SUCCESS"))

	core.SetPorts(core.InputPortWithDefaultValue(THRESHOLD_FAILURE, 1,
		"number of children that need to fail to trigger a FAILURE"))
}

const (
	THRESHOLD_SUCCESS = "success_count"
	THRESHOLD_FAILURE = "failure_count"
)

// 并行执行”所有孩子结点。直到至少M个孩子(M的值在1到N之间)结点返回成功状态或所有孩子结点返回失败状态。
type ParallelNode struct {
	*core.ControlNode
	success_threshold_ int
	failure_threshold_ int

	completedList map[int]struct{}

	successCount int
	failureCount int

	readParameterFromPorts bool
}

func NewParallelNodeWithName(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &ParallelNode{
		ControlNode:            core.NewControlNode(name, core.DefaultNodeConfig),
		completedList:          map[int]struct{}{},
		success_threshold_:     -1,
		failure_threshold_:     1,
		readParameterFromPorts: true,
	}
	if cfg == nil {
		n.SetRegistrationID("Parallel")
		n.readParameterFromPorts = false
	}
	return n
}

func NewParallelNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	return &ParallelNode{
		ControlNode:            core.NewControlNode(name, cfg),
		completedList:          map[int]struct{}{},
		success_threshold_:     -1,
		failure_threshold_:     1,
		readParameterFromPorts: true,
	}
}

func (n *ParallelNode) Tick() core.NodeStatus {
	if n.readParameterFromPorts {
		v, err := n.GetInput(THRESHOLD_SUCCESS, 0)
		n.success_threshold_ = v.(int)
		if err != nil {
			panic(fmt.Sprintf("Missing parameter [%v] in ParallelNode err:%v", THRESHOLD_SUCCESS, err))
		}
		v, err = n.GetInput(THRESHOLD_FAILURE, 0)
		n.failure_threshold_ = v.(int)
		if err != nil {
			panic(fmt.Sprintf("Missing parameter [%v] in ParallelNode :%v", THRESHOLD_FAILURE, err))
		}
	}

	childrenCount := len(n.Children)

	if childrenCount < n.successThreshold() {
		panic("Number of children is less than threshold. Can never succeed.")
	}

	if childrenCount < n.failureThreshold() {
		panic("Number of children is less than threshold. Can never fail.")
	}

	n.SetStatus(core.NodeStatus_RUNNING)

	skippedCount := 0

	// Routing the tree according to the sequence node's logic:
	for i := 0; i < childrenCount; i++ {
		if len(n.completedList) == 0 {
			childNode := n.Children[i]
			childStatus := childNode.ExecuteTick()
			switch childStatus {
			case core.NodeStatus_SKIPPED:
				skippedCount++
			case core.NodeStatus_SUCCESS:
				n.completedList[i] = struct{}{}
				n.successCount++
			case core.NodeStatus_FAILURE:
				n.completedList[i] = struct{}{}
				n.failureCount++
			case core.NodeStatus_RUNNING:
				// Still working. Check the next
			case core.NodeStatus_IDLE:
				{
					panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))
				}
			}
		}

		requiredSuccessCount := n.successThreshold()

		if n.successCount >= requiredSuccessCount ||
			(n.success_threshold_ < 0 && (n.successCount+skippedCount) >= requiredSuccessCount) {
			n.clear()
			n.ResetChildren()
			return core.NodeStatus_SUCCESS
		}

		// It fails if it is not possible to succeed anymore or if
		// number of failures are equal to failure_threshold_
		if ((childrenCount - n.failureCount) < requiredSuccessCount) ||
			(n.failureCount == n.failureThreshold()) {
			n.clear()
			n.ResetChildren()
			return core.NodeStatus_FAILURE
		}
	}
	if skippedCount == childrenCount {
		return core.NodeStatus_SKIPPED
	}
	// Skip if ALL the nodes have been skipped
	return core.NodeStatus_RUNNING
}

func (n *ParallelNode) clear() {
	n.completedList = map[int]struct{}{}
	n.successCount = 0
	n.failureCount = 0
}

func (n *ParallelNode) Halt() {
	n.clear()
	n.ControlNode.Halt()
}

func (n *ParallelNode) successThreshold() int {
	if n.success_threshold_ < 0 {
		return max(len(n.Children)+n.success_threshold_+1, 0)
	} else {
		return n.success_threshold_
	}
}

func (n *ParallelNode) failureThreshold() int {
	if n.failure_threshold_ < 0 {
		return max(len(n.Children)+n.failure_threshold_+1, 0)
	} else {
		return n.failure_threshold_
	}
}

func (n *ParallelNode) setSuccessThreshold(threshold int) {
	n.success_threshold_ = threshold
}

func (n *ParallelNode) SetFailureThreshold(threshold int) {
	n.failure_threshold_ = threshold
}

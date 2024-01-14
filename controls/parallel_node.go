package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

const (
	THRESHOLD_SUCCESS = "success_count"
	THRESHOLD_FAILURE = "failure_count"
)

// 并行执行”所有孩子结点。直到至少M个孩子(M的值在1到N之间)结点返回成功状态或所有孩子结点返回失败状态。
type ParallelNode struct {
	*core.ControlNode
	success_threshold_ int
	failure_threshold_ int

	completed_list_ map[int]struct{}

	success_count_ int
	failure_count_ int

	read_parameter_from_ports_ bool
}

func NewParallelNodeWithName(name string) *ParallelNode {
	return &ParallelNode{
		ControlNode:                core.NewControlNode(name, core.DefaultNodeConfig),
		completed_list_:            map[int]struct{}{},
		success_threshold_:         -1,
		failure_threshold_:         1,
		read_parameter_from_ports_: false,
	}
}

func NewParallelNode(name string, cfg *core.NodeConfig) *ParallelNode {
	return &ParallelNode{
		ControlNode:                core.NewControlNode(name, cfg),
		completed_list_:            map[int]struct{}{},
		success_threshold_:         -1,
		failure_threshold_:         1,
		read_parameter_from_ports_: true,
	}
}
func (n *ParallelNode) Tick() core.NodeStatus {
	if n.read_parameter_from_ports_ {
		v, err := n.GetInput(THRESHOLD_SUCCESS)
		n.success_threshold_ = v.(int)
		if err != nil {
			panic(fmt.Sprintf("Missing parameter [%v] in ParallelNode err:%v", THRESHOLD_SUCCESS, err))
		}
		v, err = n.GetInput(THRESHOLD_FAILURE)
		n.failure_threshold_ = v.(int)
		if err != nil {
			panic(fmt.Sprintf("Missing parameter [%v] in ParallelNode :%v", THRESHOLD_FAILURE, err))
		}
	}

	children_count := len(n.Children)

	if children_count < n.successThreshold() {
		panic("Number of children is less than threshold. Can never succeed.")
	}

	if children_count < n.failureThreshold() {
		panic("Number of children is less than threshold. Can never fail.")
	}

	n.SetStatus(core.NodeStatus_RUNNING)

	skipped_count := 0

	// Routing the tree according to the sequence node's logic:
	for i := 0; i < children_count; i++ {
		if len(n.completed_list_) == 0 {
			child_node := n.Children[i]
			child_status := child_node.ExecuteTick()

			switch child_status {
			case core.NodeStatus_SKIPPED:
				{
					skipped_count++
				}

			case core.NodeStatus_SUCCESS:
				{
					n.completed_list_[i] = struct{}{}
					n.success_count_++
				}

			case core.NodeStatus_FAILURE:
				{
					n.completed_list_[i] = struct{}{}
					n.failure_count_++
				}

			case core.NodeStatus_RUNNING:
				{
					// Still working. Check the next
				}

			case core.NodeStatus_IDLE:
				{
					panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))
				}
			}
		}

		required_success_count := n.successThreshold()

		if n.success_count_ >= required_success_count ||
			(n.success_threshold_ < 0 && (n.success_count_+skipped_count) >= required_success_count) {
			n.clear()
			n.ResetChildren()
			return core.NodeStatus_SUCCESS
		}

		// It fails if it is not possible to succeed anymore or if
		// number of failures are equal to failure_threshold_
		if ((children_count - n.failure_count_) < required_success_count) ||
			(n.failure_count_ == n.failureThreshold()) {
			n.clear()
			n.ResetChildren()
			return core.NodeStatus_FAILURE
		}
	}
	if skipped_count == children_count {
		return core.NodeStatus_SKIPPED
	}
	// Skip if ALL the nodes have been skipped
	return core.NodeStatus_RUNNING
}

func (n *ParallelNode) clear() {
	n.completed_list_ = map[int]struct{}{}
	n.success_count_ = 0
	n.failure_count_ = 0
}

func (n *ParallelNode) Halt() {
	n.clear()
	n.ControlNode.Halt()
}

func (n *ParallelNode) successThreshold() int {
	if n.success_threshold_ < 0 {
		return max(int(len(n.Children))+n.success_threshold_+1, 0)
	} else {
		return n.success_threshold_
	}
}

func (n *ParallelNode) failureThreshold() int {
	if n.failure_threshold_ < 0 {
		return max(int(len(n.Children))+n.failure_threshold_+1, 0)
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

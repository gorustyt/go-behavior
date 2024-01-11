package go_behavior

type NodeType int

const (
	UNDEFINED NodeType = iota
	ACTION
	CONDITION
	SUBTREE
	DECORATOR
	CONTROL
)

type NodeStatus int

const (
	IDLE NodeStatus = iota
	RUNNING
	SUCCESS
	FAILURE
	SKIPPED
)

func (status NodeStatus) IsActive() bool {
	return status != IDLE && status != SKIPPED
}

func (status NodeStatus) IsCompleted() bool {
	return status == SUCCESS || status == FAILURE
}

type TreeNode struct {
}

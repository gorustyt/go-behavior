package actions

import "github.com/gorustyt/go-behavior/core"

type AlwaysFailureNode struct {
	*core.SyncActionNode
}

func NewAlwaysFailureNode() *AlwaysFailureNode {
	return &AlwaysFailureNode{}
}

func (n *AlwaysFailureNode) Tick() core.NodeStatus {
	return core.FAILURE
}

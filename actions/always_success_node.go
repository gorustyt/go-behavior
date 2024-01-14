package actions

import "github.com/gorustyt/go-behavior/core"

type AlwaysSuccessNode struct {
	*core.SyncActionNode
}

func NewAlwaysSuccessNode() *AlwaysSuccessNode {
	return &AlwaysSuccessNode{}
}

func (n *AlwaysSuccessNode) Tick() core.NodeStatus {
	return core.SUCCESS
}

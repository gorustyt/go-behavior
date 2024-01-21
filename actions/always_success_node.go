package actions

import "github.com/gorustyt/go-behavior/core"

type AlwaysSuccessNode struct {
	*core.SyncActionNode
}

func NewAlwaysSuccessNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &AlwaysSuccessNode{SyncActionNode: core.NewSyncActionNode(name, cfg)}
	n.SetRegistrationID("AlwaysSuccess")
	return n
}

func (n *AlwaysSuccessNode) Tick() core.NodeStatus {
	return core.NodeStatus_SUCCESS
}

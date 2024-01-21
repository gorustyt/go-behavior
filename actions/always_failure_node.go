package actions

import "github.com/gorustyt/go-behavior/core"

type AlwaysFailureNode struct {
	*core.SyncActionNode
}

func NewAlwaysFailureNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &AlwaysFailureNode{
		SyncActionNode: core.NewSyncActionNode(name, cfg),
	}
	n.SetRegistrationID("AlwaysFailure")
	return n
}

func (n *AlwaysFailureNode) Tick() core.NodeStatus {
	return core.NodeStatus_FAILURE
}

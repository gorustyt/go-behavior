package actions

import "github.com/gorustyt/go-behavior/core"

func init() {
	core.SetPorts(&UnsetBlackboardNode{}, core.InputPortWithDefaultValue(
		"key", "", "Key of the entry to remove",
	))
}

type UnsetBlackboardNode struct {
	*core.SyncActionNode
}

func NewUnsetBlackboardNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &UnsetBlackboardNode{SyncActionNode: core.NewSyncActionNode(name, cfg)}
	n.SetRegistrationID("UnsetBlackboard")
	return n
}

func (n *UnsetBlackboardNode) Tick() core.NodeStatus {
	var key string
	v, err := n.GetInput("key", &key)
	if err != nil {
		panic("missing input port [key]")
	}
	key = v.(string)
	n.Config().Blackboard.Unset(key)
	return core.NodeStatus_SUCCESS
}

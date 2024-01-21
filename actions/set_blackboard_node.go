package actions

import "github.com/gorustyt/go-behavior/core"

func init() {
	core.SetPorts(&SetBlackboardNode{}, core.InputPort("value", "Value to be written int othe output_key"))
	core.SetPorts(&SetBlackboardNode{}, core.BidirectionalPort("output_key", "Name of the blackboard entry where the value should be written"))
}

type SetBlackboardNode struct {
	*core.SyncActionNode
}

func NewSetBlackboardNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &SetBlackboardNode{SyncActionNode: core.NewSyncActionNode(name, cfg)}
	n.SetRegistrationID("SetBlackboard")
	return n
}

func (n *SetBlackboardNode) tick() core.NodeStatus {
	var outputKey string
	if v, err := n.GetInput("output_key", outputKey); err != nil {
		panic("missing port [output_key]")
	} else {
		outputKey = v.(string)
	}

	valueStr := n.Config().InputPorts["value"]

	if inputKey, ok := core.IsBlackboardPointer(valueStr); ok {
		srcEntry := n.Config().Blackboard.GetEntry(inputKey)
		dstEntry := n.Config().Blackboard.GetEntry(outputKey)

		if srcEntry == nil {
			panic("Can't find the port referred by [value]")
		}
		if dstEntry == nil {
			n.Config().Blackboard.CreateEntry(outputKey, srcEntry.info)
			dstEntry = n.Config().Blackboard.GetEntry(outputKey)
		}
		dstEntry.Value = srcEntry.Value
	} else {
		n.Config().Blackboard.Set(outputKey, valueStr)
	}

	return core.NodeStatus_SUCCESS
}

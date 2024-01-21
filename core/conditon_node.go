package core

type ConditionNode struct {
	*LeafNode
}

func NewConditionNode(name string, config *NodeConfig) *ConditionNode {
	return &ConditionNode{LeafNode: NewLeafNode(name, config)}
}

func (n *ConditionNode) NodeType() NodeType {
	return NodeType_CONDITION
}

// Do nothing
func (n *ConditionNode) Halt() {
	n.ResetStatus()
}

type SimpleConditionNode struct {
	*ConditionNode
	tickFunctor TickFunctor
}

func NewSimpleConditionNode(name string, config *NodeConfig, args ...interface{}) ITreeNode {
	n := &SimpleConditionNode{
		ConditionNode: NewConditionNode(name, config),
	}
	if len(args) > 0 {
		n.tickFunctor = args[0].(TickFunctor)
	}
	return n
}

func (s *SimpleConditionNode) Tick() NodeStatus {
	return s.tickFunctor(s)
}

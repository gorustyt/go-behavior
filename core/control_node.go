package core

type ControlNode struct {
	*TreeNode
	Children []*TreeNode
}

func NewControlNode(name string, config *NodeConfig) *ControlNode {
	return &ControlNode{}
}
func (n *ControlNode) Type() NodeType {
	return NodeType_CONTROL
}

func (n *ControlNode) HaltChild(i int) {
	child := n.Children[i]
	if child.status == NodeStatus_RUNNING {
		child.HaltNode()
	}
	child.resetStatus()
}

func (n *ControlNode) HaltChildren() {
	for i := range n.Children {
		n.HaltChild(i)
	}
}

func (n *ControlNode) Halt() {
	n.ResetChildren()
	n.resetStatus() // might be redundant
}

func (n *ControlNode) AddChild(child *TreeNode) {
	n.Children = append(n.Children, child)
}

func (n *ControlNode) ResetChildren() {
	for _, child := range n.Children {
		if child.status == NodeStatus_RUNNING {
			child.HaltNode()
		}
		child.resetStatus()
	}
}

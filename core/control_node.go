package core

type ControlNode struct {
	*TreeNode
	Children []ITreeNode
}

func NewControlNode(name string, config *NodeConfig) *ControlNode {
	return &ControlNode{TreeNode: NewTreeNode(name, config)}
}
func (n *ControlNode) NodeType() NodeType {
	return NodeType_CONTROL
}

func (n *ControlNode) HaltChild(i int) {
	child := n.Children[i]
	if child.Status() == NodeStatus_RUNNING {
		child.HaltNode()
	}
	child.ResetStatus()
}

func (n *ControlNode) HaltChildren() {
	for i := range n.Children {
		n.HaltChild(i)
	}
}

func (n *ControlNode) Halt() {
	n.ResetChildren()
	n.ResetStatus() // might be redundant
}

func (n *ControlNode) AddChild(child ITreeNode) {
	n.Children = append(n.Children, child)
}

func (n *ControlNode) ResetChildren() {
	for _, child := range n.Children {
		if child.Status() == NodeStatus_RUNNING {
			child.HaltNode()
		}
		child.ResetStatus()
	}
}

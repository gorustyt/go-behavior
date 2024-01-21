package core

type LeafNode struct {
	*TreeNode
}

func NewLeafNode(name string, config *NodeConfig) *LeafNode {
	return &LeafNode{NewTreeNode(name, config)}
}

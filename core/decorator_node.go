package core

import (
	"fmt"
)

// 以自定义的方式修改孩子结点的行为。比如Invert类型的装饰结点，可以反转其孩子结点返回的状态信息。为了方便他人理解，应该尽可能使用比较常见的装饰结点。
type DecoratorNode struct {
	*TreeNode
	childNode ITreeNode
}

func NewDecoratorNode(name string, config *NodeConfig) *DecoratorNode {
	return &DecoratorNode{
		TreeNode: NewTreeNode(name, config),
	}
}
func (n *DecoratorNode) NodeType() NodeType {
	return NodeType_CONDITION
}

func (n *DecoratorNode) SetChild(child ITreeNode) error {
	if n.childNode != nil {
		return fmt.Errorf("decorator [%v] has already a child assigned", n.Name())
	}
	n.childNode = child
	return nil
}

func (n *DecoratorNode) Halt() {
	n.ResetChild()
	n.ResetStatus() // might be redundant
}

func (n *DecoratorNode) Child() ITreeNode {
	return n.childNode
}
func (n *DecoratorNode) HaltChild() {
	n.ResetChild()
}

func (n *DecoratorNode) ResetChild() {
	if n.childNode == nil {
		return
	}
	if n.childNode.Status() == NodeStatus_RUNNING {
		n.childNode.HaltNode()
	}
	n.childNode.ResetStatus()
}

type SimpleDecoratorNode struct {
	*DecoratorNode
	tickFn TickFunctor
}

func (n *DecoratorNode) ExecuteTick() NodeStatus {
	status := n.TreeNode.ExecuteTick()
	childStatus := n.Child().Status()
	if childStatus == NodeStatus_SUCCESS || childStatus == NodeStatus_FAILURE {
		n.Child().ResetStatus()
	}
	return status
}

func NewSimpleDecoratorNode(name string,
	config *NodeConfig, args ...interface{}) ITreeNode {
	n := &SimpleDecoratorNode{DecoratorNode: NewDecoratorNode(name, config)}
	if len(args) > 0 {
		n.tickFn = args[0].(TickFunctor)
	}
	return n
}

func (n *SimpleDecoratorNode) Tick() NodeStatus {
	return n.tickFn(n, n.Child().ExecuteTick())
}

package core

import (
	"sync"
	"sync/atomic"
)

type ActionNodeBase struct {
	*LeafNode
}

func NewActionNodeBase(name string, config *NodeConfig) *ActionNodeBase {
	return &ActionNodeBase{LeafNode: NewLeafNode(name, config)}
}

func (ActionNodeBase) NodeType() NodeType {
	return NodeType_ACTION
}

type SyncActionNode struct {
	*ActionNodeBase
}

func NewSyncActionNode(name string, config *NodeConfig) *SyncActionNode {
	return &SyncActionNode{
		ActionNodeBase: NewActionNodeBase(name, config),
	}
}
func (s *SyncActionNode) ExecuteTick() NodeStatus {
	stat := s.ActionNodeBase.ExecuteTick()
	if stat == NodeStatus_RUNNING {
		panic("SyncActionNode MUST never return RUNNING")
	}
	return stat
}

type SimpleActionNode struct {
	*SyncActionNode
	tickFunctor TickFunctor
}

func NewSimpleActionNode(name string, config *NodeConfig, args ...interface{}) ITreeNode {
	n := &SimpleActionNode{
		SyncActionNode: NewSyncActionNode(name, config),
	}
	if len(args) > 0 {
		n.tickFunctor = args[0].(TickFunctor)
	}
	return n
}
func (n *SimpleActionNode) Tick() NodeStatus {
	prevStatus := n.Status()

	if prevStatus == NodeStatus_IDLE {
		n.SetStatus(NodeStatus_RUNNING)
		prevStatus = NodeStatus_RUNNING
	}
	status := n.tickFunctor(n)
	if status != prevStatus {
		n.SetStatus(status)
	}
	return status
}

type ThreadedAction struct {
	*ActionNodeBase
}

type IStatefulActionNode interface {
	OnStart() NodeStatus
	OnRunning() NodeStatus
	OnHalted()
}

type StatefulActionNode struct {
	halt_requested_ atomic.Bool
	mutex_          sync.Mutex
	*ActionNodeBase
	IStatefulActionNode
}

func NewStatefulActionNode(name string, config *NodeConfig) *StatefulActionNode {
	return &StatefulActionNode{
		ActionNodeBase: NewActionNodeBase(name, config),
	}
}
func (n *StatefulActionNode) IsHaltRequested() bool {
	return n.halt_requested_.Load()
}

func (n *StatefulActionNode) Tick() NodeStatus {
	prevStatus := n.Status()

	if prevStatus == NodeStatus_IDLE {
		newStatus := n.OnStart()
		if newStatus == NodeStatus_IDLE {
			panic("StatefulActionNode::onStart() must not return IDLE")
		}
		return newStatus
	}
	//------------------------------------------
	if prevStatus == NodeStatus_RUNNING {
		newStatus := n.OnRunning()
		if newStatus == NodeStatus_IDLE {
			panic("StatefulActionNode::onRunning() must not return IDLE")
		}
		return newStatus
	}
	return prevStatus
}

func (n *StatefulActionNode) Halt() {
	n.halt_requested_.Store(true)
	if n.Status() == NodeStatus_RUNNING {
		n.OnHalted()
	}
	n.ResetStatus() // might be redundant
}

type CoroActionNode struct {
	*ActionNodeBase
}

package core

import "time"

type Subtree struct {
	TreeId       string
	InstanceName string
	Blackboard   *Blackboard
	Nodes        []ITreeNode
	Subtree      *Subtree
}

type Tree struct {
	uidCounter uint16
	Subtrees   []*Subtree
	manifests  map[string]*TreeNodeManifest
	wakeUp     *WakeUpSignal
}

func NewTree() *Tree {
	return &Tree{manifests: make(map[string]*TreeNodeManifest)}
}
func (t *Tree) GetUID() uint16 {
	t.uidCounter++
	return t.uidCounter
}
func (t *Tree) Root() ITreeNode {
	if len(t.Subtrees) == 0 {
		return nil
	}
	subtreeNodes := t.Subtrees[0].Nodes
	if len(subtreeNodes) == 0 {
		return nil
	}
	return subtreeNodes[0]
}

func (t *Tree) TickOnce() NodeStatus {
	return t.tickRoot(ONCE_UNLESS_WOKEN_UP, 0)
}

func (t *Tree) TickExactlyOnce() NodeStatus {
	return t.tickRoot(EXACTLY_ONCE, 0)
}

func (t *Tree) TickWhileRunning(sleepTimes ...time.Duration) NodeStatus {
	sleepTime := 10 * time.Millisecond
	if len(sleepTimes) > 0 {
		sleepTime = sleepTimes[0]
	}
	return t.tickRoot(WHILE_RUNNING, sleepTime)
}

func (t *Tree) TickRoot(opt TickOption, sleepTime time.Duration) NodeStatus {
	return t.tickRoot(opt, sleepTime)
}
func (t *Tree) tickRoot(opt TickOption, sleepTime time.Duration) NodeStatus {
	root := t.Root()
	status := NodeStatus_IDLE
	for status == NodeStatus_IDLE ||
		(opt == WHILE_RUNNING && status == NodeStatus_RUNNING) {
		status = root.ExecuteTick()
		for opt != EXACTLY_ONCE &&
			status == NodeStatus_RUNNING &&
			t.wakeUp.WaitFor(time.Duration(0)) {
			status = root.ExecuteTick()
		}
		if IsStatusCompleted(status) {
			root.ResetStatus()
		}
		if status == NodeStatus_RUNNING && sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}

	return status
}

func (t *Tree) Init() {
	t.wakeUp = NewWakeUpSignal()
	for _, subtree := range t.Subtrees {
		for _, node := range subtree.Nodes {
			node.SetWakeUpInstance(t.wakeUp)
		}
	}
}

func (t *Tree) haltTree() {
	root := t.Root()
	if root == nil {
		return
	}
	// the halt should propagate to all the node if the nodes
	// have been implemented correctly
	t.Root().HaltNode()
	//but, just in case.... this should be no-op
	for _, v := range t.Subtrees {
		for _, node := range v.Nodes {
			node.HaltNode()
		}
	}
	root.ResetStatus()
}

func (t *Tree) Sleep(timeout time.Duration) {
	t.wakeUp.WaitFor(timeout)
}

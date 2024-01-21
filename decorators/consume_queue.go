package decorators

import "github.com/gorustyt/go-behavior/core"

func init() {
	core.SetPorts(&ConsumeQueue{}, core.InputPortWithDefaultValue("queue", &core.ProtectedQueue{}))
	core.SetPorts(ConsumeQueue{}, core.OutputPortWithDefaultValue("popped_item", &core.ProtectedQueue{}))

}

type ConsumeQueue struct {
	*core.DecoratorNode
	runningChild bool
}

func NewConsumeQueue(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	return &ConsumeQueue{DecoratorNode: core.NewDecoratorNode(name, cfg)}
}
func (n *ConsumeQueue) Tick() core.NodeStatus {
	// by default, return SUCCESS, even if queue is empty
	statusToBeReturned := core.NodeStatus_SUCCESS

	if n.runningChild {
		childState := n.Child().ExecuteTick()
		n.runningChild = childState == core.NodeStatus_RUNNING
		if n.runningChild {
			return core.NodeStatus_RUNNING
		} else {
			n.HaltChild()
			statusToBeReturned = childState
		}
	}

	var queue *core.ProtectedQueue
	v, err := n.GetInput("queue", queue)
	if err != nil {
		panic(err)
	}
	queue = v.(*core.ProtectedQueue)
	if queue != nil {
		queue.Mtx.Lock()
		items := queue.Items

		for items.Len() != 0 {
			n.SetStatus(core.NodeStatus_RUNNING)
			val := items.Front()
			items.Remove(val)
			n.SetOutput("popped_item", val)
			queue.Mtx.Unlock()
			childState := n.Child().ExecuteTick()
			queue.Mtx.Lock()
			n.runningChild = childState == core.NodeStatus_RUNNING
			if n.runningChild {
				return core.NodeStatus_RUNNING
			} else {
				n.HaltChild()
				if childState == core.NodeStatus_FAILURE {
					return core.NodeStatus_FAILURE
				}
				statusToBeReturned = childState
			}
		}
		queue.Mtx.Unlock()
	}

	return statusToBeReturned
}

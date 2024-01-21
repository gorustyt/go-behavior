package actions

import (
	"github.com/gorustyt/go-behavior/core"
)

func init() {
	core.SetPorts(&PopFromQueue{}, core.InputPort("queue"))
	core.SetPorts(&PopFromQueue{}, core.OutPort("popped_item"))

	core.SetPorts(&QueueSize{}, core.InputPort("queue"))
	core.SetPorts(&QueueSize{}, core.OutPort("size"))
}

/*
 * Few words about why we represent the queue as std::shared_ptr<ProtectedQueue>:
 *
 * Since we will pop from the queue, the fact that the blackboard uses
 * a value semantic is not very convenient, since it would oblige us to
 * copy the entire std::list from the BB and than copy again a new one with one less element.
 *
 * We avoid this using reference semantic (wrapping the object in a shared_ptr).
 * Unfortunately, remember that this makes our access to the list not thread-safe!
 * This is the reason why we add a mutex to be used when modyfying the ProtectedQueue::items
 *
 * */

func NewPopFromQueue(name string, config *core.NodeConfig, args ...interface{}) core.ITreeNode {
	return &PopFromQueue{
		SyncActionNode: core.NewSyncActionNode(name, config),
	}
}

type PopFromQueue struct {
	*core.SyncActionNode
}

func (n *PopFromQueue) Tick() core.NodeStatus {
	var queue *core.ProtectedQueue
	v, err := n.GetInput("queue", queue)
	if err != nil {
		panic(err)
	}
	queue = v.(*core.ProtectedQueue)
	if queue != nil {
		queue.Mtx.Lock()
		defer queue.Mtx.Unlock()
		items := queue.Items
		if items.Len() == 0 {
			return core.NodeStatus_FAILURE
		} else {
			val := items.Front()
			items.Remove(val)
			n.SetOutput("popped_item", val)
			return core.NodeStatus_SUCCESS
		}
	} else {
		return core.NodeStatus_FAILURE
	}
}

/**
 * Get the size of a queue. Usefull is you want to write something like:
 *
 *  <QueueSize queue="{waypoints}" size="{wp_size}" />
 *  <Repeat num_cycles="{wp_size}" >
 *      <Sequence>
 *          <PopFromQueue  queue="{waypoints}" popped_item="{wp}" >
 *          <UseWaypoint   waypoint="{wp}" />
 *      </Sequence>
 *  </Repeat>
 */

func NewQueueSize(name string, config *core.NodeConfig, args ...interface{}) core.ITreeNode {
	return &PopFromQueue{
		SyncActionNode: core.NewSyncActionNode(name, config),
	}
}

type QueueSize struct {
	*core.SyncActionNode
}

func (n *QueueSize) Tick() core.NodeStatus {
	var queue *core.ProtectedQueue
	r, err := n.GetInput("queue", queue)
	if err != nil {
		panic(err)
	}
	queue = r.(*core.ProtectedQueue)
	if queue != nil {
		queue.Mtx.Lock()
		defer queue.Mtx.Unlock()
		items := queue.Items

		if items.Len() == 0 {
			return core.NodeStatus_FAILURE
		} else {
			n.SetOutput("size", int(items.Len()))

			return core.NodeStatus_SUCCESS
		}
	}
	return core.NodeStatus_FAILURE
}

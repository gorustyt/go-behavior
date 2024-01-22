package decorators

import (
	"container/list"
	"github.com/gorustyt/go-behavior/core"
	"reflect"
)

func init() {
	core.SetPorts(&LoopNode{}, core.OutPort("value"))
	core.SetPorts(&LoopNode{}, core.BidirectionalPort("queue"))
	core.SetPorts(&LoopNode{}, core.InputPortWithDefaultValue("if_empty", core.NodeStatus_SUCCESS, "Status to return if queue is empty: SUCCESS, FAILURE, SKIPPED"))
}

type LoopNode struct {
	*core.DecoratorNode
	child_running_ bool
	static_queue_  *list.List
	current_queue_ *list.List
	Type           reflect.Kind
}

func NewLoopNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	return &LoopNode{Type: args[0].(reflect.Kind), DecoratorNode: core.NewDecoratorNode(name, cfg)}
}

func (n *LoopNode) Tick() core.NodeStatus {
	popped := false
	if n.Status() == core.NodeStatus_IDLE {
		n.child_running_ = false
		// special case: the port contains a string that was converted to SharedQueue<T>
		if n.static_queue_ != nil {
			n.current_queue_ = n.static_queue_
		}
	}

	// Pop value from queue, if the child is not RUNNING
	if !n.child_running_ {
		// if the port is static, any_ref is empty, otherwise it will keep access to
		// port locked for thread-safety
		var v *core.Entry
		if n.static_queue_ == nil {
			any_ref := n.GetLockedPortContent("queue")
			if any_ref() != nil {
				v = any_ref()
			}
		}

		if v != nil {
			n.current_queue_ = v.Value.(*list.List)
		}

		if n.current_queue_ != nil && n.current_queue_.Len() != 0 {
			value := n.current_queue_.Front()
			n.current_queue_.Remove(value)
			popped = true
			n.SetOutput("value", value.Value)
		}
	}

	if !popped && !n.child_running_ {
		var status core.NodeStatus
		t, err := n.GetInput("if_empty", &status)
		if err == nil {
			status = t.(core.NodeStatus)
		} else {
			panic(err)
		}
		return status
	}

	if n.Status() == core.NodeStatus_IDLE {
		n.SetStatus(core.NodeStatus_RUNNING)
	}

	childState := n.Child().ExecuteTick()
	n.child_running_ = childState == core.NodeStatus_RUNNING

	if core.IsStatusCompleted(childState) {
		n.ResetChild()
	}

	if childState == core.NodeStatus_FAILURE {
		return core.NodeStatus_FAILURE
	}
	return core.NodeStatus_RUNNING
}

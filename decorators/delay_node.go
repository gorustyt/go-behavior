package decorators

import (
	"github.com/gorustyt/go-behavior/core"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	core.SetPorts(&DelayNode{}, core.InputPortWithDefaultValue("delay_msec", 1, "Tick the child after a few milliseconds"))
}

type DelayNode struct {
	*core.DecoratorNode
	timer                      *time.Timer
	delay_started_             bool
	delay_complete_            atomic.Bool
	delay_aborted_             bool
	msec_                      time.Duration
	read_parameter_from_ports_ bool
	delay_mutex_               sync.Mutex
}

func NewDelayNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	return &DelayNode{}
}

func (n *DelayNode) Halt() {
	n.delay_started_ = false
	n.timer.Stop()
	n.DecoratorNode.Halt()
}

func (n *DelayNode) Tick() core.NodeStatus {
	if n.read_parameter_from_ports_ {
		var msec_ int
		if v, err := n.GetInput("delay_msec", &msec_); err != nil {
			panic("Missing parameter [delay_msec] in DelayNode")
		} else {
			msec_ = v.(int)
		}
	}

	if !n.delay_started_ {
		n.delay_complete_.Store(false)
		n.delay_aborted_ = false
		n.delay_started_ = true
		n.SetStatus(core.NodeStatus_RUNNING)

		n.timer = time.AfterFunc(n.msec_, func() {
			n.delay_mutex_.Lock()
			n.delay_complete_.Store(true)
			n.EmitWakeUpSignal()
			n.delay_mutex_.Unlock()
		})
	}

	n.delay_mutex_.Lock()
	defer n.delay_mutex_.Unlock()
	if n.delay_aborted_ {
		n.delay_aborted_ = false
		n.delay_started_ = false
		return core.NodeStatus_FAILURE
	} else if n.delay_complete_.Load() {
		childStatus := n.Child().ExecuteTick()
		if core.IsStatusCompleted(childStatus) {
			n.delay_started_ = false
			n.delay_aborted_ = false
			n.ResetChild()
		}
		return childStatus
	} else {
		return core.NodeStatus_RUNNING
	}
}

package decorators

import (
	"github.com/gorustyt/go-behavior/core"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	core.SetPorts(&TimeoutNode{}, core.InputPortWithDefaultValue("msec", 1, "After a certain amount of time, halt() the child if it is still running."))
}

type TimeoutNode struct {
	*core.DecoratorNode
	childHalted            atomic.Bool
	timerId                int64
	msec_                  time.Duration
	readParameterFromPorts bool
	timeoutStarted         atomic.Bool
	timeoutMutex           sync.Mutex
	timer                  *time.Timer
}

func NewTimeoutNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	milliseconds := 0
	if len(args) > 0 {
		milliseconds = args[0].(int)
	}
	return &TimeoutNode{msec_: time.Duration(milliseconds) * time.Millisecond, DecoratorNode: core.NewDecoratorNode(name, cfg)}
}

func (n *TimeoutNode) Tick() core.NodeStatus {
	var msec_ int
	if n.readParameterFromPorts {
		if v, err := n.GetInput("msec", &msec_); err != nil {
			panic("Missing parameter [msec] in TimeoutNode")
		} else {
			msec_ = v.(int)
		}
	}
	if !n.timeoutStarted.Load() {
		n.timeoutStarted.Store(true)
		n.timeoutMutex.Lock()
		n.SetStatus(core.NodeStatus_RUNNING)
		n.timeoutMutex.Unlock()
		n.childHalted.Store(false)

		if n.msec_ > 0 {
			n.timer = time.AfterFunc(n.msec_, func() {
				n.timeoutMutex.Lock()
				if n.Child().Status() == core.NodeStatus_RUNNING {
					n.childHalted.Store(true)
					n.HaltChild()
					n.EmitWakeUpSignal()
				}
				n.timeoutMutex.Unlock()
			})
		}
	}

	if n.childHalted.Load() {
		n.timeoutStarted.Store(false)
		return core.NodeStatus_FAILURE
	} else {
		n.timeoutMutex.Lock()
		childStatus := n.Child().ExecuteTick()
		n.timeoutMutex.Unlock()
		if core.IsStatusCompleted(childStatus) {
			n.timeoutStarted.Store(false)
			n.timeoutMutex.Lock()
			n.timer.Stop()
			n.ResetChild()
			n.timeoutMutex.Unlock()
		}
		return childStatus
	}
}

func (n *TimeoutNode) Halt() {
	n.timeoutStarted.Store(false)
	n.timer.Stop()
	n.DecoratorNode.Halt()
}

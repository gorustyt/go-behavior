package actions

import (
	"github.com/gorustyt/go-behavior/core"
	"sync"
	"time"
)

func init() {
	core.SetPorts(&SleepNode{}, core.InputPortWithDefaultValue("msec", 0))
}

type SleepNode struct {
	*core.StatefulActionNode
	cfg          *core.NodeConfig
	timerWaiting bool
	timer        *time.Timer
	delayMutex   sync.Mutex
}

func NewSleepNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &SleepNode{
		StatefulActionNode: core.NewStatefulActionNode(name, cfg),
		cfg:                cfg,
	}
	n.StatefulActionNode.IStatefulActionNode = n
	return n
}

func (n *SleepNode) OnStart() core.NodeStatus {
	var m int
	v, err := n.GetInput("msec", &m)
	if err != nil {
		panic("Missing parameter [msec] in SleepNode")
	}
	msec := time.Duration(v.(int))
	if msec <= 0 {
		return core.NodeStatus_SUCCESS
	}

	n.SetStatus(core.NodeStatus_RUNNING)

	n.timerWaiting = true
	n.timer = time.AfterFunc(msec, func() {
		n.delayMutex.Lock()
		n.EmitWakeUpSignal()
		n.timerWaiting = false
		n.timer = nil
		n.delayMutex.Unlock()
	})

	return core.NodeStatus_RUNNING
}

func (n *SleepNode) OnRunning() core.NodeStatus {
	if n.timerWaiting {
		return core.NodeStatus_RUNNING
	}
	return core.NodeStatus_SUCCESS
}

func (n *SleepNode) OnHalt() {
	n.delayMutex.Lock()
	if n.timer != nil {
		n.timer.Stop()
	}
	n.timer = nil
	n.timerWaiting = false
	n.delayMutex.Unlock()
}

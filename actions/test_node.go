package actions

import (
	"github.com/gorustyt/go-behavior/core"
	"sync/atomic"
	"time"
)

type TestNode struct {
	*core.StatefulActionNode
	TestConfig *core.TestNodeConfig
	_completed atomic.Bool
	_executor  core.ScriptFunction
	timer      *time.Timer
}

func NewTestNode(name string, cfg *core.NodeConfig, args ...interface{}) *TestNode {
	n := &TestNode{StatefulActionNode: core.NewStatefulActionNode(name, cfg)}
	n.SetRegistrationID("TestNode")
	n.StatefulActionNode.IStatefulActionNode = n
	return n
}
func (t *TestNode) setConfig(config *core.TestNodeConfig) {
	if config.ReturnStatus == core.NodeStatus_IDLE {
		panic("TestNode can not return IDLE")
	}
	t.TestConfig = config

	if t.TestConfig.PostScript != "" {
		executor, err := core.ParseScript(t.TestConfig.PostScript)
		if err != nil {
			panic(err)
		}
		t._executor = executor
	}
}
func (t *TestNode) OnStart() core.NodeStatus {
	if t.TestConfig.PreFunc != nil {
		t.TestConfig.PreFunc()
	}

	if t.TestConfig.AsyncDelay <= 0 {
		return t.OnCompleted()
	}
	// convert this in an asynchronous operation. Use another thread to count
	// a certain amount of time.
	t._completed.Store(false)
	t.timer = time.AfterFunc(t.TestConfig.AsyncDelay, func() {
		t._completed.Store(true)
		t.EmitWakeUpSignal()
	})

	return core.NodeStatus_RUNNING
}

func (t *TestNode) OnRunning() core.NodeStatus {
	if t._completed.Load() {
		return t.OnCompleted()
	}
	return core.NodeStatus_RUNNING
}

func (t *TestNode) OnHalted() {
	t.timer.Stop()
}
func (t *TestNode) OnCompleted() core.NodeStatus {
	if t._executor != nil {
		t._executor(t.Config().Blackboard, t.Config().Enums)
	}
	if t.TestConfig.PostFunc != nil {
		t.TestConfig.PostFunc()
	}
	return t.TestConfig.ReturnStatus
}

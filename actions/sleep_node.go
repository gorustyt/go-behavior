package actions

import "github.com/gorustyt/go-behavior/core"

type SleepNode struct {
	*core.StatefulActionNode
	cfg          *core.NodeConfig
	timerWaiting bool
}

func NewSleepNode(cfg *core.NodeConfig) *SleepNode {
	return &SleepNode{
		cfg: cfg,
	}
}

func (n *SleepNode) OnStart() {

}

func (n *SleepNode) OnRunning() {

}

func (n *SleepNode) OnHalt() {
	n.timerWaiting = false
}

package decorators

import "github.com/gorustyt/go-behavior/core"

func init() {
	core.SetPorts(&RunOnceNode{}, core.InputPortWithDefaultValue(
		"then_skip", true,
		"If true, skip after the first execution, otherwise return the same NodeStatus returned once bu the child."))
}

type RunOnceNode struct {
	*core.DecoratorNode
	alreadyTicked  bool
	returnedStatus core.NodeStatus
}

func NewRunOnceNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &RunOnceNode{
		DecoratorNode:  core.NewDecoratorNode(name, cfg),
		returnedStatus: core.NodeStatus_IDLE,
	}
	n.SetRegistrationID("RunOnce")
	return n
}

// ------------ implementation ----------------------------
func (n *RunOnceNode) Tick() core.NodeStatus {
	skip := true

	if v, err := n.GetInput("then_skip", &skip); err != nil {
		panic(err)
	} else {
		skip = v.(bool)
	}

	if n.alreadyTicked {
		if skip {
			return core.NodeStatus_SKIPPED
		}
		return n.returnedStatus
	}

	n.SetStatus(core.NodeStatus_RUNNING)
	status := n.Child().ExecuteTick()

	if core.IsStatusCompleted(status) {
		n.alreadyTicked = true
		n.returnedStatus = status
		n.ResetChild()
	}
	return status
}

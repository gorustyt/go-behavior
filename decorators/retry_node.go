package decorators

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

func init() {
	core.SetPorts(core.InputPortWithDefaultValue(NUM_ATTEMPTS, 1, "Execute again a failing child up to N times. Use -1 to create an infinite loop."))
}

const (
	NUM_ATTEMPTS = "num_attempts"
)

type RetryNode struct {
	*core.DecoratorNode
	maxAttempts            int
	tryCount               int
	allSkipped             bool
	readParameterFromPorts bool
}

func NewRetryNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &RetryNode{
		allSkipped:             true,
		readParameterFromPorts: true,
		DecoratorNode:          core.NewDecoratorNode(name, cfg),
	}
	n.SetRegistrationID("RetryUntilSuccessful")
	if len(args) > 0 {
		n.readParameterFromPorts = false
		n.maxAttempts = args[0].(int)
	}
	return n
}

func (n *RetryNode) Halt() {
	n.tryCount = 0
	n.Halt()
}

func (n *RetryNode) Tick() core.NodeStatus {

	if n.readParameterFromPorts {
		if v, err := n.GetInput(NUM_ATTEMPTS, &n.maxAttempts); err != nil {
			panic(fmt.Sprintf("Missing parameter [%v] in RetryNode", NUM_ATTEMPTS))
		} else {
			n.maxAttempts = v.(int)
		}
	}

	do_loop := n.tryCount < n.maxAttempts || n.maxAttempts == -1

	if n.Status() == core.NodeStatus_IDLE {
		n.allSkipped = true
	}
	n.SetStatus(core.NodeStatus_RUNNING)

	for do_loop {
		prevStatus := n.Child().Status()
		childStatus := n.Child().ExecuteTick()
		// switch to RUNNING state as soon as you find an active child
		n.allSkipped = n.allSkipped && (childStatus == core.NodeStatus_SKIPPED)
		switch childStatus {
		case core.NodeStatus_SUCCESS:
			n.tryCount = 0
			n.ResetChild()
			return core.NodeStatus_SUCCESS
		case core.NodeStatus_FAILURE:
			n.tryCount++
			do_loop = n.tryCount < n.maxAttempts || n.maxAttempts == -1

			n.ResetChild()

			// Return the execution flow if the child is async,
			// to make this interruptable.
			if n.RequiresWakeUp() && prevStatus == core.NodeStatus_IDLE && do_loop {
				n.EmitWakeUpSignal()
				return core.NodeStatus_RUNNING
			}
			break
		case core.NodeStatus_RUNNING:
			return core.NodeStatus_RUNNING
		case core.NodeStatus_SKIPPED:
			// to allow it to be skipped again, we must reset the node
			n.ResetChild()
			// the child has been skipped. Slip this too
			return core.NodeStatus_SKIPPED
		case core.NodeStatus_IDLE:
			panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))
		}
	}

	n.tryCount = 0
	if n.allSkipped {
		return core.NodeStatus_SKIPPED
	}
	return core.NodeStatus_FAILURE
}

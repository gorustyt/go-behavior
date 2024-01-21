package decorators

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

const (
	NUM_CYCLES = "num_cycles"
)

func init() {
	core.SetPorts(core.InputPortWithDefaultValue(NUM_CYCLES, 1, "Repeat a successful child up to N times. Use -1 to create an infinite loop."))
}

type RepeatNode struct {
	*core.DecoratorNode
	num_cycles_                int
	repeat_count_              int
	all_skipped_               bool
	read_parameter_from_ports_ bool
}

func NewRepeatNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &RepeatNode{
		all_skipped_:               true,
		read_parameter_from_ports_: true,
		DecoratorNode:              core.NewDecoratorNode(name, cfg),
	}
	n.SetRegistrationID("Repeat")
	if len(args) > 0 {
		n.read_parameter_from_ports_ = false
		n.num_cycles_ = args[0].(int)
	}
	return n
}
func (n *RepeatNode) Tick() core.NodeStatus {

	if n.read_parameter_from_ports_ {
		v, err := n.GetInput(NUM_CYCLES, &n.num_cycles_)
		if err != nil {
			panic(fmt.Sprintf("Missing parameter [%v] in RepeatNode", NUM_CYCLES))
		}
		n.num_cycles_ = v.(int)
	}

	do_loop := n.repeat_count_ < n.num_cycles_ || n.num_cycles_ == -1
	if n.Status() == core.NodeStatus_IDLE {
		n.all_skipped_ = true
	}
	n.SetStatus(core.NodeStatus_RUNNING)

	for do_loop {
		prevStatus := n.Child().Status()
		childStatus := n.Child().ExecuteTick()

		// switch to RUNNING state as soon as you find an active child
		n.all_skipped_ = n.all_skipped_ && (childStatus == core.NodeStatus_SKIPPED)

		switch childStatus {
		case core.NodeStatus_SUCCESS:
			n.repeat_count_++
			do_loop = n.repeat_count_ < n.num_cycles_ || n.num_cycles_ == -1

			n.ResetChild()

			// Return the execution flow if the child is async,
			// to make this interruptable.
			if n.RequiresWakeUp() && prevStatus == core.NodeStatus_IDLE && do_loop {
				n.EmitWakeUpSignal()
				return core.NodeStatus_RUNNING
			}
		case core.NodeStatus_FAILURE:

			n.repeat_count_ = 0
			n.ResetChild()
			return core.NodeStatus_FAILURE

		case core.NodeStatus_RUNNING:
			return core.NodeStatus_RUNNING

		case core.NodeStatus_SKIPPED:
			// to allow it to be skipped again, we must reset the node
			n.ResetChild()
			// the child has been skipped. Skip the decorator too.
			// Don't reset the counter, though !
			return core.NodeStatus_SKIPPED
		case core.NodeStatus_IDLE:
			panic(fmt.Sprintf("[%v]: A children should not return IDLE", n.Name()))
		}
	}

	n.repeat_count_ = 0
	if n.all_skipped_ {
		return core.NodeStatus_SKIPPED
	}
	return core.NodeStatus_SUCCESS
}

func (n *RepeatNode) Halt() {
	n.repeat_count_ = 0
	n.DecoratorNode.Halt()
}

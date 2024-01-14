package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

type SwitchNode struct {
	*core.ControlNode
	running_child_ int
}

func (n *SwitchNode) Halt() {
	n.running_child_ = -1
	n.ControlNode.Halt()
}

func (n *SwitchNode) ProvidedPorts() map[string]*core.PortInfo {
	res := map[string]*core.PortInfo{}
	v := core.InputPortWithDefaultValue("variable", "")
	res[v.Name] = v
	for i := 0; i < NUM_CASES; i++ {
		case_str := fmt.Sprintf("case_%d", i+1)
		res[v.Name] = core.InputPortWithDefaultValue(case_str, "")
	}
	return res
}
func (n *SwitchNode) Tick() core.NodeStatus {
	if len(n.Children) != NUM_CASES+1 {
		panic("Wrong number of children in SwitchNode; must be (num_cases + default)")
	}

	match_index := int(NUM_CASES) // default index;
	variable, err := n.GetInput("variable")
	if err == nil {
		// check each case until you find a match
		for index := 0; index < int(NUM_CASES); index++ {
			case_key := fmt.Sprintf("case_%d", int(index+1))
			value, err := n.GetInput(case_key)

			if err == nil && variable == value {
				match_index = index
				break
			}
		}
	}

	// if another one was running earlier, halt it
	if n.running_child_ != -1 && n.running_child_ != match_index {
		n.HaltChild(n.running_child_)
	}

	selected_child := n.Children[match_index]
	ret := selected_child.ExecuteTick()
	if ret == core.NodeStatus_SKIPPED {
		// if the matching child is SKIPPED, should I jump to default or
		// be SKIPPED myself? Going with the former, for the time being.
		n.running_child_ = -1
		return core.NodeStatus_SKIPPED
	} else if ret == core.NodeStatus_RUNNING {
		n.running_child_ = match_index
	} else {
		n.ResetChildren()
		n.running_child_ = -1
	}
}

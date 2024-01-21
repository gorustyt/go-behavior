package controls

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
	"strconv"
)

type SwitchNode struct {
	*core.ControlNode
	runningChild int
	numCases     int
}

func NewSwitchNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	var numCases int
	if len(args) > 0 {
		numCases, _ = strconv.Atoi(args[0].(string))
	}
	n := &SwitchNode{
		numCases:     numCases,
		ControlNode:  core.NewControlNode(name, cfg),
		runningChild: -1,
	}
	core.SetPorts(&SwitchNode{}, core.InputPortWithDefaultValue("variable", ""))
	for i := 0; i < numCases; i++ {
		core.SetPorts(&SwitchNode{}, core.InputPortWithDefaultValue(fmt.Sprintf("case_%d", i+1), ""))
	}
	n.SetRegistrationID("Switch")
	return n
}
func (n *SwitchNode) Halt() {
	n.runningChild = -1
	n.ControlNode.Halt()
}

func (n *SwitchNode) ProvidedPorts() map[string]*core.PortInfo {
	res := map[string]*core.PortInfo{}
	v := core.InputPortWithDefaultValue("variable", "")
	res[v.Name] = v
	for i := 0; i < n.numCases; i++ {
		caseStr := fmt.Sprintf("case_%d", i+1)
		res[v.Name] = core.InputPortWithDefaultValue(caseStr, "")
	}
	return res
}
func (n *SwitchNode) Tick() core.NodeStatus {
	if len(n.Children) != n.numCases+1 {
		panic("Wrong number of children in SwitchNode; must be (num_cases + default)")
	}

	matchIndex := n.numCases // default index;
	var (
		variable int
		value    int
	)
	v, err := n.GetInput("variable", &variable)
	if err == nil {
		variable = v.(int)
		// check each case until you find a match
		for index := 0; index < n.numCases; index++ {
			caseKey := fmt.Sprintf("case_%d", int(index+1))
			v1, err := n.GetInput(caseKey, &value)
			if err != nil {
				panic(err)
			}
			value = v1.(int)
			if err == nil && variable == value {

				matchIndex = index
				break
			}
		}
	}

	// if another one was running earlier, halt it
	if n.runningChild != -1 && n.runningChild != matchIndex {
		n.HaltChild(n.runningChild)
	}

	selectedChild := n.Children[matchIndex]
	ret := selectedChild.ExecuteTick()
	if ret == core.NodeStatus_SKIPPED {
		// if the matching child is SKIPPED, should I jump to default or
		// be SKIPPED myself? Going with the former, for the time being.
		n.runningChild = -1
		return core.NodeStatus_SKIPPED
	} else if ret == core.NodeStatus_RUNNING {
		n.runningChild = matchIndex
	} else {
		n.ResetChildren()
		n.runningChild = -1
	}
	return ret
}

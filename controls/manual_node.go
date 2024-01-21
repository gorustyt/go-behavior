package controls

import "github.com/gorustyt/go-behavior/core"

const (
	REPEAT_LAST_SELECTION = "repeat_last_selection"
)

func init() {
	core.SetPorts(&ManualSelectorNode{}, core.InputPortWithDefaultValue(REPEAT_LAST_SELECTION, false,
		"If true, execute again the same child that was selected the last time"))
}

type NumericalStatus int

const (
	NUM_SUCCESS = 253
	NUM_FAILURE = 254
	NUM_RUNNING = 255
)

type ManualSelectorNode struct {
	running_child_idx_       int
	previously_executed_idx_ int
}

func (n *ManualSelectorNode) Halt() {

}

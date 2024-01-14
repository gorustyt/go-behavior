package controls

type NumericalStatus int

const (
	NUM_SUCCESS = 253
	NUM_FAILURE = 254
	NUM_RUNNING = 255
)

type ManualSelectorNode struct {
}

func (n *ManualSelectorNode) Halt() {

}

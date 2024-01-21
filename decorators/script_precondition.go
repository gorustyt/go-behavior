package decorators

import "github.com/gorustyt/go-behavior/core"

func init() {
	core.SetPorts(&PreconditionNode{}, core.InputPortWithDefaultValue("if", ""))
	core.SetPorts(&PreconditionNode{}, core.InputPortWithDefaultValue("else", core.NodeStatus_FAILURE, "Return status if condition is false"))
}

type PreconditionNode struct {
	*core.DecoratorNode
	_script   string
	_executor core.ScriptFunction
}

func NewPreconditionNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &PreconditionNode{
		DecoratorNode: core.NewDecoratorNode(name, cfg),
	}
	n.loadExecutor()
	return n
}
func (n *PreconditionNode) Tick() core.NodeStatus {
	n.loadExecutor()

	var else_return core.NodeStatus

	if v, err := n.GetInput("else", 0); err != nil {
		panic("Missing parameter [else] in Precondition")
	} else {
		else_return = core.NodeStatus(v.(int))
	}
	if n._executor(n.Config().Blackboard, n.Config().Enums) {
		child_status := n.Child().ExecuteTick()
		if core.IsStatusCompleted(child_status) {
			n.ResetChild()
		}
		return child_status
	} else {
		return else_return
	}
}

func (n *PreconditionNode) loadExecutor() {
	var script string
	if v, err := n.GetInput("if", script); err != nil {
		panic("Missing parameter [if] in Precondition")
	} else {
		script = v.(string)
	}
	if script == n._script {
		return
	}

	if n._executor == nil {
		panic("not impl")
	} else {

		n._script = script
	}

}

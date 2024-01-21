package actions

import "github.com/gorustyt/go-behavior/core"

func init() {
	core.SetPorts(&ScriptNode{}, core.InputPort("code", "Piece of code that can be parsed"))
}

type ScriptNode struct {
	*core.SyncActionNode
	script   string
	executor core.ScriptFunction
}

func NewScriptNode(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &ScriptNode{SyncActionNode: core.NewSyncActionNode(name, cfg)}
	n.SetRegistrationID("ScriptNode")
	n.loadExecutor()
	return n
}
func (n *ScriptNode) Tick() core.NodeStatus {
	n.loadExecutor()

	result := n.executor(n.Config().Blackboard, n.Config().Enums)
	if result {
		return core.NodeStatus_SUCCESS
	}
	return core.NodeStatus_FAILURE
}

func (n *ScriptNode) loadExecutor() {
	var script string
	v, err := n.GetInput("code", script)
	if err != nil {
		panic("Missing port [code] in Script")
	}
	script = v.(string)
	if script == n.script {
		return
	}
	executor, err := core.ParseScript(script)
	if err != nil {
		panic(err)
	} else {
		n.executor = executor
		n.script = script
	}
}

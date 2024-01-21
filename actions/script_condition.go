package actions

import "github.com/gorustyt/go-behavior/core"

func init() {

}

type ScriptCondition struct {
	*core.ConditionNode
	script   string
	executor core.ScriptFunction
}

func NewScriptCondition(name string, cfg *core.NodeConfig, args ...interface{}) core.ITreeNode {
	n := &ScriptCondition{ConditionNode: core.NewConditionNode(name, cfg)}
	n.SetRegistrationID("ScriptCondition")
	n.loadExecutor()
	return n
}

func (n *ScriptCondition) Tick() core.NodeStatus {
	n.loadExecutor()

	result := n.executor(n.Config().Blackboard, n.Config().Enums)
	if result {
		return core.NodeStatus_SUCCESS
	}
	return core.NodeStatus_FAILURE
}

func (n *ScriptCondition) loadExecutor() {
	var script string
	v, err := n.GetInput("code", script)
	if err != nil {
		panic("Missing port [code] in ScriptCondition")
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

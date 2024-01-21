package core

import (
	"errors"
	"fmt"
	"github.com/gorustyt/go-behavior/actions"
	"github.com/gorustyt/go-behavior/controls"
	"github.com/gorustyt/go-behavior/decorators"
	"log"
	"reflect"
	"strings"
	"time"
)

type TickOption int

const (
	EXACTLY_ONCE = iota
	ONCE_UNLESS_WOKEN_UP
	WHILE_RUNNING
)

type NodeBuilderFn func(name string, config *NodeConfig, args ...interface{}) ITreeNode //构造函数

type NodeBuilder struct {
	Cons        NodeBuilderFn
	DefaultArgs []any //默认参数
	*TreeNodeManifest
}

type TestNodeConfig struct {
	Id string
	/// status to return when the action is completed
	ReturnStatus NodeStatus
	/// script to execute when actions is completed
	PostScript string
	//毫秒
	AsyncDelay time.Duration
	PreFunc    func()
	PostFunc   func()
}

func NewTestNodeConfig() *TestNodeConfig {
	return &TestNodeConfig{
		ReturnStatus: NodeStatus_SUCCESS,
	}
}

type BehaviorTreeFactory struct {
	Builders                map[string]*NodeBuilder
	behaviorTreeDefinitions map[string]any
	scriptingEnums          map[string]int
	substitutionRules       map[string]*TestNodeConfig
	parser                  Parser
}

func NewBehaviorTreeFactory() *BehaviorTreeFactory {
	return &BehaviorTreeFactory{
		Builders:                map[string]*NodeBuilder{},
		behaviorTreeDefinitions: map[string]any{},
		scriptingEnums:          map[string]int{},
		substitutionRules:       map[string]*TestNodeConfig{},
		parser:                  NewXmlParser(),
	}
}

func (f *BehaviorTreeFactory) TickOnce() NodeStatus {
	return f.tickRoot(ONCE_UNLESS_WOKEN_UP, 0)
}

func (f *BehaviorTreeFactory) TickExactlyOnce() NodeStatus {
	return f.tickRoot(EXACTLY_ONCE, 0)
}

func (f *BehaviorTreeFactory) TickWhileRunning(sleepTime time.Duration) NodeStatus {
	return f.tickRoot(WHILE_RUNNING, sleepTime)
}

func (f *BehaviorTreeFactory) tickRoot(opt TickOption, sleepTime time.Duration) NodeStatus {
	return f.tickRoot(opt, sleepTime)
}

func (f *BehaviorTreeFactory) RegisterNodeType(id string, cons NodeBuilderFn, args ...any) {
	b := &NodeBuilder{
		TreeNodeManifest: NewTreeNodeManifest(cons("", nil, args...)),
		DefaultArgs:      args,
		Cons:             cons,
	}
	b.RegistrationID = id
	f.Builders[id] = b
}

func (f *BehaviorTreeFactory) InstantiateTreeNode(name, ID string, config *NodeConfig) (node ITreeNode, err error) {
	idNotFound := func() error {
		log.Printf("%v  not included in this list:", ID)
		for k, _ := range f.Builders {
			log.Println(k)
		}
		return fmt.Errorf("BehaviorTreeFactory: ID [%v] not registered", ID)
	}

	b, ok := f.Builders[ID]
	if !ok || b.TreeNodeManifest == nil {
		return node, idNotFound()
	}
	substituted := false
	for filter, rule := range f.substitutionRules {
		if filter == name || filter == ID || strings.Contains(config.Path, filter) { //TODo 验证原来的方法
			// first case: the rule is simply a string with the name of the
			// node to create instead
			if substitutedId := rule.Id; substitutedId != "" {
				builder, ok := f.Builders[substitutedId]
				if ok && builder.Cons != nil {
					node = builder.Cons(name, config)
				} else {
					return node, errors.New("substituted Node ID not found")
				}
				substituted = true
				break
			} else {
				// second case, the varian is a TestNodeConfig
				testNode := actions.NewTestNode(name, config)
				testNode.TestConfig = rule
				node = testNode
				substituted = true
				break
			}
		}
	}

	// No substitution rule applied: default behavior
	if !substituted {
		builder, ok := f.Builders[ID]
		if !ok || builder.Cons == nil {
			return node, idNotFound()
		}
		node = builder.Cons(name, config)
	}

	node.SetRegistrationID(ID)
	node.Config().Enums = f.scriptingEnums
	for condId, script := range config.PreConditions {
		node.PreConditionsScripts()[condId] = func(args ...interface{}) bool {
			log.Println("not impl", script)
			return true
		}
	}

	for condId, script := range config.PostConditions {
		node.PostConditionsScripts()[condId] = func(args ...interface{}) bool {
			log.Println("not impl", script)
			return true
		}
	}
	return node, nil
}
func (f *BehaviorTreeFactory) RegisterScriptingEnum(name string, value int) {
	f.scriptingEnums[name] = value
}

func (f *BehaviorTreeFactory) Init() {
	f.RegisterNodeType("Fallback", controls.NewFallbackNode)
	f.RegisterNodeType("AsyncFallback", controls.NewFallbackNode, true)
	f.RegisterNodeType("Sequence", controls.NewSequenceNode)
	f.RegisterNodeType("AsyncSequence", controls.NewSequenceNode, true)
	f.RegisterNodeType("SequenceWithMemory", controls.NewSequenceWithMemory)
	f.RegisterNodeType("SequenceStar", controls.NewSequenceWithMemory)
	f.RegisterNodeType("Parallel", controls.NewParallelNode)
	f.RegisterNodeType("ParallelAll", controls.NewParallelAllNode)
	f.RegisterNodeType("ReactiveSequence", controls.NewReactiveSequence)
	f.RegisterNodeType("ReactiveFallback", controls.NewReactiveFallback)
	f.RegisterNodeType("IfThenElse", controls.NewIfThenElseNode)
	f.RegisterNodeType("WhileDoElse", controls.NewWhileDoElseNode)
	f.RegisterNodeType("Inverter", decorators.NewInverterNode)
	f.RegisterNodeType("RetryUntilSuccessful", decorators.NewRetryNode)
	f.RegisterNodeType("KeepRunningUntilFailure", decorators.NewKeepRunningUntilFailureNode)
	f.RegisterNodeType("Repeat", decorators.NewRepeatNode)
	f.RegisterNodeType("Timeout", decorators.NewTimeoutNode)
	f.RegisterNodeType("Delay", decorators.NewDelayNode)
	f.RegisterNodeType("RunOnce", decorators.NewRunOnceNode)
	f.RegisterNodeType("ForceSuccess", decorators.NewForceSuccessNode)
	f.RegisterNodeType("ForceFailure", decorators.NewForceFailureNode)
	f.RegisterNodeType("AlwaysSuccess", actions.NewAlwaysSuccessNode)
	f.RegisterNodeType("AlwaysFailure", actions.NewAlwaysFailureNode)
	f.RegisterNodeType("Script", actions.NewScriptNode)
	f.RegisterNodeType("ScriptCondition", actions.NewScriptCondition)
	f.RegisterNodeType("SetBlackboard", actions.NewSetBlackboardNode)
	f.RegisterNodeType("Sleep", actions.NewSleepNode)
	f.RegisterNodeType("UnsetBlackboard", actions.NewUnsetBlackboardNode)
	f.RegisterNodeType("SubTree", decorators.NewSubTreeNode)
	f.RegisterNodeType("Precondition", decorators.NewPreconditionNode)
	f.RegisterNodeType("Switch2", controls.NewSwitchNode, 2)
	f.RegisterNodeType("Switch3", controls.NewSwitchNode, 3)
	f.RegisterNodeType("Switch4", controls.NewSwitchNode, 4)
	f.RegisterNodeType("Switch5", controls.NewSwitchNode, 5)
	f.RegisterNodeType("Switch6", controls.NewSwitchNode, 6)
	f.RegisterNodeType("LoopInt", decorators.NewLoopNode, reflect.Int)
	f.RegisterNodeType("LoopBool", decorators.NewLoopNode, reflect.Bool)
	f.RegisterNodeType("LoopDouble", decorators.NewLoopNode, reflect.Float64)
	f.RegisterNodeType("LoopString", decorators.NewLoopNode, reflect.String)
}

func (f *BehaviorTreeFactory) CreateTreeFromText(text string) (*Tree, error) {
	res := f.parser.RegisteredBehaviorTrees()
	if len(res) != 0 {
		fmt.Println("WARNING: You executed BehaviorTreeFactory::createTreeFromText ",
			"after registerBehaviorTreeFrom[File/Text].\n",
			"This is NOT, probably, what you want to do.\n",
			"You should probably use BehaviorTreeFactory::createTree, instead")
	}
	parser := NewXmlParser()
	err := parser.LoadFromText(text)
	if err != nil {
		return nil, err
	}
	blackboard := NewBlackboard(nil)
	tree, err := parser.InstantiateTree(blackboard, "")
	if err != nil {
		return nil, err
	}
	for k, v := range f.Builders {
		tree.manifests[k] = v.TreeNodeManifest
	}
	return tree, nil
}

func (f *BehaviorTreeFactory) CreateTreeFromFile(file_path string) (*Tree, error) {
	res := f.parser.RegisteredBehaviorTrees()
	if len(res) != 0 {
		log.Println(
			"WARNING: You executed BehaviorTreeFactory::createTreeFromFile ",
			"after registerBehaviorTreeFrom[File/Text].\n",
			"This is NOT, probably, what you want to do.\n",
			"You should probably use BehaviorTreeFactory::createTree, instead",
		)
	}

	parser := NewXmlParser()
	err := parser.LoadFromFile(file_path)
	if err != nil {
		return nil, err
	}
	blackboard := NewBlackboard(nil)
	tree, err := parser.InstantiateTree(blackboard, "")
	if err != nil {
		return nil, err
	}
	for k, v := range f.Builders {
		tree.manifests[k] = v.TreeNodeManifest
	}
	return tree, nil
}
func (f *BehaviorTreeFactory) CreateTree(treeName string) (*Tree, error) {
	blackboard := NewBlackboard(nil)
	tree, err := f.parser.InstantiateTree(blackboard, treeName)
	if err != nil {
		return nil, err
	}
	for k, v := range f.Builders {
		tree.manifests[k] = v.TreeNodeManifest
	}
	return tree, nil
}

func (f *BehaviorTreeFactory) RegisterBehaviorTreeFromFile(filename string) error {
	return f.parser.LoadFromFile(filename)
}

func (f *BehaviorTreeFactory) RegisterBehaviorTreeFromText(xml_text string) error {
	return f.parser.LoadFromText(xml_text)
}

func (f *BehaviorTreeFactory) RegisteredBehaviorTrees() []string {
	return f.parser.RegisteredBehaviorTrees()
}

func (f *BehaviorTreeFactory) RegisterSimpleCondition(
	ID string, tickFunctor TickFunctor,
	ports map[string]*PortInfo) {
	f.RegisterNodeType(ID, NewSimpleConditionNode, tickFunctor)
}

func (f *BehaviorTreeFactory) RegisterSimpleAction(ID string, tickFunctor TickFunctor,
	ports map[string]*PortInfo) {
	f.RegisterNodeType(ID, NewSimpleActionNode, tickFunctor)
}

func (f *BehaviorTreeFactory) RegisterSimpleDecorator(
	ID string, tickFunctor TickFunctor,
	ports map[string]*PortInfo) {
	f.RegisterNodeType(ID, NewSimpleDecoratorNode, tickFunctor)
}

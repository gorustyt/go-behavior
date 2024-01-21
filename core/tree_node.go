package core

import (
	"errors"
	"fmt"
	"sync"
)

type NodeType int

func (p NodeType) String() string {
	switch p {
	case NodeType_ACTION:
		return "Action"
	case NodeType_CONDITION:
		return "Condition"
	case NodeType_DECORATOR:
		return "Decorator"
	case NodeType_CONTROL:
		return "Control"
	case NodeType_SUBTREE:
		return "SubTree"
	default:
		return "Undefined"
	}
}

const (
	NodeType_UNDEFINED NodeType = iota
	NodeType_ACTION
	NodeType_CONDITION
	NodeType_SUBTREE
	NodeType_DECORATOR
	NodeType_CONTROL
)

type NodeStatus int

func (n NodeStatus) String() string {
	switch n {
	case NodeStatus_SUCCESS:
		return "SUCCESS"
	case NodeStatus_FAILURE:
		return "FAILURE"
	case NodeStatus_RUNNING:
		return "RUNNING"
	case NodeStatus_IDLE:
		return "IDLE"
	case NodeStatus_SKIPPED:
		return "SKIPPED"
	}
	return ""
}
func (n NodeStatus) StringColor(colored bool) string {
	if !colored {
		return n.String()
	} else {
		switch n {
		case NodeStatus_SUCCESS:
			return "\x1b[32mSUCCESS\x1b[0m" // RED
		case NodeStatus_FAILURE:
			return "\x1b[31mFAILURE\x1b[0m" // GREEN
		case NodeStatus_RUNNING:
			return "\x1b[33mRUNNING\x1b[0m" // YELLOW
		case NodeStatus_SKIPPED:
			return "\x1b[34mSKIPPED\x1b[0m" // BLUE
		case NodeStatus_IDLE:
			return "\x1b[36mIDLE\x1b[0m" // CYAN
		}
	}
	return "Undefined"
}

const (
	NodeStatus_IDLE NodeStatus = iota
	NodeStatus_RUNNING
	NodeStatus_SUCCESS
	NodeStatus_FAILURE
	NodeStatus_SKIPPED
)

type PreCond int

func (p PreCond) String() string {
	switch p {
	case PreCond_SUCCESS_IF:
		return "_successIf"
	case PreCond_FAILURE_IF:
		return "_failureIf"
	case PreCond_SKIP_IF:
		return "_skipIf"
	case PreCond_WHILE_TRUE:
		return "_while"
	default:
		return "Undefined"
	}
}

const (
	// order of the enums also tell us the execution order
	PreCond_FAILURE_IF PreCond = iota
	PreCond_SUCCESS_IF
	PreCond_SKIP_IF
	PreCond_WHILE_TRUE
	PreCond_COUNT_
)

type PostCond int

func (p PostCond) String() string {
	switch p {
	case PostCond_ON_SUCCESS:
		return "_onSuccess"
	case PostCond_ON_FAILURE:
		return "_onFailure"
	case PostCond_ALWAYS:
		return "_post"
	case PostCond_ON_HALTED:
		return "_onHalted"
	default:
		return "Undefined"
	}
}

const (
	// order of the enums also tell us the execution order
	PostCond_ON_HALTED PostCond = iota
	PostCond_ON_FAILURE
	PostCond_ON_SUCCESS
	PostCond_ALWAYS
	PostCond_COUNT_
)

func (n NodeStatus) IsActive() bool {
	return n != NodeStatus_IDLE && n != NodeStatus_SKIPPED
}

func (n NodeStatus) IsCompleted() bool {
	return n == NodeStatus_SUCCESS || n == NodeStatus_FAILURE
}

type PortDirection int

func (p PortDirection) String() string {
	switch p {
	case PortDirection_INPUT:
		return "Input"
	case PortDirection_OUTPUT:
		return "Output"
	case PortDirection_INOUT:
		return "InOut"
	}
	return "InOut"
}

const (
	PortDirection_INPUT PortDirection = iota
	PortDirection_OUTPUT
	PortDirection_INOUT
)

type PortInfo struct {
	direction       PortDirection
	description     string
	defaultValue    any
	defaultValueStr string
	Name            string
}

func (p *PortInfo) Description() string {
	return p.description
}
func (p *PortInfo) Direction() PortDirection {
	return p.direction
}
func (p *PortInfo) DefaultValue() any {
	return p.defaultValue
}

func (p *PortInfo) DefaultValueString() string {
	return p.defaultValueStr
}

func (p *PortInfo) SetDefaultValue(value any) {
	if v, ok := value.(fmt.Stringer); ok {
		p.defaultValueStr = v.String()
	}
	p.defaultValue = value
}

func (p *PortInfo) SetDescription(description string) {
	p.description = description
}

func NewPortInfo(direction PortDirection, name string, desc ...string) *PortInfo {
	p := &PortInfo{direction: direction,
		Name: name}
	if len(desc) > 0 {
		p.description = desc[0]
	}
	return p
}

func InputPort(name string, desc ...string) *PortInfo {
	return NewPortInfo(PortDirection_INPUT, name, desc...)
}

func OutPort(name string, desc ...string) *PortInfo {
	return NewPortInfo(PortDirection_OUTPUT, name, desc...)
}

func BidirectionalPort(name string, desc ...string) *PortInfo {
	return NewPortInfo(PortDirection_INOUT, name, desc...)
}

func InputPortWithDefaultValue(name string, defaultValue any, desc ...string) *PortInfo {
	p := NewPortInfo(PortDirection_INPUT, name, desc...)
	p.SetDefaultValue(defaultValue)
	return p
}

func OutputPortWithDefaultValue(name string, defaultValue any, desc ...string) *PortInfo {
	p := NewPortInfo(PortDirection_OUTPUT, name, desc...)
	p.SetDefaultValue(defaultValue)
	return p
}

func BidirectionalPortWithDefaultValue(name string, defaultValue any, desc ...string) *PortInfo {
	p := NewPortInfo(PortDirection_INOUT, name, desc...)
	p.SetDefaultValue(defaultValue)
	return p
}

type TreeNodeManifest struct {
	Type           NodeType
	RegistrationID string
	Ports          map[string]*PortInfo
	Metadata       []map[string]string
}

func NewTreeNodeManifest(value any) *TreeNodeManifest {
	v := &TreeNodeManifest{
		Ports:    map[string]*PortInfo{},
		Metadata: make([]map[string]string, 0),
		Type:     NodeType_UNDEFINED,
	}
	t, ok := value.(INodeType)
	if ok {
		v.Type = t.NodeType()
	}
	return v
}

var (
	DefaultNodeConfig = &NodeConfig{}
)

type NodeConfig struct {
	Blackboard     *Blackboard
	Enums          map[string]int
	InputPorts     map[string]string
	OutputPorts    map[string]string
	Manifest       *TreeNodeManifest
	Uid            uint16
	Path           string
	PreConditions  map[PreCond]string  //有序
	PostConditions map[PostCond]string //有序

}

type TreeNode struct {
	name                   string
	status                 NodeStatus
	mutex                  *sync.Mutex
	callbackInjectionMutex *sync.Mutex
	config                 *NodeConfig
	substitutionCallback   PreTickCallback
	postConditionCallback  PostTickCallback
	cond                   *sync.Cond
	wake_up                *WakeUpSignal
	registrationID         string
	pre_parsed             []ScriptFunction
	post_parsed            []ScriptFunction
	state_change_signal    *Signal
}

func NewTreeNode(name string, cfg *NodeConfig) *TreeNode {
	mu := &sync.Mutex{}
	return &TreeNode{
		callbackInjectionMutex: &sync.Mutex{},
		mutex:                  mu,
		config:                 cfg,
		name:                   name,
		pre_parsed:             make([]ScriptFunction, PreCond_COUNT_),
		post_parsed:            make([]ScriptFunction, PostCond_COUNT_),
		cond:                   sync.NewCond(mu),
		state_change_signal:    NewSignal(),
	}
}

func (n *TreeNode) SetStatus(status NodeStatus) {
	if status == NodeStatus_IDLE {
		panic(fmt.Sprintf("Node [", "", "]: you are not allowed to set manually the status to IDLE. If you know what you are doing (?) use resetStatus() instead."))
	}
	var prev_status NodeStatus
	n.mutex.Lock()
	prev_status = n.status
	n.mutex.Unlock()
	if prev_status != status {
		n.cond.Broadcast()
		n.state_change_signal.Notify(prev_status, status)
	}
}

func (n *TreeNode) EmitWakeUpSignal() {
	if n.wake_up != nil {
		n.wake_up.EmitSignal()
	}
}

func (n *TreeNode) RequiresWakeUp() bool {
	return n.wake_up != nil
}
func (n *TreeNode) Tick() NodeStatus {
	panic("not ok")
	return NodeStatus_RUNNING
}
func (n *TreeNode) ExecuteTick() NodeStatus {
	new_status := n.status

	// a pre-condition may return the new status.
	// In this case it override the actual tick()
	if precond, err := n.checkPreConditions(); err == nil {
		new_status = precond
	} else {
		// injected pre-callback
		substituted := false
		if !IsStatusCompleted(n.status) {
			var callback PreTickCallback
			n.callbackInjectionMutex.Lock()
			callback = n.substitutionCallback
			n.callbackInjectionMutex.Unlock()
			if callback != nil {
				override_status := callback(n)
				if IsStatusCompleted(override_status) {
					// don't execute the actual tick()
					substituted = true
					new_status = override_status
				}
			}
		}

		// Call the ACTUAL tick
		if !substituted {
			new_status = n.Tick()
		}
	}

	n.checkPostConditions(new_status)

	// injected post callback
	if IsStatusCompleted(new_status) {
		var callback PostTickCallback
		n.callbackInjectionMutex.Lock()
		callback = n.postConditionCallback
		n.callbackInjectionMutex.Unlock()

		if callback != nil {
			override_status := callback(n, new_status)
			if IsStatusCompleted(override_status) {
				new_status = override_status
			}
		}
	}

	// preserve the IDLE state if skipped, but communicate SKIPPED to parent
	if new_status != NodeStatus_SKIPPED {
		n.SetStatus(new_status)
	}
	return new_status
}
func (n *TreeNode) Status() NodeStatus {
	var prevStatus NodeStatus
	n.mutex.Lock()
	prevStatus = n.status
	n.mutex.Unlock()
	return prevStatus
}

func (n *TreeNode) SetWakeUpInstance(instance *WakeUpSignal) {
	n.wake_up = instance
}

// / The method used to interrupt the execution of a RUNNING node.
// / Only Async nodes that may return RUNNING should implement it.
func (n *TreeNode) Halt() {

}
func (n *TreeNode) HaltNode() {
	n.Halt()
	ex := n.post_parsed[PostCond_ON_HALTED]
	if ex != nil {
		ex(n.Config().Blackboard, n.Config().Enums)
	}
}

func (n *TreeNode) PreConditionsScripts() []ScriptFunction {
	return n.pre_parsed
}

func (n *TreeNode) SetRegistrationID(ID string) {
	n.registrationID = ID
}
func (n *TreeNode) PostConditionsScripts() []ScriptFunction {
	return n.post_parsed
}

func (n *TreeNode) Name() string {
	return n.name
}
func (n *TreeNode) UID() uint16 {
	return n.config.Uid
}
func (n *TreeNode) checkPreConditions() (NodeStatus, error) {
	// check the pre-conditions
	for index := 0; index < int(PreCond_COUNT_); index++ {
		parse_executor := n.pre_parsed[index]
		if parse_executor == nil {
			continue
		}
		args := []interface{}{n.Config().Blackboard, n.Config().Enums}
		preID := PreCond(index)
		// Some preconditions are applied only when the node state is IDLE or SKIPPED
		if n.status == NodeStatus_IDLE ||
			n.status == NodeStatus_SKIPPED {
			// what to do if the condition is true
			if parse_executor(args...) {
				if preID == PreCond_FAILURE_IF {
					return NodeStatus_FAILURE, nil
				} else if preID == PreCond_SUCCESS_IF {
					return NodeStatus_SUCCESS, nil
				} else if preID == PreCond_SKIP_IF {
					return NodeStatus_SKIPPED, nil
				}
			} else if preID == PreCond_WHILE_TRUE { // if the conditions is false
				return NodeStatus_SKIPPED, nil
			}
		} else if n.status == NodeStatus_RUNNING && preID == PreCond_WHILE_TRUE {
			// what to do if the condition is false
			if !parse_executor(args...) {
				n.HaltNode()
				return NodeStatus_SKIPPED, nil
			}
		}
	}
	return NodeStatus_FAILURE, errors.New("")
}

func (n *TreeNode) checkPostConditions(status NodeStatus) {
	executeScript := func(cond PostCond) {
		parse_executor := n.post_parsed[cond]
		if parse_executor != nil {
			parse_executor(n.Config().Blackboard, n.Config().Enums)
		}
	}

	if status == NodeStatus_SUCCESS {
		executeScript(PostCond_ON_SUCCESS)
	} else if status == NodeStatus_FAILURE {
		executeScript(PostCond_ON_FAILURE)
	}
	executeScript(PostCond_ALWAYS)
}

func (n *TreeNode) ResetStatus() {
	var prev_status NodeStatus
	n.mutex.Lock()
	prev_status = n.status
	n.mutex.Unlock()
	if prev_status != NodeStatus_IDLE {
		n.cond.Broadcast()
		n.state_change_signal.Notify(prev_status, NodeStatus_IDLE)
	}
}

func (n *TreeNode) ModifyPortsRemapping(newRemapping map[string]string) {
	for k, v := range newRemapping {
		if _, ok := n.config.InputPorts[k]; ok {
			n.config.InputPorts[k] = v
		}
		if _, ok := n.config.OutputPorts[k]; ok {
			n.config.InputPorts[k] = v
		}
	}
}

func (n *TreeNode) SetPreTickFunction(callback func(node *TreeNode) NodeStatus) {
	n.callbackInjectionMutex.Lock()
	n.substitutionCallback = callback
	n.callbackInjectionMutex.Unlock()
}

func (n *TreeNode) SetPostTickFunction(callback func(node *TreeNode, status NodeStatus) NodeStatus) {
	n.callbackInjectionMutex.Lock()
	n.postConditionCallback = callback
	n.callbackInjectionMutex.Unlock()
}

func (n *TreeNode) FullPath() string {
	return n.config.Path
}

func (n *TreeNode) registrationName() string {
	return n.registrationID
}

func (n *TreeNode) Config() *NodeConfig {
	return n.config
}

func (n *TreeNode) SetOutput(key string, value any) {
	if n.config.Blackboard == nil {
		panic("setOutput() failed: trying to access a Blackboard(BB) entry, but BB is invalid")
	}

	remappedKey, ok := n.config.OutputPorts[key]
	if !ok {
		panic(fmt.Sprintf("setOutput() failed:  NodeConfig::output_ports does not contain the key: [%v]", key))
	}
	if remappedKey == "=" {
		n.config.Blackboard.Set(key, value)
		return
	}

	if _, ok = IsBlackboardPointer(remappedKey); !ok {
		panic("setOutput requires a blackboard pointer. Use {}")
	}

	if value == nil {
		panic("setOutput<Any> is not allowed, unless the port was declared using OutputPort<Any>")
	}

	remappedKey = StripBlackboardPointer(remappedKey)
	n.config.Blackboard.Set(remappedKey, value)
}

func (n *TreeNode) GetRawPortValue(key string) string {
	remap, ok := n.config.InputPorts[key]
	if !ok {
		remap, ok = n.config.OutputPorts[key]
		if !ok {
			panic(fmt.Sprintf("[%v] not found", key))
		}
	}
	return remap
}

func IsBlackboardPointer(str string) (res string, ok bool) {
	if len(str) < 3 {
		return str, false
	}
	// strip leading and following spaces
	front_index := 0
	last_index := len(str) - 1
	for str[front_index] == ' ' && front_index <= last_index {
		front_index++
	}
	for str[last_index] == ' ' && front_index <= last_index {
		last_index--
	}
	size := (last_index - front_index) + 1
	valid := size >= 3 && str[front_index] == '{' && str[last_index] == '}'
	if valid {
		res = str[front_index+1:]
	}
	return res, valid
}

func StripBlackboardPointer(str string) string {
	v, _ := IsBlackboardPointer(str)
	return v
}

func AssignDefaultRemapping(config *NodeConfig, value IGetProvidedPorts) {
	for port_name, v := range value.GetProvidedPorts() {
		direction := v.Direction()
		if direction != PortDirection_OUTPUT {
			// PortDirection::{INPUT,INOUT}
			config.InputPorts[port_name] = "="
		}
		if direction != PortDirection_INPUT {
			// PortDirection::{OUTPUT,INOUT}
			config.OutputPorts[port_name] = "="
		}
	}
}
func (n *TreeNode) GetLockedPortContent(key string) func() *Entry {
	if remapped_key, err := GetRemappedKey(key, n.GetRawPortValue(key)); err == nil {
		return n.Config().Blackboard.GetAnyLocked(remapped_key)
	}
	return nil
}

func GetRemappedKey(portName string, remappedPort string) (res string, err error) {
	if remappedPort == "=" {
		return portName, nil
	}
	stripped, ok := IsBlackboardPointer(remappedPort)
	if ok {
		return stripped, nil
	}
	return res, errors.New("not a blackboard pointer")
}
func (n *TreeNode) GetInput(key string, destination any) (res any, err error) {
	ParseString := func(str string) any {
		switch destination.(type) {
		case NodeType, PortDirection:
			it, ok := n.config.Enums[str]
			if ok {
				return it
			}
		}
		return ConvFromString(str, destination)
	}
	portValueStr, ok := n.config.InputPorts[key]
	if !ok {
		return res, fmt.Errorf("getInput() of node `%v` failed because NodeConfig::input_ports does not contain the key: [%v]", n.FullPath(), key)
	}

	// special case. Empty port value, we should use the default value,
	// if available in the model.
	// BUT, it the port type is a string, then an empty string might be
	// a valid value
	if portValueStr == "" && n.config.Manifest != nil {
		portManifest := n.config.Manifest.Ports[key]
		defaultValue := portManifest.DefaultValue()
		_, ok := defaultValue.(string)
		if defaultValue != nil && !ok {
			destination = defaultValue
			return destination, nil
		}
	}

	remappedRes, err := GetRemappedKey(key, portValueStr)
	if err != nil {
		return res, err
	}
	// pure string, not a blackboard key
	if remappedRes != "" {
		destination = ParseString(portValueStr)
		return res, nil
	}
	remappedKey := remappedRes

	if n.config.Blackboard == nil {
		return res, fmt.Errorf("getInput(): trying to access an invalid Blackboard")
	}

	if anyRef := n.config.Blackboard.GetAnyLocked(remappedKey); anyRef != nil {
		val := anyRef()
		v, ok := val.Value.(string)
		if ok {
			destination = ParseString(v)
		} else {
			return destination, nil
		}
	}

	return res, fmt.Errorf("getInput() failed because it was unable to find the key [%v] remapped to [%v]", key,
		remappedKey)

}

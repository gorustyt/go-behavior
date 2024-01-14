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
	registrationID string
	ports          map[string]*PortInfo
	metadata       []map[string]string
}

var (
	DefaultNodeConfig = &NodeConfig{}
)

type NodeConfig struct {
	blackboard     *Blackboard
	enums          map[string]int
	inputPorts     map[string]string
	outputPorts    map[string]string
	manifest       *TreeNodeManifest
	uid            uint16
	path           string
	preConditions  map[PreCond]string  //有序
	postConditions map[PostCond]string //有序
}

type StatusChangeEvent struct {
	Pre NodeStatus
	Now NodeStatus
}
type TreeNode struct {
	name                   string
	status                 NodeStatus
	mutex                  *sync.Mutex
	callbackInjectionMutex *sync.Mutex
	config                 *NodeConfig
	substitutionCallback   func(node *TreeNode) NodeStatus
	postConditionCallback  func(node *TreeNode, status NodeStatus) NodeStatus
	cond                   *sync.Cond
	statusCh               chan *StatusChangeEvent
	registrationID         string
}

func NewTreeNode(name string, cfg *NodeConfig) *TreeNode {
	mu := &sync.Mutex{}
	return &TreeNode{
		callbackInjectionMutex: &sync.Mutex{},
		mutex:                  mu,
		config:                 cfg,
		name:                   name,
		cond:                   sync.NewCond(mu),
		statusCh:               make(chan *StatusChangeEvent, 1),
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
		n.statusCh <- &StatusChangeEvent{Pre: prev_status, Now: status}
	}
}

func (n *TreeNode) ExecuteTick() NodeStatus {

}
func (n *TreeNode) Status() NodeStatus {
	var prev_status NodeStatus
	n.mutex.Lock()
	prev_status = n.status
	n.mutex.Unlock()
	return prev_status
}
func (n *TreeNode) HaltNode() {

}
func (n *TreeNode) Type() {

}
func (n *TreeNode) Name() string {
	return n.name
}

func (n *TreeNode) checkPreConditions() error {
	return errors.New("")
}

func (n *TreeNode) checkPostConditions() {

}

func (n *TreeNode) ResetStatus() {
	var prev_status NodeStatus
	n.mutex.Lock()
	prev_status = n.status
	n.mutex.Unlock()
	if prev_status != NodeStatus_IDLE {
		n.cond.Broadcast()
		n.statusCh <- &StatusChangeEvent{Pre: prev_status, Now: NodeStatus_IDLE}
	}
}

func (n *TreeNode) ModifyPortsRemapping(newRemapping map[string]string) {
	for k, v := range newRemapping {
		if _, ok := n.config.inputPorts[k]; ok {
			n.config.inputPorts[k] = v
		}
		if _, ok := n.config.outputPorts[k]; ok {
			n.config.inputPorts[k] = v
		}
	}
}

func (n *TreeNode) resetStatus() {
	var prevStatus NodeStatus

	n.mutex.Lock()
	prevStatus = n.status
	n.status = NodeStatus_IDLE
	n.mutex.Lock()

	if prevStatus != NodeStatus_IDLE {
		n.cond.Broadcast()
		n.statusCh.notify(
			prevStatus, NodeStatus_IDLE)
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
	return n.config.path
}

func (n *TreeNode) registrationName() string {
	return n.registrationID
}

func (n *TreeNode) Config() *NodeConfig {
	return n.config
}

func (n *TreeNode) SetOutput(key string, value any) (err error) {
	if n.config.blackboard == nil {
		return errors.New("setOutput() failed: trying to access a Blackboard(BB) entry, but BB is invalid")
	}

	remapped_key, ok := n.config.outputPorts[key]
	if !ok {
		return fmt.Errorf("setOutput() failed:  NodeConfig::output_ports does not contain the key: [%v]", key)
	}
	if remapped_key == "=" {
		n.config.blackboard.Set(key, value)
		return nil
	}

	if _, ok = n.IsBlackboardPointer(remapped_key); !ok {
		return errors.New("setOutput requires a blackboard pointer. Use {}")
	}

	if value == nil {
		panic("setOutput<Any> is not allowed, unless the port was declared using OutputPort<Any>")

	}

	remapped_key = n.stripBlackboardPointer(remapped_key)
	n.config.blackboard.Set(remapped_key, value)

	return nil
}

func (n *TreeNode) GetRawPortValue(key string) string {
	remap, ok := n.config.inputPorts[key]
	if !ok {
		remap, ok = n.config.outputPorts[key]
		if !ok {
			panic(fmt.Sprintf("[%v] not found", key))
		}
	}
	return remap
}

func (n *TreeNode) IsBlackboardPointer(str string) (res string, ok bool) {
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

func (n *TreeNode) stripBlackboardPointer(str string) string {
	v, _ := n.IsBlackboardPointer(str)
	return v
}

type IGetProvidedPorts interface {
	GetProvidedPorts() map[string]*PortInfo
}

func AssignDefaultRemapping(config *NodeConfig, value IGetProvidedPorts) {
	for port_name, v := range value.GetProvidedPorts() {
		direction := v.Direction()
		if direction != PortDirection_OUTPUT {
			// PortDirection::{INPUT,INOUT}
			config.inputPorts[port_name] = "="
		}
		if direction != PortDirection_INPUT {
			// PortDirection::{OUTPUT,INOUT}
			config.outputPorts[port_name] = "="
		}
	}
}

func (n *TreeNode) getRemappedKey(port_name string, remapped_port string) (res string, err error) {
	if remapped_port == "=" {
		return port_name, nil
	}
	stripped, ok := n.IsBlackboardPointer(remapped_port)
	if ok {
		return stripped, nil
	}
	return res, errors.New("Not a blackboard pointer")
}
func (n *TreeNode) GetInput(key string, destination any) (res any, err error) {

	ParseString := func(str string) any {
		switch destination.(type) {
		case NodeType, PortDirection:
			it, ok := n.config.enums[str]
			if ok {
				return it
			}
		default:

		}
		return ConvFromString(str, destination)
	}
	remap_it, ok := n.config.inputPorts[key]
	if !ok {
		return res, fmt.Errorf("getInput() of node `%v` failed because NodeConfig::input_ports does not contain the key: [%v]", n.FullPath(), key)
	}

	// special case. Empty port value, we should use the default value,
	// if available in the model.
	// BUT, it the port type is a string, then an empty string might be
	// a valid value
	port_value_str := remap_it
	if port_value_str == "" && n.config.manifest != nil {
		port_manifest := n.config.manifest.ports[key]
		default_value := port_manifest.DefaultValue()
		_, ok := default_value.(string)
		if default_value != nil && !ok {
			destination = default_value
			return destination, nil
		}
	}

	remapped_res, err := n.getRemappedKey(key, port_value_str)
	if err != nil {
		return res, err
	}
	// pure string, not a blackboard key
	if remapped_res != "" {
		destination = ParseString(port_value_str)
		return res, nil
	}
	remapped_key := remapped_res

	if n.config.blackboard == nil {
		return res, fmt.Errorf("getInput(): trying to access an invalid Blackboard")
	}

	if any_ref := n.config.blackboard.getAnyLocked(remapped_key); any_ref != nil {
		val := any_ref()
		v, ok := val.Value.(string)
		if ok {
			destination = ParseString(v)
		} else {
			return val, nil
		}
	}

	return res, fmt.Errorf("getInput() failed because it was unable to find the key [%v] remapped to [%v]", key,
		remapped_key)

}

package core

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func IsAlpha(c uint8) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

type ProtectedQueue struct {
	Items list.List
	Mtx   sync.Mutex
}
type IMetadata interface {
	Metadata() [][]string
}

type IGetProvidedPorts interface {
	GetProvidedPorts() map[string]*PortInfo
}
type PostTickCallback func(node *TreeNode, status NodeStatus) NodeStatus
type PreTickCallback func(node *TreeNode) NodeStatus
type ScriptFunction func(args ...interface{}) bool
type TickFunctor func(node ITreeNode, status ...NodeStatus) NodeStatus
type INodeType interface {
	NodeType() NodeType
}

type ITreeNode interface {
	Tick() NodeStatus
	Config() *NodeConfig
	PreConditionsScripts() []ScriptFunction
	PostConditionsScripts() []ScriptFunction
	SetRegistrationID(ID string)

	NodeType() NodeType
	UID() uint16
	HaltNode()

	SetWakeUpInstance(instance *WakeUpSignal)
	ExecuteTick() NodeStatus
	ResetStatus()
	Status() NodeStatus
}

func IsAllowedPortName(str string) bool {
	if str == "_autoremap" {
		return true
	}
	if str == "" {
		return false
	}

	if !IsAlpha(str[0]) {
		return false
	}
	if str == "name" || str == "ID" {
		return false
	}
	return true
}
func IsStatusCompleted(status NodeStatus) bool {
	return status == NodeStatus_SUCCESS || status == NodeStatus_FAILURE
}
func ConvFromString(str string, value any) error {
	switch v := value.(type) {
	case *int:
		res, err := convertInt64FromString(str)
		if err != nil {
			return err
		}
		*v = int(res)
	case *float64:
		res, err := convertFloat64FromString(str)
		if err != nil {
			return err
		}
		*v = res
	case *[]float64:
		res, err := convertFloat64sFromString(str)
		if err != nil {
			return err
		}
		*v = append(*v, res...)
	case *int32:
		res, err := convertInt64FromString(str)
		if err != nil {
			return err
		}
		*v = int32(res)
	case *int64:
		res, err := convertInt64FromString(str)
		if err != nil {
			return err
		}
		*v = res
	case *[]int64:
		res, err := convertInt64sFromString(str)
		if err != nil {
			return err
		}
		*v = append(*v, res...)
	case *bool:
		*v = convertBoolFromString(str)
	case *NodeStatus:
		*v = convertNodeStatusFromString(str)
	case *PortDirection:
		*v = convertPortDirectionFromString(str)
	case *NodeType:
		*v = convertNodeTypeFromString(str)
	}
	return errors.New("invalid format")
}

func convertBoolFromString(str string) bool {
	if len(str) == 1 {
		if str[0] == '0' {
			return false
		}
		if str[0] == '1' {
			return true
		}
	} else if len(str) == 4 {
		if str == "true" || str == "TRUE" || str == "True" {
			return true
		}
	} else if len(str) == 5 {
		if str == "false" || str == "FALSE" || str == "False" {
			return false
		}
	}
	panic("convertFromString(): invalid bool conversion")
}

func convertNodeStatusFromString(str string) NodeStatus {
	if str == "IDLE" {
		return NodeStatus_IDLE
	}

	if str == "RUNNING" {
		return NodeStatus_RUNNING
	}

	if str == "SUCCESS" {
		return NodeStatus_SUCCESS
	}

	if str == "FAILURE" {
		return NodeStatus_FAILURE
	}

	if str == "SKIPPED" {
		return NodeStatus_SKIPPED
	}

	panic(fmt.Sprintf("Cannot convert this to NodeStatus:%v ", str))
}

func convertNodeTypeFromString(str string) NodeType {
	if str == "Action" {
		return NodeType_ACTION
	}

	if str == "Condition" {
		return NodeType_CONDITION
	}

	if str == "Control" {
		return NodeType_CONTROL
	}

	if str == "Decorator" {
		return NodeType_DECORATOR
	}

	if str == "SubTree" {
		return NodeType_SUBTREE
	}

	return NodeType_UNDEFINED
}

func convertPortDirectionFromString(str string) PortDirection {
	if str == "Input" || str == "INPUT" {
		return PortDirection_INPUT
	}

	if str == "Output" || str == "OUTPUT" {
		return PortDirection_OUTPUT
	}

	if str == "InOut" || str == "INOUT" {
		return PortDirection_INOUT
	}

	panic(fmt.Sprintf("Cannot convert this to PortDirection: %v", str))
}

func convertInt64sFromString(str string) (res []int64, err error) {
	parts := strings.Split(str, ";")
	for _, v := range parts {
		r, err := convertInt64FromString(v)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func convertFloat64sFromString(str string) (res []float64, err error) {
	parts := strings.Split(str, ";")
	for _, v := range parts {
		r, err := convertFloat64FromString(v)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func convertFloat64FromString(str string) (res float64, err error) {
	return strconv.ParseFloat(str, 64)
}

func convertInt64FromString(str string) (res int64, err error) {
	return strconv.ParseInt(str, 10, 64)
}

func ParseScript(s string) (ScriptFunction, error) {
	return nil, nil
}

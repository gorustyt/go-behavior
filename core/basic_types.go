package core

import (
	"container/list"
	"errors"
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

type IFromStr interface {
	FromString(str string) error
}

func ConvFromString(str string, value any) error {
	switch v := value.(type) {
	case *int:
		res, err := ConvertInt64FromString(str)
		if err != nil {
			return err
		}
		*v = int(res)
	case *float64:
		res, err := ConvertFloat64FromString(str)
		if err != nil {
			return err
		}
		*v = res
	case *[]float64:
		res, err := ConvertFloat64sFromString(str)
		if err != nil {
			return err
		}
		*v = append(*v, res...)
	case *int32:
		res, err := ConvertInt64FromString(str)
		if err != nil {
			return err
		}
		*v = int32(res)
	case *int64:
		res, err := ConvertInt64FromString(str)
		if err != nil {
			return err
		}
		*v = res
	case *[]int64:
		res, err := ConvertInt64sFromString(str)
		if err != nil {
			return err
		}
		*v = append(*v, res...)
	case *bool:
		*v = ConvertBoolFromString(str)
	default:
		if tmp, ok := value.(IFromStr); ok {
			return tmp.FromString(str)
		}
		return errors.New("invalid format")
	}
	return nil
}

func ConvertBoolFromString(str string) bool {
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

func ConvertInt64sFromString(str string) (res []int64, err error) {
	parts := strings.Split(str, ";")
	for _, v := range parts {
		r, err := ConvertInt64FromString(v)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func ConvertFloat64sFromString(str string) (res []float64, err error) {
	parts := strings.Split(str, ";")
	for _, v := range parts {
		r, err := ConvertFloat64FromString(v)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func ConvertFloat64FromString(str string) (res float64, err error) {
	return strconv.ParseFloat(str, 64)
}

func ConvertInt64FromString(str string) (res int64, err error) {
	return strconv.ParseInt(str, 10, 64)
}

func ParseScript(s string) (ScriptFunction, error) {
	return nil, nil
}

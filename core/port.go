package core

import (
	"sync"
)

var (
	ports sync.Map
)

func SetPorts(descType any, ps ...*PortInfo) {
	for _, v := range ps {
		ports.Store(v.Name, v)
	}
}

func GetPorts(name string) (res *PortInfo, ok bool) {
	v, ok := ports.Load(name)
	if !ok {
		return res, ok
	}
	return v.(*PortInfo), ok
}

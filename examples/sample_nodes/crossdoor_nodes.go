package sample_nodes

import (
	"github.com/gorustyt/go-behavior/core"
	"time"
)

type CrossDoor struct {
	_door_open     bool
	_door_locked   bool
	_pick_attempts int
}

func NewCrossDoor() *CrossDoor {
	return &CrossDoor{
		_door_locked: true,
	}
}

func SleepMS(ms int64) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (n *CrossDoor) IsDoorClosed() core.NodeStatus {
	SleepMS(200)
	if !n._door_open {
		return core.NodeStatus_SUCCESS
	}
	return core.NodeStatus_FAILURE
}

func (n *CrossDoor) PassThroughDoor() core.NodeStatus {
	SleepMS(500)
	if n._door_open {
		return core.NodeStatus_SUCCESS
	}
	return core.NodeStatus_FAILURE
}

func (n *CrossDoor) OpenDoor() core.NodeStatus {
	SleepMS(500)
	if n._door_locked {
		return core.NodeStatus_FAILURE
	} else {
		n._door_open = true
		return core.NodeStatus_SUCCESS
	}
}

func (n *CrossDoor) PickLock() core.NodeStatus {
	SleepMS(500)
	// succeed at 3rd attempt
	if n._pick_attempts > 3 {
		n._door_locked = false
		n._door_open = true
	}
	n._pick_attempts++
	if n._door_open {
		return core.NodeStatus_SUCCESS
	}
	return core.NodeStatus_FAILURE
}

func (n *CrossDoor) SmashDoor() core.NodeStatus {
	n._door_locked = false
	n._door_open = true
	// smash always works
	return core.NodeStatus_SUCCESS
}

func (n *CrossDoor) registerNodes(factory core.BehaviorTreeFactory) {
	factory.RegisterSimpleCondition(
		"IsDoorClosed", n.IsDoorClosed)

	factory.RegisterSimpleAction(
		"PassThroughDoor", n.PassThroughDoor)

	factory.RegisterSimpleAction(
		"OpenDoor", n.OpenDoor)

	factory.RegisterSimpleAction(
		"PickLock", n.PickLock)

	factory.RegisterSimpleCondition(
		"SmashDoor", n.SmashDoor)
}

func (n *CrossDoor) Reset() {
	n._door_open = false
	n._door_locked = true
	n._pick_attempts = 0
}

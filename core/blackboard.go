package core

import (
	"fmt"
	"reflect"
	"sync"
)

type Blackboard struct {
	mutex_             sync.Mutex
	recursiveMutex     sync.Mutex
	storage            map[string]*Entry
	parentBb           *Blackboard
	internalToExternal map[string]string
	automapping        bool
}

func NewBlackboard(parent *Blackboard) *Blackboard {
	return &Blackboard{
		internalToExternal: map[string]string{},
		storage:            map[string]*Entry{},
		parentBb:           parent,
	}
}

type Entry struct {
	entryMutex sync.Mutex
	Value      any
}

func (n *Blackboard) Clear() {
	n.mutex_.Lock()
	n.storage = map[string]*Entry{}
	n.mutex_.Unlock()
}

func IsPrivateKey(str string) bool {
	return len(str) >= 1 && str[0] == '_'
}

func (n *Blackboard) enableAutoRemapping(remapping bool) {
	n.automapping = remapping
}
func (n *Blackboard) Get(key string) any {
	if any_ref := n.GetAnyLocked(key); any_ref != nil {
		a := any_ref()
		if a == nil {
			panic(fmt.Sprintf("Blackboard::get() error. Entry [%v] hasn't been initialized, yet", key))
		}
		return a
	} else {
		panic(fmt.Sprintf("Blackboard::get() error. Missing key [%v]", key))
	}
}

func (n *Blackboard) Unset(key string) {
	n.mutex_.Lock()
	defer n.mutex_.Unlock()
	delete(n.storage, key)
}

func (n *Blackboard) GetAnyLocked(key string) func() *Entry {
	return func() (entry *Entry) {
		n.mutex_.Lock()
		entry = n.GetEntry(key)
		n.mutex_.Unlock()
		return entry
	}
}

func (n *Blackboard) AddSubtreeRemapping(internal, external string) {
	n.internalToExternal[internal] = external
}

func (n *Blackboard) DebugMessage() {
	for key, entry := range n.storage {
		fmt.Printf("%v (%v)", key, reflect.ValueOf(entry.Value).String())
	}

	for from, to := range n.internalToExternal {
		fmt.Printf("[%v] remapped to port of parent tree [%v]", from, to)
		continue
	}
}

func (n *Blackboard) GetKeys() (res []string) {
	if len(n.storage) == 0 {
		return
	}
	for k := range n.storage {
		res = append(res, k)
	}
	return
}

func (n *Blackboard) GetEntry(key string) *Entry {
	n.mutex_.Lock()
	defer n.mutex_.Unlock()
	it, ok := n.storage[key]
	if ok {
		return it
	}

	// not found. Try autoremapping
	if parent := n.parentBb; parent != nil {
		newKey, ok := n.internalToExternal[key]
		if ok {
			entry := parent.GetEntry(newKey)
			if entry != nil {
				n.storage[key] = entry
			}
			return entry
		}
		if n.automapping && !IsPrivateKey(key) {
			entry := parent.GetEntry(key)
			if entry != nil {
				n.storage[key] = entry
			}
			return entry
		}
	}
	return nil
}

func (n *Blackboard) CreateEntry(key string, info *PortInfo) *Entry {
	return n.createEntryImpl(key, info)
}

func (n *Blackboard) createEntryImpl(key string, info *PortInfo) *Entry {
	n.mutex_.Lock()
	defer n.mutex_.Unlock()
	// This function might be called recursively, when we do remapping, because we move
	// to the top scope to find already existing  entries

	// search if exists already
	storageIt, ok := n.storage[key]
	if ok {
		return storageIt
	}

	var entry *Entry

	// manual remapping first
	remappedKey, ok := n.internalToExternal[key]
	if ok {
		if n.parentBb != nil {
			entry = n.parentBb.createEntryImpl(remappedKey, info)
		}
	} else if n.automapping && !IsPrivateKey(key) {
		if n.parentBb != nil {
			entry = n.parentBb.createEntryImpl(key, info)
		}
	} else {
		// not remapped, not found. Create locally.
		entry = &Entry{}
		// even if empty, let's assign to it a default type
		entry.Value = info.defaultValue

	}

	n.storage[key] = entry
	return entry
}

func (n *Blackboard) Set(key string, value any) {
	n.mutex_.Lock()
	defer n.mutex_.Unlock()
	entry, ok := n.storage[key]
	if !ok {
		n.mutex_.Unlock()
		// if a new generic port is created with a string, it's type should be AnyTypeAllowed
		s, ok := value.(string)
		p := NewPortInfo(PortDirection_INOUT, "")
		p.defaultValueStr = s
		if ok {
			entry = n.createEntryImpl(key, p)
		} else {
			p.SetDefaultValue(value)
			entry = n.createEntryImpl(key, p)
		}
		n.mutex_.Lock()
		n.storage[key] = entry
		entry.Value = value
	} else {
		// this is not the first time we set this entry, we need to check
		// if the type is the same or not.

		entry.entryMutex.Lock()
		defer entry.entryMutex.Unlock()
		if entry == nil || reflect.TypeOf(entry.Value) == nil {
			entry.Value = value
			entry.entryMutex.Unlock()
			return
		}

		previousType := reflect.TypeOf(entry.Value)

		// check type mismatch
		if previousType != reflect.TypeOf(value) {
			mismatching := true
			if v, ok := value.(fmt.Stringer); ok {
				anyFromString := v.String()
				if anyFromString != "" {
					mismatching = false
					entry.Value = anyFromString
				}
			}
			if mismatching {
				n.DebugMessage()
				panic(fmt.Sprintf("Blackboard::set(%v", key, "): once declared, the type of a port shall not change. Previously declared type [, reflect.TypeOf(previous_type),], current type [", "]"))
			}
		}
	}
}

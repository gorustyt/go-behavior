package core

import (
	"fmt"
	"reflect"
	"sync"
)

type Blackboard struct {
	mutex_                sync.Mutex
	recursive_mutex       sync.Mutex
	storage_              map[string]*Entry
	parent_bb_            *Blackboard
	internal_to_external_ map[string]string
	autoremapping_        bool
}

func NewBlackboard() *Blackboard {
	return &Blackboard{
		internal_to_external_: map[string]string{},
		storage_:              map[string]*Entry{},
	}
}

type Entry struct {
	Value any
}

func (n *Blackboard) Clear() {
	n.mutex_.Lock()
	n.storage_ = map[string]*Entry{}
	n.mutex_.Unlock()
}

func IsPrivateKey(str string) bool {
	return len(str) >= 1 && str[0] == '_'
}

func (n *Blackboard) enableAutoRemapping(remapping bool) {
	n.autoremapping_ = remapping
}
func (n *Blackboard) Get(key string) any {
	if any_ref := n.getAnyLocked(key); any_ref != nil {
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
	_, ok := n.storage_[key]
	if !ok {
		return
	}
	delete(n.storage_, key)
}

func (n *Blackboard) getAnyLocked(key string) func() *Entry {
	return func() (entry *Entry) {
		n.mutex_.Lock()
		entry = n.GetEntry(key)
		n.mutex_.Unlock()
		return entry
	}
}

func (n *Blackboard) AddSubtreeRemapping(internal, external string) {
	n.internal_to_external_[internal] = external
}

func (n *Blackboard) DebugMessage() {
	for key, entry := range n.storage_ {
		fmt.Printf("%v (%v)", key, reflect.ValueOf(entry.Value).String())
	}

	for from, to := range n.internal_to_external_ {
		fmt.Printf("[%v] remapped to port of parent tree [%v]", from, to)
		continue
	}
}

func (n *Blackboard) GetKeys() (res []string) {
	if len(n.storage_) == 0 {
		return
	}
	for k := range n.storage_ {
		res = append(res, k)
	}
	return
}

func (n *Blackboard) GetEntry(key string) *Entry {
	n.mutex_.Lock()
	defer n.mutex_.Unlock()
	it, ok := n.storage_[key]
	if ok {
		return it
	}

	// not found. Try autoremapping
	if parent := n.parent_bb_; parent != nil {
		new_key, ok := n.internal_to_external_[key]
		if ok {

			entry := parent.GetEntry(new_key)
			if entry != nil {
				n.storage_[key] = entry
			}
			return entry
		}
		if n.autoremapping_ && !IsPrivateKey(key) {
			entry := parent.GetEntry(key)
			if entry != nil {
				n.storage_[key] = entry
			}
			return entry
		}
	}
	return nil
}

func (n *Blackboard) CreateEntry(key string) *Entry {
	return n.createEntryImpl(key)
}

func (n *Blackboard) createEntryImpl(key string) *Entry {
	n.mutex_.Lock()
	defer n.mutex_.Unlock()
	// This function might be called recursively, when we do remapping, because we move
	// to the top scope to find already existing  entries

	// search if exists already
	storageIt, ok := n.storage_[key]
	if ok {
		return storageIt
	}

	var entry *Entry

	// manual remapping first
	remappedKey, ok := n.internal_to_external_[key]
	if ok {
		if n.parent_bb_ != nil {
			entry = n.parent_bb_.createEntryImpl(remappedKey)
		}
	} else if n.autoremapping_ && !IsPrivateKey(key) {
		if n.parent_bb_ != nil {
			entry = n.parent_bb_.createEntryImpl(key)
		}
	}
	n.storage_[key] = entry
	return entry
}

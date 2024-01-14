package core

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core/parser"
)

type NodeBuilder func(name string, config *NodeConfig) *TreeNode
type BehaviorTreeFactory struct {
	builders                map[string]*NodeBuilder
	manifests               map[string]*TreeNodeManifest
	builtinIDs              map[string]struct{}
	behaviorTreeDefinitions map[string]any
	scriptingEnums          map[string]int
	substitutionRules       map[string]SubstitutionRule
	parser                  parser.Parser
}

func (f *BehaviorTreeFactory) registerBuilder(manifest *TreeNodeManifest,
	builder *NodeBuilder) {
	_, ok := f.builders[manifest.registrationID]
	if ok {
		panic(fmt.Sprintf("ID [%v] already registered", manifest.registrationID))
	}
	f.builders[manifest.registrationID] = builder
	f.manifests[manifest.registrationID] = manifest
}

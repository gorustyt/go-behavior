package core

type Parser interface {
	LoadFromFile(fileName string, addIncludes ...bool) error
	LoadFromText(xmlText string, addIncludes ...bool) error
	RegisteredBehaviorTrees() []string
	InstantiateTree(rootBlackboard *Blackboard, mainTreeId string) (tree *Tree, err error)
	ClearInternalState()
}

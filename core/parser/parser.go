package parser

type Parser interface {
	LoadFromFile(fileName string, addIncludes ...bool)
	LoadFromText(xmlText string, addIncludes ...bool)
	RegisteredBehaviorTrees() []string
	InstantiateTree()
	ClearInternalState()
}

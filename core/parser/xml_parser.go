package parser

type xmlParser struct {
}

func NewParser() Parser {
	return &xmlParser{}
}

func (p *xmlParser) LoadFromFile(fileName string, addIncludes ...bool) {
	addInclude := true
	if len(addIncludes) > 0 {
		addInclude = addIncludes[0]
	}
}

func (p *xmlParser) LoadFromText(xmlText string, addIncludes ...bool) {
	addInclude := true
	if len(addIncludes) > 0 {
		addInclude = addIncludes[0]
	}
}

func (p *xmlParser) RegisteredBehaviorTrees() []string {

}

func (p *xmlParser) InstantiateTree() {

}

func (p *xmlParser) ClearInternalState() {

}

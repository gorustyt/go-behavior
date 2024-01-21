package core

import (
	"bytes"
	"container/list"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gorustyt/go-behavior/decorators"
	"io"
	"os"
	"reflect"
	"strings"
)

var (
	ErrorsInvalidTag = errors.New("invalid tag")
)

type SubtreeModel struct {
	ports map[string]*PortInfo
}

type xmlParser struct {
	rootTag       []*XmlTag
	treesRoot     *SortMap[string, *XmlTag] //所有的树
	stack         *stack
	factory       *BehaviorTreeFactory
	currentPath   string
	subtreeModels map[string]*SubtreeModel
	suffixCount   int
}

func NewXmlParser() Parser {
	return &xmlParser{
		stack:     &stack{},
		treesRoot: NewSortMap[string, *XmlTag](),
	}
}

func (p *xmlParser) LoadFromFile(fileName string, addIncludes ...bool) error {
	addInclude := true
	if len(addIncludes) > 0 {
		addInclude = addIncludes[0]
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return p.Parse(bytes.NewReader(data), addInclude)
}

func (p *xmlParser) LoadFromText(xmlText string, addIncludes ...bool) error {
	addInclude := true
	if len(addIncludes) > 0 {
		addInclude = addIncludes[0]
	}
	return p.Parse(strings.NewReader(xmlText), addInclude)
}

func (p *xmlParser) RegisteredBehaviorTrees() []string {
	var res []string
	p.treesRoot.Range(func(key string, value *XmlTag) (stop bool) {
		res = append(res, key)
		return false
	})
	return res
}

func (p *xmlParser) createNodeFromXML(element *XmlTag,
	blackboard *Blackboard,
	nodeParent ITreeNode,
	prefixPath string,
	outputTree *Tree) (tree ITreeNode, err error) {
	elementName := element.TagName()
	elementId := element.GetAttr("ID")
	var nodeType NodeType
	err = ConvFromString(elementName, &nodeType)
	if err != nil {
		return
	}
	// name used by the factory
	var typeId string

	if nodeType == NodeType_UNDEFINED {
		// This is the case of nodes like <MyCustomAction>
		// check if the factory has this name
		if _, ok := p.factory.Builders[elementName]; !ok {
			return tree, fmt.Errorf("%v  is not a registered node", elementName)
		}
		typeId = elementName

		if elementId != "" {
			return tree, fmt.Errorf("attribute [ID] is not allowed in <%v>", typeId)
		}
	} else {
		// in this case, it is mandatory to have a field "ID"
		if elementId == "" {
			return tree, fmt.Errorf("attribute [ID] is mandatory in <%v>", typeId)
		}
		typeId = elementId
	}

	// By default, the instance name is equal to ID, unless the
	// attribute [name] is present.
	attrName := element.GetAttr("name")
	instanceName := attrName
	if attrName == "" {
		instanceName = typeId
	}

	b, ok := p.factory.Builders[typeId]
	portRemap := map[string]string{}
	for _, att := range element.GetAttrs() {
		if IsAllowedPortName(att.Name.Local) {
			attributeName := att.Name.Local
			portRemap[attributeName] = att.Value
		}
	}

	var config NodeConfig
	config.Blackboard = blackboard
	config.Path = prefixPath + instanceName
	config.Uid = outputTree.GetUID()
	config.Manifest = b.TreeNodeManifest

	if typeId == instanceName {
		config.Path += fmt.Sprintf("::%v", config.Uid)
	}

	for i := 0; i < int(PreCond_COUNT_); i++ {
		pre := PreCond(i)
		if script := element.GetAttr(pre.String()); script != "" {
			config.PreConditions[pre] = script
		}
	}
	for i := 0; i < int(PostCond_COUNT_); i++ {
		post := PostCond(i)
		if script := element.GetAttr(post.String()); script != "" {
			config.PostConditions[post] = script
		}
	}

	//---------------------------------------------
	var newNode ITreeNode

	if nodeType == NodeType_SUBTREE {
		config.InputPorts = portRemap
		newNode, err = p.factory.InstantiateTreeNode(instanceName,
			NodeType_SUBTREE.String(),
			&config)
		if err != nil {
			return nil, err
		}

		subtreeNode, ok := newNode.(*decorators.SubTreeNode)
		if ok {
			subtreeNode.SetSubtreeID(typeId)
		}
	} else {
		if !ok || b == nil || b.TreeNodeManifest == nil {
			return tree, fmt.Errorf("missing manifest for element_ID: %v . It shouldn't happen. Please report this issue", elementId)
		}
		manifest := b.TreeNodeManifest
		//Check that name in remapping can be found in the manifest
		for nameInSubtree := range portRemap {
			_, ok = manifest.Ports[nameInSubtree]
			if !ok {
				return tree, fmt.Errorf("possible typo? In the XML, you tried to remap port \"%v\" in node [%v(type :%v)], but the manifest of this node does not contain a port with this name",
					nameInSubtree, config.Path, typeId)
			}
		}

		// Initialize the ports in the BB to set the type
		for portName, portInfo := range manifest.Ports {
			remappedPort, ok := portRemap[portName]
			if !ok {
				continue
			}
			portKey, err := GetRemappedKey(portName, remappedPort)
			if err != nil {
				return nil, err
			}
			if portKey != "" {
				// port_key will contain the key to find the entry in the blackboard
				// if the entry already exists, check that the type is the same
				prevInfo := blackboard.GetEntry(portKey)
				if prevInfo != nil {
					// Check consistency of types.
					portTypeMismatch := reflect.TypeOf(prevInfo.Value) != reflect.TypeOf(portInfo.defaultValue)

					// special case related to convertFromString
					_, stringInput := prevInfo.Value.(string)
					if portTypeMismatch && !stringInput {
						blackboard.DebugMessage()

						return tree, fmt.Errorf("The creation of the tree failed because the port [%v] was initially created with type [] and, later type [] was used somewhere else.", portKey)
					}
				} else {
					// not found, insert for the first time.
					blackboard.CreateEntry(portKey, portInfo)
				}
			}
		}

		// Set the port direction in config
		for portName, v := range portRemap {
			portIt, ok := manifest.Ports[portName]
			if ok {
				direction := portIt.direction
				if direction != PortDirection_OUTPUT {
					config.InputPorts[portName] = v
				}
				if direction != PortDirection_INPUT {
					config.OutputPorts[portName] = v
				}
			}
		}

		// use default value if available for empty ports. Only inputs
		for portName, portInfo := range manifest.Ports {
			direction := portInfo.direction
			_, ok := config.InputPorts[portName]
			if direction != PortDirection_OUTPUT && !ok && portInfo.defaultValue != nil {
				config.InputPorts[portName] = portInfo.DefaultValueString()
			}
		}

		newNode, err = p.factory.InstantiateTreeNode(instanceName, typeId, &config)
		if err != nil {
			return nil, err
		}
	}

	// add the pointer of this node to the parent
	if nodeParent != nil {

		if controlParent, ok := nodeParent.(*ControlNode); ok {
			controlParent.AddChild(newNode)
		} else if decoratorParent, ok := nodeParent.(*DecoratorNode); ok {
			err = decoratorParent.SetChild(newNode)
			if err != nil {
				return nil, err
			}
		}
	}
	return newNode, nil
}

func (p *xmlParser) recursivelyCreateSubtree(parentNode ITreeNode,
	subtree *Subtree,
	outputTree *Tree,
	blackboard *Blackboard,
	prefix string,
	element *XmlTag) error {
	// create the node
	node, err := p.createNodeFromXML(element, blackboard, parentNode, prefix, outputTree)
	if err != nil {
		return err
	}
	subtree.Nodes = append(subtree.Nodes, node)

	// common case: iterate through all children
	if node.NodeType() != NodeType_SUBTREE {
		for _, childElement := range element.Children {
			err = p.recursivelyCreateSubtree(node, subtree, outputTree, blackboard, prefix, childElement)
			if err != nil {
				return err
			}
		}
	} else { // special case: SubTreeNode
		newBb := NewBlackboard(blackboard)
		subtreeId := element.GetAttr("ID")
		remapping := map[string]string{}
		doAutoRemap := false
		for _, attr := range element.GetAttrs() {
			attrName := attr.Name.Local
			attrValue := attr.Value
			if attrName == "_autoremap" {
				err = ConvFromString(attrValue, &doAutoRemap)
				if err != nil {
					return err
				}
				newBb.enableAutoRemapping(doAutoRemap)
				continue
			}
			if !IsAllowedPortName(attrName) {
				continue
			}
			remapping[attrName] = attrValue
		}
		// check if this subtree has a model. If it does,
		// we want o check if all the mandatory ports were remapped and
		// add default ones, if necessary
		subtreeModelPorts, ok := p.subtreeModels[subtreeId]
		if ok {
			// check if:
			// - remapping contains mondatory ports
			// - if any of these has default value
			for portName, portInfo := range subtreeModelPorts.ports {
				_, ok := remapping[portName]
				// don't override existing remapping
				if !ok && !doAutoRemap {
					// remapping is not explicitly defined in the XML: use the model
					if portInfo.DefaultValueString() == "" {
						return fmt.Errorf("in the <TreeNodesModel> the <Subtree ID=\"%v\"> is defining a mandatory port called [%v], but you are not remapping it", subtreeId,
							portName)
					} else {
						remapping[portName] = portInfo.DefaultValueString()
					}
				}
			}
		}

		for attrName, attrValue := range remapping {
			if _, ok = IsBlackboardPointer(attrValue); ok {
				// do remapping
				portName := StripBlackboardPointer(attrValue)
				newBb.AddSubtreeRemapping(attrName, portName)
			} else {
				// constant string: just set that constant value into the BB
				// IMPORTANT: this must not be autoremapped!!!
				newBb.enableAutoRemapping(false)
				newBb.Set(attrName, attrValue)
				newBb.enableAutoRemapping(doAutoRemap)
			}
		}

		subtreePath := subtree.InstanceName
		if subtreePath != "" {
			subtreePath += "/"
		}
		if name := element.GetAttr("name"); name != "" {
			subtreePath += name
		} else {
			subtreePath += fmt.Sprintf("%v::%v", subtreeId, node.UID())
		}

		return p.createSubtree(subtreeId,
			subtreePath,     // name
			subtreePath+"/", //prefix
			outputTree, newBb, node)
	}
	return nil
}

func (p *xmlParser) createSubtree(
	treeId string,
	treePath string,
	prefixPath string,
	outputTree *Tree,
	blackboard *Blackboard, rootNode ITreeNode) (err error) {
	rootElement, ok := p.treesRoot.Get(treeId)
	if !ok {
		return fmt.Errorf("can't find a tree with name: %v", treeId)
	}

	// Append a new subtree to the list
	newTree := &Subtree{}
	newTree.Blackboard = blackboard
	newTree.InstanceName = treePath
	newTree.TreeId = treeId
	outputTree.Subtrees = append(outputTree.Subtrees, newTree)
	return p.recursivelyCreateSubtree(rootNode, newTree, outputTree, blackboard, prefixPath, rootElement)
}

func (p *xmlParser) InstantiateTree(rootBlackboard *Blackboard, mainTreeId string) (tree *Tree, err error) {
	tree = NewTree()
	if mainTreeId == "" {
		root := p.rootTag[0]
		mainTreeId = root.GetAttr("main_tree_to_execute")
		if mainTreeId == "" && len(p.rootTag) == 1 {
			// special case: there is only one registered BT.
			v, _ := p.treesRoot.Front()
			mainTreeId = v.TagName()
		} else {
			return tree, fmt.Errorf("[main_tree_to_execute] was not specified correctly")
		}
	}
	if rootBlackboard == nil {
		return tree, fmt.Errorf("XMLParser::instantiateTree needs a non-empty root_blackboard")
	}

	err = p.createSubtree(mainTreeId, "", "",
		tree, rootBlackboard, nil)
	if err != nil {
		return nil, err
	}
	tree.Init()
	return tree, nil
}

func (p *xmlParser) ClearInternalState() {

}

func (p *xmlParser) parseSubtreeModel(tag *XmlTag) error {
	var portMap = map[string]PortDirection{"input_port": PortDirection_INPUT,
		"output_port": PortDirection_OUTPUT,
		"inout_port":  PortDirection_INOUT}
	for _, v := range tag.Children {
		if v.IsTag("SubTree") {
			id := v.GetAttr("ID")
			models, ok := p.subtreeModels[id]
			if !ok {
				p.subtreeModels[id] = &SubtreeModel{ports: map[string]*PortInfo{}}
			}
			for name, direction := range portMap {
				if v.IsTag(name) {
					n := v.GetAttr(name)
					if n == "" {
						return errors.New("missing attribute [name] in port (SubTree model)")
					}
					port := NewPortInfo(direction, n, v.GetAttr("description"))
					port.SetDefaultValue(v.GetAttr("default"))
					models.ports[n] = port
				}
			}
		}
	}
	return nil
}

type XmlTag struct {
	element  xml.StartElement
	Children []*XmlTag
}

func (p *XmlTag) TagName() string {
	return p.element.Name.Local
}
func (p *XmlTag) IsTag(key string) bool {
	return p.element.Name.Local == key
}

func (p *XmlTag) GetAttrs() (res []xml.Attr) {
	for _, v := range p.element.Attr {
		res = append(res, v)
	}
	return res
}

func (p *XmlTag) GetAttr(key string) string {
	for _, v := range p.element.Attr {
		if v.Name.Local == key {
			return v.Value
		}
	}
	return ""
}

func (p *xmlParser) parse(tags []*XmlTag) error {
	for _, v := range tags {
		var err error
		switch v.TagName() {
		case "BehaviorTree":
			p.parseBehaviorTree(v)
		case "TreeNodesModel":
			err = p.parseSubtreeModel(v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *xmlParser) verify(roots []*XmlTag) error {
	if roots == nil || len(roots) == 0 {
		return errors.New("the XML must have a root node called <root>")
	}
	if len(roots) != 1 {
		return errors.New("only a single node <TreeNodesModel> is supported")
	}
	root := roots[0]
	var model *XmlTag
	for _, child := range root.Children {
		if model != nil {
			return errors.New("only a single node <TreeNodesModel> is supported")
		}
		if child.IsTag("TreeNodesModel") {
			model = child
		}
	}
	if model != nil {
		// not having a MetaModel is not an error. But consider that the
		// Graphical editor needs it.
		for _, child := range root.Children {
			name := child.TagName()
			if name == "Action" || name == "Decorator" ||
				name == "SubTree" || name == "Condition" ||
				name == "Control" {
				if child.GetAttr("ID") == "" {
					return fmt.Errorf("Error at line  . The attribute [ID] is mandatory")
				}
			}
		}
	}
	for _, v := range root.Children {
		if v.IsTag("BehaviorTree") {
			err := p.verifyNode(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (p *xmlParser) GetPortsRecursively(element *XmlTag, res *[]string) {
	for _, attr := range element.GetAttrs() {
		attrName := attr.Name.Local
		attrValue := attr.Value
		_, ok := IsBlackboardPointer(attrValue)
		if IsAllowedPortName(attrName) && ok {
			portName := StripBlackboardPointer(attrValue)
			*res = append(*res, portName)
		}
	}
	for _, v := range element.Children {
		p.GetPortsRecursively(v, res)
	}
}
func (p *xmlParser) verifyNode(tag *XmlTag) error {
	childrenCount := len(tag.Children)
	name := tag.TagName()
	if name == "Decorator" {
		if childrenCount != 1 {
			return fmt.Errorf("the node <Decorator> must have exactly 1 child")
		}
		if tag.GetAttr("ID") == "" {
			return fmt.Errorf("the node <Decorator> must have the attribute [ID]")
		}
	} else if name == "Action" {
		if childrenCount != 0 {
			return fmt.Errorf("the node <Action> must not have any child")
		}
		if tag.GetAttr("ID") == "" {
			return fmt.Errorf("the node <Action> must have the attribute [ID]")
		}
	} else if name == "Condition" {
		if childrenCount != 0 {
			return fmt.Errorf("the node <Condition> must not have any child")
		}
		if tag.GetAttr("ID") == "" {
			return fmt.Errorf("the node <Condition> must have the attribute [ID]")
		}
	} else if name == "Control" {
		if childrenCount == 0 {
			return fmt.Errorf("the node <Control> must have at least 1 child")
		}
		if tag.GetAttr("ID") == "" {
			return fmt.Errorf("the node <Control> must have the attribute [ID]")
		}
	} else if name == "Sequence" || name == "ReactiveSequence" ||
		name == "SequenceWithMemory" || name == "Fallback" {
		if childrenCount == 0 {
			var nameAttr string
			if tag.GetAttr("name") != "" {
				nameAttr = "(`" + tag.GetAttr("name") + "`)"
			}
			return fmt.Errorf("A Control node must have at least 1 child, error in XML node %v ` %v`", tag.TagName(), nameAttr)
		}
	} else if name == "SubTree" {
		if len(tag.Children) > 0 {
			if tag.Children[0].TagName() == "remap" {
				return fmt.Errorf("<remap> was deprecated")
			} else {
				return fmt.Errorf("<SubTree> should not have any child")
			}
		}

		if tag.GetAttr("ID") == "" {
			return fmt.Errorf("the node <SubTree> must have the attribute [ID]")
		}
	} else if name == "BehaviorTree" {
		if childrenCount != 1 {
			return fmt.Errorf("the node <BehaviorTree> must have exactly 1 child")
		}
	} else {
		// search in the factory and the list of subtrees
		search, ok := p.factory.Builders[name]
		if !ok {
			return fmt.Errorf("node not recognized: %v", name)
		}

		if search.Type == NodeType_DECORATOR {
			if childrenCount != 1 {
				return fmt.Errorf("The node <%v> must have exactly 1 child ", name)
			}
		}
	}
	//recursion
	if name != "SubTree" {
		for _, v := range tag.Children {
			err := p.verifyNode(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *xmlParser) parseRoot(roots []*XmlTag, addInclude bool) error {
	if len(roots) == 0 {
		return errors.New("not found roots ")
	}
	tag := roots[0]
	format := tag.GetAttr("BTCPP_format")
	if format == "" {
		fmt.Println("Warnings: The first tag of the XML (<root>) should contain the attribute [BTCPP_format=\"4\"]\nPlease check if your XML is compatible with version 4.x of BT.CPP")
	}
	for _, v := range tag.Children {
		if !addInclude {
			break
		}
		if v.IsTag("include") {
			fpath := v.GetAttr("path")
			rospath := v.GetAttr("ros_pkg")
			if rospath != "" {
				fpath = rospath
			}
			data, err := os.ReadFile(fpath)
			if err != nil {
				return err
			}
			tags, err := p.buildTags(bytes.NewReader(data))
			if err != nil {
				return err
			}
			err = p.parseRoot(tags, addInclude)
			if err != nil {
				return err
			}
		}
	}
	err := p.verify(roots)
	if err != nil {
		return err
	}
	p.rootTag = append(p.rootTag, roots[0])
	return p.parse(roots[0].Children)
}

func (p *xmlParser) parseBehaviorTree(tag *XmlTag) {
	id := tag.GetAttr("ID")
	if id == "" {
		id = fmt.Sprintf("BehaviorTree_%v", p.suffixCount)
		p.suffixCount++
	}
	p.treesRoot.Set(id, tag)
}

func (p *xmlParser) Parse(reader io.Reader, add_includes bool) error {
	roots, err := p.buildTags(reader)
	if err != nil {
		return err
	}
	return p.parseRoot(roots, add_includes) //TODO
}

func (p *xmlParser) buildTags(reader io.Reader) (roots []*XmlTag, err error) {
	d := xml.NewDecoder(reader)
	for t, err := d.Token(); err == nil; t, err = d.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			tmp := &XmlTag{element: token}
			if tmp.IsTag("root") {
				roots = append(roots, tmp)
			}
			p.stack.Push(tmp)
		case xml.EndElement:
			e := p.stack.Peek()
			if e.element.Name.Local == token.Name.Local { //一个闭合标签
				p.stack.Pop()
				if p.stack.Len() > 0 {
					p.stack.Peek().Children = append(p.stack.Peek().Children, e)
				}
			}
		}
	}
	if p.stack.Len() > 0 {
		p.stack = &stack{}
		return roots, ErrorsInvalidTag
	}
	return roots, nil
}

type stack struct {
	list.List
}

func (s *stack) Peek() *XmlTag {
	b := s.List.Back()
	return b.Value.(*XmlTag)
}

func (s *stack) Pop() *XmlTag {
	b := s.List.Back()
	s.List.Remove(b)
	return b.Value.(*XmlTag)
}
func (s *stack) Len() int {
	return s.List.Len()
}
func (s *stack) Push(v *XmlTag) {
	s.PushBack(v)
}

package main

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
	"strings"
)

type Point3D struct {
	x, y, z float64
}

// We want to create an ActionNode that calls the method MyLegacyMoveTo::go
type MyLegacyMoveTo struct {
}

func (n *MyLegacyMoveTo) Go(goal Point3D) bool {
	fmt.Printf("Going to: %f %f %f\n", goal.x, goal.y, goal.z)
	return true // true means success in my legacy code
}

// Similarly to the previous tutorials, we need to implement this parsing method,
// providing a specialization of BT::convertFromString

func convertFromString(key string) *Point3D {
	// three real numbers separated by semicolons
	parts := strings.Split(key, ";")
	if parts.size() != 3 {
		panic("invalid input)")
	} else {
		var output Point3D
		output.x = convertFromString < double > (parts[0])
		output.y = convertFromString < double > (parts[1])
		output.z = convertFromString < double > (parts[2])
		return &output
	}
}

var xml_text = `(
 <root BTCPP_format="4">
     <BehaviorTree>
        <MoveTo  goal="-1;3;0.5" />
     </BehaviorTree>
 </root>
 )`

func main() {

	var move_to MyLegacyMoveTo

	// Here we use a lambda that captures the reference of move_to
	MoveToWrapperWithLambda := func(parentNode *core.TreeNode, status ...core.NodeStatus) core.NodeStatus {
		var goal Point3D
		// thanks to paren_node, you can access easily the input and output ports.
		parentNode.GetInput("goal", goal)

		res := move_to.Go(goal)
		if res {
			return core.NodeStatus_SUCCESS
		}
		// convert bool to NodeStatus
		return core.NodeStatus_FAILURE
	}

	factory := core.NewBehaviorTreeFactory()

	// Register the lambda with BehaviorTreeFactory::registerSimpleAction

	ports := core.InputPortWithDefaultValue("goal", &Point3D{})
	factory.RegisterSimpleAction("MoveTo", MoveToWrapperWithLambda, ports)

	tree, err := factory.CreateTreeFromText(xml_text)
	if err != nil {
		panic(err)
	}
	tree.TickWhileRunning()

}

/* Expected output:

Going to: -1.000000 3.000000 0.500000

*/

package main

import (
	"fmt"
	"github.com/gorustyt/go-behavior/core"
)

/** In the CrossDoor example we did not exchange any information
 * between the Maintree and the DoorClosed subtree.
 *
 * If we tried to do that, we would have noticed that it can't be done, because
 * each of the tree/subtree has its own Blackboard, to avoid the problem of name
 * clashing in very large trees.
 *
 * But a SubTree can have its own input/output ports.
 * In practice, these ports are nothing more than "soft links" between the
 * ports inside the SubTree (called "internal") and those in the parent
 * Tree (called "external").
 *
 */

// clang-format off

var xml_text = `(
<root BTCPP_format="4">

    <BehaviorTree ID="MainTree">
        <Sequence>
            <Script code=" move_goal='1;2;3' " />
            <SubTree ID="MoveRobot" target="{move_goal}" result="{move_result}" />
            <SaySomething message="{move_result}"/>
        </Sequence>
    </BehaviorTree>

    <BehaviorTree ID="MoveRobot">
        <Fallback>
            <Sequence>
                <MoveBase  goal="{target}"/>
                <Script code=" result:='goal reached' " />
            </Sequence>
            <ForceFailure>
                <Script code=" result:='error' " />
            </ForceFailure>
        </Fallback>
    </BehaviorTree>

</root>
 )`

func main() {
	var factory core.BehaviorTreeFactory

	factory.RegisterNodeType < SaySomething > ("SaySomething")
	factory.RegisterNodeType < MoveBaseAction > ("MoveBase")

	err := factory.RegisterBehaviorTreeFromText(xml_text)
	if err != nil {
		panic(err)
	}
	tree, err := factory.CreateTree("MainTree")
	if err != nil {
		panic(err)
	}
	tree.TickWhileRunning()

	// let's visualize some information about the current state of the blackboards.
	fmt.Printf("\n------ First BB ------")
	tree.Subtrees[0].Blackboard.DebugMessage()
	fmt.Printf("\n------ Second BB------")
	tree.Subtrees[1].Blackboard.DebugMessage()

}

/* Expected output:

------ First BB ------
move_result (std::string)
move_goal (Pose2D)

------ Second BB------
[result] remapped to port of parent tree [move_result]
[target] remapped to port of parent tree [move_goal]

*/

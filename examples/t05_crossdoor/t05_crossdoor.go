package main

import "github.com/gorustyt/go-behavior/core"

/** This is a more complex example that uses Fallback,
 * Decorators and Subtrees
 *
 * For the sake of simplicity, we aren't focusing on ports remapping to the time being.
 */

// clang-format off

var  xml_text = `(
<root BTCPP_format="4">

    <BehaviorTree ID="MainTree">
        <Sequence>
            <Fallback>
                <Inverter>
                    <IsDoorClosed/>
                </Inverter>
                <SubTree ID="DoorClosed"/>
            </Fallback>
            <PassThroughDoor/>
        </Sequence>
    </BehaviorTree>

    <BehaviorTree ID="DoorClosed">
        <Fallback>
            <OpenDoor/>
            <RetryUntilSuccessful num_attempts="5">
                <PickLock/>
            </RetryUntilSuccessful>
            <SmashDoor/>
        </Fallback>
    </BehaviorTree>

</root>
 )`

// clang-format on

func main() {
  var  factory core.BehaviorTreeFactory;

  var cross_door  CrossDoor ;
  cross_door.RegisterNodes(factory);

  // In this example a single XML contains multiple <BehaviorTree>
  // To determine which one is the "main one", we should first register
  // the XML and then allocate a specific tree, using its ID

  err:=factory.RegisterBehaviorTreeFromText(xml_text);
  if err!=nil{
      panic(err)
  }
  tree ,err:= factory.CreateTree("MainTree");
    if err!=nil{
        panic(err)
    }
  // helper function to print the tree
  BT::printTreeRecursively(tree.rootNode());

  // Tick multiple times, until either FAILURE of SUCCESS is returned
  tree.TickWhileRunning();

}

package main

import (
    "fmt"
    "github.com/gorustyt/go-behavior/core"
)

/** Show the use of the TreeObserver.
 */

// clang-format off

var  xml_text = `(
<root BTCPP_format="4">

    <BehaviorTree ID="MainTree">
        <Sequence>
            <Fallback>
                <AlwaysFailure name="failing_action"/>
                <SubTree ID="SubTreeA" name="mysub"/>
            </Fallback>
            <AlwaysSuccess name="last_action"/>
        </Sequence>
    </BehaviorTree>

    <BehaviorTree ID="SubTreeA">
        <Sequence>
            <AlwaysSuccess name="action_subA"/>
            <SubTree ID="SubTreeB" name="sub_nested"/>
            <SubTree ID="SubTreeB" />
        </Sequence>
    </BehaviorTree>

    <BehaviorTree ID="SubTreeB">
        <AlwaysSuccess name="action_subB"/>
    </BehaviorTree>

</root>
 )`

// clang-format on

func main() {
  var factory core.BehaviorTreeFactory ;

 err:= factory.RegisterBehaviorTreeFromText(xml_text);
 if err!=nil{
     panic(err)
 }
   tree ,err:= factory.CreateTree("MainTree");
    if err!=nil{
        panic(err)
    }
  // Helper function to print the tree.
  BT::printTreeRecursively(tree.rootNode());

  // The purpose of the observer is to save some statistics about the number of times
  // a certain node returns SUCCESS or FAILURE.
  // This is particularly useful to create unit tests and to check if
  // a certain set of transitions happened as expected
  BT::TreeObserver observer(tree);

  // Print the unique ID and the corresponding human readable path
  // Path is also expected to be unique.
  ordered_UID_to_path:=map[uint16]string{}
  for name, uid:=range  observer.pathToUID() {
    ordered_UID_to_path[uid] = name;
  }

  for uid, name:=range  ordered_UID_to_path{
    std::cout << uid << " -> " << name << std::endl;
  }


  tree.TickWhileRunning();

  // You can access a specific statistic, using is full path or the UID
  const auto& last_action_stats = observer.getStatistics("last_action");
  assert(last_action_stats.transitions_count > 0);
  fmt.Println("----------------" )
  // print all the statistics
  for uid, name:=range ordered_UID_to_path {
    const auto& stats = observer.getStatistics(uid);

    std::cout << "[" << name
              << "] \tT/S/F:  " << stats.transitions_count
              << "/" << stats.success_count
              << "/" << stats.failure_count
              << std::endl;
  }

  return 0;
}

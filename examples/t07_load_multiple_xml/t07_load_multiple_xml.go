package main

import (
    "fmt"
    "github.com/gorustyt/go-behavior/core"
)

/** This example show how it is possible to:
 * - load BehaviorTrees from multiple files manually (without the <include> tag)
 * - instantiate a specific tree, instead of the one specified by [main_tree_to_execute]
 */

// clang-format off

var  xml_text_main =`(
<root BTCPP_format="4">
    <BehaviorTree ID="MainTree">
        <Sequence>
            <SaySomething message="starting MainTree" />
            <SubTree ID="SubA"/>
            <SubTree ID="SubB"/>
        </Sequence>
    </BehaviorTree>
</root>  )";

static const char* xml_text_subA = R"(
<root BTCPP_format="4">
    <BehaviorTree ID="SubA">
        <SaySomething message="Executing SubA" />
    </BehaviorTree>
</root>  )`

var  xml_text_subB = `(
<root BTCPP_format="4">
    <BehaviorTree ID="SubB">
        <SaySomething message="Executing SubB" />
    </BehaviorTree>
</root>  )`

// clang-format on



func main() {
  var factory core.BehaviorTreeFactory ;
  factory.RegisterNodeType<DummyNodes::SaySomething>("SaySomething");

  // Register the behavior tree definitions, but do not instantiate them yet.
  // Order is not important.
  factory.RegisterBehaviorTreeFromText(xml_text_subA);
  factory.RegisterBehaviorTreeFromText(xml_text_subB);
  factory.RegisterBehaviorTreeFromText(xml_text_main);

  //Check that the BTs have been registered correctly
  fmt.Printf("Registered BehaviorTrees:")
  for  _,bt_name :=range  factory.RegisteredBehaviorTrees(){
      fmt.Printf(" -%v\n ",bt_name )
  }

  // You can create the MainTree and the subtrees will be added automatically.
  fmt.Println("----- MainTree tick ----")
  main_tree ,err:= factory.CreateTree("MainTree");
  if err!=nil{
      panic(err)
  }
  main_tree.TickWhileRunning();

  // ... or you can create only one of the subtree
  fmt.Println("----- SubA tick ----")

  subA_tree,err := factory.CreateTree("SubA");
  if err!=nil{
      panic(err)
  }
  subA_tree.TickWhileRunning();

}
/* Expected output:

Registered BehaviorTrees:
 - MainTree
 - SubA
 - SubB
----- MainTree tick ----
Robot says: starting MainTree
Robot says: Executing SubA
Robot says: Executing SubB
----- SubA tick ----
Robot says: Executing SubA

*/

package main

import (
    "fmt"
    "github.com/gorustyt/go-behavior/core"
    "time"
)

/** This tutorial will teach you:
 *
 *  - The difference between Sequence and ReactiveSequence
 *
 *  - How to create an asynchronous ActionNode.
*/

// clang-format off

var  xml_text_sequence = `(

 <root BTCPP_format="4" >

     <BehaviorTree ID="MainTree">
        <Sequence name="root">
            <BatteryOK/>
            <SaySomething   message="mission started..." />
            <MoveBase       goal="1;2;3"/>
            <SaySomething   message="mission completed!" />
        </Sequence>
     </BehaviorTree>

 </root>
 )`

var  xml_text_reactive = `(

 <root BTCPP_format="4" >

     <BehaviorTree ID="MainTree">
        <ReactiveSequence name="root">
            <BatteryOK/>
            <Sequence>
                <SaySomething   message="mission started..." />
                <MoveBase       goal="1;2;3"/>
                <SaySomething   message="mission completed!" />
            </Sequence>
        </ReactiveSequence>
     </BehaviorTree>

 </root>
 )`


func  main() {
  var factory  core.BehaviorTreeFactory ;

  factory.RegisterSimpleCondition("BatteryOK", std::bind(CheckBattery));
  factory.RegisterNodeType<MoveBaseAction>("MoveBase");
  factory.RegisterNodeType<SaySomething>("SaySomething");

  // Compare the state transitions and messages using either
  // xml_text_sequence and xml_text_reactive.

  // The main difference that you should notice is:
  //  1) When Sequence is used, the ConditionNode is executed only __once__ because it returns SUCCESS.
  //  2) When ReactiveSequence is used, BatteryOK is executed at __each__ tick()

  for _,xml_text :=range  []string {xml_text_sequence, xml_text_reactive}{
    fmt.Printf("\n------------ BUILDING A NEW TREE ------------\n\n")
     tree ,err:= factory.CreateTreeFromText(xml_text);
    if err!=nil{
        panic(err)
    }
   status := core.NodeStatus_IDLE;

    // Tick the root until we receive either SUCCESS or RUNNING
    // same as: tree.tickRoot(Tree::WHILE_RUNNING)
    fmt.Printf("--- ticking\n")
    status = tree.TickWhileRunning();
    fmt.Printf("--- status:%v \n\n",status.String())
    // If we need to run code between one tick() and the next,
    // we can implement our own while loop
   for (status != core.NodeStatus_SUCCESS){
       fmt.Printf("--- ticking\n")
      status = tree.TickOnce();
      fmt.Printf( "--- status:%v \n\n",status.)

      // if still running, add some wait time
      if (status == core.NodeStatus_RUNNING) {
        tree.Sleep(time.Millisecond*(100));
      }
    }

  }

}

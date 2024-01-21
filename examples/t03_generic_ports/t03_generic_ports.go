package main

import (
  "fmt"
  "github.com/gorustyt/go-behavior/core"
  "strings"
)

/* This tutorial will teach you how to deal with ports when its
 *  type is not std::string.
*/

// We want to be able to use this custom type
type Position2D struct {
 x, y float64
};

// It is recommended (or, in some cases, mandatory) to define a template
// specialization of convertFromString that converts a string to Position2D.
func  convertFromString( str string)*Position2D {
fmt.Printf("Converting string: \"%s\"\n", str);

// real numbers separated by semicolons
 parts := strings.Split(str, ";");
if (len(parts)!= 2) {
panic("invalid input)");
} else{
var  output Position2D;
output.x = convertFromString<double>(parts[0]);
output.y = convertFromString<double>(parts[1]);
return &output;
}
}
func NewCalculateGoal( name string,  config *core.NodeConfig,args ...interface{}) core.ITreeNode{
  return&CalculateGoal{SyncActionNode:core.NewSyncActionNode(name,config)}
}
func (n*CalculateGoal)Tick()core.NodeStatus {
var mygoal=Position2D {1.1, 2.3};
n.SetOutput("goal", mygoal);
return core.NodeStatus_SUCCESS;
}
static PortsList providedPorts()
{
return {OutputPort<Position2D>("goal")};
}
type  CalculateGoal struct {
  *core.SyncActionNode


};

func NewPrintTarget(name string,  config *core.NodeConfig,args ...interface{})*PrintTarget  {
  return &PrintTarget{
    SyncActionNode:core.NewSyncActionNode(name,config),
  }
}
func (*PrintTarget)Tick() core.NodeStatus {
 res := getInput<Position2D>("target");
if (!res) {
panic("error reading port [target]:", res.error());
}
var goal  Position2D
fmt.Printf("Target positions: [ %.1f, %.1f ]\n", goal.x, goal.y);
return core.NodeStatus_SUCCESS;
}
type  PrintTarget struct {
  *core.SyncActionNode



  static PortsList providedPorts()
  {
    // Optionally, a port can have a human readable description
    const char* description = "Simply print the target on console...";
    return {InputPort<Position2D>("target", description)};
  }
};

//----------------------------------------------------------------

/** The tree is a Sequence of 4 actions

*  1) Store a value of Position2D in the entry "GoalPosition"
*     using the action CalculateGoal.
*
*  2) Call PrintTarget. The input "target" will be read from the Blackboard
*     entry "GoalPosition".
*
*  3) Use the built-in action Script to write the key "OtherGoal".
*     A conversion from string to Position2D will be done under the hood.
*
*  4) Call PrintTarget. The input "goal" will be read from the Blackboard
*     entry "OtherGoal".
*/

// clang-format off
var  xml_text = `(

 <root BTCPP_format="4" >
     <BehaviorTree ID="MainTree">
        <Sequence name="root">
            <CalculateGoal   goal="{GoalPosition}" />
            <PrintTarget     target="{GoalPosition}" />
            <Script          code="OtherGoal='-1;3'" />
            <PrintTarget     target="{OtherGoal}" />
        </Sequence>
     </BehaviorTree>
 </root>
 )`

// clang-format on

func main() {


  var factory core.BehaviorTreeFactory ;
  factory.RegisterNodeType<CalculateGoal>("CalculateGoal");
  factory.RegisterNodeType<PrintTarget>("PrintTarget");

  tree ,err:= factory.CreateTreeFromText(xml_text);
  if err!=nil{
    panic(err)
  }
  tree.TickWhileRunning();

  /* Expected output:
 *
    Target positions: [ 1.1, 2.3 ]
    Converting string: "-1;3"
    Target positions: [ -1.0, 3.0 ]
*/

}

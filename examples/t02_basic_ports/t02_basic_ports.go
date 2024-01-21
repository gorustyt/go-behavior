package main

import "github.com/gorustyt/go-behavior/core"

/** This tutorial will teach you how basic input/output ports work.
 *
 * Ports are a mechanism to exchange information between Nodes using
 * a key/value storage called "Blackboard".
 * The type and number of ports of a Node is statically defined.
 *
 * Input Ports are like "argument" of a functions.
 * Output ports are, conceptually, like "return values".
 *
 * In this example, a Sequence of 5 Actions is executed:
 *
 *   - Actions 1 and 4 read the input "message" from a static string.
 *
 *   - Actions 3 and 5 read the input "message" from an entry in the
 *     blackboard called "the_answer".
 *
 *   - Action 2 writes something into the entry of the blackboard
 *     called "the_answer".
*/

// clang-format off
var  xml_text = `(

 <root BTCPP_format="4" >

     <BehaviorTree ID="MainTree">
        <Sequence name="root">
            <SaySomething     message="hello" />
            <SaySomething2    message="this works too" />
            <ThinkWhatToSay   text="{the_answer}"/>
            <SaySomething2    message="{the_answer}" />
        </Sequence>
     </BehaviorTree>

 </root>
 )`

func NewThinkWhatToSay(name string,cfg*core.NodeConfig,args ...interface{})*ThinkWhatToSay  {
return &ThinkWhatToSay{
    SyncActionNode:core.NewSyncActionNode(name,cfg),
}
}
type  ThinkWhatToSay struct {
    *core.SyncActionNode


  // A node having ports MUST implement this STATIC method
  static BT::PortsList providedPorts()
  {
    return {BT::OutputPort<std::string>("text")};
  }
};
// This Action simply write a value in the port "text"
func  (n*ThinkWhatToSay)Tick() core.NodeStatus {
n.SetOutput("text", "The answer is 42");
return core.NodeStatus_SUCCESS;
}
func  main() {


  var factory core.BehaviorTreeFactory ;

  // The class SaySomething has a method called providedPorts() that define the INPUTS.
  // In this case, it requires an input called "message"
  factory.RegisterNodeType<SaySomething>("SaySomething");

  // Similarly to SaySomething, ThinkWhatToSay has an OUTPUT port called "text"
  // Both these ports are std::string, therefore they can connect to each other
  factory.RegisterNodeType<ThinkWhatToSay>("ThinkWhatToSay");

  // SimpleActionNodes can not define their own method providedPorts(), therefore
  // we have to pass the PortsList explicitly if we want the Action to use getInput()
  // or setOutput();
    say_something_ports := core.InputPortWithDefaultValue("message","");
  factory.RegisterSimpleAction("SaySomething2", NewSaySomethingSimple, say_something_ports);

  /* An INPUT can be either a string, for instance:
     *
     *     <SaySomething message="hello" />
     *
     * or contain a "pointer" to a type erased entry in the Blackboard,
     * using this syntax: {name_of_entry}. Example:
     *
     *     <SaySomething message="{the_answer}" />
     */

  tree ,err:= factory.CreateTreeFromText(xml_text);
    if err!=nil{
        panic(err)
    }
  tree.TickWhileRunning();

  /*  Expected output:
     *
        Robot says: hello
        Robot says: this works too
        Robot says: The answer is 42
    *
    * The way we "connect" output ports to input ports is to "point" to the same
    * Blackboard entry.
    *
    * This means that ThinkSomething will write into the entry with key "the_answer";
    * SaySomething and SaySomething will read the message from the same entry.
    *
    */

}

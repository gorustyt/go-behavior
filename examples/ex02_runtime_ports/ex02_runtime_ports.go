package main

import (
    "fmt"
    "github.com/gorustyt/go-behavior/core"
)

// clang-format off
var  xml_text = `(
 <root BTCPP_format="4" >
     <BehaviorTree ID="MainTree">
        <Sequence name="root">
            <ThinkRuntimePort   text="{the_answer}"/>
            <SayRuntimePort     message="{the_answer}" />
        </Sequence>
     </BehaviorTree>
 </root>
 )`

func NewThinkRuntimePort(name string,config *core.NodeConfig,args ...interface{})core.ITreeNode  {
    return &ThinkRuntimePort{
        SyncActionNode:core.NewSyncActionNode(name,config),
    }
}

type ThinkRuntimePort struct {
    *core.SyncActionNode
};
func (n*ThinkRuntimePort)Tick() core.NodeStatus {
n.SetOutput("text", "The answer is 42");
return core.NodeStatus_SUCCESS;
}
func NewSayRuntimePort(name string,config *core.NodeConfig,args ...interface{})core.ITreeNode  {
    return &SayRuntimePort{
        SyncActionNode:core.NewSyncActionNode(name,config),
    }
}
type  SayRuntimePort struct {
    *core.SyncActionNode
};

// You must override the virtual function tick()
 func (n*ThinkRuntimePort)Tick()core.NodeStatus {
auto msg = getInput<std::string>("message");
if (!msg) {
panic("missing required input [message]: ", msg.error());
}
fmt.Println("Robot says: " , msg.value())
return core.NodeStatus_SUCCESS;
}
func  main() {
  var factory core.BehaviorTreeFactory ;

  //-------- register ports that might be defined at runtime --------
  // more verbose way
   think_ports := core.OutputPortWithDefaultValue("text","");
  factory.RegisterBuilder(
      CreateManifest<ThinkRuntimePort>("ThinkRuntimePort", think_ports),
      CreateBuilder<ThinkRuntimePort>());
  // less verbose way
  say_ports :=core.InputPortWithDefaultValue("message","")
  factory.RegisterNodeType("SayRuntimePort",NewSayRuntimePort, say_ports);

  err:=factory.RegisterBehaviorTreeFromText(xml_text);
  if err!=nil{
      panic(err)
  }
  tree,err := factory.CreateTree("MainTree");
  if err!=nil{
      panic(err)
  }
  tree.TickWhileRunning();

  return
}

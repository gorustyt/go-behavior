package main

import (
  "fmt"
  "github.com/gorustyt/go-behavior/core"
  "github.com/gorustyt/go-behavior/examples/sample_nodes"
  "time"
)

/** We are using the same example in Tutorial 5,
 *  But this time we also show how to connect
 */

// A custom structuree that I want to visualize in Groot2
type Position2D struct{
  x float64
   y float64
};

// Allows Position2D to be visualized in Groot2
// You still need BT::RegisterJsonDefinition<Position2D>(PositionToJson)
func PositionToJson( j map[string]interface{},  p *Position2D) {
  j["x"] = p.x;
  j["y"] = p.y;
}

func NewUpdatePosition( name string, config *core.NodeConfig,args ...interface{})core.ITreeNode{
  return &UpdatePosition{
    SyncActionNode:core.NewSyncActionNode(name,config),
  }
}
// Simple Action that updates an instance of Position2D in the blackboard
type  UpdatePosition struct {
  *core.SyncActionNode
  _pos Position2D


  static BT::PortsList providedPorts()
  {
    return {BT::OutputPort<Position2D>("pos")};
  }


};

func(n*UpdatePosition)  Tick() core.NodeStatus {
n._pos.x += 0.2;
n._pos.y += 0.1;
n.SetOutput("pos", n._pos);
return core.NodeStatus_SUCCESS;
}
// clang-format off

var  xml_text = `(
<root BTCPP_format="4">

  <BehaviorTree ID="MainTree">
    <Sequence>
      <Script code="door_open:=false" />
      <UpdatePosition pos="{pos_2D}" />
      <Fallback>
        <Inverter>
          <IsDoorClosed/>
        </Inverter>
        <SubTree ID="DoorClosed" _autoremap="true" door_open="{door_open}"/>
      </Fallback>
      <PassThroughDoor/>
    </Sequence>
  </BehaviorTree>

  <BehaviorTree ID="DoorClosed">
    <Fallback name="tryOpen" _onSuccess="door_open:=true">
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
  factory:=core.NewBehaviorTreeFactory()

  // Nodes registration, as usual
  sample_nodes.RegisterNodes(factory);
  factory.RegisterNodeType("UpdatePosition",NewUpdatePosition);

  // Groot2 editor requires a model of your registered Nodes.
  // You don't need to write that by hand, it can be automatically
  // generated using the following command.
  std::string xml_models = BT::writeTreeNodesModelXML(factory);

 err:= factory.RegisterBehaviorTreeFromText(xml_text);
if err!=nil{
  panic(err)
}
  // Add this to allow Groot2 to visualize your custom type
  BT::RegisterJsonDefinition<Position2D>(PositionToJson);

  tree,err := factory.CreateTree("MainTree");
  if err!=nil{
    panic(err)
  }
  std::cout << "----------- XML file  ----------\n"
            << BT::WriteTreeToXML(tree, false, false)
            << "--------------------------------\n";

  // Connect the Groot2Publisher. This will allow Groot2 to
  // get the tree and poll status updates.
   port := 1667;
  BT::Groot2Publisher publisher(tree, port);

  // Add two more loggers, to save the transitions into a file.
  // Both formats are compatible with Groot2

  // Lightweight serialization
  BT::FileLogger2 logger2(tree, "t12_logger2.btlog");
  // SQLite logger can save multiple sessions into the same database
  bool append_to_database = true;
  BT::SqliteLogger sqlite_logger(tree, "t12_sqlitelog.db3", append_to_database);

  for {
    fmt.Println("Start")
    cross_door.reset();
    tree.TickWhileRunning();
    time.Sleep(2000*time.Millisecond);
  }
}

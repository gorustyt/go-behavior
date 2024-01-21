package main

import (
  "fmt"
  "github.com/gorustyt/go-behavior/core"
  "time"
)

/*
 * In this example we will show how a common design pattern could be implemented.
 * We want to iterate through the elements of a queue, for instance a list of waypoints.
 */

type   Pose2D  struct {
   x, y, theta  float64
};

func NewGenerateWaypoints(name string,cfg *core.NodeConfig,args ...interface{}) *GenerateWaypoints{
  return &GenerateWaypoints{
    SyncActionNode:core.NewSyncActionNode(name,cfg),
  }
}
func (n*GenerateWaypoints)Tick()  core.NodeStatus{
  auto shared_queue = std::make_shared<std::deque<Pose2D>>();
  for  i := 0; i < 5; i++{
  shared_queue->push_back(Pose2D{double(i), double(i), 0});
  }
  n.SetOutput("waypoints", shared_queue);
  return core.NodeStatus_SUCCESS;
}

type  GenerateWaypoints struct {

  *core.SyncActionNode

  static PortsList providedPorts()
  {
    return {OutputPort<SharedQueue<Pose2D>>("waypoints")};
  }
};

func NewPrintNumber(name string,cfg *core.NodeConfig,args ...interface{}) *PrintNumber{
  return &PrintNumber{
    SyncActionNode:core.NewSyncActionNode(name,cfg),
  }
}
func (n*PrintNumber)Tick()  core.NodeStatus{
  var  value float64
  if (getInput("value", value)) {
    fmt.Printf("PrintNumber:%v \n",value)
    return  core.NodeStatus_SUCCESS;
  }
  return  core.NodeStatus_FAILURE;
}
//--------------------------------------------------------------
type  PrintNumber struct {
  *core.SyncActionNode
  static PortsList providedPorts()
  {
    return {InputPort<double>("value")};
  }
};

//--------------------------------------------------------------

/**
 * @brief Simple Action that uses the output of PopFromQueue<Pose2D> or ConsumeQueue<Pose2D>
 */
func NewUseWaypoint(name string,cfg *core.NodeConfig,args ...interface{}) core.ITreeNode {
  return &UseWaypoint{

  }
}
type  UseWaypoint struct {
  *core.ThreadedAction
};

func (n*UseWaypoint)Tick()core.NodeStatus  {
  var wp Pose2D
  if (getInput("waypoint", wp)) {
  time.Sleep(100*time.Millisecond);
  fmt.Printf("Using waypoint:%v/%v ",  wp.x,wp.y)
    return core.NodeStatus_SUCCESS;
  } else{
    return core.NodeStatus_FAILURE;
  }
}

// clang-format off
var  xml_tree = `(
 <root BTCPP_format="4" >
     <BehaviorTree ID="TreeA">
        <Sequence>
            <LoopDouble queue="1;2;3"  value="{number}">
              <PrintNumber value="{number}" />
            </LoopDouble>

            <GenerateWaypoints waypoints="{waypoints}" />
            <LoopPose queue="{waypoints}"  value="{wp}">
              <UseWaypoint waypoint="{wp}" />
            </LoopPose>
        </Sequence>
     </BehaviorTree>
 </root>
 )`

// clang-format on

func  main() {
  var factory core.BehaviorTreeFactory ;
  
  factory.RegisterNodeType<LoopNode<Pose2D>>("LoopPose");

  factory.RegisterNodeType<UseWaypoint>("UseWaypoint");
  core.SetPorts(&UseWaypoint{},core.InputPortWithDefaultValue("waypoint",&Pose2D{}))
  factory.RegisterNodeType<PrintNumber>("PrintNumber");
  factory.RegisterNodeType<GenerateWaypoints>("GenerateWaypoints");

  tree ,err:= factory.CreateTreeFromText(xml_tree);
if err!=nil{
  panic(err)
}

  tree.TickWhileRunning();

  return
}

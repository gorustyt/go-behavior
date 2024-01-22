
package sample_nodes

import (
    "fmt"
    "github.com/gorustyt/go-behavior/core"
    "time"
)


type  GripperInterface struct {

     _opened bool
};

//--------------------------------------
func NewApproachObject( name string, config *core.NodeConfig,args ...interface{} )core.ITreeNode{
    return &ApproachObject{SyncActionNode:core.NewSyncActionNode(name, config)}
}
// Example of custom SyncActionNode (synchronous action)
// without ports.
type  ApproachObject struct {
    *core.SyncActionNode
};

func (t*ApproachObject)Tick()  core.NodeStatus{
    panic("You must override the virtual function tick()")
}

// Example of custom SyncActionNode (synchronous action)
// with an input port.
func NewSaySomething( name string, config *core.NodeConfig,args ...interface{} )core.ITreeNode{
    return &SaySomething{SyncActionNode:core.NewSyncActionNode(name, config)}
}
type SaySomething struct {
    *core.SyncActionNode

    // It is mandatory to define this static method.
    static BT::PortsList providedPorts()
    {
        return{ BT::InputPort<std::string>("message") };
    }
};



func NewSleepNode( name string, config *core.NodeConfig,args ...interface{} )*SleepNode{
    return &SleepNode{
        StatefulActionNode:core.NewStatefulActionNode(name, config),
}
}
// Example os Asynchronous node that use StatefulActionNode as base class
type  SleepNode struct {
    *core.StatefulActionNode
    deadline_ time.Time


    static BT::PortsList providedPorts()
    {
        // amount of milliseconds that we want to sleep
        return{ BT::InputPort<int>("msec") };
    }
}
 func (s*SleepNode)onStart() core.NodeStatus {
 msec := 0;
getInput("msec", msec);
if( msec <= 0 ) {
// no need to go into the RUNNING state
return core.NodeStatus_SUCCESS;
} else {

// once the deadline is reached, we will return SUCCESS.
s.deadline_ = time.Now().Add(time.Duration(msec)*time.Millisecond)
return core.NodeStatus_RUNNING;
}
}

/// method invoked by an action in the RUNNING state.
func (s*SleepNode) OnRunning() core.NodeStatus {
    now:=time.Now()
if now.After(s.deadline_ ){
return core.NodeStatus_SUCCESS;
}else {
return core.NodeStatus_RUNNING;
}
}

func (s*SleepNode)OnHalted() {
// nothing to do here...
fmt.Println("SleepNode interrupted" )
}

func RegisterNodes(factory *core.BehaviorTreeFactory ) {
    var  grip_singleton GripperInterface ;

    factory.RegisterSimpleCondition("CheckBattery", std::bind(CheckBattery));
    factory.RegisterSimpleCondition("CheckTemperature", std::bind(CheckTemperature));
    factory.RegisterSimpleAction("SayHello", std::bind(SayHello));
    factory.RegisterSimpleAction("OpenGripper", std::bind(&GripperInterface::open, &grip_singleton));
    factory.RegisterSimpleAction("CloseGripper", std::bind(&GripperInterface::close, &grip_singleton));
    factory.RegisterNodeType<ApproachObject>("ApproachObject");
    factory.RegisterNodeType<SaySomething>("SaySomething");
}

func  CheckBattery()core.NodeStatus {
fmt.Println("[ Battery: OK ]")
return core.NodeStatus_SUCCESS;
}

func CheckTemperature()core.NodeStatus {
fmt.Println("[ Temperature: OK ]")
return core.NodeStatus_SUCCESS;
}
func SayHello()core.NodeStatus {
fmt.Println("Robot says: Hello World")
return core.NodeStatus_SUCCESS;
}

func (n*GripperInterface)Open()core.NodeStatus{
n._opened = true;
fmt.Println("GripperInterface::open")
return core.NodeStatus_SUCCESS;
}

func (n*GripperInterface)Close()core.NodeStatus {
fmt.Println("GripperInterface::close" )
n._opened = false;
return core.NodeStatus_SUCCESS;
}

func (n*GripperInterface)Tick()core.NodeStatus {
fmt.Printf("ApproachObject:%v " , n.Name())
return core.NodeStatus_SUCCESS;
}

func  (n*SaySomething)Tick()core.NodeStatus {
auto msg = getInput<std::string>("message");
if (!msg) {
throw BT::RuntimeError( "missing required input [message]: ", msg.error() );
}

std::cout << "Robot says: " << msg.value() << std::endl;
return core.NodeStatus_SUCCESS;
}

 func SaySomethingSimple(BT::TreeNode &self)core.NodeStatus {
auto msg = self.getInput<std::string>("message");
if (!msg)
{
panic( "missing required input [message]: ", msg.error() );
}

std::cout << "Robot says: " << msg.value() << std::endl;
return core.NodeStatus_SUCCESS;
}


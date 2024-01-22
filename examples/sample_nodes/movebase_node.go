

package sample_nodes

import (
    "fmt"
    "github.com/gorustyt/go-behavior/core"
    "strings"
    "time"
)

// Custom type
type Pose2D struct {
     x, y, theta float64
};

// Use this to register this function into JsonExporter:
//
// BT::JsonExporter::get().addConverter<Pose2D>();
func  to_json( dest map[string]interface{}, pose* Pose2D ) {
    dest["x"] = pose.x;
    dest["y"] = pose.y;
    dest["theta"] = pose.theta;
}



 func convertFromString(key string)*Pose2D {
    // three real numbers separated by semicolons
   parts := strings.Split(key, ";");
    if (len(parts) != 3) {
        panic("invalid input)");
    } else {
        var output Pose2D ;
        output.x     = convertFromString<double>(parts[0]);
        output.y     = convertFromString<double>(parts[1]);
        output.theta = convertFromString<double>(parts[2]);
        return &output;
    }
}



// Any TreeNode with ports must have a constructor with this signature
func NewMoveBaseAction( name string, config *core.NodeConfig )*MoveBaseAction{
    return &MoveBaseAction{
        StatefulActionNode:core.NewStatefulActionNode(name, config),
    }
}
// This is an asynchronous operation
type   MoveBaseAction struct {
    _goal Pose2D
    _completion_time time.Time
    *core.StatefulActionNode

    // It is mandatory to define this static method.
    static BT::PortsList providedPorts()
    {
        return{ BT::InputPort<Pose2D>("goal") };
    }

};


func (n*MoveBaseAction)OnStart()core.NodeStatus  {
if ( !getInput<Pose2D>("goal", _goal)){
throw BT::RuntimeError("missing required input [goal]");
}
fmt.Printf("[ MoveBase: SEND REQUEST ]. goal: x=%.1f y=%.1f theta=%.1f\n",
_goal.x, _goal.y, _goal.theta);

// We use this counter to simulate an action that takes a certain
// amount of time to be completed (220 ms)
n._completion_time = time.Now().Add(220*time.Millisecond)

return core.NodeStatus_RUNNING;
}

func (n*MoveBaseAction)OnRunning()core.NodeStatus {
// Pretend that we are checking if the reply has been received
// you don't want to block inside this function too much time.
time.Sleep(10*time.Millisecond);

// Pretend that, after a certain amount of time,
// we have completed the operation
if(time.Now().After(n._completion_time) ) {
fmt.Println("[ MoveBase: FINISHED ]")
return core.NodeStatus_SUCCESS;
}
return core.NodeStatus_RUNNING;
}

func (n*MoveBaseAction)OnHalted() {
fmt.Println("[ MoveBase: ABORTED ]")
}

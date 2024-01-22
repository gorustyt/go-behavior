package main

import (
  "fmt"
  "github.com/gorustyt/go-behavior/core"
  "time"
)

// clang-format off

var  xml_text = `(
<root BTCPP_format="4">

  <BehaviorTree ID="MainTree">
    <Sequence>
      <SaySomething name="talk" message="hello world"/>
        <Fallback>
          <AlwaysFailure name="failing_action"/>
          <SubTree ID="MySub" name="mysub"/>
        </Fallback>
        <SaySomething message="before last_action"/>
        <Script code="msg:='after last_action'"/>
        <AlwaysSuccess name="last_action"/>
        <SaySomething message="{msg}"/>
    </Sequence>
  </BehaviorTree>

  <BehaviorTree ID="MySub">
    <Sequence>
      <AlwaysSuccess name="action_subA"/>
      <AlwaysSuccess name="action_subB"/>
    </Sequence>
  </BehaviorTree>

</root>
 )`

var  json_text = `(
{
  "TestNodeConfigs": {
    "MyTest": {
      "async_delay": 2000,
      "return_status": "SUCCESS",
      "post_script": "msg ='message SUBSTITUED'"
    }
  },

  "SubstitutionRules": {
    "mysub/action_*": "TestAction",
    "talk": "TestSaySomething",
    "last_action": "MyTest"
  }
}
 )`

// clang-format on

func main() {

 var  factory core.BehaviorTreeFactory ;

  factory.RegisterNodeType<SaySomething>("SaySomething");

  // We use lambdas and registerSimpleAction, to create
  // a "dummy" node, that we want to create instead of a given one.

  // Simple node that just prints its name and return SUCCESS
  factory.RegisterSimpleAction("DummyAction", func(node core.ITreeNode, status ...core.NodeStatus) core.NodeStatus {
  fmt.Printf("DummyAction substituting:%v ",node.Name())
    return core.NodeStatus_SUCCESS;
  });

  // Action that is meant to substitute SaySomething.
  // It will try to use the input port "message"
  factory.RegisterSimpleAction("TestSaySomething", func(node core.ITreeNode, status ...core.NodeStatus) core.NodeStatus {
    auto msg = self.getInput<std::string>("message");
    if (!msg) {
    throw BT::RuntimeError( "missing required input [message]: ", msg.error() );
    }
    std::cout << "TestSaySomething: " << msg.value() << std::endl;
    return core.NodeStatus_SUCCESS;
  });

  //----------------------------
  // pass "no_sub" as first argument to avoid adding rules
  bool skip_substitution = (argc == 2) && std::string(argv[1]) == "no_sub";

  if(!skip_substitution) {
    // we can use a JSON file to configure the substitution rules
    // or do it manually
   USE_JSON := true;

    if(USE_JSON) {
      factory.loadSubstitutionRuleFromJSON(json_text);
    } else {
      // Substitute nodes which match this wildcard pattern with TestAction
      factory.addSubstitutionRule("mysub/action_*", "TestAction");

      // Substitute the node with name [talk] with TestSaySomething
      factory.addSubstitutionRule("talk", "TestSaySomething");

      // This configuration will be passed to a TestNode
     var  test_config core.TestNodeConfig
      // Convert the node in asynchronous and wait 2000 ms
      test_config.AsyncDelay = (2000)*time.Millisecond;
      // Execute this postcondition, once completed
      test_config.PostScript = "msg ='message SUBSTITUED'";

      // Substitute the node with name [last_action] with a TestNode,
      // configured using test_config
      factory.AddSubstitutionRule("last_action", test_config);
    }
  }

  err:=factory.RegisterBehaviorTreeFromText(xml_text);
  if err!=nil{
      panic(err)
  }
  // During the construction phase of the tree, the substitution
  // rules will be used to instantiate the test nodes, instead of the
  // original ones.
 tree,err := factory.CreateTree("MainTree");
 if err!=nil{
   panic(err)
 }
  tree.TickWhileRunning();
}

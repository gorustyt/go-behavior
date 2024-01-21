package main

import "github.com/gorustyt/go-behavior/core"

/*
 * Demonstrate how to use a SubTree Model (since version 4.4)
 */

/**
 * You can optionally add a model to a SubTrees, in this case, "MySub".
 * We are telling the factory that the callee should remap
 * two mandatory inputs, and two outputs:
 *
 * - sub_in_value (that has the default value 42)
 * - sub_in_name (no default)
 * - sub_out_result (default remapping to port {output})
 * - sub_out_state (no default)
 *
 * The callee (parent tree, including the subtree) MUST specify those
 * remapping which have no default value.
 */

// clang-format off
var xml_subtree = `(
<root BTCPP_format="4">

  <TreeNodesModel>
    <SubTree ID="MySub">
      <input_port name="sub_in_value" default="42"/>
      <input_port name="sub_in_name"/>
      <output_port name="sub_out_result" default="{out_result}"/>
      <output_port name="sub_out_state"/>
    </SubTree>
  </TreeNodesModel>

  <BehaviorTree ID="MySub">
    <Sequence>
      <ScriptCondition code="sub_in_value==42 && sub_in_name=='john'" />
      <Script code="sub_out_result:=69; sub_out_state:='ACTIVE'" />
    </Sequence>
  </BehaviorTree>
</root>
 )`

/**
 * Here, when calling "MySub", only `sub_in_name` and `sub_out_state` are explicitly
 * remapped. We will use the default values for the other two.
 */

var xml_maintree = `(
<root BTCPP_format="4">

  <BehaviorTree ID="MainTree">
    <Sequence>
      <Script code="in_name:= 'john' "/>
      <SubTree ID="MySub" sub_in_name="{in_name}"
                          sub_out_state="{out_state}"/>
      <ScriptCondition code=" out_result==69 && out_state=='ACTIVE' " />
    </Sequence>
  </BehaviorTree>

</root>
 )`

// clang-format on

func main() {
	var factory core.BehaviorTreeFactory
	err := factory.RegisterBehaviorTreeFromText(xml_subtree)
	if err != nil {
		panic(err)
	}
	err = factory.RegisterBehaviorTreeFromText(xml_maintree)
	if err != nil {
		panic(err)
	}
	tree, err := factory.CreateTree("MainTree")
	if err != nil {
		panic(err)
	}
	tree.TickWhileRunning()

	// We expect the sequence to be successful.

	// The full remapping was:
	//
	// - sub_in_name    <-> {in_name}      // specified by callee 'MainTree'
	// - sub_in_value   <-> 42             // default remapping, see model
	// - sub_out_result <-> {out_result}   // default remapping, see model
	// - sub_out_state  <-> {out_state}    // specified by callee 'MainTree'
}

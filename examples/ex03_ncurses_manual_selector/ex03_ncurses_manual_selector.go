package main

/* Try also
*      <ManualSelector repeat_last_selection="1">
*  to see the difference.
*/

var  xml_text= `(
 <root BTCPP_format="4" >
     <BehaviorTree ID="MainTree">
        <Repeat num_cycles="3">
            <ManualSelector repeat_last_selection="0">
                <SaySomething name="Option1"    message="Option1" />
                <SaySomething name="Option2"    message="Option2" />
                <SaySomething name="Option3"    message="Option3" />
                <SaySomething name="Option4"    message="Option4" />
                <ManualSelector name="YouChoose" />
            </ManualSelector>
        </Repeat>
     </BehaviorTree>
 </root>
 )`

func main() {
  BehaviorTreeFactory factory;
  factory.registerNodeType<DummyNodes::SaySomething>("SaySomething");

  auto tree = factory.createTreeFromText(xml_text);
  auto ret = tree.tickWhileRunning();

  std::cout << "Result: " << ret << std::endl;

  return 0;
}

package main

import "github.com/gorustyt/go-behavior/core"

// clang-format off
    var xml_text = `(
 <root BTCPP_format="4">
     <BehaviorTree>
        <Sequence>
            <Script code=" msg:='hello world' " />
            <Script code=" A:=THE_ANSWER; B:=3.14; color:=RED " />
            <Precondition if="A>B && color != BLUE" else="FAILURE">
                <Sequence>
                  <SaySomething message="{A}"/>
                  <SaySomething message="{B}"/>
                  <SaySomething message="{msg}"/>
                  <SaySomething message="{color}"/>
                </Sequence>
            </Precondition>
        </Sequence>
     </BehaviorTree>
 </root>
 )`

// clang-format on
const (
    RED = 1
        BLUE = 2
    GREEN = 3
)
func  main() {
  var factory core.BehaviorTreeFactory ;
  factory.RegisterNodeType<DummyNodes::SaySomething>("SaySomething");


  // We can add these enums to the scripting language
  factory.RegisterScriptingEnums<Color>();

  // Or we can do it manually
  factory.RegisterScriptingEnum("THE_ANSWER", 42);

   tree ,err:= factory.CreateTreeFromText(xml_text);
   if err!=nil{
       panic(err)
   }
  tree.TickWhileRunning();

}

/* Expected output:

Robot says: 42.000000
Robot says: 3.140000
Robot says: hello world
Robot says: 1.000000

*/

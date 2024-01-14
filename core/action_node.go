package core

type ActionNodeBase struct {
	*LeafNode
}

type SyncActionNode struct {
	*ActionNodeBase
}

type SimpleActionNode struct {
	*SyncActionNode
}

type ThreadedAction struct {
	*ActionNodeBase
}

type StatefulActionNode struct {
	*ActionNodeBase
}

type CoroActionNode struct {
	*ActionNodeBase
}

package cluster

// Kind 集群实例类型
type Kind string

const (
	Gate Kind = "gate" // 网关服
	Node Kind = "node" // 节点服
)

// State 集群实例状态
type State string

const (
	Work State = "work" // 工作（节点正常工作，可以分配更多玩家到该节点）
	Busy State = "busy" // 繁忙（节点资源紧张，不建议分配更多玩家到该节点上）
	Hang State = "hang" // 挂起（节点即将关闭，正处于资源回收中）
	Shut State = "shut" // 关闭（节点已经关闭，无法正常访问该节点）
)

// Event 事件
type Event int16

const (
	Send       Event = 0 //  正常收发数据
	Connect    Event = 1 //   打开连接
	Reconnect  Event = 2 //   断线重连
	Disconnect Event = 3 //   断开连接
)

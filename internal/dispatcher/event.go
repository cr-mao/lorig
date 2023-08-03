package dispatcher

import (
	"github.com/cr-mao/lorig/cluster"
)

type Event struct {
	abstract
	event cluster.Event // 路由ID
}

func newEvent(dispatcher *Dispatcher, event cluster.Event) *Event {
	return &Event{
		abstract: abstract{
			dispatcher:  dispatcher,
			endpointMap: make(map[string]*serviceEndpoint),
			endpointArr: make([]*serviceEndpoint, 0),
		},
		event: event,
	}
}

// Event 获取事件
func (e *Event) Event() cluster.Event {
	return e.event
}

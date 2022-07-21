package amethyst

import (
	"errors"
	"fmt"
	"github.com/phrara/amethyst/config"
	"github.com/phrara/amethyst/interval"
)

type routerSet map[uint32]Router

type RouterManager struct {
	routers routerSet
	wp      *interval.WorkerPool
}

func NewRouterManager() *RouterManager {
	rm := &RouterManager{
		routers: make(routerSet, 10),
	}
	wp := interval.NewWorkPool(config.Global.MaxPoolSize, config.Global.MaxQueLen, rm, interval.RR)
	wp.Start()
	rm.wp = wp
	return rm
}

func (m *RouterManager) HandleRequest(request *Request) error {
	_, ok := m.routers[request.Tag()]
	if !ok {
		return errors.New(fmt.Sprintf("tag <%d> unrecognized", request.Tag()))
	}

	// put the request into WorkerPool's TaskQueue
	m.wp.AppendTask(request, int(request.Conn().ConnID()))
	return nil
}

func (m *RouterManager) Handle(task any) {
	req := task.(*Request)
	r, _ := m.routers[req.Tag()]
	r.Before(req)
	r.Cope(req)
	r.After(req)
}

func (m *RouterManager) AppendRouter(tag uint32, r Router) {
	m.routers[tag] = r
}

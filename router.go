package amethyst

type Router interface {
	Before(req *Request)
	Cope(req *Request)
	After(req *Request)
}

// DefaultRouter
// Users should extend this struct to design their own routers
type DefaultRouter struct {
}

func (d *DefaultRouter) Before(req *Request) {

}

func (d *DefaultRouter) Cope(req *Request) {

}

func (d *DefaultRouter) After(req *Request) {

}

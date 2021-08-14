package router

type RouterItem struct {
	Route      string
	Controller ILionController
}

// Register is to registe these URL into default route,
// then when linsen the default route.
func (r *LionRouter) Register(c *BaseController) {
	c.router = r
	c.route = addRoute(r.route, c.route)
	item := &RouterItem{
		Route:      c.route,
		Controller: c,
	}
	r.t.add(item)
}

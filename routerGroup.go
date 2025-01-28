package gfx

// RouteGroup represents a group of routes
type RouteGroup struct {
	engine     *GFXEngine
	parent     *RouteGroup
	basePath   string
	middleware []MiddlewareFunc
}

func (rg *RouteGroup) Group(basePath string) *RouteGroup {
	newGroup := &RouteGroup{
		engine:   rg.engine,
		parent:   rg,
		basePath: basePath,
	}

	newGroup.middleware = append(newGroup.middleware, rg.middleware...)

	return newGroup
}

// UseMiddleware adds a middleware to the group
func (rg *RouteGroup) UseMiddleware(middleware MiddlewareFunc) {
	rg.middleware = append(rg.middleware, middleware)
}

// Get adds a GET route to the group
func (rg *RouteGroup) Get(path string, handler HandlerFunc) {
	rg.engine.addRoute("GET", path, handler, rg.middleware, rg)
}

// Post adds a POST route to the group
func (rg *RouteGroup) Post(path string, handler HandlerFunc) {
	rg.engine.addRoute("POST", path, handler, rg.middleware, rg)
}

// Patch adds a PATCH route to the engine
func (rg *RouteGroup) Patch(path string, handler HandlerFunc) {
	rg.engine.addRoute("PATCH", path, handler, rg.middleware, rg)
}

// Put adds a PUT route to the group
func (rg *RouteGroup) Put(path string, handler HandlerFunc) {
	rg.engine.addRoute("PUT", path, handler, rg.middleware, rg)
}

// Delete adds a DELETE route to the group
func (rg *RouteGroup) Delete(path string, handler HandlerFunc) {
	rg.engine.addRoute("DELETE", path, handler, rg.middleware, rg)
}

// Options adds a OPTIONS route to the group
func (rg *RouteGroup) Options(path string, handler HandlerFunc) {
	rg.engine.addRoute("OPTIONS", path, handler, rg.middleware, rg)
}

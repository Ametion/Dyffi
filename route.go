package gfx

// Route represents a HTTP route
type Route struct {
	method      string
	handler     HandlerFunc
	middleware  []MiddlewareFunc
	parts       []string
	paramsIndex []int
}

package gfx

import (
	"fmt"
	"net/http"
	"strings"
)

// GFXEngine represents the main engine
type GFXEngine struct {
	routes         []Route
	middleware     []MiddlewareFunc
	development    bool
	isCors         bool
	allowedMethods []string
	allowedOrigins []string
	allowedHeaders []string
}

// NewGFXEngine creates a new GFXEngine
func NewGFXEngine() *GFXEngine {
	return &GFXEngine{
		development: false,
	}
}

// ServeHTTP handles the request
func (g *GFXEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	methodAllowed := false

	for _, allowedMethod := range g.allowedMethods {
		if r.Method == allowedMethod {
			methodAllowed = true
			break
		}
	}

	for _, allowedOrigin := range g.allowedOrigins {
		if (allowedOrigin == origin || allowedOrigin == "*") && methodAllowed {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(g.allowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(g.allowedHeaders, ", "))
			break
		}
	}

	if !methodAllowed {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("Method not allowed"))
		if err != nil {
			return
		}
		return
	}

	requestParts := strings.Split(r.URL.Path, "/")
	statusCode := http.StatusNotFound // Default status code
	wrappedWriter := &LoggingResponseWriter{
		ResponseWriter: w,
		development:    g.development,
		statusCode:     http.StatusOK,
		method:         r.Method,
		route:          r.URL.Path,
	}

	for _, route := range g.routes {
		if r.Method == route.method && len(requestParts) == len(route.parts) {
			if c := g.processRoute(route, wrappedWriter, r, requestParts); c != nil {
				c.Next()
				statusCode = http.StatusOK
				break
			}
		}
	}

	if statusCode == http.StatusNotFound {
		http.NotFound(wrappedWriter, r)
	}
}

func (g *GFXEngine) IsDevelopment() {
	g.development = true
}

// Get adds a GET route to the engine
func (g *GFXEngine) Get(path string, handler HandlerFunc) {
	g.addRoute("GET", path, handler, nil, nil)
}

// Post adds a POST route to the engine
func (g *GFXEngine) Post(path string, handler HandlerFunc) {
	g.addRoute("POST", path, handler, nil, nil)
}

// Patch adds a PATCH route to the engine
func (g *GFXEngine) Patch(path string, handler HandlerFunc) {
	g.addRoute("PATCH", path, handler, nil, nil)
}

// Put adds a PUT route to the engine
func (g *GFXEngine) Put(path string, handler HandlerFunc) {
	g.addRoute("PUT", path, handler, nil, nil)
}

// Delete adds a DELETE route to the engine
func (g *GFXEngine) Delete(path string, handler HandlerFunc) {
	g.addRoute("DELETE", path, handler, nil, nil)
}

// Options adds a OPTIONS route to the engine
func (g *GFXEngine) Options(path string, handler HandlerFunc) {
	g.addRoute("OPTIONS", path, handler, nil, nil)
}

// Group creates a new RouteGroup
func (g *GFXEngine) Group(basePath string) *RouteGroup {
	return &RouteGroup{
		engine:     g,
		basePath:   basePath,
		middleware: g.middleware,
	}
}

// Func which use for add middleware to whole engine
func (g *GFXEngine) UseMiddleware(middleware MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware)
}

// Run starts the web
func (g *GFXEngine) Run(addr string) error {
	fmt.Println("GFXEngine starting with the following routes:")
	for _, route := range g.routes {
		if route.method != "OPTIONS" {
			path := strings.Join(route.parts, "/")
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			fmt.Printf("%s %s\n", route.method, path)
		}
	}
	fmt.Printf("Listening on %s\n", addr)
	return http.ListenAndServe(addr, g)
}

func (g *GFXEngine) processRoute(route Route, w http.ResponseWriter, r *http.Request, requestParts []string) *Context {
	params := make(map[string]string)
	match := true

	for _, i := range route.paramsIndex {
		params[route.parts[i]] = requestParts[i]
	}

	for i, part := range requestParts {
		if i >= len(route.parts) || (part != route.parts[i] && !contains(route.paramsIndex, i)) {
			match = false
			break
		}
	}

	if match {
		c := &Context{
			Writer:     w,
			Request:    r,
			params:     params,
			Headers:    r.Header,
			index:      0,
			middleware: append(route.middleware, handlerToMiddleware(route.handler)),
		}
		return c
	}

	return nil
}

// addRoute adds a route to the engine
func (g *GFXEngine) addRoute(method string, path string, handler HandlerFunc, middleware []MiddlewareFunc, group *RouteGroup) {
	fullPath := path
	var fullMiddleware []MiddlewareFunc

	// Handle Group-based path and middleware merging
	if group != nil {
		// If the group has a parent, walk up the hierarchy and prepend each parent's basePath
		for p := group; p != nil; p = p.parent {
			fullPath = p.basePath + fullPath
		}

		// Prepend parent middleware in order from topmost parent to current group
		for p := group; p != nil; p = p.parent {
			fullMiddleware = append(p.middleware, fullMiddleware...)
		}
	}

	// Merge middleware from the engine itself
	fullMiddleware = append(g.middleware, fullMiddleware...)

	// Finally, include the route-specific middleware
	fullMiddleware = append(fullMiddleware, middleware...)

	// Filter out duplicate middleware
	fullMiddleware = removeDuplicateMiddleware(fullMiddleware)

	// Split the path into its components
	parts := strings.Split(fullPath, "/")
	var paramsIndex []int

	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			paramsIndex = append(paramsIndex, i)
			parts[i] = part[1:] // Remove the ":" prefix
		}
	}

	// Create and add the new Route
	route := Route{
		method:      method,
		handler:     handler,
		middleware:  fullMiddleware,
		parts:       parts,
		paramsIndex: paramsIndex,
	}
	g.routes = append(g.routes, route)
}

// Helper function to remove duplicate middleware
func removeDuplicateMiddleware(middleware []MiddlewareFunc) []MiddlewareFunc {
	seen := make(map[string]struct{})
	result := []MiddlewareFunc{}

	for _, m := range middleware {
		key := fmt.Sprintf("%p", m) // Use the memory address as a unique key
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, m)
		}
	}

	return result
}

func handlerToMiddleware(h HandlerFunc) MiddlewareFunc {
	return func(c *Context) {
		h(c)
	}
}

func contains(arr []int, value int) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

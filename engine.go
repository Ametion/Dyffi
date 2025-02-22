package dyffi

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Engine represents the main engine of the web server
type Engine struct {
	routes         []Route
	middleware     []MiddlewareFunc
	development    bool
	isCors         bool
	AllowedMethods []string
	allowedOrigins []string
	AllowedHeaders []string
}

// NewDyffiEngine creates a new Engine
func NewDyffiEngine() *Engine {
	return &Engine{
		development: false,
	}
}

// ServeHTTP handles the request
func (g *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	methodAllowed := false

	for _, allowedMethod := range g.AllowedMethods {
		if r.Method == allowedMethod {
			methodAllowed = true
			break
		}
	}

	for _, allowedOrigin := range g.allowedOrigins {
		if (allowedOrigin == origin || allowedOrigin == "*") && methodAllowed {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(g.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(g.AllowedHeaders, ", "))
			break
		}
	}

	requestParts := strings.Split(r.URL.Path, "/")
	statusCode := http.StatusNotFound
	params := make(map[string]pathPart)

	for _, route := range g.routes {
		if r.Method == route.method && len(requestParts) == len(route.parts) {
			if ctx := g.processRoute(route, w, r, requestParts); ctx != nil {
				statusCode = http.StatusOK
				params = ctx.params
				g.logRequest(r.Method, statusCode, r.URL.Path, params)
				ctx.Next()
				return
			}
		}
	}

	if !methodAllowed {
		statusCode = http.StatusMethodNotAllowed
		w.WriteHeader(statusCode)
		w.Write([]byte("Method not allowed"))
	} else {
		http.NotFound(w, r)
	}
}

func (g *Engine) IsDevelopment() {
	g.development = true
}

// Get adds a GET route to the engine
func (g *Engine) Get(path string, handler HandlerFunc) {
	g.addRoute("GET", path, handler, nil, nil)
}

// Post adds a POST route to the engine
func (g *Engine) Post(path string, handler HandlerFunc) {
	g.addRoute("POST", path, handler, nil, nil)
}

// Patch adds a PATCH route to the engine
func (g *Engine) Patch(path string, handler HandlerFunc) {
	g.addRoute("PATCH", path, handler, nil, nil)
}

// Put adds a PUT route to the engine
func (g *Engine) Put(path string, handler HandlerFunc) {
	g.addRoute("PUT", path, handler, nil, nil)
}

// Delete adds a DELETE route to the engine
func (g *Engine) Delete(path string, handler HandlerFunc) {
	g.addRoute("DELETE", path, handler, nil, nil)
}

// Options adds a OPTIONS route to the engine
func (g *Engine) Options(path string, handler HandlerFunc) {
	g.addRoute("OPTIONS", path, handler, nil, nil)
}

// Group creates a new RouteGroup
func (g *Engine) Group(basePath string) *RouteGroup {
	return &RouteGroup{
		engine:     g,
		basePath:   basePath,
		middleware: g.middleware,
	}
}

// UseMiddleware Func which use for add middleware to whole engine
func (g *Engine) UseMiddleware(middleware MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware)
}

// Run starts the web
func (g *Engine) Run(addr string) error {
	fmt.Println("\n\033[1;32mDyffi Engine starting with the following routes:\033[0m\n")

	if len(g.routes) == 0 {
		fmt.Println("\033[1;31mNo routes registered!\033[0m") // Red warning if no routes exist
	} else {
		for _, route := range g.routes {
			if route.method != "OPTIONS" {
				path := formatRoute(route.parts, route.paramsIndex)
				fmt.Printf("  \033[1;35m%-7s\033[0m \033[1;34m%s\033[0m\n", route.method, path)
			}
		}
	}

	fmt.Printf("\n\033[1;36mListening on %s\033[0m\n\n", addr) // Cyan color for "Listening"
	return http.ListenAndServe(addr, g)
}

func (g *Engine) processRoute(route Route, w http.ResponseWriter, r *http.Request, requestParts []string) *Context {
	params := make(map[string]pathPart)

	for _, part := range route.parts {
		for _, i := range route.paramsIndex {
			if part.isParam {
				temp := params[part.part]
				temp.part = part.part
				temp.isParam = true
				temp.index = part.index

				if part.index == i {
					temp.value = requestParts[i]

					if part.regexPattern != "" {
						temp.regexPattern = "^" + part.regexPattern + "$"

						matched, regexErr := regexp.MatchString(temp.regexPattern, requestParts[i])

						if !matched || regexErr != nil {
							w.WriteHeader(http.StatusBadRequest)
							w.Write([]byte("Regex not matched"))
							return nil
						}
					}
				}

				params[part.part] = temp
			}
		}
	}

	fmt.Println(params)

	// Collect all middleware (engine -> group -> route)
	middlewareQueue := []MiddlewareFunc{}
	middlewareQueue = append(middlewareQueue, g.middleware...)                    // Engine-level middleware
	middlewareQueue = append(middlewareQueue, route.middleware...)                // Route-specific middleware
	middlewareQueue = append(middlewareQueue, handlerToMiddleware(route.handler)) // Final handler

	// Create Context with middleware queue
	ctx := &Context{
		writer:     w,
		request:    r,
		params:     params,
		Headers:    r.Header,
		index:      0,
		middleware: middlewareQueue, // Middleware queue including the handler
	}
	return ctx
}

// addRoute adds a route to the engine
func (g *Engine) addRoute(method string, path string, handler HandlerFunc, middleware []MiddlewareFunc, group *RouteGroup) {
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

	pathParts := []pathPart{}

	var paramsIndex []int

	for i, part := range parts {
		pathPart := pathPart{
			index: i,
			part:  part,
		}

		staticPart, regexPattern := extractPartAndRegex(part)
		pathPart.part = staticPart

		if regexPattern != "" {
			pathPart.regexPattern = regexPattern
		}

		if strings.HasPrefix(part, ":") {
			pathPart.isParam = true
			paramsIndex = append(paramsIndex, i)
			parts[i] = part[1:]
			pathPart.part = part[1:]
		}

		pathParts = append(pathParts, pathPart)
	}

	// Create and add the new Route
	route := Route{
		method:      method,
		handler:     handler,
		middleware:  fullMiddleware,
		parts:       pathParts,
		paramsIndex: paramsIndex,
	}
	g.routes = append(g.routes, route)
}

func extractPartAndRegex(part string) (string, string) {
	regexSymbols := `.*+?^$|{}[]()`

	for i, char := range part {
		if strings.ContainsRune(regexSymbols, char) {
			return part[:i], part[i:]
		}
	}
	return part, ""
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

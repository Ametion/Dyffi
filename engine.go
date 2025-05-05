package dyffi

import (
	"fmt"
	dyffiBroker "github.com/Ametion/dyffi/broker"
	"github.com/graphql-go/graphql"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

// Engine represents the main engine of the web server
type Engine struct {
	routes         []Route
	middleware     []MiddlewareFunc
	development    bool
	isCors         bool
	graphqlUsage   bool
	graphqlSchemas GraphQLSchemas
	AllowedMethods []string
	allowedOrigins []string
	AllowedHeaders []string

	broker       *dyffiBroker.Queue
	brokerConfig dyffiBroker.BrokerConfig
	authConf     any
	services map[reflect.Type]reflect.Value
}

// NewDyffiEngine creates a new Engine
func NewDyffiEngine() *Engine {
	return &Engine{
		development:  false,
		isCors:       false,
		graphqlUsage: false,
		services:     make(map[reflect.Type]reflect.Value),
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

	for _, route := range g.routes {
		if r.Method == route.method && len(requestParts) == len(route.parts) {
			if g.matchRoute(route, requestParts) {
				if ctx := g.processRoute(route, w, r, requestParts); ctx != nil {
					statusCode = http.StatusOK
					g.logRequest(r.Method, statusCode, r.URL.Path, ctx.params)
					ctx.Next()
					return
				}
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

func (g *Engine) matchRoute(route Route, requestParts []string) bool {
	// Must have same number of segments
	if len(route.parts) != len(requestParts) {
		return false
	}

	for i, routePart := range route.parts {
		if routePart.isParam {
			// Parameter segment always matches, but no static check
			continue
		}

		// If it's a static segment, it must match exactly
		if routePart.part != requestParts[i] {
			return false
		}
	}
	return true
}

func (g *Engine) Provide(svc any) {
    v := reflect.ValueOf(svc)
    t := v.Type()

    if t.Kind() == reflect.Struct {
        ptr := reflect.New(t)
        ptr.Elem().Set(v)
        g.services[ptr.Type()] = ptr
        return
    }

    g.services[t] = v
}

// SetDevelopment sets development mode
func (g *Engine) SetDevelopment() {
	g.development = true
}

// Get adds a GET route to the engine
func (g *Engine) Get(path string, handler any) {
	g.addInjectedRoute("GET", path, handler)
}

// Post adds a POST route to the engine
func (g *Engine) Post(path string, handler any) {
	g.addInjectedRoute("POST", path, handler)
}

// Patch adds a PATCH route to the engine
func (g *Engine) Patch(path string, handler any) {
	g.addInjectedRoute("PATCH", path, handler)
}

// Put adds a PUT route to the engine
func (g *Engine) Put(path string, handler any) {
	g.addInjectedRoute("PUT", path, handler)
}

// Delete adds a DELETE route to the engine
func (g *Engine) Delete(path string, handler any) {
	g.addInjectedRoute("DELETE", path, handler)
}

// Options adds a OPTIONS route to the engine
func (g *Engine) Options(path string, handler any) {
	g.addInjectedRoute("OPTIONS", path, handler)
}

// GraphQLModel creates a GraphQL route
func (g *Engine) GraphQLModel(modelName string, model interface{}, resolvers GraphQLResolvers) {
	g.graphqlUsage = true

	if g.graphqlSchemas == nil {
		g.graphqlSchemas = make(GraphQLSchemas)
	}

	// Store the model and its resolvers
	g.graphqlSchemas[modelName] = &GraphQLModel{
		ModelType: reflect.TypeOf(model),
		Resolvers: resolvers,
	}

	// Auto-register the /graphql route the first time a model is added
	if len(g.graphqlSchemas) == 1 {
		g.Post("/graphql", func(c *Context) {
			g.processGraphQLRequest(c)
		})
	}
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

func (g *Engine) UseBroker(brokerConfig dyffiBroker.BrokerConfig) {
	g.brokerConfig = brokerConfig

	broker, err := dyffiBroker.CreateNewBroker(brokerConfig)

	if err != nil {
		fmt.Println("Error creating broker:", err)
		return
	}

	g.broker = &broker
}

// Run starts the web
func (g *Engine) Run(addr string) error {
	fmt.Println("\n\033[1;32mDyffi Engine starting with the following routes:\033[0m\n")

	if len(g.routes) == 0 && len(g.graphqlSchemas) == 0 {
		fmt.Println("\033[1;31mNo routes registered!\033[0m")
	} else {
		for _, route := range g.routes {
			if route.method != "OPTIONS" {
				path := formatRoute(route.parts, route.paramsIndex)
				fmt.Printf("  \033[1;35m%-7s\033[0m \033[1;34m%s\033[0m\n", route.method, path)
			}
		}

		// Print registered GraphQL models
		if len(g.graphqlSchemas) > 0 {
			fmt.Println("\n\033[1;33mGraphQL Models Registered:\033[0m")
			for modelName := range g.graphqlSchemas {
				fmt.Printf("  \033[1;36m/graphql\033[0m → Model: \033[1;32m%s\033[0m\n", modelName)
			}
		}
	}

	fmt.Printf("\n\033[1;36mListening on %s\033[0m\n\n", addr) // Cyan color for "Listening"
	return http.ListenAndServe(addr, g)
}

func (g *Engine) wrapResolver(fn interface{}) graphql.FieldResolveFn {
    hv := reflect.ValueOf(fn)
    ht := hv.Type()
    ctxT := reflect.TypeOf(QLContext{})

    if ht.Kind() != reflect.Func || ht.NumIn() == 0 || ht.In(0) != ctxT {
        panic("resolver must be func(QLContext, …) (T, error)")
    }
    if ht.NumOut() != 2 || !ht.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
        panic("resolver must return (T, error)")
    }

    // PRE-RESOLVE extra deps beyond QLContext
    deps := make([]reflect.Value, ht.NumIn()-1)
    for i := 1; i < ht.NumIn(); i++ {
        depT := ht.In(i)
        svc, ok := g.services[depT]
        if !ok && depT.Kind() == reflect.Interface {
            for _, cand := range g.services {
                if cand.Type().Implements(depT) {
                    svc = cand
                    ok = true
                    break
                }
            }
        }
        if !ok {
            panic(fmt.Sprintf("no service registered for resolver dep %v", depT))
        }
        deps[i-1] = svc
    }

    return func(p graphql.ResolveParams) (interface{}, error) {
        qctx := QLContext{Params: p}
        args := make([]reflect.Value, 1+len(deps))
        args[0] = reflect.ValueOf(qctx)
        copy(args[1:], deps)

        out := hv.Call(args)
        result := out[0].Interface()
        var err error
        if e, _ := out[1].Interface().(error); e != nil {
            err = e
        }
        return result, err
    }
}

// generateGraphQLSchema builds your schema and plugs in wrapResolver for each resolver
func (g *Engine) generateGraphQLSchema() *graphql.Schema {
    queryFields := graphql.Fields{}
    mutationFields := graphql.Fields{}

    for modelName, model := range g.graphqlSchemas {
        modelType := model.ModelType

        // build object type
        fields := graphql.Fields{}
        for i := 0; i < modelType.NumField(); i++ {
            f := modelType.Field(i)
            fields[f.Name] = &graphql.Field{Type: g.goTypeToGraphQL(f.Type)}
        }
        objectType := graphql.NewObject(graphql.ObjectConfig{
            Name:   modelName,
            Fields: fields,
        })

        // Query resolver
        if model.Resolvers.Query != nil {
            queryFields["get"+modelName] = &graphql.Field{
                Type:    objectType,
                Args:    graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.Int}},
                Resolve: g.wrapResolver(model.Resolvers.Query),
            }
        }

        // Create mutation
        if model.Resolvers.MutationCreate != nil {
            mutationFields["create"+modelName] = &graphql.Field{
                Type:    objectType,
                Args:    fieldsToArgs(fields),
                Resolve: g.wrapResolver(model.Resolvers.MutationCreate),
            }
        }

        // Delete mutation
        if model.Resolvers.MutationDelete != nil {
            mutationFields["delete"+modelName] = &graphql.Field{
                Type: graphql.NewObject(graphql.ObjectConfig{
                    Name: "DeleteResponse",
                    Fields: graphql.Fields{
                        "message": &graphql.Field{Type: graphql.String},
                        "ID":      &graphql.Field{Type: graphql.Int},
                    },
                }),
                Args:    graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.Int}},
                Resolve: g.wrapResolver(model.Resolvers.MutationDelete),
            }
        }
    }

    schema, _ := graphql.NewSchema(graphql.SchemaConfig{
        Query:    graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: queryFields}),
        Mutation: graphql.NewObject(graphql.ObjectConfig{Name: "Mutation", Fields: mutationFields}),
    })
    return &schema
}

// processGraphQLRequest processes a GraphQL request
func (g *Engine) processGraphQLRequest(c *Context) {
	var request struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}

	if err := c.SetBody(&request); err != nil {
		c.SendJSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	// Use merged schema
	schema := g.generateGraphQLSchema()
	result := graphql.Do(graphql.Params{
		Schema:         *schema,
		RequestString:  request.Query,
		VariableValues: request.Variables,
	})

	c.SendJSON(http.StatusOK, result)
}

// goTypeToGraphQL converts a Go type to a GraphQL scalar
func (g *Engine) goTypeToGraphQL(t reflect.Type) *graphql.Scalar {
	switch t.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		return graphql.Int
	case reflect.String:
		return graphql.String
	case reflect.Bool:
		return graphql.Boolean
	default:
		return graphql.String
	}
}

// fieldToArgs converts a map of fields to a GraphQL argument map
func fieldsToArgs(fields graphql.Fields) graphql.FieldConfigArgument {
	args := graphql.FieldConfigArgument{}
	for name, field := range fields {
		args[name] = &graphql.ArgumentConfig{Type: field.Type}
	}
	return args
}

// processRoute processes a route
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
						temp.regexPattern = part.regexPattern

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

	// Collect all middleware (engine -> group -> route)
	middlewareQueue := []MiddlewareFunc{}
	middlewareQueue = append(middlewareQueue, g.middleware...)                    // Engine-level middleware
	middlewareQueue = append(middlewareQueue, route.middleware...)                // Route-specific middleware
	middlewareQueue = append(middlewareQueue, handlerToMiddleware(route.handler)) // Final handler

	var broker dyffiBroker.Queue
	if g.broker != nil {
		broker = *g.broker
	}

	ctx := &Context{
		writer:                     w,
		request:                    r,
		params:                     params,
		Headers:                    r.Header,
		index:                      0,
		middleware:                 middlewareQueue,
		queueBroker:                broker,
		authorizationConfiguration: g.authConf,
	}

	return ctx
}

func (g *Engine) addInjectedRoute(method, path string, handler interface{}) {
    hv := reflect.ValueOf(handler)
    ht := hv.Type()
    ctxT := reflect.TypeOf(&Context{})

    if ht.Kind() != reflect.Func ||
       ht.NumIn() == 0 ||
       ht.In(0) != ctxT {
        panic("handler must be func(*Context, …)")
    }

    deps := make([]reflect.Value, ht.NumIn()-1)
    for i := 1; i < ht.NumIn(); i++ {
        depT := ht.In(i)

        svc, ok := g.services[depT]

        if !ok && depT.Kind() == reflect.Interface {
            for _, candidate := range g.services {
                if candidate.Type().Implements(depT) {
                    svc = candidate
                    ok = true
                    break
                }
            }
        }

        if !ok {
            panic(fmt.Sprintf("no service registered for %v", depT))
        }
        deps[i-1] = svc
    }

    wrapper := func(c *Context) {
        args := make([]reflect.Value, 1+len(deps))
        args[0] = reflect.ValueOf(c)
        copy(args[1:], deps)
        hv.Call(args)
    }

    g.addRoute(method, path, wrapper, nil, nil)
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
			if pathPart.regexPattern == "" {
				pathPart.part = part[1:]
			}

			pathPart.part = pathPart.part[1:]
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

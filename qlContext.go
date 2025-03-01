package dyffi

import (
	"fmt"
	"github.com/graphql-go/graphql"
)

// QLContext is a wrapper for GraphQL ResolveParams
type QLContext struct {
	Params graphql.ResolveParams
}

// ArgString safely gets a string argument
func (c *QLContext) ArgString(name string) string {
	if val, ok := c.Params.Args[name].(string); ok {
		return val
	}
	return ""
}

// ArgInt safely gets an int argument
func (c *QLContext) ArgInt(name string) int {
	if val, ok := c.Params.Args[name].(int); ok {
		return val
	}
	return 0
}

// ArgBool safely gets a boolean argument
func (c *QLContext) ArgBool(name string) bool {
	if val, ok := c.Params.Args[name].(bool); ok {
		return val
	}
	return false
}

// Log prints GraphQL arguments (for debugging)
func (c *QLContext) Log() {
	fmt.Println("GraphQL Args:", c.Params.Args)
}

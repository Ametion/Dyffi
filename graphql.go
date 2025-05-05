package dyffi

import (
	"github.com/graphql-go/graphql"
	"reflect"
)

// GraphQLResolverFunc is a function that resolves a GraphQL query or mutation
type GraphQLResolverFunc func(QLContext) (interface{}, error)

// GraphQLResolvers holds all GraphQL resolvers
type GraphQLResolvers struct {
	Query          any
	MutationCreate any
	MutationUpdate any
	MutationDelete any
}

// GraphQLRoute holds a GraphQL schema and resolvers
type GraphQLRoute struct {
	Schema    *graphql.Schema
	Resolvers GraphQLResolvers
}

// GraphQLSchemas holds all GraphQL routes
type GraphQLSchemas map[string]*GraphQLModel

// GraphQLModel holds a GraphQL model type and resolvers
type GraphQLModel struct {
	ModelType reflect.Type
	Resolvers GraphQLResolvers
}

package dyffi

import (
	"github.com/graphql-go/graphql"
	"reflect"
)

// GraphQLResolverFunc is a function that resolves a GraphQL query or mutation
type GraphQLResolverFunc func(QLContext) (interface{}, error)

// GraphQLResolvers holds all GraphQL resolvers
type GraphQLResolvers struct {
	Query          GraphQLResolverFunc
	MutationCreate GraphQLResolverFunc
	MutationUpdate GraphQLResolverFunc
	MutationDelete GraphQLResolverFunc
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

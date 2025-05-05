package main

import (
	"github.com/Ametion/dyffi"
)

func main() {
	engine := dyffi.NewDyffiEngine()

	engine.SetDevelopment()

	//Register GraphQL Model (no need to manually register /graphql)
	engine.GraphQLModel("Post", Post{}, dyffi.GraphQLResolvers{
		Query:          PostQueryResolver,  // Query resolver
		MutationCreate: PostMutationCreate, // Mutation resolver for creating a post
		MutationUpdate: PostMutationUpdate, // Mutation resolver for updating a post
		MutationDelete: PostMutationDelete, // Mutation resolver for deleting a post
	})

	engine.Run(":8080")
}

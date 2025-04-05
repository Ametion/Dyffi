package main

import (
	"github.com/Ametion/dyffi"
)

// Define Post struct
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Define GraphQL resolvers for Post
func PostQueryResolver(ctx dyffi.QLContext) (interface{}, error) {
	id := ctx.ArgInt("id")
	return Post{ID: id, Title: "GraphQL in Go", Content: "Building GraphQL APIs with Dyffi."}, nil
}

// Define GraphQL resolver for Post Creation
func PostMutationCreate(ctx dyffi.QLContext) (interface{}, error) {
	return Post{ID: ctx.ArgInt("id"), Title: ctx.ArgString("Title"), Content: ctx.ArgString("Content")}, nil
}

// Define GraphQL resolver for Post Update
func PostMutationUpdate(ctx dyffi.QLContext) (interface{}, error) {
	return Post{ID: ctx.ArgInt("id"), Title: ctx.ArgString("Title"), Content: ctx.ArgString("Content")}, nil
}

// Define GraphQL resolver for Post Deletion
func PostMutationDelete(ctx dyffi.QLContext) (interface{}, error) {
	return Post{ID: ctx.ArgInt("id"), Title: ctx.ArgString("Title"), Content: ctx.ArgString("Content")}, nil
}

package main

import (
	"fmt"
	"github.com/Ametion/dyffi"
)

func main() {
	engine := dyffi.NewDyffiEngine()

	engine.IsDevelopment()

	api := engine.Group("/api")

	api.UseMiddleware(APIMiddleware) // Enable custom middleware for this group

	// Path for this route will be /api/test, and all middleware used for this group, will be applied for this group only.
	api.Get("/test", func(context *dyffi.Context) {
		context.SendJSON(200, "Hello World")
	})

	engine.Run(":8080")
}

func APIMiddleware(context *dyffi.Context) {
	fmt.Println("API Middleware executed")

	context.Next()
}

package main

import (
	"fmt"
	"github.com/Ametion/dyffi"
)

func main() {
	engine := dyffi.NewDyffiEngine()

	engine.IsDevelopment()

	engine.UseMiddleware(ContextStorageMiddleware) // Enable custom middleware

	engine.Get("/test", func(context *dyffi.Context) {
		item := context.GetItem("key")

		fmt.Println(item) // Retrieve data from context storage

		context.SendJSON(200, "Hello World")
	})

	engine.Run(":8080")
}

func ContextStorageMiddleware(context *dyffi.Context) {
	// Store data in context storage
	context.SetItem("key", "value")

	context.Next() // Call the next middleware or handler
}

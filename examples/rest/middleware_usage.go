package main

import (
	"fmt"
	"github.com/Ametion/dyffi"
)

func main() {
	engine := dyffi.NewDyffiEngine()

	engine.SetDevelopment()

	engine.UseMiddleware(Middleware) // Enable custom middleware

	engine.Get("/test", func(context *dyffi.Context) {
		context.SendJSON(200, "Hello World")
	})

	engine.Run(":8080")
}

// Custom middleware function
func Middleware(context *dyffi.Context) {

	fmt.Println("Middleware executed")

	context.Next() // Call the next middleware or handler
}

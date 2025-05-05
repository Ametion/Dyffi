package main

import (
	"github.com/Ametion/dyffi"
)

func main() {
	cors := dyffi.CorsConfig{
		AllowedOrigins: []string{"http://localhost:3000"},        // Allow specific origin or all with "*"
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}, // Allow specific methods
		AllowedHeaders: []string{"Content-Type"},                 // Allow specific headers
	}

	engine := dyffi.NewDyffiEngine()

	engine.UseCors(cors) // Enable CORS middleware

	engine.SetDevelopment()

	engine.Get("/test", func(context *dyffi.Context) {
		context.SendJSON(200, "Hello World")
	})

	engine.Run(":8080")
}

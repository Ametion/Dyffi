package main

import (
	"github.com/Ametion/dyffi"
)

func main() {
	engine := dyffi.NewDyffiEngine()

	engine.IsDevelopment()

	engine.Get("/test", func(context *dyffi.Context) {
		context.SendJSON(200, "Hello World")
	})

	engine.Run(":8080")
}

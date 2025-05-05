package main

import (
	"github.com/Ametion/dyffi"
	"net/http"
)

func main() {
	engine := dyffi.NewDyffiEngine()

	engine.SetDevelopment()

	// Example of using regex in route parameters, this route will only match if the id is a number, and you need to write whole regex pattern inside brackets
	engine.Get("/user/:id(^/d+$)", func(c *dyffi.Context) {
		id := c.Param("id")
		c.SendJSON(http.StatusOK, map[string]string{"id": id})
	})

	engine.Run(":8080")
}

package main

import (
	"fmt"
	"github.com/Ametion/dyffi"
)

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	engine := dyffi.NewDyffiEngine()

	engine.IsDevelopment()

	engine.Post("/test", func(context *dyffi.Context) {
		var body User

		// Parse the JSON body into the User struct
		if err := context.SetBody(&body); err != nil {
			context.SendJSON(400, "invalid request")
			return
		}

		// Access the parsed data
		firstName := body.FirstName
		fmt.Println("First Name:", firstName)

		context.SendJSON(201, "user created")
	})

	// Start the server on port 8080
	engine.Run(":8080")
}

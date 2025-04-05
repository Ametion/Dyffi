package main

import (
	"github.com/Ametion/dyffi"
)

func main() {
	engine := dyffi.NewDyffiEngine()

	// Define authentication configurations
	auth := dyffi.APIAuthorization{
		AuthorizationType: dyffi.BASIC,
		ExcludedRoutes:    []string{"/login", "/register"}, // Excluded routes for Basic Auth
	}

	// Define Basic authentication configuration
	authConf := dyffi.BasicAuth{
		Username: "admin",
		Password: "admin",
	}

	engine.Authorization(auth, authConf)

	// This route will be protected by Basic Auth
	engine.Get("/test", AnotherTest)

	engine.Run(":8080")
}

// This handler will be protected by Basic Auth
func AnotherTest(context *dyffi.Context) {
	context.SendJSON(200, "Hello World")
}

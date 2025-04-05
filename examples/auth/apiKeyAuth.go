package main

import (
	"github.com/Ametion/dyffi"
)

func main() {
	engine := dyffi.NewDyffiEngine()

	// Define authentication configurations
	auth := dyffi.APIAuthorization{
		AuthorizationType: dyffi.APIKEY,
		ExcludedRoutes:    []string{"/login", "/register"}, // Excluded routes for APIKey Auth
	}

	// Define APIKey authentication configuration
	authConf := dyffi.APIKeyAuth{
		APIKey: "some_api_key",
	}

	engine.Authorization(auth, authConf)

	// This route will be protected by APIKey Auth
	engine.Get("/test", AnotherTest)

	engine.Run(":8080")
}

package main

import (
	"fmt"
	"github.com/Ametion/dyffi"
	"time"
)

// IMPORTANT NOTE: WHEN USING JWT AUTH, YOU CAN NOT USE KEY "claims" IN CONTEXT, ITS RESERVED FOR JWT AUTH
func main() {
	engine := dyffi.NewDyffiEngine()

	// Define authentication configurations
	auth := dyffi.APIAuthorization{
		AuthorizationType: dyffi.JWT,
		ExcludedRoutes:    []string{"/login", "/register"}, // Excluded routes for JWT Auth
	}

	// Define JWT authentication configuration
	authConf := dyffi.JWTAuth{
		JWTSecret: "some_secret_key",
		ExpireAt:  24 * time.Hour,
	}

	engine.Authorization(auth, authConf)

	// Route for login
	engine.Post("/login", Login)

	// Route for test with JWT authentication and claims inside context
	engine.Get("/test", Test)

	engine.Run(":8080")
}

func Login(context *dyffi.Context) {
	username := context.PostForm("username")
	password := context.PostForm("password")

	// (Assume credentials are valid for now)
	fmt.Println("Username and password received: ", username, password)

	accessToken, err := context.LoginJWT(map[string]interface{}{"data": "some_data"})
	if err != nil {
		context.SendJSON(500, "Internal Server Error")
		return
	}

	context.SendJSON(200, map[string]string{"access_token": accessToken})
}

func Test(context *dyffi.Context) {
	claims := context.GetItem("claims") //here is saved whole data from token

	fmt.Println(claims) //will be printed map with data

	context.SendJSON(200, "Hello World")
}

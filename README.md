# Dyffi Router

Dyffi is a lightweight, modular, and developer-friendly HTTP router for building scalable web servers in Go. Designed with simplicity and flexibility in mind, it supports middleware, route grouping, GraphQL, and advanced CORS handling to help you develop robust web applications.

---

# Features

- 🌎 **Simple REST Routing** – Define routes for common HTTP methods (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`).
- 🔐 **Auto Authentication System** – Supports JWT, API Key, and Basic Auth with middleware.
- 🔍 **Regex in Routes** – Use regex constraints in dynamic path parameters.
- 📦 **GraphQL Support** – Built-in GraphQL API handling with automatic schema generation.
- 🔀 **Route Grouping** – Organize routes logically with route groups.
- 🔧 **Middleware System** – Easily extend functionality with middleware.
- 🌍 **CORS Support** – Advanced configuration for Cross-Origin Resource Sharing.
- 🛠 **Developer Mode** – Real-time color-coded request logging for debugging.
- ⚡ **High Performance** – Optimized for speed and low memory footprint.

---

# Installation

Install Dyffi using `go get`:

```bash
go get github.com/Ametion/dyffi@latest
```

---

# Quick Start

### **REST API Example**

```go
package main

import (
	"github.com/Ametion/dyffi"
	"net/http"
)

func main() {
	engine := dyffi.NewDyffiEngine()

	// Enable development mode for logging
	engine.IsDevelopment()

	// Define REST routes
	engine.Get("/hello", func(c *dyffi.Context) {
		c.SendJSON(http.StatusOK, map[string]string{"message": "Hello, world!"})
	})

	engine.Post("/submit", func(c *dyffi.Context) {
		var input map[string]interface{}
		if err := c.SetBody(&input); err != nil {
			c.SendJSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
			return
		}
		c.SendJSON(http.StatusOK, input)
	})

	// Start the server
	engine.Run(":8080")
}
```

---

# GraphQL Support

Dyffi natively supports **GraphQL APIs**, allowing you to define schemas and resolvers easily.

### **🔹 Quick GraphQL Example**

```go
package main

import (
	"github.com/Ametion/dyffi"
)

// Define User struct
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Define Post struct
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// GraphQL Query Resolvers

func UserQueryResolver(ctx dyffi.QLContext) (interface{}, error) {
	id := ctx.ArgInt("id")
	return User{ID: id, Name: "John Doe", Age: 25}, nil
}

func PostQueryResolver(ctx dyffi.QLContext) (interface{}, error) {
	id := ctx.ArgInt("id")
	return Post{ID: id, Title: "GraphQL in Go", Content: "Building GraphQL APIs with Dyffi."}, nil
}

func PostMutationCreate(ctx dyffi.QLContext) (interface{}, error) {
	return Post{ID: ctx.ArgInt("ID"), Title: ctx.ArgString("Title"), Content: ctx.ArgString("Content")}, nil
}


func main() {
	engine := dyffi.NewDyffiEngine()

	// Register GraphQL API
	engine.GraphQLModel("User", User{}, dyffi.GraphQLResolvers{
		Query:          UserQueryResolver,
	})

	// Register GraphQL API
	engine.GraphQLModel("Post", Post{}, dyffi.GraphQLResolvers{
		Query:          PostQueryResolver,
		MutationCreate: PostMutationCreate,
	})

	// Start the server
	engine.Run(":8080")
}
```

## **🔹 Making a GraphQL Request**

#### **Query Example**
```graphql
query {
  getUser(id: 1) {
    ID
    Name
    Age
  }
}
```
### **Expected Response:**
```json
{
  "data": {
    "getUser": {
      "ID": 1,
      "Name": "John Doe",
      "Age": 25
    }
  }
}
```

### **Mutation Example**
```graphql
mutation CreatePost {
    createPost(ID: 1, Title: "New Post", Content: "Some content") {
        Content
        ID
        Title
    }
}
```
### **Expected Response:**
```json
{
  "data": {
    "createPost": {
      "ID": 1,
      "Title": "New Post",
      "Content": "Some content"
    }
  }
}
```

### **✨ Now Dyffi seamlessly supports both REST and GraphQL APIs in a single router!**

---

#  ✨Advanced Features✨

---

# **Auto Authentication System**

#### Dyffi now supports **automatic authentication** middleware, allowing you to secure routes with **JWT tokens**, **API keys**, and **Basic Auth**.

## **Supported Authentication Methods**
1. **JWT Authentication** – Uses JSON Web Tokens (JWT) for user sessions.
2. **API Key Authentication** – Uses an API key sent via headers.
3. **Basic Authentication** – Requires a username and password.

### **Setup Example**

```go
package main

import (
    "github.com/Ametion/dyffi"
    "time"
)

func main() {
    engine := dyffi.NewDyffiEngine()

    // Define authentication configurations
    auth := dyffi.APIAuthorization{
        AuthorizationTypes: "JWT", // can be "JWT", "APIKey", "Basic"
        ExcludedRoutes:    []string{"/login", "/register"},
    }

    //Configuration for JWT
    authConf := dyffi.JWTAuth{
        JWTSecret: "some_secret_key",
        ExpireAt:  24 * time.Hour,
    }

    //Configuration for API Key
    authConf := dyffi.APIKeyAuth{
        APIKey: "some_api_key",
    }

    //Configuration for Basic Auth
    authConf := dyffi.BasicAuth{
        Username: "admin",
        Password: "admin",
    }
	
    //apply authentication to engine
    engine.Authorization(auth, authConf)
}
```

---

## JWT Authentication Example

## **Login Route (Token Generation)**

```go
engine.Post("/login", func(c *dyffi.Context) {
    username := c.PostForm("username")
    password := c.PostForm("password")

    // (Assume credentials are valid for now) 
    accessToken, err := context.LoginJWT(map[string]interface{}{"data": "some_data"})
    if err != nil {
        context.SendJSON(500, "Internal Server Error")
        return
    }

    c.SendJSON(200, map[string]string{"access_token": token})
})
```

## **Get Data from Token on Protected Route**

```go
engine.Get("/some_route", func(c *dyffi.Context) {
    claims := context.GetItem("claims") //here is saved whole data from token
    
    fmt.Println(claims) //will be printed map with data
	
	//some logic here
	
    c.SendJSON(200, "Protected content")
})
```

## **Making a Secure Request**

Include the token in the `Authorization` header:

```
Authorization: Bearer <your-token>
```

---

# **Route Protection & Exclusions**

By default, authentication is **required for all routes** except those **explicitly excluded** in `ExcludedRoutes`.

```go
auth := dyffi.APIAuthorization{
    AuthorizationTypes: "JWT",
    ExcludedRoutes:    []string{"/login", "/register"},
}
```

---

# **Middleware**

Use middleware to extend functionality, such as authentication, logging, or modifying requests:

```go
engine.UseMiddleware(func(c *dyffi.Context) {
	c.writer.Header().Set("X-Powered-By", "Dyffi")
	c.Next()
})
```

# **Route Grouping**

Group related routes together for better organization:

```go
api := engine.Group("/api")
api.UseMiddleware(func(c *dyffi.Context) {
	c.writer.Header().Set("X-API-Version", "1.0")
	c.Next()
})

api.Get("/users", func(c *dyffi.Context) {
	c.SendJSON(http.StatusOK, []string{"user1", "user2", "user3"})
})
```

# **CORS Configuration**

Enable CORS to control access for different origins:

```go
engine.UseCors(dyffi.CorsConfig{
	AllowedOrigins: []string{"*"}, // Allow all origins
	AllowedMethods: []string{"GET", "POST", "PUT"},
	AllowedHeaders: []string{"Authorization", "Content-Type"},
})
```

---

# Regex in Path Parameters

Dyffi supports regex-based path parameters to enforce constraints on dynamic segments. *You need to put regex into "()" brackets.*

```go
api.Get("/user/:id(^\d+$)", func(c *dyffi.Context) {
	id := c.Param("id")
	c.SendJSON(http.StatusOK, map[string]string{"id": id})
})
```

---

# Development Logging

When **development mode** is enabled (`engine.IsDevelopment()`), Dyffi logs:
- **Date & Time**
- **HTTP Method**
- **Status Code** (color-coded)
- **Request Path**

---

# Contributing

Contributions are welcome! If you’d like to improve Dyffi:
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request.

---

# Feedback

Have ideas, suggestions, or found an issue?  
📢 **Open an issue on the GitHub repository!**

---

# Authors

- **Yehor Kochetov** - [GitHub](https://github.com/Ametion)

---

## 🚀 Happy coding with Dyffi!
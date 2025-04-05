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

All REST examples you can find [here](https://github.com/Ametion/Dyffi/tree/dev/examples/rest)


---

# GraphQL Support

Dyffi natively supports **GraphQL APIs**, allowing you to define schemas and resolvers easily.

#### Full example [here](https://github.com/Ametion/Dyffi/tree/dev/examples/graphql).

### **🔹 Quick GraphQL Example**

```go
package main

import (
	"github.com/Ametion/dyffi"
)

func main() {
	engine := dyffi.NewDyffiEngine()

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
query GetPost {
    getPost(id: 1) {
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
    "getPost": {
      "Content": "Building GraphQL APIs with Dyffi.",
      "ID": 1,
      "Title": "GraphQL in Go"
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
1. **JWT Authentication** – Uses JSON Web Tokens (JWT) for user sessions, full example [here](https://github.com/Ametion/Dyffi/tree/dev/examples/auth/jwtAuth.go).
2. **API Key Authentication** – Uses an API key sent via headers, full example [here](https://github.com/Ametion/Dyffi/tree/dev/examples/auth/apiKeyAuth.go).
3. **Basic Authentication** – Requires a username and password, full example [here](https://github.com/Ametion/Dyffi/tree/dev/examples/auth/basicAuth.go).

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
        AuthorizationType: dyffi.JWT,
        ExcludedRoutes:    []string{"/login", "/register"},
    }

    //Configuration for JWT
    authConf := dyffi.JWTAuth{
        JWTSecret: "some_secret_key",
        ExpireAt:  24 * time.Hour,
    }
	
    //apply authentication to engine
    engine.Authorization(auth, authConf)
}
```

---

# **Route Protection & Exclusions**

By default if you are using any auth type, authentication is **required for all routes** except those **explicitly excluded** in `ExcludedRoutes`.

```go
auth := dyffi.APIAuthorization{
    AuthorizationTypes: "JWT",
    ExcludedRoutes:    []string{"/login", "/register"},
}
```

---

# **Brokers**
Dyffi supports **brokers** for message queuing and event-driven architectures. You can use brokers like **RabbitMQ**, **Kafka**, or **Nats** to send message to your brokers/queues.

For start using brokers inside your router you need to create broker config and pass it to `engine.UseBroker()` method. Otherwise, it will not work.

All examples for all supported brokers are available [here](https://github.com/Ametion/Dyffi/tree/dev/examples/brokers)

### **Example with Kafka**

```go
package main

import (
	"fmt"
	"github.com/Ametion/dyffi"
	dyffiBroker "github.com/Ametion/dyffi/broker"
)

func main() {
	brokerConf := dyffiBroker.BrokerConfig{
		BrokerType: dyffiBroker.Kafka,
		Host:       "localhost",
		Port:       "9092",
	}

	engine := dyffi.NewDyffiEngine()
	

	engine.UseBroker(brokerConf)

	engine.IsDevelopment()

	engine.Post("/someRoute", func(context *dyffi.Context) {
		//This line will send message "test message" to Kafka "test_topic" topic, note, it will not automatically create topic!.
		err := context.PublishToQueue("test_topic", []byte("test message"), dyffiBroker.KafkaParams{})

		if err != nil {
			fmt.Println(err.Error())
		}

		context.SendJSON(200, "Data")
	})
	
	engine.Run(":8080")
}

````

### Important Note: You need to provide right params type for publish function. For example, if you are using Kafka, you need to provide `dyffiBroker.KafkaParams{}` as params. You can use any params inside this struct which is supported for chosen broker.

### Supported Brokers:
- **Kafka** - [Kafka](https://kafka.apache.org/)
- **RabbitMQ** - [RabbitMQ](https://www.rabbitmq.com/)
- **Nats** - [Nats](https://nats.io/)

---

# **Middleware**

Use middleware to extend functionality, such as authentication, logging, or modifying requests:

[Full Example](https://github.com/Ametion/Dyffi/tree/dev/examples/rest/middleware_usage.go)

```go
engine.UseMiddleware(func(c *dyffi.Context) {
	c.writer.Header().Set("X-Powered-By", "Dyffi")
	c.Next()
})
```

# **Route Grouping**

Group related routes together for better organization:

[Full Example](https://github.com/Ametion/Dyffi/tree/dev/examples/rest/routerGrouping_usage.go)

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

[Full Example](https://github.com/Ametion/Dyffi/tree/dev/examples/rest/cors_usage.go)


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

[Full Example](https://github.com/Ametion/Dyffi/tree/dev/examples/rest/regex_usage.go)

```go
api.Get("/user/:id(^\d+$)", func(c *dyffi.Context) {
	id := c.Param("id")
	c.SendJSON(http.StatusOK, map[string]string{"id": id})
})
```

# Context Storage

Dyffi provides a context storage mechanism to store and retrieve data during request processing. This is useful for sharing data between middleware and handlers.

[Full Example](https://github.com/Ametion/Dyffi/tree/dev/examples/rest/contextStorage_usage.go)

```go
engine.Get("/user/:id", func(c *dyffi.Context) {
    userID := c.GetItem("userID")
	
    c.SendJSON(http.StatusOK, map[string]string{"userID": userID})
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
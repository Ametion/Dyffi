# Dyffi Router

Dyffi is a lightweight, modular, and developer-friendly HTTP router for building scalable web servers in Go. Designed with simplicity and flexibility in mind, it supports middleware, route grouping, and advanced CORS handling to help you develop robust web applications.

---

## Features

- **Simple Routing**: Easily define routes for common HTTP methods (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`).
- **Middleware Support**: Add global, group-specific, or route-specific middleware for extensibility.
- **Route Grouping**: Organize routes logically with route groups.
- **CORS Support**: Built-in configuration for Cross-Origin Resource Sharing.
- **Context Handling**: Intuitive request/response utilities for JSON, query params, headers, and more.
- **Developer Mode**: Real-time color-coded request logging for debugging.
- **Flexible Extensibility**: Easily customize with minimal configuration.

---

## Installation

Install Dyffi using `go get`:

```bash
go get github.com/Ametion/dyffi
```

---

## Quick Start

Here’s an example to get you started:

```go
package main

import (
	"github.com/Ametion/dyffi"
	"net/http"
)

func main() {
	// Create a new Dyffi engine
	engine := dyffi.NewDyffiEngine()

	// Enable development mode for detailed request logging
	engine.IsDevelopment()

	// Add a GET route
	engine.Get("/hello", func(c *dyffi.Context) {
		c.SendJSON(http.StatusOK, map[string]string{"message": "Hello, world!"})
	})

	// Add a POST route
	engine.Post("/submit", func(c *dyffi.Context) {
		var input map[string]interface{}
		if err := c.SetBody(&input); err != nil {
			c.SendJSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
			return
		}
		c.SendJSON(http.StatusOK, input)
	})

	// Enable CORS
	engine.UseCors(dyffi.CorsConfig{
		AllowedOrigins: []string{"http://example.com"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Start the server
	engine.Run(":8080")
}
```

---

## Advanced Usage

### Middleware

Dyffi allows you to define middleware to extend functionality, such as authentication, logging, or input validation:

```go
engine.UseMiddleware(func(c *dyffi.Context) {
	// Custom middleware logic
	c.writer.Header().Set("X-Powered-By", "Dyffi")
	c.Next() // Proceed to the next middleware or route handler
})
```

### Route Grouping

Group related routes together for better organization:

```go
api := engine.Group("/api")
api.UseMiddleware(func(c *dyffi.Context) {
	// Group-specific middleware
	c.writer.Header().Set("X-API-Version", "1.0")
	c.Next()
})

api.Get("/users", func(c *dyffi.Context) {
	c.SendJSON(http.StatusOK, []string{"user1", "user2", "user3"})
})
```

### CORS Configuration

Enable CORS to control access for different origins:

```go
engine.UseCors(dyffi.CorsConfig{
	AllowedOrigins: []string{"*"}, // Allow all origins
	AllowedMethods: []string{"GET", "POST", "PUT"},
	AllowedHeaders: []string{"Authorization", "Content-Type"},
})
```

---

## Development Logging

When development mode is enabled using `engine.IsDevelopment()`, Dyffi logs detailed information about incoming requests:
- **Date & Time**
- **HTTP Method**
- **Status Code** (color-coded)
- **Request Path**

Example log output:
```
Date: Tue, 28 Jan 2025 14:05:00 GMT, Method: GET, Status code: 200, Route: /hello
```

---

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests to improve Dyffi. To get started:
1. Fork the repository.
2. Create a new branch for your feature/bug fix.
3. Submit a pull request.

---

## Feedback

We’d love to hear your feedback! If you have any suggestions, feature requests, or issues, please open an issue on the GitHub repository.

---

## Authors

- **Yehor Kochetov** - [GitHub](https://github.com/Ametion)

---

Happy coding with Dyffi! 🚀

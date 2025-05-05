package dyffi

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthType string

const (
	JWT    AuthType = "JWT"
	APIKEY AuthType = "APIKey"
	BASIC  AuthType = "Basic"
)

type Authorization interface {
    authorize(c *Context)
}

// APIAuthorization holds your global auth settings.
type APIAuthorization struct {
	AuthorizationType AuthType
	ExcludedRoutes    []string
}

// JWTAuth holds JWT-specific config.
type JWTAuth struct {
	JWTSecret string
	ExpireAt  time.Duration
}

// APIKeyAuth holds API-key-specific config.
type APIKeyAuth struct {
	APIKey string
}

// BasicAuth holds basic auth-specific config.
type BasicAuth struct {
	Username string
	Password string
}

// Authorization sets up a global middleware for the chosen auth type.
// You call this once in your main setup, like:
//
//	engine.Authorization(auth, authConf)
//
// where 'auth' is an APIAuthorization, and 'authConf' is either JWTAuth or APIKeyAuth, etc.
func (g *Engine) Authorization(auth APIAuthorization, conf Authorization) {
	g.authConf = conf

	// Attach the global middleware
	g.UseMiddleware(func(c *Context) {
		requestPath := c.request.URL.Path
		for _, excluded := range auth.ExcludedRoutes {
			if strings.HasPrefix(requestPath, excluded) {
				c.Next()
				return
			}
		}

		conf.authorize(c)
	})
}

// handleJWTAuth is the actual JWT-check logic.
func (auth JWTAuth) authorize(c *Context) {
	tokenHeader := c.request.Header.Get("Authorization")
	if tokenHeader == "" {
		http.Error(c.writer, "Missing Authorization header", http.StatusUnauthorized)
		c.Abort()
	}

	parts := strings.Split(tokenHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(c.writer, "Invalid Authorization header format", http.StatusUnauthorized)
		c.Abort()
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(c.writer, "Invalid token", http.StatusUnauthorized)
		c.Abort()
	}

	c.SetItem("claims", token.Claims)

	c.Next()
}

// handleAPIKeyAuth is the actual API key check logic.
func (auth APIKeyAuth) authorize(c *Context) {
	apiKey := c.request.Header.Get("X-API-KEY")
	if apiKey == "" || apiKey != auth.APIKey {
		http.Error(c.writer, "Invalid API Key", http.StatusUnauthorized)
		c.Abort()
	}
	c.Next()
}

// handleBasicAuth is the actual basic auth check logic.
func (auth BasicAuth) authorize(c *Context) {
	username, password, ok := c.request.BasicAuth()
	if !ok || username != auth.Username || password != auth.Password {
		http.Error(c.writer, "Invalid credentials", http.StatusUnauthorized)
		c.Abort()
	}
	c.Next()
}

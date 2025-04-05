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
func (g *Engine) Authorization(auth APIAuthorization, conf interface{}) {
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

		switch auth.AuthorizationType {
		case JWT:
			jwtConf, ok := conf.(JWTAuth)
			if !ok {
				http.Error(c.writer, "Invalid JWT configuration", http.StatusInternalServerError)
				return
			}
			handleJWTAuth(c, jwtConf)

		case APIKEY:
			apiKeyConf, ok := conf.(APIKeyAuth)
			if !ok {
				http.Error(c.writer, "Invalid APIKey configuration", http.StatusInternalServerError)
				return
			}
			handleAPIKeyAuth(c, apiKeyConf)

		case BASIC:
			basicConf, ok := conf.(BasicAuth)
			if !ok {
				http.Error(c.writer, "Invalid BasicAuth configuration", http.StatusInternalServerError)
				return
			}
			handleBasicAuth(c, basicConf)

		default:
			http.Error(c.writer, "Unsupported authorization type", http.StatusForbidden)
			return
		}
	})
}

// handleJWTAuth is the actual JWT-check logic.
func handleJWTAuth(c *Context, conf JWTAuth) {
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
		return []byte(conf.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(c.writer, "Invalid token", http.StatusUnauthorized)
		c.Abort()
	}

	c.SetItem("claims", token.Claims)

	c.Next()
}

// handleAPIKeyAuth is the actual API key check logic.
func handleAPIKeyAuth(c *Context, conf APIKeyAuth) {
	apiKey := c.request.Header.Get("X-API-KEY")
	if apiKey == "" || apiKey != conf.APIKey {
		http.Error(c.writer, "Invalid API Key", http.StatusUnauthorized)
		c.Abort()
	}
	c.Next()
}

// handleBasicAuth is the actual basic auth check logic.
func handleBasicAuth(c *Context, conf BasicAuth) {
	username, password, ok := c.request.BasicAuth()
	if !ok || username != conf.Username || password != conf.Password {
		http.Error(c.writer, "Invalid credentials", http.StatusUnauthorized)
		c.Abort()
	}
	c.Next()
}

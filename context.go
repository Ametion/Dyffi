package dyffi

import (
	"encoding/json"
	"fmt"
	dyffiBroker "github.com/Ametion/dyffi/broker"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type JWTClaims struct {
	Data any `json:"data"`
	jwt.StandardClaims
}

// Context represents request context
type Context struct {
	writer                     http.ResponseWriter
	request                    *http.Request
	Headers                    http.Header
	aborted                    bool
	params                     map[string]pathPart
	index                      int
	middleware                 []MiddlewareFunc
	items                      map[string]any
	queueBroker                dyffiBroker.Queue
	authorizationConfiguration interface{}
}

func (c *Context) LoginJWT(data any) (string, error) {
	jwtAuth := c.authorizationConfiguration.(JWTAuth)

	expirationTime := time.Now().Add(jwtAuth.ExpireAt)

	claims := &JWTClaims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtAuth.JWTSecret))
}

func (c *Context) PublishToQueue(topic string, message []byte, params dyffiBroker.DefaultParams) error {
	if c.queueBroker == nil {
		return fmt.Errorf("queue broker is not initialized")
	}

	return c.queueBroker.Publish(topic, message, params)
}

// Set choosed item by choosed index
func (c *Context) SetItem(index string, item any) {
	if len(c.items) <= 0 {
		c.items = make(map[string]any)
	}

	c.items[index] = item
}

// Return choosed item by choosed index from param
func (c *Context) GetItem(index string) any {
	return c.items[index]
}

// Set abort variable to true
func (c *Context) Abort() {
	c.aborted = true
}

// Next proceeds to the next middleware
func (c *Context) Next() {
	for c.index < len(c.middleware) && !c.aborted {
		middleware := c.middleware[c.index]
		c.index++
		middleware(c)
	}
}

// Redirect redirects to the specific url with chosen status code
func (c *Context) Redirect(url string, statusCode int) {
	http.Redirect(c.writer, c.request, url, statusCode)
}

// Query gets a query value
func (c *Context) Query(key string) string {
	return c.request.URL.Query().Get(key)
}

// Param gets a path parameter
func (c *Context) Param(key string) string {
	return c.params[key].value
}

// PostForm gets a post form value with presented key
func (c *Context) PostForm(key string) string {
	if err := c.request.ParseForm(); err != nil {
		return ""
	}

	return c.request.PostFormValue(key)
}

func (c *Context) SetBody(v interface{}) error {
	decoder := json.NewDecoder(c.request.Body)
	defer c.request.Body.Close()
	return decoder.Decode(v)
}

// SendJSON sends a SendJSON response
func (c *Context) SendJSON(statusCode int, v interface{}) {
	c.writer.Header().Set("Content-Type", "application/json")
	c.writer.WriteHeader(statusCode)
	err := json.NewEncoder(c.writer).Encode(v)

	if err != nil {
		c.writer.WriteHeader(http.StatusInternalServerError)
		c.writer.Write([]byte("error while trying to encode json: " + err.Error()))
	}
}

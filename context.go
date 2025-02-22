package dyffi

import (
	"encoding/json"
	"net/http"
)

// Context represents request context
type Context struct {
	writer     http.ResponseWriter
	request    *http.Request
	Headers    http.Header
	aborted    bool
	params     map[string]pathPart
	index      int
	middleware []MiddlewareFunc
	items      map[string]any
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
	json.NewEncoder(c.writer).Encode(v)
}

package gfx

import (
	"encoding/json"
	"net/http"
)

// Context represents request context
type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Headers    http.Header
	aborted    bool
	params     map[string]string
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
	if c.aborted {
		return
	}

	for c.index < len(c.middleware) {
		middleware := c.middleware[c.index]
		c.index++ // Increment the index before calling middleware
		middleware(c)

		if c.aborted {
			return
		}
	}
}

// Redirect redirects to the specific url with choosed status code
func (c *Context) Redirect(url string, statusCode int) {
	http.Redirect(c.Writer, c.Request, url, statusCode)
}

// Query gets a query value
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// Param gets a path parameter
func (c *Context) Param(key string) string {
	return c.params[key]
}

// PostForm gets a post form value with presented key
func (c *Context) PostForm(key string) string {
	if err := c.Request.ParseForm(); err != nil {
		return ""
	}

	return c.Request.PostFormValue(key)
}

func (c *Context) SetBody(v interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	defer c.Request.Body.Close()
	return decoder.Decode(v)
}

// SendJSON sends a SendJSON response
func (c *Context) SendJSON(statusCode int, v interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(statusCode)
	json.NewEncoder(c.Writer).Encode(v)
}

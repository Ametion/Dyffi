package dyffi

// Route represents a HTTP route
type Route struct {
	method      string
	handler     HandlerFunc
	middleware  []MiddlewareFunc
	parts       []pathPart
	paramsIndex []int
}

type pathPart struct {
	index        int
	part         string
	value        string
	isParam      bool
	regexPattern string
}

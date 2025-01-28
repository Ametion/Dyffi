package gfx

import (
	"fmt"
	"net/http"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	development bool
	statusCode  int
	method      string
	route       string
}

func (l *LoggingResponseWriter) WriteHeader(statusCode int) {
	l.statusCode = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

func (l *LoggingResponseWriter) Write(b []byte) (int, error) {
	n, err := l.ResponseWriter.Write(b)

	if err == nil && l.development {
		statusColor := "\033[33m" // Default to yellow
		if l.statusCode == 200 || l.statusCode == 201 {
			statusColor = "\033[32m" // Green for success
		} else if l.statusCode == 500 || l.statusCode == 400 || l.statusCode == 401 || l.statusCode == 402 || l.statusCode == 403 || l.statusCode == 404 {
			statusColor = "\033[31m" // Red for errors
		}

		methodColor := "\033[35m" // Violet for method

		// Bold text
		bold := "\033[1m"

		// Reset formatting after printing
		reset := "\033[0m"

		fmt.Printf("%sDate: %s, Method: %s%s%s, Status code: %s%d%s, Route: %s%s%s\n",
			bold, time.Now().Format(time.RFC1123),
			methodColor, l.method, reset, // Method in violet
			statusColor, l.statusCode, reset, // Status code in green/red/yellow
			"\033[34m", l.route, reset) // Full route in blue

	}

	return n, err
}

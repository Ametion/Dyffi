package dyffi

import (
	"fmt"
	"strings"
	"time"
)

func (g *Engine) logRequest(method string, statusCode int, route string, params map[string]pathPart) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	if g.development {
		// **Colorized Logging (For Dev Mode)**
		statusColor := getStatusColor(statusCode)
		methodColor := "\033[1;35m" // Magenta for method
		routeColor := "\033[1;34m"  // Blue for route
		reset := "\033[0m"

		// Format parameters only if they exist
		var paramsString string
		if params != nil && len(params) > 0 {
			paramParts := []string{}
			for key, value := range params {
				paramParts = append(paramParts, fmt.Sprintf("\033[1;33m%s: \033[1;32m%s\033[0m", key, value.value))
			}
			paramsString = " | Params: " + strings.Join(paramParts, ", ")
		}

		// Print colorized log
		fmt.Printf("\033[1;31m%s\033[0m | Method: %s%s%s | Status: %s%d%s | Route: %s%s%s%s\n",
			timestamp,
			methodColor, method, reset,
			statusColor, statusCode, reset,
			routeColor, route, reset,
			paramsString,
		)
	} else {
		// **Plain Logging (Production Mode)**
		if params != nil && len(params) > 0 {
			fmt.Printf("[%s] Method: %s | Status: %d | Route: %s | Params: %v\n",
				timestamp, method, statusCode, route, params)
		} else {
			fmt.Printf("[%s] Method: %s | Status: %d | Route: %s\n",
				timestamp, method, statusCode, route)
		}
	}
}

func formatRoute(parts []pathPart, paramsIndex []int) string {
	formattedParts := make([]string, len(parts))

	for i, part := range parts {
		if contains(paramsIndex, i) {
			formattedParts[i] = fmt.Sprintf("\033[1;33m:%s\033[0m", part.part) // Yellow for params
		} else {
			formattedParts[i] = fmt.Sprintf("\033[1;34m%s\033[0m", part.part) // Blue for static parts
		}
	}

	return strings.Join(formattedParts, "/")
}

// Helper function to colorize status codes
func getStatusColor(status int) string {
	if status >= 200 && status < 300 {
		return "\033[1;32m" // Green for 2xx Success
	}
	if status >= 400 && status < 500 {
		return "\033[1;33m" // Yellow for 4xx Client Errors
	}
	return "\033[1;31m" // Red for 5xx Server Errors
}

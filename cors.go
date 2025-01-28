package dyffi

type CorsConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

func (g *Engine) UseCors(corsCFG CorsConfig) {
	g.isCors = true
	g.allowedOrigins = corsCFG.AllowedOrigins
	g.AllowedMethods = corsCFG.AllowedMethods
	g.AllowedHeaders = corsCFG.AllowedHeaders
}

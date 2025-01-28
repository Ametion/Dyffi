package gfx

type CorsConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

func (g *GFXEngine) UseCors(corsCFG CorsConfig) {
	g.isCors = true
	g.allowedOrigins = corsCFG.AllowedOrigins
	g.allowedMethods = corsCFG.AllowedMethods
	g.allowedHeaders = corsCFG.AllowedHeaders
}

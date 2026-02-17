package config

type CorsConfig struct {
	Origins []string
	Headers []string `default:"Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control"`
	Methods []string `default:"POST,OPTIONS,GET,PUT,DELETE,PATCH"`
}

package config

var AllowedOrigins = []string{
	"http://localhost:1010",
	"http://localhost:5173",
}

var GetAllowedOrigins = func() []string {
	return AllowedOrigins
}

package config

var AllowedOrigins = []string{
	"http://localhost:6969",
	"http://localhost:5173",
}

var GetAllowedOrigins = func() []string {
	return AllowedOrigins
}

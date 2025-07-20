package config

var AllowedOrigins = []string{
	"http://localhost:1010",
	"http://localhost:5173",
	"https://embeck.onrender.com",
}

func GetAllowedOrigins() []string {
	return AllowedOrigins
}

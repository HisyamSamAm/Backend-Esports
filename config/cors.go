package config


var AllowedOrigins = []string{
	"http://localhost:6969",
}

var GetAllowedOrigins = func() []string {
	return AllowedOrigins
}

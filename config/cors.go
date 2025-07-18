package config

var AllowedOrigins = []string {
		"http://localhost:1010",
		"http://localhost:5173",
	}

func GetAllowedOrigins() []string {
	return AllowedOrigins
}	



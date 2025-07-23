package config

var AllowedOrigins = []string{
	"http://localhost:1010",
	"http://localhost:5173",
	// "https://embeck.onrender.com",
	"https://backend-esports.up.railway.app", // deploy
	"https://esports-app.netlify.app", // deploy

}

func GetAllowedOrigins() []string {
	return AllowedOrigins
}

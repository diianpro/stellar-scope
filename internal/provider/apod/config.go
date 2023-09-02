package apod

type Config struct {
	Address string `env:"APOD_ADDRESS,required" envDefault:"https://api.nasa.gov/planetary/apod"`
	ApiKey  string `env:"APOD_API_KEY,required" envDefault:"hexAwcWxOrEDePOQf3NrtOo78460rq8WplXhl3K9"`
}

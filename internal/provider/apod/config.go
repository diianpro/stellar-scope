package apod

type Config struct {
	Address string `env:"APOD_ADDRESS,required" envDefault:"https://api.nasa.gov/planetary/apod"`
	ApiKey  string `env:"APOD_API_KEY,required" envDefault:"uqe4UgO54kToru7pU6PqJDdmlhuMXxpkzEbpHQZV"`
}

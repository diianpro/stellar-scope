package apod

type Config struct {
	Address string `env:"ADDRESS,required"`
	ApiKey  string `env:"API_KEY,required"`
}

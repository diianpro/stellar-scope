package postgres

type Config struct {
	URL      string `env:"DB_URL,required" envDefault:"postgres://su:su@postgres:5432/image?sslmode=disable"`
	MinConns int32  `env:"MIN_CONNS,required" envDefault:"1"`
	MaxConns int32  `env:"MAX_CONNS,required" envDefault:"3"`
	//Username string `env:"USERNAME,required" envDefault:"su"`
	//Password string `env:"PASSWORD,required" envDefault:"su"`
}

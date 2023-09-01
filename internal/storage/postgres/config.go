package postgres

type Config struct {
	URL        string `env:"URL,required"`
	MinConns   int32  `env:"MIN_CONNS,required"`
	MaxConns   int32  `env:"MAX_CONNS,required"`
	SecretName string `env:"SECRET_NAME,required"`
	Username   string `env:"USERNAME,required"`
	Password   string `env:"PASSWORD,required"`
}

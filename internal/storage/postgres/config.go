package postgres

//type Config struct {
//	Port     int    `env:"APP_DB_PORT" envDefault:"5433"`
//	DBName   string `env:"APP_DB_NAME" envDefault:"postgres"`
//	Password string `env:"APP_DB_PASSWORD"`
//	Username string `env:"APP_DB_USERNAME" envDefault:"postgres"`
//	URL      string
//	Host     string `envDefault:"localhost"`
//}

type Config struct {
	URL        string `mapstructure:"url" valid:"required"`
	MinConns   int32  `mapstructure:"min_connections"`
	MaxConns   int32  `mapstructure:"max_connections" valid:"required"`
	SecretName string `mapstructure:"secret_name" valid:"required"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

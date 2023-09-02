package s3

type Config struct {
	Bucket string `env:"BUCKET,required" envDefault:"images"`
	Region string `env:"REGION,required" envDefault:"us-east-1"`
}

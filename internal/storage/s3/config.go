package s3

type Config struct {
	Bucket string `env:"BUCKET,required"`
	Region string `env:"REGION,required"`
}

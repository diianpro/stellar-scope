package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type IntegrationTestSuite struct {
	container *Container
	storage   *Storage
	suite.Suite
}

func (i *IntegrationTestSuite) SetupSuite() {
	var err error
	i.container, err = NewContainer()
	i.Require().NoError(err)

	const defaultRegion = "us-east-1"
	awsConfig := aws.Config{
		Region:      defaultRegion,
		Credentials: credentials.NewStaticCredentialsProvider("minioadmin", "minioadmin", ""),
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               fmt.Sprintf("http://%s", i.container.addr),
				SigningRegion:     defaultRegion,
				HostnameImmutable: true,
			}, nil
		}),
	}

	s3Storage := s3.NewFromConfig(awsConfig)
	cfg := Config{
		Bucket: "images",
		Region: defaultRegion,
	}

	_, err = s3Storage.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		log.Errorf("Could not create bucket: %v", err)
	}

	i.storage = New(cfg.Bucket, s3Storage)
}

func (i *IntegrationTestSuite) TearDownSuite() {
	err := i.container.Purge()
	i.Assert().NoError(err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (i *IntegrationTestSuite) TestIntegrationUpload() {
	ctx := context.Background()

	date := time.Now()
	imageData := []byte("test-image-data")

	err := i.storage.Upload(ctx, date, imageData)
	i.Assert().NoError(err)

	retData, err := i.storage.Get(ctx, date.Format(time.DateOnly))
	i.Assert().NoError(err)
	i.Assert().NotNil(retData)

	i.Require().Equal(imageData, retData)

}

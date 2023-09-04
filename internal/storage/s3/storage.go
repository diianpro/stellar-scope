package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/diianpro/stellar-scope/internal/domain/image"
)

type Storage struct {
	bucket   string
	cli      *s3.Client
	uploader *manager.Uploader
}

func New(bucket string, cli *s3.Client) *Storage {
	return &Storage{
		bucket:   bucket,
		cli:      cli,
		uploader: manager.NewUploader(cli),
	}
}

func (s *Storage) Get(ctx context.Context, date string) ([]byte, error) {
	output, err := s.cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(date),
	})
	if err != nil {
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			return nil, image.ErrImageNotExists
		}
		return nil, fmt.Errorf("storage error: %w", err)
	}
	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("storage error read output body")
	}
	return data, nil
}

func (s *Storage) Upload(ctx context.Context, date time.Time, imageData []byte) error {
	body := bytes.NewBuffer(imageData)
	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(date.Format(time.DateOnly)),
		Body:   body,
	})
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}
	return nil
}

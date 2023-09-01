package service

import (
	"context"
	"time"

	"github.com/diianpro/stellar-scope/internal/domain/apod"
	"github.com/diianpro/stellar-scope/internal/domain/image"
)

type Picture interface {
	GetByDate(ctx context.Context, date time.Time) (*apod.Data, error)
	GetAll(ctx context.Context, limit, offset int) ([]apod.Data, error)
	ObserveDailyImage(ctx context.Context) error
}

type ImageProvider interface {
	GetMetadata(ctx context.Context, date time.Time) (*apod.Data, error)
}

type Service struct {
	apodRps       apod.Repository
	imageRps      image.Repository
	imageProvider ImageProvider
}

func New(apodRps apod.Repository, imageRps image.Repository, imageProvider ImageProvider) *Service {
	return &Service{
		apodRps:       apodRps,
		imageRps:      imageRps,
		imageProvider: imageProvider,
	}
}

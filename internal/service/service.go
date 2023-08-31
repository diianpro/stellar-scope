package service

import (
	"context"
	"github.com/diianpro/stellar-scope/internal/domain"
	"time"
)

type Picture interface {
	GetByDate(ctx context.Context, date time.Time) (*domain.ApodData, error)
	GetAll(ctx context.Context) ([]domain.ApodData, error)
}

type ImageStorage interface {
	Upload(ctx context.Context, image []byte, filename string) (string, error)
	Download(ctx context.Context, date string) ([]byte, error)
}

type Storage interface {
	domain.Repository
	ImageStorage
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

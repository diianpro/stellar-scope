package service

import (
	"context"
	"github.com/diianpro/stellar-scope/internal/domain"
	"time"
)

func (s *Service) GetByDate(ctx context.Context, date time.Time) (*domain.ApodData, error) {
	return s.storage.GetByDate(ctx, date)
}

func (s *Service) GetAll(ctx context.Context) ([]domain.ApodData, error) {
	return s.storage.GetAll(ctx)
}

package domain

import (
	"context"
	"time"
)

type Repository interface {
	GetByDate(ctx context.Context, date time.Time) (*ApodData, error)
	GetAll(ctx context.Context) ([]ApodData, error)
	Create(ctx context.Context, data *ApodData, imageLink string) (int64, error)
}

type ApodData struct {
	Title       string
	Explanation string
	Date        string
	ImageLink   string
	Copyright   string
}

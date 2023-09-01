package apod

import (
	"context"
	"time"

	"github.com/diianpro/stellar-scope/internal/domain/image"
)

type Repository interface {
	GetByDate(ctx context.Context, date time.Time) (*Data, error)
	GetAll(ctx context.Context, limit, offset int) ([]Data, error)
	Create(ctx context.Context, data *Data) error
}

type Data struct {
	Explanation string
	Date        string
	Image       *image.Image
	Copyright   string
}

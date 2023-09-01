package image

import (
	"context"
	"errors"
	"time"
)

var (
	ErrImageNotExists = errors.New("image not exists")
)

type Repository interface {
	Get(ctx context.Context, date string) ([]byte, error)
	Upload(ctx context.Context, date time.Time, imageData []byte) error
}

type Image struct {
	Title     string
	Extension string
	Data      []byte
}

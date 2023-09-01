package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"

	"github.com/diianpro/stellar-scope/internal/domain/apod"
	"github.com/diianpro/stellar-scope/internal/domain/image"
)

func (s *Service) GetByDate(ctx context.Context, date time.Time) (*apod.Data, error) {
	data, err := s.apodRps.GetByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	data.Image.Data, err = s.imageRps.Get(ctx, date.Format(time.DateOnly))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) GetAll(ctx context.Context, limit, offset int) ([]apod.Data, error) {
	metadata, err := s.apodRps.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, data := range metadata {
		data.Image.Data, err = s.imageRps.Get(ctx, data.Date)
		if err != nil {
			if errors.Is(err, image.ErrImageNotExists) {
				log.Error(fmt.Errorf("%w: %s", err, data.Image.Title))
				continue
			}
			return nil, err
		}
	}
	return metadata, nil
}

func (s *Service) ObserveDailyImage(ctx context.Context) error {
	sheduler := gocron.NewScheduler(time.UTC)
	defer sheduler.Stop()

	errChan := make(chan error, 1)
	defer close(errChan)
	job, err := sheduler.Every(1).Day().At("00:00").Do(func(errorChan chan<- error) {
		err := linearBackOff(11, time.Hour, func() error {
			now := time.Now()
			data, err := s.imageProvider.GetMetadata(ctx, now)
			if err != nil {
				return err
			}
			err = s.imageRps.Upload(ctx, now, data.Image.Data)
			if err != nil {
				return err
			}
			err = s.apodRps.Create(ctx, data)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			errorChan <- err
		}
	}, errChan)
	if err != nil {
		return err
	}
	defer sheduler.Remove(job)

	sheduler.StartAsync()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err = <-errChan:
			log.Error(fmt.Errorf("picture of the day by %s not uploaded: %w", time.Now().Format(time.DateOnly), err))
		}
	}
}

func linearBackOff(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
		} else {
			return nil
		}

		log.Warnf("Retry: %d.", i)

		time.Sleep(delay * time.Duration(i+1))
	}
	return err
}

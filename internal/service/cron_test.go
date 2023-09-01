package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/require"
)

func TestCron(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sheduler := gocron.NewScheduler(time.UTC)
	defer sheduler.Stop()

	errChan := make(chan error, 1)
	defer close(errChan)
	job, err := sheduler.Every(10).Second().Do(func(errorChan chan<- error) {
		fmt.Println(time.Now())
		errorChan <- errors.New("error")
	}, errChan)
	require.NoError(t, err)
	defer sheduler.Remove(job)

	sheduler.StartAsync()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("app canceled")
				return
			case err = <-errChan:
				fmt.Println(err)
			}
		}
	}()
	<-time.After(time.Minute)
	cancel()
}

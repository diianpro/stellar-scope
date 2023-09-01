package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/diianpro/stellar-scope/app/waiter"
	"github.com/diianpro/stellar-scope/internal/config"
	"github.com/diianpro/stellar-scope/internal/provider/apod"
	"github.com/diianpro/stellar-scope/internal/service"
	"github.com/diianpro/stellar-scope/internal/storage/postgres"
	imagestor "github.com/diianpro/stellar-scope/internal/storage/s3"
	appServer "github.com/diianpro/stellar-scope/internal/transport/http"
)

const (
	defaultAddr    = ":8080"
	defaultAppName = "stellar_scope"
)

type LoadConfigFn func() (config.Config, error)

type App struct {
	cfg      config.Config
	server   *appServer.Server
	waiter   waiter.Waiter
	ctx      context.Context
	cancelFn context.CancelFunc
	name     string
}

func New(loadConfigFn LoadConfigFn) *App {
	ctx, cancelFn := context.WithCancel(context.Background())
	cfg, err := loadConfigFn()
	if err != nil {
		log.Fatal("failed to load config")
	}

	w := waiter.NewWaiter(ctx, cancelFn)

	return &App{
		cfg:      cfg,
		waiter:   w,
		ctx:      ctx,
		cancelFn: cancelFn,
		name:     defaultAppName,
	}
}

func (a *App) Start() {
	defer a.cancelFn()

	repo, err := postgres.New(a.ctx, &a.cfg.Postgres)
	if err != nil {
		log.Fatal("Could not setup storage.")
	}
	defer repo.Close()

	awsConfig, err := awsconf.LoadDefaultConfig(a.ctx, awsconf.WithRegion(a.cfg.ImageStorage.Region))
	if err != nil {
		log.Fatal(err)
	}

	s3Storage := s3.NewFromConfig(awsConfig)
	imageRps := imagestor.New(a.cfg.ImageStorage.Bucket, s3Storage)

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	httpClient := &http.Client{
		Transport: t,
		Timeout:   time.Second * 5,
	}

	imageProvider := apod.NewApodClient(a.cfg.APOD, httpClient)

	srv := service.New(repo, imageRps, imageProvider)
	handler := appServer.NewHandler(srv)

	a.server = appServer.New(handler)

	a.waitForServer()
	a.waitForWorker(srv.ObserveDailyImage)

	if err = a.waiter.Wait(); err != nil {
		log.Error(fmt.Errorf("app crash with error: %w", err))
	}
}

func (a *App) Stop() {
	a.cancelFn()
}

func (a *App) waitForServer() {
	a.waiter.Add(func(ctx context.Context) error {
		defer log.Info("server has been shutdown")

		group, gCtx := errgroup.WithContext(ctx)
		group.Go(func() error {
			defer log.Info("public server exited")
			log.Infof("starting server at: %s", defaultAddr)
			err := a.server.ServePublic(defaultAddr)
			if err != nil && err != http.ErrServerClosed {
				return err
			}
			return nil
		})

		group.Go(func() error {
			<-gCtx.Done()
			log.Info("shutting down the server")
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			if err := a.server.ShutdownPublic(ctx); err != nil {
				log.Warn(fmt.Errorf("error while shutting down the server: %w", err))
			}
			return nil
		})

		return group.Wait()
	})
}

func (a *App) waitForWorker(workerFn func(ctx context.Context) error) {
	a.waiter.Add(func(ctx context.Context) error {
		group, gCtx := errgroup.WithContext(ctx)
		group.Go(func() error {
			return workerFn(gCtx)
		})
		return group.Wait()
	})
}

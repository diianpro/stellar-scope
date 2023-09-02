package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	db        *Repository
	container *Container
	suite.Suite
}

func (i *IntegrationTestSuite) SetupSuite() {
	var err error

	ctx := context.Background()

	cfg := &Config{
		MinConns: 1,
		MaxConns: 2,
	}
	i.container, err = NewContainer(cfg, func() error {
		i.db, err = New(ctx, cfg)
		if err != nil {
			return err
		}
		return ApplyMigrate(cfg.URL, "../../../migration")
	})
	i.Require().NoError(err)
}

func (i *IntegrationTestSuite) TearDownSuite() {
	err := i.container.Purge()
	i.Assert().NoError(err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

package postgres

import (
	"context"
	"github.com/diianpro/stellar-scope/internal/domain"
	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
	"testing"
	"time"
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
		Username: "postgres",
		Password: "postgres",
	}
	i.container, err = NewContainer(cfg, func() error {
		i.db, err = New(ctx, cfg)
		if err != nil {
			return err
		}

		ApplyMigrate(cfg.URL)
		return nil
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

func (i *IntegrationTestSuite) TestIntegrationCreate() {
	ctx := context.Background()

	data := domain.ApodData{
		Title:       "test",
		Explanation: "test",
		Date:        time.Now().Format("2006-01-02"),
		ImageLink:   "",
		Copyright:   "test",
	}

	_, err := i.db.Create(ctx, &data, "test_link")

	i.Require().NoError(err)

	dataFound, err := i.db.GetByDate(ctx, time.Now())
	i.Require().NoError(err)

	assert.Equal(i.T(), data.Title, dataFound.Title)
	assert.Equal(i.T(), data.Explanation, dataFound.Explanation)
	assert.Equal(i.T(), data.Date, dataFound.Date)
	// assert.Equal(i.T(), data.ImageLink, dataFound.ImageLink)
	assert.Equal(i.T(), data.Copyright, dataFound.Copyright)

}

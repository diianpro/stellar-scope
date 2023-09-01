package postgres

import (
	"context"
	"time"

	"github.com/diianpro/stellar-scope/internal/domain/apod"
)

func (i *IntegrationTestSuite) TestIntegrationCreate() {
	ctx := context.Background()

	data := apod.Data{
		Title:       "test",
		Explanation: "test",
		Date:        time.Now().Format("2006-01-02"),
		ImageLink:   "test_link",
		Copyright:   "test",
	}

	err := i.db.Create(ctx, &data, data.ImageLink)
	i.Require().NoError(err)

	dataFound, err := i.db.GetByDate(ctx, time.Now())
	i.Require().NoError(err)

	i.Require().Equal(data.Title, dataFound.Title)
	i.Require().Equal(data.Explanation, dataFound.Explanation)
	i.Require().Equal(data.Date, dataFound.Date)
	i.Require().Equal(data.ImageLink, dataFound.ImageLink)
	i.Require().Equal(data.Copyright, dataFound.Copyright)
}

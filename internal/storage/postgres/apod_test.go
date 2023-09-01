package postgres

import (
	"context"
	"time"

	"github.com/diianpro/stellar-scope/internal/domain/apod"
)

func (i *IntegrationTestSuite) TestIntegrationCreate() {
	ctx := context.Background()

	data := apod.Data{
		Explanation: "test",
		Date:        "test",
		Image:       nil,
		Copyright:   "test",
	}

	err := i.db.Create(ctx, &data)
	i.Require().NoError(err)

	dataFound, err := i.db.GetByDate(ctx, time.Now())
	i.Require().NoError(err)

	i.Require().Equal(data.Explanation, dataFound.Explanation)
	i.Require().Equal(data.Date, dataFound.Date)
	i.Require().Equal(data.Image, dataFound.Image)
	i.Require().Equal(data.Copyright, dataFound.Copyright)
}

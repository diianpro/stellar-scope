package postgres

import (
	"context"
	"time"

	"github.com/diianpro/stellar-scope/internal/domain/apod"
	"github.com/diianpro/stellar-scope/internal/domain/image"
)

func (i *IntegrationTestSuite) TestIntegrationCreate() {
	ctx := context.Background()

	data := apod.Data{
		Explanation: "test",
		Date:        time.Now().Format(time.DateOnly),
		Image:       &image.Image{},
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

func (i *IntegrationTestSuite) TestIntegrationGetAll() {
	ctx := context.Background()

	testData := []apod.Data{
		{
			Explanation: "test",
			Date:        time.Now().Format(time.DateOnly),
			Image:       &image.Image{},
			Copyright:   "test"},
		{
			Explanation: "test-2",
			Date:        "2023-08-07",
			Image:       &image.Image{},
			Copyright:   "test",
		},
		{
			Explanation: "test-3",
			Date:        "2023-08-06",
			Image:       &image.Image{},
			Copyright:   "test",
		},
	}

	for _, data := range testData {
		err := i.db.Create(ctx, &data)
		i.Require().NoError(err)
	}

	dataFound, err := i.db.GetAll(ctx, 10, 0)
	i.Require().NoError(err)
	i.Require().Len(dataFound, len(testData))

	for idx, found := range dataFound {
		i.Require().Equal(testData[idx].Explanation, found.Explanation)
		i.Require().Equal(testData[idx].Date, found.Date)
		i.Require().Equal(testData[idx].Image.Title, found.Image.Title)
		i.Require().Equal(testData[idx].Copyright, found.Copyright)
	}
}

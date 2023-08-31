package postgres

import (
	"context"
	"github.com/diianpro/stellar-scope/internal/domain"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/gommon/log"
	"time"
)

func (r Repository) GetByDate(ctx context.Context, date time.Time) (*domain.ApodData, error) {
	rows := r.db.QueryRow(ctx, `SELECT  title, explanation, "date", image_link, copyright  FROM images WHERE "date" = $1`, date.Format("2006-01-02"))
	return r.getByDate(rows)
}

func (r Repository) getByDate(rows pgx.Row) (*domain.ApodData, error) {
	apod := &domain.ApodData{}
	err := rows.Scan(&apod.Title, &apod.Explanation, &apod.Date, &apod.ImageLink, &apod.Copyright)
	if err != nil {
		log.Errorf("images db: getByDate error: %v", err)
		return nil, err
	}
	return apod, nil
}

func (r Repository) GetAll(ctx context.Context) ([]domain.ApodData, error) {
	row, err := r.db.Query(ctx, `SELECT  title, explanation,"date", image_link, copyright FROM images`)
	if err != nil {
		log.Errorf("images db: getAll error: %v", err)
		return nil, err
	}
	defer row.Close()

	var result []domain.ApodData

	for row.Next() {
		newRow, err := r.getByDate(row)
		if err != nil {
			log.Errorf("images db get data: getAll error: %v", err)
			return nil, err
		}

		result = append(result, *newRow)
	}

	return result, err
}

func (r Repository) Create(ctx context.Context, data *domain.ApodData, imageLink string) (int64, error) {
	res, err := r.db.Exec(ctx, `INSERT INTO images ( title, explanation,"date",  image_link, copyright) 
		VALUES ($1, $2, $3, $4, $5)`, data.Title, data.Explanation, time.Now().Format("2006-01-02"), imageLink, data.Copyright)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil

}

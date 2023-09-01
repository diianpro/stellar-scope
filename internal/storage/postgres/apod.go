package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/gommon/log"

	"github.com/diianpro/stellar-scope/internal/domain/apod"
	"github.com/diianpro/stellar-scope/internal/domain/image"
)

func (r *Repository) GetByDate(ctx context.Context, date time.Time) (*apod.Data, error) {
	rows := r.db.QueryRow(ctx, `SELECT  title, explanation, "date", image_extension, copyright
		FROM image WHERE "date" = $1`, date.Format(time.DateOnly))
	return r.getByDate(rows)
}

func (r *Repository) GetAll(ctx context.Context, limit, offset int) ([]apod.Data, error) {
	row, err := r.db.Query(ctx, `SELECT  title, explanation,"date", image_extension, copyright
		FROM image LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		log.Errorf("images db: getAll error: %v", err)
		return nil, err
	}
	defer row.Close()

	var result []apod.Data
	for row.Next() {
		apod, err := r.getByDate(row)
		if err != nil {
			log.Errorf("images db get data: getAll error: %v", err)
			return nil, err
		}

		result = append(result, *apod)
	}

	return result, err
}

func (r *Repository) Create(ctx context.Context, data *apod.Data) error {
	_, err := r.db.Exec(ctx, `INSERT INTO image (title, explanation,"date",  image_extension, copyright) 
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
		data.Image.Title,
		data.Explanation,
		time.Now().Format(time.DateOnly),
		data.Image.Extension,
		data.Copyright)
	if err != nil {
		return err
	}
	return nil

}
func (r *Repository) getByDate(row pgx.Row) (*apod.Data, error) {
	apod := &apod.Data{
		Image: &image.Image{},
	}
	err := row.Scan(&apod.Image.Title, &apod.Explanation, &apod.Date, &apod.Image.Extension, &apod.Copyright)
	if err != nil {
		log.Errorf("images db: getByDate error: %v", err)
		return nil, err
	}
	return apod, nil
}

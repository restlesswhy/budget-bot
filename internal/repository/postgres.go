package repository

import (
	"bot/internal"
	"bot/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

type repos struct {
	pool *pgxpool.Pool
}

func NewRepos(pool *pgxpool.Pool) internal.Repository {
	return &repos{
		pool: pool,
	}
}

func (r *repos) WriteButton(btn *models.Buttons) error {
	q := `INSERT INTO buttons (button_id, message_id, amount, firstname, lastname, username)
			VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := r.pool.Exec(context.Background(),
		q,
		btn.ID,
		btn.MessageID,
		btn.Amount,
		btn.Firstname,
		btn.Lastname,
		btn.Username,
	)
	if err != nil {
		return errors.Wrap(err, "can't exec button")
	}

	return nil
}

func (r *repos) WriteTransaction(tx *models.Transaction) error {
	q := `INSERT INTO transactions (button_id, amount, category, time) 
			VALUES (
				$1,
				(select amount from buttons where button_id=$1),
				$2,
				$3
			);`

	_, err := r.pool.Exec(context.Background(),
		q,
		tx.ButtonID,
		tx.Category,
		tx.Time,
	)
	if err != nil {
		return errors.Wrap(err, "can't exec button")
	}

	return nil
}

func (r *repos) GetMonthReport() (*models.TotalReport, error) {
	res := &models.TotalReport{
		SpendsSet: make([]*models.Report, 0),
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	q := `SELECT category, SUM(amount)
			FROM transactions
			WHERE time BETWEEN $1 and $2
			GROUP BY category;`

	rows, err := r.pool.Query(context.Background(), q, firstOfMonth.Format(TIME_FORMAT), lastOfMonth.Format(TIME_FORMAT))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		scanRes := &models.Report{}

		scanRes.Category = values[0].(string)
		scanRes.Amount = values[1].(int64)

		res.SpendsSet = append(res.SpendsSet, scanRes)
	}

	q = `SELECT SUM(amount)
	FROM transactions
	WHERE time BETWEEN $1 and $2;`

	err = r.pool.QueryRow(context.Background(), q, firstOfMonth.Format(TIME_FORMAT), lastOfMonth.Format(TIME_FORMAT)).Scan(&res.TotalSpend)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repos) GetDayReport() error {
	return nil
}

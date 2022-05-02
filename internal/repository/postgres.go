package repository

import (
	"bot/internal"
	"bot/internal/models"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
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
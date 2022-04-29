package repository

import (
	"bot/internal"
	"bot/internal/models"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type dataRepo struct {
	pool *pgxpool.Pool
}

func NewDataRepo(pool *pgxpool.Pool) internal.Repository {
	return &dataRepo{
		pool: pool,
	}
}

func (d *dataRepo) WriteMessage(msg *models.Message) error {
	q := `INSERT INTO messages (id, text, firstname, lastname, username) VALUES ($1, $2, $3, $4, $5);`

	_, err := d.pool.Exec(context.Background(),
		q,
		msg.ID,
		msg.Text,
		msg.Firstname,
		msg.Lastname,
		msg.Username,
	)
	if err != nil {
		return errors.Wrap(err, "can't exec message")
	}

	return nil
}

func (d *dataRepo) WriteButton(msg *models.Buttons) error {
	q := `INSERT INTO buttons (id, message_relation_id) VALUES ($1, $2);`

	_, err := d.pool.Exec(context.Background(),
		q,
		msg.ID,
		msg.MessageRelationID,
	)
	if err != nil {
		return errors.Wrap(err, "can't exec button")
	}

	return nil
}

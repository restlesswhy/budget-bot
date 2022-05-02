package internal

import "bot/internal/models"

type Repository interface {
	// WriteMessage(msg *models.Message) error
	WriteButton(btn *models.Buttons) error
	WriteTransaction(tx *models.Transaction) error
}
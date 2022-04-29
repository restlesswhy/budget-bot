package internal

import "bot/internal/models"

type Repository interface {
	WriteMessage(msg *models.Message) error
	WriteButton(msg *models.Buttons) error
}
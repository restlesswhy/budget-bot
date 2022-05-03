package internal

import "bot/internal/models"

type Repository interface {
	WriteButton(btn *models.Buttons) error
	WriteTransaction(tx *models.Transaction) error
	GetMonthReport() ([]*models.Report, error)
	GetDayReport() error
}
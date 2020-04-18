package service

import "github.com/poofik33/db-technopark/internal/models"

type Usecase interface {
	GetStatus() (*models.Status, error)
	DeleteAll() error
}

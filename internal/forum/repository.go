package forum

import "github.com/poofik33/db-technopark/internal/models"

type Repository interface {
	InsertInto(f *models.Forum) error
	GetBySlug(slug string) (*models.Forum, error)
}

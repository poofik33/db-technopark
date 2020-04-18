package post

import "github.com/poofik33/db-technopark/internal/models"

type Repository interface {
	InsertInto([]*models.Post) error
	GetByThread(string) ([]*models.Post, error)
	GetCountByForumSlug(string) (uint64, error)
	Update(*models.Post) error
}

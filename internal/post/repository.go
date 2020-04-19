package post

import "github.com/poofik33/db-technopark/internal/models"

type Repository interface {
	InsertInto([]*models.Post) error
	GetByThread(uint64, uint64, uint64, string, bool) ([]*models.Post, error)
	GetByID(uint64) (*models.Post, error)
	GetCountByForumSlug(string) (uint64, error)
	Update(*models.Post) error
}

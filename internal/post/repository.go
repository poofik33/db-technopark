package post

import "github.com/poofik33/db-technopark/internal/models"

type Repository interface {
	InsertInto([]*models.Post) error
	GetByThread(uint64, uint64, uint64, string, bool) ([]*models.Post, error)
	CheckParentPosts([]*models.Post, uint64) (bool, error)
	GetByID(uint64) (*models.Post, error)
	GetCountByForumID(uint64) (uint64, error)
	Update(*models.Post) error
}

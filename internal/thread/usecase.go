package thread

import (
	"github.com/poofik33/db-technopark/internal/models"
)

type Usecase interface {
	AddThread(*models.Thread) (*models.Thread, error)
	CreatePosts(string, []*models.Post) ([]*models.Post, error)
	GetBySlugOrID(string) (*models.Thread, error)
	GetThreadPosts(string, uint64, uint64, string, bool) ([]*models.Post, error)
	Update(string, *models.Thread) (*models.Thread, error)
	Vote(string, *models.Vote) (*models.Thread, error)
}

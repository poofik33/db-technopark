package thread

import (
	"github.com/poofik33/db-technopark/internal/models"
)

type Usecase interface {
	AddThread(*models.Thread) (*models.Thread, error)
	GetBySlugOrID(string) (*models.Thread, error)
	GetThreadPosts(string) ([]*models.Post, error)
	Update(*models.Thread) error
	Vote(string, models.Vote) (*models.Thread, error)
}

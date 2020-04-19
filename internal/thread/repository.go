package thread

import (
	"github.com/poofik33/db-technopark/internal/models"
)

type Repository interface {
	InsertInto(*models.Thread) error
	GetByID(uint64) (*models.Thread, error)
	GetBySlug(string) (*models.Thread, error)
	GetByForumSlug(string, uint64, string, bool) ([]*models.Thread, error)
	GetCountByForumSlug(string) (uint64, error)
	Update(*models.Thread) error
}

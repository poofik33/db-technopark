package vote

import "github.com/poofik33/db-technopark/internal/models"

type Repository interface {
	GetThreadVotes(uint64) (int64, error)
	InsertInto(*models.Vote) error
}

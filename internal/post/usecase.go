package post

import "github.com/poofik33/db-technopark/internal/models"

type Usecase interface {
	GetPostDetails(uint64, []string) (*models.PostFull, error)
	UpdatePost(*models.Post) (*models.Post, error)
}

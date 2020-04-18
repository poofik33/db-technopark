package forum

import "github.com/poofik33/db-technopark/internal/models"

type Usecase interface {
	AddForum(*models.Forum) (*models.Forum, error)
	GetForumBySlug(string) (*models.Forum, error)
	GetForumUsers(string, uint64, string, bool) ([]*models.User, error)
	GetForumThreads(string, uint64, string, bool) ([]*models.Thread, error)
}

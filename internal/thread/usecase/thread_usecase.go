package usecase

import (
	"github.com/poofik33/db-technopark/internal/forum"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/thread"
	"github.com/poofik33/db-technopark/internal/tools"
	"github.com/poofik33/db-technopark/internal/user"
)

type ThreadUsecase struct {
	forumRepo  forum.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
}

func NewThreadUsecase(tr thread.Repository, ur user.Repository, fr forum.Repository) thread.Usecase {
	return &ThreadUsecase{
		forumRepo:  fr,
		threadRepo: tr,
		userRepo:   ur,
	}
}

func (tu *ThreadUsecase) AddThread(t *models.Thread) (*models.Thread, error) {
	if _, err := tu.forumRepo.GetBySlug(t.Forum); err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}

		return nil, err
	}

	if _, err := tu.userRepo.GetByNickname(t.Author); err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrUserDoesntExists
		}
	}

	if t.Slug != "" {
		returnThread, err := tu.threadRepo.GetBySlug(t.Slug)
		if err != nil && err != tools.ErrDoesntExists {
			return nil, err
		}
		if returnThread != nil {
			return returnThread, tools.ErrExistWithSlug
		}
	}

	if err := tu.threadRepo.InsertInto(t); err != nil {
		return nil, err
	}

	return t, nil
}

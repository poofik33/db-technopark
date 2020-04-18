package usecase

import (
	"github.com/poofik33/db-technopark/internal/forum"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/post"
	"github.com/poofik33/db-technopark/internal/thread"
	"github.com/poofik33/db-technopark/internal/tools"
	"github.com/poofik33/db-technopark/internal/user"
)

type ForumUsecase struct {
	forumRepo  forum.Repository
	postRepo   post.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
}

func NewForumUsecase(fr forum.Repository, ur user.Repository) forum.Usecase {
	return &ForumUsecase{
		forumRepo: fr,
		userRepo:  ur,
	}
}

func (fu *ForumUsecase) AddForum(forum *models.Forum) (*models.Forum, error) {
	returnForum, err := fu.forumRepo.GetBySlug(forum.Slug)
	if err != nil && err != tools.ErrDoesntExists {
		return nil, err
	}

	if returnForum != nil {
		return returnForum, tools.ErrExistWithSlug
	}

	if _, err := fu.userRepo.GetByNickname(forum.AdminNickname); err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrUserDoesntExists
		}
	}

	if err = fu.forumRepo.InsertInto(forum); err != nil {
		return nil, err
	}

	return forum, nil
}

func (fu *ForumUsecase) GetForumBySlug(slug string) (*models.Forum, error) {
	returnForum, err := fu.forumRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	if returnForum == nil {
		return nil, tools.ErrDoesntExists
	}

	postCount, err := fu.postRepo.GetCountByForumSlug(slug)
	threadCount, err := fu.threadRepo.GetCountByForumSlug(slug)
	if err != nil {
		return nil, err
	}

	returnForum.PostsCount = postCount
	returnForum.ThreadsCount = threadCount

	return returnForum, nil
}

func (fu *ForumUsecase) GetForumThreads(
	slug string, limit uint64, since string, desc bool) ([]*models.Thread, error) {
	if _, err := fu.forumRepo.GetBySlug(slug); err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}
	}

	returnThreads, err := fu.threadRepo.GetByForumSlug(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}

	return returnThreads, nil
}

func (fu *ForumUsecase) GetForumUsers(
	slug string, limit uint64, since string, desc bool) ([]*models.User, error) {
	if _, err := fu.forumRepo.GetBySlug(slug); err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}
	}

	returnUsers, err := fu.userRepo.GetUsersByForum(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}

	return returnUsers, nil
}

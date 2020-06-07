package usecase

import (
	"github.com/poofik33/db-technopark/internal/forum"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/post"
	"github.com/poofik33/db-technopark/internal/thread"
	"github.com/poofik33/db-technopark/internal/tools"
	"github.com/poofik33/db-technopark/internal/user"
	"github.com/poofik33/db-technopark/internal/vote"
)

type PostUsecase struct {
	forumRepo  forum.Repository
	postRepo   post.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
	voteRepo   vote.Repository
}

func NewPostUsecase(pr post.Repository, fr forum.Repository, vr vote.Repository,
	tr thread.Repository, ur user.Repository) post.Usecase {
	return &PostUsecase{
		forumRepo:  fr,
		postRepo:   pr,
		threadRepo: tr,
		userRepo:   ur,
		voteRepo:   vr,
	}
}

func (pUC *PostUsecase) GetPostDetails(id uint64, related []string) (*models.PostFull, error) {
	p, err := pUC.postRepo.GetByID(id)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrPostDoesntExists
		}

		return nil, err
	}

	returnPost := &models.PostFull{PostData: p}

	for _, rel := range related {
		switch rel {
		case "user":
			u, err := pUC.userRepo.GetByNickname(p.Author)
			if err != nil {
				if err == tools.ErrDoesntExists {
					return nil, tools.ErrUserDoesntExists
				}

				return nil, err
			}

			returnPost.Author = u
		case "thread":
			t, err := pUC.threadRepo.GetByID(p.ThreadID)
			if err != nil {
				if err == tools.ErrDoesntExists {
					return nil, tools.ErrThreadDoesntExists
				}

				return nil, err
			}

			returnPost.Thread = t
		case "forum":
			f, err := pUC.forumRepo.GetBySlug(p.Forum)
			if err != nil {
				if err == tools.ErrDoesntExists {
					return nil, tools.ErrForumDoesntExists
				}

				return nil, err
			}

			postCount, err := pUC.postRepo.GetCountByForumID(f.ID)
			threadCount, err := pUC.threadRepo.GetCountByForumID(f.ID)
			if err != nil {
				return nil, err
			}

			f.ThreadsCount = threadCount
			f.PostsCount = postCount
			returnPost.Forum = f
		}
	}

	return returnPost, nil
}

func (pUC *PostUsecase) UpdatePost(p *models.Post) (*models.Post, error) {
	updPost, err := pUC.postRepo.GetByID(p.ID)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrPostDoesntExists
		}

		return nil, err
	}
	if p.Message == "" || p.Message == updPost.Message {
		return updPost, nil
	}
	updPost.Message = p.Message
	updPost.IsEdited = true

	if err = pUC.postRepo.Update(updPost); err != nil {
		return nil, err
	}

	return updPost, nil
}

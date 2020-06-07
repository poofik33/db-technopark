package usecase

import (
	"github.com/poofik33/db-technopark/internal/forum"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/post"
	"github.com/poofik33/db-technopark/internal/thread"
	"github.com/poofik33/db-technopark/internal/tools"
	"github.com/poofik33/db-technopark/internal/user"
	"github.com/poofik33/db-technopark/internal/vote"
	"strconv"
)

type ThreadUsecase struct {
	forumRepo  forum.Repository
	postRepo   post.Repository
	threadRepo thread.Repository
	userRepo   user.Repository
	voteRepo   vote.Repository
}

func NewThreadUsecase(tr thread.Repository, ur user.Repository,
	fr forum.Repository, pr post.Repository, vr vote.Repository) thread.Usecase {
	return &ThreadUsecase{
		forumRepo:  fr,
		postRepo:   pr,
		threadRepo: tr,
		userRepo:   ur,
		voteRepo:   vr,
	}
}

func (tUC *ThreadUsecase) AddThread(t *models.Thread) (*models.Thread, error) {
	f, err := tUC.forumRepo.GetBySlug(t.Forum)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrForumDoesntExists
		}

		return nil, err
	}

	u, err := tUC.userRepo.GetByNickname(t.Author)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrUserDoesntExists
		}
	}

	if t.Slug != "" {
		returnThread, err := tUC.threadRepo.GetBySlug(t.Slug)
		if err != nil && err != tools.ErrDoesntExists {
			return nil, err
		}
		if returnThread != nil {
			return returnThread, tools.ErrExistWithSlug
		}
	}

	t.Forum = f.Slug
	t.Author = u.Nickname
	t.AuthorID = u.ID
	t.ForumID = f.ID

	if err := tUC.threadRepo.InsertInto(t); err != nil {
		return nil, err
	}

	return t, nil
}

func (tUC *ThreadUsecase) CreatePosts(slugOrID string, posts []*models.Post) ([]*models.Post, error) {
	t := &models.Thread{}

	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		t, err = tUC.threadRepo.GetBySlug(slugOrID)
	} else {
		t, err = tUC.threadRepo.GetByID(id)
	}

	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}

		return nil, err
	}

	_, err = tUC.userRepo.CheckNicknames(posts)
	if err != nil {
		return nil, err
	}

	_, err = tUC.postRepo.CheckParentPosts(posts, t.ID)
	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return []*models.Post{}, nil
	}

	for _, p := range posts {
		p.ThreadID = t.ID
		p.Forum = t.Forum
		p.ForumID = t.ForumID
	}
	if err = tUC.postRepo.InsertInto(posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (tUC *ThreadUsecase) GetBySlugOrID(slugOrID string) (*models.Thread, error) {
	t := &models.Thread{}

	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		t, err = tUC.threadRepo.GetBySlug(slugOrID)
	} else {
		t, err = tUC.threadRepo.GetByID(id)
	}

	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}

		return nil, err
	}

	return t, nil
}

func (tUC *ThreadUsecase) Update(slugOrID string, thread *models.Thread) (*models.Thread, error) {
	t := &models.Thread{}

	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		t, err = tUC.threadRepo.GetBySlug(slugOrID)
	} else {
		t, err = tUC.threadRepo.GetByID(id)
	}

	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}

		return nil, err
	}

	if thread.About != "" {
		t.About = thread.About
	}
	if thread.Title != "" {
		t.Title = thread.Title
	}

	if err = tUC.threadRepo.Update(t); err != nil {
		return nil, err
	}

	t.Votes, err = tUC.voteRepo.GetThreadVotes(t.ID)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (tUC *ThreadUsecase) Vote(slugOrID string, v *models.Vote) (*models.Thread, error) {
	t := &models.Thread{}

	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		t, err = tUC.threadRepo.GetBySlug(slugOrID)
	} else {
		t, err = tUC.threadRepo.GetByID(id)
	}

	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}

		return nil, err
	}

	u, err := tUC.userRepo.GetByNickname(v.Nickname)
	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrUserDoesntExists
		}
	}

	v.ThreadID = t.ID
	v.UserID = u.ID
	if err = tUC.voteRepo.InsertInto(v); err != nil {
		return nil, err
	}

	t.Votes, err = tUC.voteRepo.GetThreadVotes(t.ID)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (tUC *ThreadUsecase) GetThreadPosts(slugOrID string, limit uint64,
	since uint64, sort string, desc bool) ([]*models.Post, error) {
	t := &models.Thread{}

	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		t, err = tUC.threadRepo.GetBySlug(slugOrID)
	} else {
		t, err = tUC.threadRepo.GetByID(id)
	}

	if err != nil {
		if err == tools.ErrDoesntExists {
			return nil, tools.ErrThreadDoesntExists
		}

		return nil, err
	}

	returnPosts, err := tUC.postRepo.GetByThread(t.ID, limit, since, sort, desc)
	if err != nil {
		return nil, err
	}

	return returnPosts, err
}

package repository

import (
	"database/sql"
	"github.com/poofik33/db-technopark/internal/post"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) post.Repository {
	return &PostRepository{
		db: db,
	}
}

func (pr *PostRepository) GetCountByForumSlug(slug string) (uint64, error) {
	var count uint64
	if err := pr.db.QueryRow("SELECT count(*) from posts WHERE forum = $1", slug).
		Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

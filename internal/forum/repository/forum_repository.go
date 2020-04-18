package repository

import (
	"database/sql"
	"github.com/poofik33/db-technopark/internal/forum"
	"github.com/poofik33/db-technopark/internal/models"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) forum.Repository {
	return &ForumRepository{
		db: db,
	}
}

func (fr *ForumRepository) InsertInto(forum *models.Forum) error {
	if _, err := fr.db.Exec("INSERT INTO forums (slug, admin, title) "+
		"VALUES ($1, $2, $3)", forum.Slug, forum.AdminNickname, forum.Title); err != nil {
		return err
	}

	return nil
}

func (fr *ForumRepository) GetBySlug(slug string) (*models.Forum, error) {
	returnForum := &models.Forum{}
	if err := fr.db.QueryRow("SELECT slug, admin, title FROM forums "+
		"WHERE slug = $1", slug).Scan(&returnForum.Slug, &returnForum.AdminNickname, &returnForum.Title); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return returnForum, nil
}

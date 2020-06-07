package repository

import (
	"github.com/jackc/pgx"
	"github.com/poofik33/db-technopark/internal/forum"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/tools"
)

type ForumRepository struct {
	db *pgx.ConnPool
}

func NewForumRepository(db *pgx.ConnPool) forum.Repository {
	return &ForumRepository{
		db: db,
	}
}

func (fr *ForumRepository) InsertInto(forum *models.Forum) error {
	if _, err := fr.db.Exec("INSERT INTO forums (slug, admin, title) "+
		"VALUES ($1, $2, $3)", forum.Slug, forum.AdminID, forum.Title); err != nil {
		return err
	}

	return nil
}

func (fr *ForumRepository) GetBySlug(slug string) (*models.Forum, error) {
	returnForum := &models.Forum{}
	if err := fr.db.QueryRow("SELECT f.id, f.slug, u.nickname, f.title FROM forums as f "+
		"JOIN users as u ON (u.id = f.admin) "+
		"WHERE lower(slug) = lower($1)", slug).Scan(&returnForum.ID, &returnForum.Slug, &returnForum.AdminNickname,
		&returnForum.Title); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}

	return returnForum, nil
}

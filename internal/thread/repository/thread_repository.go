package repository

import (
	"database/sql"
	"fmt"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/thread"
	"github.com/poofik33/db-technopark/internal/tools"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) thread.Repository {
	return &ThreadRepository{
		db: db,
	}
}

func (tr *ThreadRepository) GetCountByForumSlug(slug string) (uint64, error) {
	var count uint64
	if err := tr.db.QueryRow("SELECT count(*) FROM threads WHERE forum = $1", slug).
		Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (tr *ThreadRepository) InsertInto(t *models.Thread) error {
	if err := tr.db.QueryRow("INSERT INTO threads "+
		"(slug, author, title, message, forum, created) "+
		"VALUES (NULLIF ($1, ''), $2, $3, $4, $5, $6) RETURNING id",
		t.Slug, t.Author, t.Title, t.About, t.Forum, t.CreationDate).
		Scan(&t.ID); err != nil {
		return err
	}

	return nil
}

func (tr *ThreadRepository) GetBySlug(slug string) (*models.Thread, error) {
	t := &models.Thread{}
	if err := tr.db.QueryRow("SELECT id, author, created, forum, message, "+
		"coalesce (slug, ''), title WHERE slug = $1", slug).
		Scan(&t.ID, &t.Author, &t.CreationDate, &t.Forum, &t.About, &t.Slug, &t.Title); err != nil {
		if err == sql.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}

		return nil, err
	}

	return t, nil
}

func (tr *ThreadRepository) GetByID(id uint64) (*models.Thread, error) {
	t := &models.Thread{}
	if err := tr.db.QueryRow("SELECT id, author, created, forum, message, "+
		"coalesce (slug, ''), title WHERE id = $1", id).
		Scan(&t.ID, &t.Author, &t.CreationDate, &t.Forum, &t.About, &t.Slug, &t.Title); err != nil {
		if err == sql.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}

		return nil, err
	}

	return t, nil
}

func (tr *ThreadRepository) GetByForumSlug(
	slug string, limit uint64, since string, desc bool) ([]*models.Thread, error) {
	returnThreads := []*models.Thread{}

	queryString := "SELECT id, author, created, forum, message, " +
		"coalesce (slug, ''), title WHERE forum = $1"

	orderString := " ORDER BY created"

	if desc {
		orderString += " DESC"
	}

	if limit != 0 {
		orderString += fmt.Sprintf(" LIMIT %d", limit)
	}

	var rows *sql.Rows
	var err error

	if since != "" {
		queryString += " AND created >= $2"
		rows, err = tr.db.Query(queryString+orderString, slug, since)
	} else {
		rows, err = tr.db.Query(queryString+orderString, slug)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := &models.Thread{}
		err = rows.Scan(&t.ID, &t.Author, &t.CreationDate, &t.Forum, &t.About, &t.Slug, &t.Title)
		if err != nil {
			return nil, err
		}

		returnThreads = append(returnThreads, t)
	}

	return returnThreads, nil
}

func (tr *ThreadRepository) Update(t *models.Thread) error {
	if _, err := tr.db.Exec("UPDATE threads SET message = $2, title = $3 WHERE id = $1",
		t.ID, t.About, t.Title); err != nil {
		return err
	}

	return nil
}

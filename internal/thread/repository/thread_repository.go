package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/thread"
	"github.com/poofik33/db-technopark/internal/tools"
)

type ThreadRepository struct {
	db *pgx.ConnPool
}

func NewThreadRepository(db *pgx.ConnPool) thread.Repository {
	return &ThreadRepository{
		db: db,
	}
}

func (tr *ThreadRepository) GetCountByForumID(id uint64) (uint64, error) {
	var count uint64
	if err := tr.db.QueryRow("SELECT count(*) FROM threads WHERE forum = $1 ", id).
		Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (tr *ThreadRepository) InsertInto(t *models.Thread) error {
	if err := tr.db.QueryRow("INSERT INTO threads "+
		"(slug, author, title, message, forum, created) "+
		"VALUES (NULLIF ($1, ''), $2, $3, $4, $5, $6) RETURNING id",
		t.Slug, t.AuthorID, t.Title, t.About, t.ForumID, t.CreationDate).
		Scan(&t.ID); err != nil {
		return err
	}

	return nil
}

func (tr *ThreadRepository) GetBySlug(slug string) (*models.Thread, error) {
	t := &models.Thread{}
	if err := tr.db.QueryRow("SELECT t.id, u.nickname, t.created, t.forum, f.slug, t.message, "+
		"coalesce (t.slug, ''), t.title, t.votes FROM threads AS t "+
		"JOIN users AS u ON (t.author = u.id) "+
		"JOIN forums AS f ON (f.id = t.forum) WHERE lower(t.slug) = lower($1)", slug).
		Scan(&t.ID, &t.Author, &t.CreationDate, &t.ForumID, &t.Forum, &t.About, &t.Slug,
			&t.Title, &t.Votes); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}

		return nil, err
	}

	return t, nil
}

func (tr *ThreadRepository) GetByID(id uint64) (*models.Thread, error) {
	t := &models.Thread{}
	if err := tr.db.QueryRow("SELECT t.id, u.nickname, t.created, t.forum, f.slug, t.message, "+
		"coalesce (t.slug, ''), t.title, t.votes FROM threads AS t "+
		"JOIN users AS u ON (t.author = u.id) "+
		"JOIN forums AS f ON (f.id = t.forum) WHERE t.id = $1", id).
		Scan(&t.ID, &t.Author, &t.CreationDate, &t.ForumID, &t.Forum, &t.About, &t.Slug,
			&t.Title, &t.Votes); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}

		return nil, err
	}

	return t, nil
}

func (tr *ThreadRepository) GetByForumSlug(
	slug string, limit uint64, since string, desc bool) ([]*models.Thread, error) {
	returnThreads := []*models.Thread{}

	queryString := "SELECT t.id, u.nickname, t.created, f.slug, t.message, " +
		"coalesce (t.slug, ''), t.title, t.votes FROM threads AS t " +
		"JOIN users AS u ON (t.author = u.id) " +
		"JOIN forums AS f ON (f.id = t.forum) WHERE lower(f.slug) = lower($1)"

	orderString := " ORDER BY t.created"

	if desc {
		orderString += " DESC"
	}

	if limit != 0 {
		orderString += fmt.Sprintf(" LIMIT %d", limit)
	}

	var rows *pgx.Rows
	var err error

	if since != "" {
		if desc {
			queryString += " AND t.created <= $2"
		} else {
			queryString += " AND t.created >= $2"
		}
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
		err = rows.Scan(&t.ID, &t.Author, &t.CreationDate, &t.Forum, &t.About, &t.Slug, &t.Title, &t.Votes)
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

package repository

import (
	"database/sql"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/post"
	"github.com/poofik33/db-technopark/internal/tools"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) post.Repository {
	return &PostRepository{
		db: db,
	}
}

func (pr *PostRepository) InsertInto(posts []*models.Post) error {
	tx, err := pr.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	for _, p := range posts {
		var forumID, userID uint64
		if err := pr.db.QueryRow("SELECT id FROM forums "+
			"WHERE lower(slug)=lower($1)", p.Forum).Scan(&forumID); err != nil {
			return err
		}

		if err := pr.db.QueryRow("SELECT id FROM users "+
			"WHERE lower(nickname)=lower($1)", p.Author).Scan(&userID); err != nil {
			return err
		}

		if err := tx.QueryRow(
			"INSERT INTO posts (id, author, forum, created, message, path, thread) "+
				"VALUES ((select nextval('posts_id_seq')::integer), $1, $2, $3, $4, "+
				"(SELECT path FROM posts WHERE id = $5) || (select currval('posts_id_seq')::integer), $6) "+
				"RETURNING id", userID, forumID, p.CreationDate,
			p.Message, p.ParentID, p.ThreadID).Scan(&p.ID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pr *PostRepository) GetCountByForumSlug(slug string) (uint64, error) {
	var count uint64
	if err := pr.db.QueryRow("SELECT count(*) from posts "+
		"JOIN forum as f (posts.forum = f.id) WHERE f.slug = $1", slug).
		Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (pr *PostRepository) GetByID(id uint64) (*models.Post, error) {
	p := &models.Post{}
	if err := pr.db.QueryRow(
		"SELECT id, author, forum, thread, message, created, isEdited, "+
			"coalesce(path[array_length(path, 1) - 1], 0) WHERE id = $1", id).
		Scan(&p.ID, &p.Author, &p.Forum, &p.ThreadID, &p.Message,
			&p.CreationDate, &p.IsEdited, &p.ParentID); err != nil {
		if err == sql.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}

	return p, nil
}

func (pr *PostRepository) Update(post *models.Post) error {
	if _, err := pr.db.Exec("UPDATE posts SET message = $2, isEdited = TRUE "+
		"WHERE id = $1", post.ID, post.Message); err != nil {
		return err
	}

	return nil
}

func (pr *PostRepository) GetByThread(id uint64, limit uint64,
	since uint64, sort string, desc bool) ([]*models.Post, error) {
	queryString := "SELECT id, author, forum, thread, created, message, isEdited, thread, " +
		"coalesce(path[array_length(path, 1) - 1], 0) " +
		"FROM posts WHERE "

	if sort == "parent_tree" {
		queryString += "path[1] IN (SELECT id FROM posts WHERE thread = $1 AND " +
			"array_length(path, 1) = 1"
		if since != 0 {
			queryString += " AND id > $2"
			if limit != 0 {
				queryString += " LIMIT $3"
			}
		} else {
			if limit != 0 {
				queryString += " LIMIT $2"
			}
		}

	} else {
		queryString += "thread = $1"
		if since != 0 {
			queryString += " AND id > $2"
		}
	}

	orderString := " ORDER BY "
	if sort == "flat" {
		orderString += "created"
	} else {
		orderString += "path[1]"
	}
	if desc {
		orderString += " DESC"
	}

	if sort != "parent_tree" && limit != 0 {
		if since != 0 {
			queryString += " LIMIT $3"
		} else {
			queryString += " LIMIT $2"
		}
	}

	rows := &sql.Rows{}
	var err error
	if since != 0 {
		if limit != 0 {
			rows, err = pr.db.Query(queryString+orderString, id, since, limit)
		} else {
			rows, err = pr.db.Query(queryString+orderString, id, since)
		}
	} else {
		if limit != 0 {
			rows, err = pr.db.Query(queryString+orderString, id, limit)
		} else {
			rows, err = pr.db.Query(queryString+orderString, id)
		}
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	returnPosts := []*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err = rows.Scan(&p.ID, &p.Author, &p.Forum, &p.ThreadID,
			&p.CreationDate, &p.Message, &p.IsEdited, &p.ParentID)
		if err != nil {
			return nil, err
		}

		returnPosts = append(returnPosts, p)
	}

	return returnPosts, nil
}

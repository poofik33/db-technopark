package repository

import (
	"database/sql"
	"fmt"
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
	if err := pr.db.QueryRow("SELECT count(*) from posts AS p "+
		"JOIN forums as f ON (p.forum = f.id) WHERE f.slug = $1", slug).
		Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (pr *PostRepository) GetByID(id uint64) (*models.Post, error) {
	p := &models.Post{}
	if err := pr.db.QueryRow(
		"SELECT p.id, u.nickname, f.slug, p.thread, p.message, p.created, p.isEdited, "+
			"coalesce(path[array_length(path, 1) - 1], 0) FROM posts AS p "+
			"JOIN users AS u ON (u.id = p.author) "+
			"JOIN forums AS f ON (f.id = p.forum) WHERE p.id = $1", id).
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
	queryStringFmt := "SELECT p.id, u.nickname, f.slug, p.thread, p.created, p.message, p.isEdited, " +
		"coalesce(p.path[array_length(p.path, 1) - 1], 0) " +
		"FROM posts AS p " +
		"JOIN users AS u ON (u.id = p.author) " +
		"JOIN forums AS f ON (f.id = p.forum) " +
		"WHERE %s %s"

	var whereString string
	var orderString string

	switch sort {
	case "flat", "":
		whereString = "p.thread = $1"
		if since != 0 {
			if desc {
				whereString += " AND p.id < $2"
			} else {
				whereString += " AND p.id > $2"
			}
		}
		orderString = "ORDER BY "
		if sort == "flat" {
			orderString += "p.created"
			if desc {
				orderString += " DESC"
			}
			orderString += ", p.id"
			if desc {
				orderString += " DESC"
			}
		} else {
			orderString += "p.id"
			if desc {
				orderString += " DESC"
			}
		}
		if limit != 0 {
			if since != 0 {
				orderString += " LIMIT $3"
			} else {
				orderString += " LIMIT $2"
			}
		}
	case "tree":
		whereString = "p.thread = $1"
		if since != 0 {
			if desc {
				whereString += " AND coalesce(path < (select path FROM posts where id = $2), true)"
			} else {
				whereString += " AND coalesce(path > (select path FROM posts where id = $2), true)"
			}
		}
		orderString = "ORDER BY p.path[1]"
		if desc {
			orderString += " DESC"
		}
		orderString += ", p.path[2:]"
		if desc {
			orderString += " DESC"
		}
		orderString += " NULLS FIRST"
		if limit != 0 {
			if since != 0 {
				orderString += " LIMIT $3"
			} else {
				orderString += " LIMIT $2"
			}
		}
	case "parent_tree":
		whereString = "p.path[1] IN (SELECT path[1] FROM posts WHERE thread = $1 AND " +
			"array_length(path, 1) = 1"
		if since != 0 {
			if desc {
				whereString += " AND id < (SELECT path[1] FROM posts WHERE id = $2)"
			} else {
				whereString += " AND id > (SELECT path[1] FROM posts WHERE id = $2)"
			}
		}

		whereString += " ORDER BY id"
		if desc {
			whereString += " DESC"
		}

		if limit != 0 {
			if since != 0 {
				whereString += " LIMIT $3"
			} else {
				whereString += " LIMIT $2"
			}
		}
		whereString += ")"
		orderString = "ORDER BY p.path[1]"
		if desc {
			orderString += " DESC"
		}
		orderString += ", p.path[2:] NULLS FIRST"
	}

	rows := &sql.Rows{}
	var err error
	if since != 0 {
		if limit != 0 {
			rows, err = pr.db.Query(fmt.Sprintf(queryStringFmt, whereString, orderString), id, since, limit)
		} else {
			rows, err = pr.db.Query(fmt.Sprintf(queryStringFmt, whereString, orderString), id, since)
		}
	} else {
		if limit != 0 {
			rows, err = pr.db.Query(fmt.Sprintf(queryStringFmt, whereString, orderString), id, limit)
		} else {
			rows, err = pr.db.Query(fmt.Sprintf(queryStringFmt, whereString, orderString), id)
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

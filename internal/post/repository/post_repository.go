package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/post"
	"github.com/poofik33/db-technopark/internal/tools"
	"strconv"
)

type PostRepository struct {
	db *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) post.Repository {
	return &PostRepository{
		db: db,
	}
}

func (pr *PostRepository) InsertInto(posts []*models.Post) error {
	sqlRow := "INSERT INTO posts (author, forum, created, message, parent, thread) VALUES "

	var val []interface{}
	id := uint64(1)
	for _, p := range posts {
		sqlRow += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", id, id+1, id+2, id+3, id+4, id+5)
		val = append(val, p.AuthorID, p.ForumID, p.CreationDate,
			p.Message, p.ParentID, p.ThreadID)
		id += 6
	}
	sqlRow = sqlRow[0 : len(sqlRow)-1]
	sqlRow += " RETURNING id"
	rows, err := pr.db.Query(sqlRow, val...)
	if err != nil {
		return err
	}

	defer rows.Close()

	postIndex := 0
	for rows.Next() {
		if err := rows.Scan(&posts[postIndex].ID); err != nil {
			return err
		}

		postIndex++
	}

	return nil
}

func (pr *PostRepository) GetCountByForumID(id uint64) (uint64, error) {
	var count uint64
	if err := pr.db.QueryRow("SELECT count(*) from posts WHERE forum = $1", id).
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
		if err == pgx.ErrNoRows {
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

	rows := &pgx.Rows{}
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

func (pr *PostRepository) CheckParentPosts(posts []*models.Post, threadID uint64) (bool, error) {
	parents := map[uint64]uint64{}
	vals := []interface{}{threadID}

	sqlRow := "SELECT count(*) FROM posts WHERE thread = $1 AND id in ("

	i := 2
	for _, p := range posts {
		if p.ParentID > 0 {
			sqlRow += "$" + strconv.Itoa(i) + ","
			parents[p.ParentID] += 1
			vals = append(vals, p.ParentID)
			i++
		}
	}

	if len(parents) == 0 {
		return true, nil
	}

	sqlRow = sqlRow[0:len(sqlRow)-1] + ")"
	var count int
	if err := pr.db.QueryRow(sqlRow, vals...).Scan(&count); err != nil {
		return false, err
	}

	if count != len(parents) {
		return false, tools.ErrParentPostDoesntExists
	}

	return true, nil
}

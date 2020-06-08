package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/tools"
	"github.com/poofik33/db-technopark/internal/user"
	"strings"
)

type UserRepository struct {
	db *pgx.ConnPool
}

func NewUserRepository(db *pgx.ConnPool) user.Repository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) InsertInto(user *models.User) error {
	if _, err := ur.db.Exec("INSERT INTO users (nickname, email, fullname, about) "+
		"VALUES ($1, $2, $3, $4)", user.Nickname, user.Email, user.Fullname, user.About); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetByNickname(nickname string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, fullname, about FROM users "+
		"WHERE lower(nickname) = lower($1)", nickname).Scan(&u.ID, &u.Nickname, &u.Email, &u.Fullname,
		&u.About); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT id, nickname, email, fullname, about FROM users "+
		"WHERE lower(email) = lower($1)", email).Scan(&u.ID, &u.Nickname, &u.Email, &u.Fullname,
		&u.About); err != nil {
		if err == pgx.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) Update(user *models.User) error {
	if _, err := ur.db.Exec("UPDATE users SET email = $2, fullname = $3, about = $4 "+
		"WHERE lower(nickname) = lower($1)", user.Nickname, user.Email, user.Fullname, user.About); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetUsersByForum(
	id uint64, limit uint64, since string, desc bool) ([]*models.User, error) {
	returnUsers := []*models.User{}

	queryString := "SELECT u.nickname, u.email, u.fullname, u.about FROM forums_users fu " +
		"JOIN users u ON (fu.user_id = u.id) " +
		"WHERE fu.forum_id = $1"
	groupbyString := " ORDER BY lower(u.nickname)"
	if desc {
		groupbyString += " DESC"
	}

	if limit != 0 {
		groupbyString += fmt.Sprintf(" LIMIT %d", limit)
	}

	var rows *pgx.Rows
	var err error
	if since != "" {
		if desc {
			queryString += " AND lower(u.nickname) < lower($2)"
		} else {
			queryString += " AND lower(u.nickname) > lower($2)"
		}
		rows, err = ur.db.Query(queryString+groupbyString, id, since)
	} else {
		rows, err = ur.db.Query(queryString+groupbyString, id)
	}
	if err != nil {
		//if err == sql.ErrNoRows {
		//	return nil, nil
		//}

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := &models.User{}
		err = rows.Scan(&u.Nickname, &u.Email, &u.Fullname, &u.About)
		if err != nil {
			return nil, err
		}

		returnUsers = append(returnUsers, u)
	}

	return returnUsers, nil
}

func (ur *UserRepository) CheckNicknames(posts []*models.Post) (bool, error) {
	rows, err := ur.db.Query("SELECT id, lower(nickname) FROM users")
	if err != nil {
		return false, err
	}

	defer rows.Close()

	nicknames := make(map[string]uint64)
	for rows.Next() {
		n := ""
		var id uint64
		if err := rows.Scan(&id, &n); err != nil {
			return false, err
		}

		nicknames[n] = id
	}

	for _, p := range posts {
		id := nicknames[strings.ToLower(p.Author)]
		if id == 0 {
			return false, tools.ErrUserDoesntExists
		}
		p.AuthorID = id
	}

	return true, nil
}

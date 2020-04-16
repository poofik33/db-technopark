package repository

import (
	"database/sql"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/tools"
	"github.com/poofik33/db-technopark/internal/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &UserRepository{
		db:db,
	}
}

func (ur *UserRepository) InsertInto(user *models.User) error {
	if _, err := ur.db.Exec("INSERT INTO users (nickname, email, fullname, about) " +
		"VALUES ($1, $2, $3, $4)", user.Nickname, user.Email, user.Fullname, user.About); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetByNickname(nickname string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT nickname, email, fullname, about FROM users " +
		"WHERE nickname = lower($1)", nickname).Scan(&u.Nickname, &u.Email, &u.Fullname, &u.About); err != nil {
		if err == sql.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	u := &models.User{}

	if err := ur.db.QueryRow("SELECT nickname, email, fullname, about FROM users " +
		"WHERE email = lower($1)", email).Scan(&u.Nickname, &u.Email, &u.Fullname, &u.About); err != nil {
		if err == sql.ErrNoRows {
			return nil, tools.ErrDoesntExists
		}
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) Update(user *models.User) error {
	if _, err := ur.db.Exec("UPDATE users SET email = $2, fullname = $3, about = $4 " +
		"WHERE nickname = lower($1)", user.Nickname, user.Email, user.Fullname, user.About); err != nil {
		return err
	}

	return nil
}

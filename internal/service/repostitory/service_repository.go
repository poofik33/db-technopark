package repostitory

import (
	"database/sql"
	"github.com/poofik33/db-technopark/internal/service"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) service.Repository {
	return &ServiceRepository{
		db: db,
	}
}

func (sr *ServiceRepository) GetCountForum() (uint64, error) {
	var count uint64
	if err := sr.db.QueryRow("SELECT count(*) from forums").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (sr *ServiceRepository) GetCountPost() (uint64, error) {
	var count uint64
	if err := sr.db.QueryRow("SELECT count(*) from posts").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (sr *ServiceRepository) GetCountThread() (uint64, error) {
	var count uint64
	if err := sr.db.QueryRow("SELECT count(*) from threads").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (sr *ServiceRepository) GetCountUser() (uint64, error) {
	var count uint64
	if err := sr.db.QueryRow("SELECT count(*) from users").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (sr *ServiceRepository) DeleteAllForum() error {
	if _, err := sr.db.Exec("DELETE FROM forums"); err != nil {
		return err
	}
	return nil
}

func (sr *ServiceRepository) DeleteAllPost() error {
	if _, err := sr.db.Exec("DELETE FROM posts"); err != nil {
		return err
	}
	return nil
}

func (sr *ServiceRepository) DeleteAllThread() error {
	if _, err := sr.db.Exec("DELETE FROM threads"); err != nil {
		return err
	}
	return nil
}

func (sr *ServiceRepository) DeleteAllUser() error {
	if _, err := sr.db.Exec("DELETE FROM users"); err != nil {
		return err
	}
	return nil
}

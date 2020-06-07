package repostitory

import (
	"github.com/jackc/pgx"
	"github.com/poofik33/db-technopark/internal/service"
)

type ServiceRepository struct {
	db *pgx.ConnPool
}

func NewServiceRepository(db *pgx.ConnPool) service.Repository {
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
	if _, err := sr.db.Exec("TRUNCATE TABLE forums CASCADE"); err != nil {
		return err
	}
	return nil
}

func (sr *ServiceRepository) DeleteAllPost() error {
	if _, err := sr.db.Exec("TRUNCATE TABLE posts CASCADE"); err != nil {
		return err
	}
	return nil
}

func (sr *ServiceRepository) DeleteAllThread() error {
	if _, err := sr.db.Exec("TRUNCATE TABLE threads CASCADE"); err != nil {
		return err
	}
	return nil
}

func (sr *ServiceRepository) DeleteAllUser() error {
	if _, err := sr.db.Exec("TRUNCATE TABLE users CASCADE"); err != nil {
		return err
	}
	return nil
}

func (sr *ServiceRepository) DeleteAllVotes() error {
	if _, err := sr.db.Exec("TRUNCATE TABLE votes CASCADE"); err != nil {
		return err
	}
	return nil
}

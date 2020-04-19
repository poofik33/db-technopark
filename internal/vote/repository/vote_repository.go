package repository

import (
	"database/sql"
	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/vote"
)

type VoteRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) vote.Repository {
	return &VoteRepository{
		db: db,
	}
}

func (vr *VoteRepository) GetThreadVotes(id uint64) (int64, error) {
	var votes int64
	if err := vr.db.QueryRow("SELECT coalesce (sum(CASE WHEN voice THEN 1 WHEN NOT voice THEN -1 END), 0) "+
		"WHERE thread = $1", id).Scan(&votes); err != nil {
		return 0, err
	}

	return votes, nil
}

func (vr *VoteRepository) InsertInto(v *models.Vote) error {
	if _, err := vr.db.Exec("INSERT INTO votes (nickname, thread, vote) "+
		"VALUES ($1, $2, CASE WHEN $3 = 1 TRUE WHEN $3 = -1 FALSE END)",
		v.Nickname, v.ThreadID, v.Voice); err != nil {
		return err
	}

	return nil
}

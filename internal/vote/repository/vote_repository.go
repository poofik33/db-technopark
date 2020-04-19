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
	if err := vr.db.QueryRow("SELECT coalesce (sum(CASE WHEN vote THEN 1 WHEN NOT vote THEN -1 END), 0) "+
		"FROM votes WHERE thread = $1", id).Scan(&votes); err != nil {
		return 0, err
	}

	return votes, nil
}

func (vr *VoteRepository) InsertInto(v *models.Vote) error {
	userID := uint64(0)
	if err := vr.db.QueryRow("SELECT id FROM users "+
		"WHERE lower(nickname) = lower($1)", v.Nickname).Scan(&userID); err != nil {
		return err
	}

	curVoteID := uint64(0)
	err := vr.db.QueryRow("SELECT id FROM votes "+
		"WHERE thread = $1 and author = $2", v.ThreadID, userID).Scan(&curVoteID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		if _, err := vr.db.Exec("INSERT INTO votes (author, thread, vote) "+
			"VALUES ($1, $2, CASE WHEN $3 = 1 THEN TRUE WHEN $3 = -1 THEN FALSE END)",
			userID, v.ThreadID, v.Voice); err != nil {
			return err
		}

		return nil
	}

	if _, err = vr.db.Exec("UPDATE votes SET vote = (CASE WHEN $2 = 1 THEN TRUE WHEN $2 = -1 THEN FALSE END) "+
		"WHERE id = $1", curVoteID, v.Voice); err != nil {
		return err
	}

	return nil

}

package repository

import (
	"github.com/jackc/pgx"

	"github.com/poofik33/db-technopark/internal/models"
	"github.com/poofik33/db-technopark/internal/vote"
)

type VoteRepository struct {
	db *pgx.ConnPool
}

func NewVoteRepository(db *pgx.ConnPool) vote.Repository {
	return &VoteRepository{
		db: db,
	}
}

func (vr *VoteRepository) GetThreadVotes(id uint64) (int64, error) {
	var votes int64
	if err := vr.db.QueryRow("SELECT votes "+
		"FROM threads WHERE id = $1", id).Scan(&votes); err != nil {
		return 0, err
	}

	return votes, nil
}

func (vr *VoteRepository) InsertInto(v *models.Vote) error {
	curVoteID := uint64(0)
	err := vr.db.QueryRow("SELECT id FROM votes "+
		"WHERE thread = $1 and author = $2", v.ThreadID, v.UserID).Scan(&curVoteID)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	if err == pgx.ErrNoRows {
		if _, err := vr.db.Exec("INSERT INTO votes (author, thread, vote) "+
			"VALUES ($1, $2, $3)",
			v.UserID, v.ThreadID, v.Voice); err != nil {
			return err
		}

		return nil
	}

	if _, err = vr.db.Exec("UPDATE votes SET vote = $2 "+
		"WHERE id = $1", curVoteID, v.Voice); err != nil {
		return err
	}

	return nil

}

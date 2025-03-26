package postgres

import (
	"database/sql"
	"github.com/Arh0rn/absoluteCinema/pkg/models"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

func (r *TokenRepo) CreateToken(rt *models.RefreshSession) error {
	_, err := r.db.Exec(
		"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		rt.UserID, rt.Token, rt.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *TokenRepo) PopToken(oldRt string) (*models.RefreshSession, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var newRt models.RefreshSession
	err = tx.QueryRow(
		"SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1", oldRt).
		Scan(&newRt.ID, &newRt.UserID, &newRt.Token, &newRt.ExpiresAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("DELETE FROM refresh_tokens WHERE token = $1", oldRt)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &newRt, nil
}

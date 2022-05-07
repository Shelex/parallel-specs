package postgres

import (
	"fmt"

	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/errors"
)

func (pg *Postgres) AddApiKey(userID string, key entities.ApiKey) error {
	query := `INSERT INTO apiKey (id, userId, name, expireAt) VALUES ($1, $2, $3, $4)`

	if _, err := pg.db.Exec(pg.ctx, query, key.ID, key.UserID, key.Name, key.ExpireAt); err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) DeleteApiKey(userID string, keyID string) error {
	query := `DELETE FROM apiKey WHERE apiKey.id = $1 AND apiKey.userID = $2`

	command, err := pg.db.Exec(pg.ctx, query, keyID, userID)
	if err != nil {
		return err
	}

	if command.RowsAffected() == 0 {
		return errors.ApiKeyNotFound
	}

	return nil
}

func (pg *Postgres) GetApiKeys(userID string) ([]entities.ApiKey, error) {
	var keys []entities.ApiKey

	query := `SELECT * FROM apiKey WHERE userid = $1`

	rows, err := pg.db.Query(pg.ctx, query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var key entities.ApiKey
		if err := rows.Scan(&key.ID, &key.UserID, &key.Name, &key.ExpireAt); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func (pg *Postgres) IsApiKeyValid(userID string, keyID string) (bool, error) {
	apiKeyQuery := `SELECT 1 FROM apiKey WHERE id = $1 AND userid = $2`
	existsQuery := fmt.Sprintf("SELECT EXISTS (%s)", apiKeyQuery)

	var exists bool

	if err := pg.db.QueryRow(pg.ctx, existsQuery, keyID, userID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

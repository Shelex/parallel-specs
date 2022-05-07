package postgres

import (
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func (pg *Postgres) DeleteSpecs(projectID string) error {
	deleteSpecsQuery := `DELETE FROM spec WHERE spec.projectId = $1`

	if _, err := pg.db.Exec(pg.ctx, deleteSpecsQuery, projectID); err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) AddSpecsMaybe(projectID string, names []string) ([]entities.Spec, error) {
	specs := make([]entities.Spec, len(names))
	for index, name := range names {
		spec, err := pg.IsSpecAvailable(projectID, name)
		if err != nil {
			if err == errors.SpecNotFound {
				newSpec, err := pg.AddSpec(projectID, name)
				if err != nil {
					return nil, err
				}
				spec = newSpec
			} else {
				return nil, err
			}
		}
		specs[index] = spec
	}
	return specs, nil
}

func (pg *Postgres) IsSpecAvailable(projectID string, name string) (entities.Spec, error) {
	specQuery := `SELECT * FROM spec WHERE projectId = $1 AND filePath = $2`

	var spec entities.Spec

	err := pg.db.QueryRow(pg.ctx, specQuery, projectID, name).Scan(&spec.ID, &spec.FilePath, &spec.ProjectID)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.SpecNotFound
		}
		return spec, err
	}

	return spec, nil
}

func (pg *Postgres) AddSpec(projectID string, name string) (entities.Spec, error) {
	query := `INSERT INTO spec (id, filePath, projectId) VALUES ($1, $2, $3)`
	id := uuid.NewString()
	spec := entities.Spec{
		ID:        id,
		FilePath:  name,
		ProjectID: projectID,
	}
	if _, err := pg.db.Exec(pg.ctx, query, spec.ID, spec.FilePath, spec.ProjectID); err != nil {
		return spec, err
	}

	return spec, nil
}

func (pg *Postgres) GetSpec(specID string) (entities.Spec, error) {
	query := `SELECT * FROM spec WHERE id = $1`

	var spec entities.Spec

	if err := pg.db.QueryRow(pg.ctx, query, specID).Scan(&spec.ID, &spec.FilePath, &spec.ProjectID); err != nil {
		return spec, err
	}

	return spec, nil
}

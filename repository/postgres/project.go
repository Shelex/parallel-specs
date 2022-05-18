package postgres

import (
	"fmt"

	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/errors"
	"github.com/Shelex/split-specs-v2/repository"
	"github.com/jackc/pgx/v4"
)

func (pg *Postgres) GetProjectByID(ID string) (*entities.Project, error) {
	query := `SELECT * FROM project WHERE id = $1`

	var project entities.Project

	if err := pg.db.QueryRow(pg.ctx, query, ID).Scan(&project.ID, &project.Name); err != nil {
		return nil, err
	}

	return &project, nil
}

func (pg *Postgres) GetProjectUsers(projectID string) ([]string, error) {
	var users []string
	query := `SELECT userproject.userId FROM userproject WHERE projectId = $1`

	rows, err := pg.db.Query(pg.ctx, query, projectID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (pg *Postgres) AddProject(project entities.Project) error {
	query := `INSERT INTO project (id, name) VALUES ($1, $2)`

	if _, err := pg.db.Exec(pg.ctx, query, project.ID, project.Name); err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) DeleteProject(userID string, projectID string) error {
	deleteUserProjectQuery := `DELETE FROM userproject WHERE userproject.userId = $1 AND userproject.projectId = $2`

	operation, err := pg.db.Exec(pg.ctx, deleteUserProjectQuery, userID, projectID)
	if err != nil {
		return err
	}

	if operation.RowsAffected() == 0 {
		return errors.ProjectNotFound
	}

	deleteProjectQuery := `DELETE FROM project WHERE project.id = $1`
	deleteSpecsQuery := `DELETE FROM spec WHERE spec.projectId = $1`
	deleteSessionsQuery := `DELETE FROM session_execution WHERE session_execution.projectId = $1`
	deleteSpecExecutionsQuery := `
	DELETE  FROM spec_execution
	USING session_execution
	WHERE session_execution.projectId = $1
	AND sessionId = session_execution.id`

	return pg.db.BeginTxFunc(pg.ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
	}, func(tx pgx.Tx) error {
		if _, err := tx.Exec(pg.ctx, deleteSpecExecutionsQuery, projectID); err != nil {
			return err
		}

		if _, err := tx.Exec(pg.ctx, deleteSessionsQuery, projectID); err != nil {
			return err
		}

		if _, err := tx.Exec(pg.ctx, deleteSpecsQuery, projectID); err != nil {
			return err
		}

		deleteCmd, err := tx.Exec(pg.ctx, deleteProjectQuery, projectID)
		if err != nil {
			return err
		}

		if deleteCmd.RowsAffected() == 0 {
			return errors.ProjectNotFound
		}

		return nil
	})
}

func (pg *Postgres) IsProjectAccessible(userID string, projectID string) (bool, error) {
	users, err := pg.GetProjectUsers(projectID)
	if err != nil {
		return false, err
	}
	hasAccess, _ := repository.Contains(users, userID)
	return hasAccess, nil
}

func (pg *Postgres) GetProjectSessions(projectID string, pagination *entities.Pagination) ([]entities.Session, int, error) {
	var sessions []entities.Session

	limit := 20
	offset := 0
	if pagination.Limit != 0 {
		limit = pagination.Limit
	}
	if pagination.Offset != 0 {
		offset = pagination.Offset
	}

	fromAndWhere := `FROM session_execution WHERE projectId = $1`
	query := fmt.Sprintf("SELECT * %s ORDER BY createdAt DESC OFFSET $2 LIMIT $3", fromAndWhere)
	countQuery := fmt.Sprintf("SELECT COUNT(id) %s", fromAndWhere)

	var total int
	if err := pg.db.QueryRow(pg.ctx, countQuery, projectID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := pg.db.Query(pg.ctx, query, projectID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var session entities.Session
		if err := rows.Scan(&session.ID, &session.ProjectID, &session.StartedAt, &session.FinishedAt, &session.CreatedAt); err != nil {
			return nil, 0, err
		}
		executions, err := pg.GetExecutions(session.ID)
		if err != nil {
			return nil, 0, err
		}
		session.Executions = executions
		sessions = append(sessions, session)
	}

	return sessions, total, nil
}

package postgres

import (
	"github.com/Shelex/parallel-specs/internal/entities"
	"github.com/Shelex/parallel-specs/internal/errors"
	"github.com/Shelex/parallel-specs/internal/events"
	"github.com/Shelex/parallel-specs/repository"
	"github.com/jackc/pgx/v4"
)

func (pg *Postgres) AddSession(sessionExecution entities.Session) error {
	query := `INSERT INTO session_execution (id, projectId, createdAt) VALUES ($1, $2, $3)`

	if _, err := pg.db.Exec(pg.ctx, query, sessionExecution.ID, sessionExecution.ProjectID, repository.GetTimestamp()); err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) GetSession(ID string) (*entities.Session, error) {
	query := `SELECT * FROM session_execution WHERE id = $1`

	var session entities.Session

	if err := pg.db.QueryRow(pg.ctx, query, ID).Scan(&session.ID, &session.ProjectID, &session.StartedAt, &session.FinishedAt, &session.CreatedAt); err != nil {
		return nil, err
	}

	return &session, nil
}

func (pg *Postgres) GetSessionWithExecution(sessionID string) (*entities.Session, error) {
	session, err := pg.GetSession(sessionID)
	if err != nil {
		return nil, errors.SessionNotFound
	}

	executions, err := pg.GetExecutions(sessionID)
	if err != nil {
		return nil, err
	}

	session.Executions = executions

	return session, nil
}

func (pg *Postgres) StartSession(ID string) error {
	session, err := pg.GetSession(ID)
	if err != nil {
		return err
	}
	query := `UPDATE session_execution SET startedAt = $1 WHERE id = $2`

	start := repository.GetTimestamp()

	cmd, err := pg.db.Exec(pg.ctx, query, start, ID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() > 0 {
		events.Handler.Publish(events.Session, events.SessionEvent{
			Event: events.BasicEvent{
				Topic: events.Session,
				Kind:  events.Started,
				ID:    ID,
			},
			ProjectID: session.ProjectID,
			Time:      start,
		})
	}

	return nil
}
func (pg *Postgres) EndSession(ID string) error {
	session, err := pg.GetSession(ID)
	if err != nil {
		return err
	}
	query := `UPDATE session_execution SET finishedAt = $1 WHERE id = $2 AND finishedAt = 0`

	end := repository.GetTimestamp()

	cmd, err := pg.db.Exec(pg.ctx, query, end, ID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() > 0 {
		events.Handler.Publish(events.Session, events.SessionEvent{
			Event: events.BasicEvent{
				Topic: events.Session,
				Kind:  events.Finished,
				ID:    ID,
			},
			ProjectID: session.ProjectID,
			Time:      end,
		})
	}

	return nil
}

func (pg *Postgres) DeleteSession(sessionID string) error {
	deleteSessionsQuery := `DELETE FROM session_execution WHERE session_execution.id = $1`
	deleteSpecExecutionsQuery := `DELETE FROM spec_execution WHERE spec_execution.sessionId = $1`

	return pg.db.BeginTxFunc(pg.ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.Deferrable,
	}, func(tx pgx.Tx) error {
		if _, err := tx.Exec(pg.ctx, deleteSpecExecutionsQuery, sessionID); err != nil {
			return err
		}

		deleteCmd, err := tx.Exec(pg.ctx, deleteSessionsQuery, sessionID)
		if err != nil {
			return err
		}

		if deleteCmd.RowsAffected() == 0 {
			return errors.SessionNotFound
		}

		events.Handler.Publish(events.Session, events.SessionEvent{
			Event: events.BasicEvent{
				Topic: events.Session,
				Kind:  events.Deleted,
				ID:    sessionID,
			},
		})

		return nil
	})
}

package postgres

import (
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/events"
	"github.com/Shelex/split-specs-v2/repository"
)

func (pg *Postgres) AddExecutions(sessionID string, executions []entities.Execution) error {
	for _, execution := range executions {
		if err := pg.AddExecution(sessionID, execution); err != nil {
			return err
		}
	}
	return nil
}

func (pg *Postgres) AddExecution(sessionID string, execution entities.Execution) error {
	query := `INSERT INTO spec_execution (id, specId, specName, sessionId, estimatedDuration) VALUES ($1, $2, $3, $4, $5)`

	if _, err := pg.db.Exec(pg.ctx, query, execution.ID, execution.SpecID, execution.SpecName, sessionID, execution.EstimatedDuration); err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) GetExecutions(sessionID string) ([]entities.Execution, error) {
	var executions []entities.Execution

	query := `SELECT * FROM spec_execution WHERE sessionId = $1`

	rows, err := pg.db.Query(pg.ctx, query, sessionID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var execution entities.Execution

		if err := rows.Scan(&execution.ID,
			&execution.SpecID,
			&execution.SpecName,
			&execution.SessionID,
			&execution.MachineID,
			&execution.StartedAt,
			&execution.FinishedAt,
			&execution.EstimatedDuration,
			&execution.Status); err != nil {
			return nil, err
		}
		execution.Duration = uint32(execution.FinishedAt - execution.StartedAt)
		executions = append(executions, execution)
	}

	return executions, nil
}

func (pg *Postgres) GetExecutionHistory(specID string, limit int) ([]entities.Execution, error) {
	var executions []entities.Execution

	query := `SELECT * FROM spec_execution WHERE specId = $1 AND finishedAt > 0 GROUP BY id ORDER BY finishedAt DESC LIMIT $2`

	rows, err := pg.db.Query(pg.ctx, query, specID, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var execution entities.Execution
		if err := rows.Scan(&execution.ID,
			&execution.SpecID,
			&execution.SpecName,
			&execution.SessionID,
			&execution.MachineID,
			&execution.StartedAt,
			&execution.FinishedAt,
			&execution.EstimatedDuration,
			&execution.Status); err != nil {
			return nil, err
		}
		execution.Duration = uint32(execution.FinishedAt - execution.StartedAt)
		executions = append(executions, execution)
	}

	return executions, nil
}

func (pg *Postgres) StartExecution(sessionID string, machineID string, specID string) error {
	query := `UPDATE spec_execution SET startedAt = $1, machineId = $2 WHERE sessionId = $3 AND specId = $4 AND startedAt = 0`

	start := repository.GetTimestamp()

	if _, err := pg.db.Exec(pg.ctx, query, start, machineID, sessionID, specID); err != nil {
		return err
	}

	events.Handler.Publish(events.Execution, events.ExecutionEvent{
		Event: events.BasicEvent{
			Topic: events.Execution,
			Kind:  events.Started,
			ID:    specID,
		},
		Time:      start,
		SessionID: sessionID,
	})

	return nil
}

func (pg *Postgres) EndExecution(sessionID string, machineID string, status string) error {
	query := `UPDATE spec_execution SET finishedAt = $1, status = $2 WHERE sessionId = $3 AND startedAt > 0 AND finishedAt = 0 AND machineID = $4`

	end := repository.GetTimestamp()

	command, err := pg.db.Exec(pg.ctx, query, end, status, sessionID, machineID)
	if err != nil {
		return err
	}

	if command.RowsAffected() != 0 {
		events.Handler.Publish(events.Execution, events.ExecutionEvent{
			Event: events.BasicEvent{
				Topic: events.Execution,
				Kind:  events.Finished,
				ID:    machineID,
			},
			Time:      end,
			SessionID: sessionID,
		})
	}

	return nil
}

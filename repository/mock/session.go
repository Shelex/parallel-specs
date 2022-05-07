package mock

import (
	"fmt"

	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/errors"
	"github.com/Shelex/split-specs-v2/repository"
)

func (i *MockStorage) AddSession(sessionExecution entities.Session) error {
	if _, ok := i.Sessions[sessionExecution.ID]; ok {
		return fmt.Errorf("[repository]: session id already in use for project %s", sessionExecution.ProjectID)
	}

	i.Sessions[sessionExecution.ID] = &sessionExecution

	return nil
}

func (i *MockStorage) GetSession(ID string) (*entities.Session, error) {
	session, ok := i.Sessions[ID]
	if !ok {
		return nil, errors.SessionNotFound
	}
	return session, nil
}

func (i *MockStorage) GetSessionWithExecution(sessionID string) (*entities.Session, error) {
	session, err := i.GetSession(sessionID)
	if err != nil {
		return nil, errors.SessionNotFound
	}

	executions, err := i.GetExecutions(sessionID)
	if err != nil {
		return nil, err
	}

	session.Executions = executions

	return session, nil
}

func (i *MockStorage) StartSession(ID string) error {
	session, err := i.GetSession(ID)
	if err != nil {
		return err
	}

	if session.StartedAt == 0 {
		session.StartedAt = repository.GetTimestamp()
		i.Sessions[ID] = session
	}
	return nil
}
func (i *MockStorage) EndSession(ID string) error {
	session, err := i.GetSession(ID)
	if err != nil {
		return err
	}
	session.FinishedAt = repository.GetTimestamp()
	i.Sessions[ID] = session
	return nil
}

func (i *MockStorage) DeleteSession(sessionID string) error {
	for _, execution := range i.Executions {
		if execution.SessionID == sessionID {
			delete(i.Executions, execution.ID)
		}
	}

	delete(i.Sessions, sessionID)
	return nil
}

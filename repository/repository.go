package repository

import (
	"context"
	"time"

	"github.com/Shelex/parallel-specs/internal/entities"
)

var DB Storage

type Storage interface {
	AddUser(userInput entities.User) error
	UpdatePassword(userID string, newPassword string) error
	GetUserByEmail(email string) (*entities.User, error)

	AddUserProject(userID string, projectID string) error
	GetUserProjectIDs(userID string) ([]string, error)
	GetUserProjectByName(userID string, projectName string) (*entities.Project, error)
	IsProjectAccessible(userID string, projectID string) (bool, error)
	GetProjectUsers(projectID string) ([]string, error)
	DeleteUserProject(userID string, projectID string) error

	AddProject(project entities.Project) error
	GetProjectByID(ID string) (*entities.Project, error)
	GetProjectSessions(projectID string, pagination *entities.Pagination) ([]entities.Session, int, error)
	DeleteProject(userID string, projectID string) error

	AddSession(sessionExecution entities.Session) error
	GetSession(ID string) (*entities.Session, error)
	GetSessionWithExecution(sessionID string) (*entities.Session, error)
	StartSession(ID string) error
	EndSession(ID string) error
	DeleteSession(sessionID string) error

	AddSpecsMaybe(projectID string, names []string) ([]entities.Spec, error)
	AddSpec(projectID string, name string) (entities.Spec, error)
	IsSpecAvailable(projectID string, name string) (entities.Spec, error)
	DeleteSpecs(projectID string) error
	GetSpec(specID string) (entities.Spec, error)

	AddExecutions(sessionID string, executions []entities.Execution) error
	GetExecutions(sessionID string) ([]entities.Execution, error)
	GetExecutionHistory(specID string, limit int) ([]entities.Execution, error)
	StartExecution(sessionID string, machineID string, specID string) error
	EndExecution(sessionID string, machineID string, status string) error

	AddApiKey(userID string, key entities.ApiKey) error
	GetApiKeys(userID string) ([]entities.ApiKey, error)
	IsApiKeyValid(userID string, keyID string) (bool, error)
	DeleteApiKey(userID string, keyID string) error

	// close db connection
	ShutDown(ctx context.Context)
}

func Contains(input []string, query string) (bool, int) {
	for index, item := range input {
		if item == query {
			return true, index
		}
	}
	return false, -1
}

func GetTimestamp() uint64 {
	return uint64(time.Now().UnixMilli())
}

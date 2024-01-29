package mock

import (
	"context"

	"github.com/Shelex/parallel-specs/internal/entities"
	"github.com/Shelex/parallel-specs/repository"
)

type MockStorage struct {
	Users        map[string]*entities.User
	ApiKeys      map[string]*entities.ApiKey
	UserProjects map[string]*entities.UserProject
	Projects     map[string]*entities.Project
	Sessions     map[string]*entities.Session
	Specs        map[string]*entities.Spec
	Executions   map[string]entities.Execution
}

var initial = MockStorage{
	Users:        map[string]*entities.User{},
	ApiKeys:      map[string]*entities.ApiKey{},
	UserProjects: map[string]*entities.UserProject{},
	Projects:     map[string]*entities.Project{},
	Sessions:     map[string]*entities.Session{},
	Specs:        map[string]*entities.Spec{},
	Executions:   map[string]entities.Execution{},
}

func NewMockStorage(ctx context.Context) (repository.Storage, error) {
	repository.DB = &initial
	return repository.DB, nil
}

func (i *MockStorage) ShutDown(ctx context.Context) {
	// close the db connections
}

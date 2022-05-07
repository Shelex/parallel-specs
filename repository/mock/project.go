package mock

import (
	"fmt"

	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/errors"
	"github.com/Shelex/split-specs-v2/repository"
)

func (i *MockStorage) GetProjectByID(ID string) (*entities.Project, error) {
	project, ok := i.Projects[ID]
	if !ok {
		return nil, errors.ProjectNotFound
	}
	return project, nil
}

func (i *MockStorage) GetProjectUsers(projectID string) ([]string, error) {
	var users []string
	for _, userProject := range i.UserProjects {
		if userProject.ProjectID == projectID {
			users = append(users, userProject.UserID)
		}
	}
	return users, nil
}

func (i *MockStorage) IsProjectAccessible(userID string, projectID string) (bool, error) {
	users, err := i.GetProjectUsers(projectID)
	if err != nil {
		return false, err
	}
	hasAccess, _ := repository.Contains(users, userID)
	return hasAccess, nil
}

func (i *MockStorage) AddProject(project entities.Project) error {
	i.Projects[project.ID] = &project
	return nil
}

func (i *MockStorage) DeleteProject(userID string, projectID string) error {
	users, err := i.GetProjectUsers(projectID)
	if err != nil {
		return err
	}

	hasAccess, _ := repository.Contains(users, userID)
	if !hasAccess {
		return errors.ProjectNotFound
	}

	projectSessions, _, err := i.GetProjectSessions(projectID, nil)
	if err != nil {
		return err
	}

	// clear out sessions
	for _, session := range projectSessions {
		err := i.DeleteSession(session.ID)
		if err != nil {
			return err
		}
	}
	delete(i.Projects, projectID)

	// clear out project access
	for _, userID := range users {
		if err := i.DeleteUserProject(userID, projectID); err != nil {
			return fmt.Errorf("could not unassign project from user")
		}
	}

	if err := i.DeleteSpecs(projectID); err != nil {
		return fmt.Errorf("could not remove project specs")
	}

	return nil
}

func (i *MockStorage) GetProjectSessions(projectID string, pagination *entities.Pagination) ([]entities.Session, int, error) {
	var sessions []entities.Session
	for _, session := range i.Sessions {
		if session.ProjectID == projectID {
			sessionWithSpecs, err := i.GetSessionWithExecution(session.ID)
			if err != nil {
				return sessions, 0, err
			}
			sessions = append(sessions, *sessionWithSpecs)
		}
	}

	return sessions, len(sessions), nil
}

package mock

import (
	"fmt"

	"github.com/Shelex/split-specs-v2/internal/appError"
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/google/uuid"
)

func (i *MockStorage) AddUser(userInput entities.User) error {
	i.Users[userInput.ID] = &userInput
	return nil
}

func (i *MockStorage) UpdatePassword(userID string, newPassword string) error {
	i.Users[userID].Password = newPassword
	return nil
}

func (i *MockStorage) GetUserByEmail(email string) (*entities.User, error) {
	for _, user := range i.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (i *MockStorage) AddUserProject(userID string, projectID string) error {
	id := uuid.NewString()
	i.UserProjects[id] = &entities.UserProject{
		ID:        id,
		UserID:    userID,
		ProjectID: projectID,
	}
	return nil
}

func (i *MockStorage) DeleteUserProject(userID string, projectID string) error {
	for _, userProject := range i.UserProjects {
		if userProject.UserID == userID && userProject.ProjectID == projectID {
			delete(i.UserProjects, userProject.ID)
		}
	}
	return nil
}

func (i *MockStorage) GetUserProjectIDs(userID string) ([]string, error) {
	var projectIds []string
	for _, userProject := range i.UserProjects {
		if userProject.UserID == userID {
			projectIds = append(projectIds, userProject.ProjectID)
		}
	}
	return projectIds, nil
}

func (i *MockStorage) GetUserProjectByName(userID string, projectName string) (*entities.Project, error) {
	var empty *entities.Project
	projectIds, err := i.GetUserProjectIDs(userID)
	if err != nil {
		return empty, err
	}

	for _, id := range projectIds {
		project, err := i.GetProjectByID(id)
		if err != nil {
			return empty, err
		}
		if project.Name == projectName {
			return project, nil
		}
	}

	return empty, appError.ProjectNotFound
}

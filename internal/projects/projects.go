package projects

import (
	"github.com/Shelex/split-specs-v2/internal/appError"
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func GetUserProjects(userID string) ([]*entities.Project, error) {
	projectIds, err := repository.DB.GetUserProjectIDs(userID)
	if err != nil {
		return nil, err
	}

	projects := make([]*entities.Project, len(projectIds))

	for index, id := range projectIds {
		project, err := repository.DB.GetProjectByID(id)
		if err != nil {
			return nil, err
		}
		projects[index] = project
	}
	return projects, nil
}

func GetByNameOrCreateNew(userID string, projectName string) (*entities.Project, bool, error) {
	var project *entities.Project

	isNew := false

	project, err := repository.DB.GetUserProjectByName(userID, projectName)
	if err != nil {
		if err.Error() == appError.ProjectNotFound.Error() || err.Error() == pgx.ErrNoRows.Error() {
			ID := uuid.NewString()
			newProject := entities.Project{
				ID:   ID,
				Name: projectName,
			}
			if err := repository.DB.AddProject(newProject); err != nil {
				return nil, isNew, err
			}

			if err := repository.DB.AddUserProject(userID, newProject.ID); err != nil {
				return nil, isNew, err
			}

			isNew = true

			project = &newProject
		} else {
			return nil, isNew, err
		}
	}
	return project, isNew, nil
}

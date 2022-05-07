package postgres

import (
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/errors"
	"github.com/google/uuid"
)

func (pg *Postgres) AddUser(user entities.User) error {
	query := `INSERT INTO users (id, email, password) VALUES ($1, $2, $3)`

	if _, err := pg.db.Exec(pg.ctx, query, user.ID, user.Email, user.Password); err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) UpdatePassword(userID string, newPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`

	if _, err := pg.db.Exec(pg.ctx, query, newPassword, userID); err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) GetUserByEmail(email string) (*entities.User, error) {
	query := `SELECT users.id, users.email, users.password FROM users WHERE email = $1`

	var user entities.User

	if err := pg.db.QueryRow(pg.ctx, query, email).Scan(&user.ID, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (pg *Postgres) AddUserProject(userID string, projectID string) error {
	id := uuid.NewString()

	query := `INSERT INTO userproject (id, userId, projectId) VALUES ($1, $2, $3)`

	if _, err := pg.db.Exec(pg.ctx, query, id, userID, projectID); err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) DeleteUserProject(userID string, projectID string) error {
	query := `DELETE FROM userproject WHERE userproject.userID = $1 AND userproject.projectID = $2`

	if _, err := pg.db.Exec(pg.ctx, query, userID, projectID); err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) GetUserProjectIDs(userID string) ([]string, error) {
	var projectIds []string

	query := `SELECT userproject.projectId FROM userproject WHERE userId = $1`

	rows, err := pg.db.Query(pg.ctx, query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		if len(id) > 0 {
			projectIds = append(projectIds, id)
		}
	}

	return projectIds, nil
}

func (pg *Postgres) GetUserProjectByName(userID string, projectName string) (*entities.Project, error) {
	query := `SELECT project.id, project.name FROM project
	INNER JOIN userproject
	ON userproject.projectId = project.id AND userproject.userid = $1
	WHERE project.name = $2`

	var project entities.Project

	if err := pg.db.QueryRow(pg.ctx, query, userID, projectName).Scan(&project.ID, &project.Name); err != nil {
		return nil, err
	}

	if project.ID == "" {
		return nil, errors.ProjectNotFound
	}

	return &project, nil
}

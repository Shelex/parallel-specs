package users

import "github.com/Shelex/parallel-specs/internal/entities"

func UserToEntityUser(user User) entities.User {
	return entities.User{
		Email:    user.Email,
		Password: user.Password,
		ID:       user.ID,
	}
}

func EntityUserToUser(user entities.User) User {
	return User{
		Email:    user.Email,
		Password: user.Password,
		ID:       user.ID,
	}
}

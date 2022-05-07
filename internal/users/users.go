package users

import (
	"github.com/Shelex/split-specs-v2/internal/entities"
	"github.com/Shelex/split-specs-v2/internal/errors"
	"github.com/Shelex/split-specs-v2/repository"
	"golang.org/x/crypto/bcrypt"
)

type User entities.User

func (user *User) Create() error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	if err := repository.DB.AddUser(UserToEntityUser(*user)); err != nil {
		return err
	}
	return nil
}

func (user *User) Authenticate() (*entities.User, error) {
	dbUser, err := repository.DB.GetUserByEmail(user.Email)
	if err != nil {
		return nil, errors.InvalidEmailOrPassord
	}

	authenticated := CheckPasswordHash(user.Password, dbUser.Password)
	if !authenticated {
		return nil, errors.InvalidEmailOrPassord
	}

	return dbUser, nil
}

func (user *User) Exist() bool {
	if _, err := repository.DB.GetUserByEmail(user.Email); err != nil {
		return false
	}
	return true
}

func (user *User) ChangePassword(password string, newPassword string) error {
	dbUser, err := repository.DB.GetUserByEmail(user.Email)
	if err != nil {
		return errors.AccessDenied
	}
	if match := CheckPasswordHash(password, dbUser.Password); !match {
		return errors.AccessDenied
	}
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	return repository.DB.UpdatePassword(user.ID, hashedPassword)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

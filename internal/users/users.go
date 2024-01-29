package users

import (
	"log"

	"github.com/Shelex/parallel-specs/internal/entities"
	"github.com/Shelex/parallel-specs/internal/errors"
	"github.com/Shelex/parallel-specs/repository"
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
	db, err := repository.DB.GetUserByEmail(user.Email)
	if err != nil {
		log.Println(err)
		return false
	}

	if db.Email == user.Email {
		return true
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

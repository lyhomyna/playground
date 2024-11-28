package controllers

import (
	"errors"
	"fmt"
	"log"
	"qqweq/siglog/model/database"
	"qqweq/siglog/model/models"

	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
    db database.SiglogDao
}

var userController *UserController
func NewUserController(db database.SiglogDao) *UserController {
    if userController == nil {
	userController = &UserController{ db }
    }
    return userController
}

func (c *UserController) GetUserByUsername(username string) (*models.User) {
    user, err := c.db.ReadUserByUsername(username)
    if err != nil {
	log.Println(err)
    }
    return user
}

func (c *UserController) AddUser(user *models.User) (string, error) {
    encryptedPassword, err := encryptPassword(user.Password)
    if err != nil {
	return "", errors.New(fmt.Sprintf("Failed to add new user. %s", err)) 
    }
    user.Password = encryptedPassword


    newUserId, err := c.db.CreateUser(user)
    if err != nil {
	return "", errors.New(fmt.Sprintf("Failed to add new user. %s", err))
    }

    return newUserId, nil
}

func (*UserController) ComparePasswords(user *models.User, possiblePassword string) error {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(possiblePassword))
    if err != nil {
	return errors.New(fmt.Sprintf("Password don't match. %s", err))
    }
    return nil
}

func (c *UserController) DeleteUser(username string) {
    user := c.GetUserByUsername(username)
    if err := c.db.DeleteUser(user); err != nil {
	log.Fatalf("Couldn't delete user. %s", err)
    }
}

func encryptPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
	return "", errors.New(fmt.Sprintf("Cant encrypt password. %s", err)) 
    }
    return string(bytes), nil
}

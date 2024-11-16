package controllers

import (
	"errors"
	"fmt"
	"qqweq/siglog/models"

	"golang.org/x/crypto/bcrypt"
)

var users = map[string]*models.User{
    "jamesbond": {
	Username: "jamesbond",
	Password: "hashedpassword",
	Firstname: "James",
	Lastname: "Bond",
	Role: "user",
    },
}

type UserController struct {}


var userController *UserController
func NewUserController() *UserController {
    if userController == nil {
	userController = &UserController{}
    }
    return userController
}

func (*UserController) GetUserByUsername(username string) (*models.User) {
    user := users[username]
    return user
}

func (*UserController) AddUser(user *models.User) error {
    encryptedPassword, err := encryptPassword(user.Password)
    if err != nil {
	return errors.New(fmt.Sprintf("Failed to add new user. %s", err)) 
    }
    user.Password = encryptedPassword
    users[user.Username] = user

    return nil
}

func (*UserController) ComparePasswords(user *models.User, possiblePassword string) error {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(possiblePassword))
    if err != nil {
	return errors.New(fmt.Sprintf("Password don't match. %s", err))
    }
    return nil
}

func encryptPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
	return "", errors.New(fmt.Sprintf("Cant encrypt password. %s", err)) 
    }
    return string(bytes), nil
}

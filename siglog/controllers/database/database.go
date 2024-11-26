package database 

import (
	"log"

	"qqweq/siglog/controllers/database/mongoDb"
	"qqweq/siglog/models"
)

type DatabaseController interface {
    CreateUser(user *models.User) (string, error)
    ReadUserByUsername(username string) (*models.User, error)
    // TODO: Delete user
}

var dbController DatabaseController

func NewDatabase() DatabaseController {
    if err := mongoDb.ConnectToMongoDb(); err != nil {
	log.Println(err)
	return nil
    }
    
    if dbController == nil {
	dbController = &mongoDb.MongoDbController{}
    }

    return dbController 
}

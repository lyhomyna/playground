package database 

import (
	"log"

	"qqweq/siglog/model/database/mongoDb"
	"qqweq/siglog/model/models"
)

type SiglogDao interface {
    CreateUser(user *models.User) (string, error)
    ReadUserByUsername(username string) (*models.User, error)
    // TODO: Delete user
}

var dbController SiglogDao
func NewDatabase() SiglogDao {
    if err := mongoDb.ConnectToMongoDb(); err != nil {
	log.Println(err)
	return nil
    }
    
    if dbController == nil {
	dbController = &mongoDb.MongoDbDao{}
    }

    return dbController 
}

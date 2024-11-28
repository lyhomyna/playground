package database 

import (
	"log"

	"qqweq/siglog/model/database/mongoDb"
	"qqweq/siglog/model/models"
)

type SiglogDao interface {
    // Users
    CreateUser(user *models.User) (string, error)
    ReadUserByUsername(username string) (*models.User, error)
    DeleteUser(user *models.User) error 

    // Sessions
    CreateSession(username string) (string, error)
    DeleteSession(sessionId string) error
    UsernameFromSessionId(sessionId string) (string, error)
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

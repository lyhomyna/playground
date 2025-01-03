package database

import (
	"context"
	"fmt"
	"qqweq/siglog/model/database/postgres"
	"qqweq/siglog/model/models"
)

type SiglogDao interface {
    // Users
    CreateUser(user *models.User) (string, error)
    DeleteUser(user *models.User) error 
    ReadUserByUsername(username string) (*models.User, error)

    // Sessions
    CreateSession(username string) (string, error)
    DeleteSession(sessionId string) error
    UsernameFromSessionId(sessionId string) (string, error)
}

var (
    ctx context.Context
    sd SiglogDao
)

func GetDao() (SiglogDao, error) {
   // if err := mongo.ConnectToMongoDb(ctx); err != nil {
   //     log.Println(err)
   //     return nil
   // }
   var err error
    if sd == nil {
	// dbController = &mongo.MongoDbDao{}
	sd, err = postgres.GetDao(ctx)
	if err != nil {
	    return nil, fmt.Errorf("Failure connecting to the PostgreSQL. %w", err)
	}
    }

    return sd, err
}

// Validate if Dao has implemented interface SiglogDao
//var _ SiglogDao = (*mongo.MongoDbDao)(nil)
var _ SiglogDao = (*postgres.PostgresDao)(nil)

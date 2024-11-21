package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"qqweq/siglog/models"
)

type DatabaseController struct {}

var dbController *DatabaseController

var mongoDbClient *mongo.Client
var mongoDb *mongo.Database
var ctx context.Context // TODO: what a heck context is?

func NewDatabaseController() *DatabaseController {
    if err := connectToDb(); err != nil {
	log.Println(err)
	return nil
    }
    
    if dbController == nil {
	dbController = &DatabaseController{}
    }

    return dbController 
}

func connectToDb() error {
    env := getEnv()

    var err error
    mongoDbClient, err = mongo.Connect(ctx, options.Client().ApplyURI(env.ConnString))
    if err != nil {
	return errors.New(fmt.Sprintf("Can't get mongo client. %s", err))
    }
    defer mongoDisconnect()

    mongoDb = mongoDbClient.Database("users")
    
    return nil
}

// TODO: refactor
func getEnv() *models.EnvModel {
    f, err := os.Open("./variables.env")
    if err != nil { log.Fatalf("Can't open env file. %s", err) }
    f.Close()
    
    var env models.EnvModel
    if err := json.NewDecoder(f).Decode(&env); err != nil {
	log.Fatalf("Can't decode env object. %s", err)
    }
    return &env
}

func mongoDisconnect() {
    if err := mongoDbClient.Disconnect(ctx); err != nil {
	log.Fatalf("Can't disconnect from mongo db. %s", err)
    }
}

package mongoDb

import (
	"context"
	"errors"
	"fmt"
	"qqweq/siglog/models"
	"qqweq/siglog/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbController struct {}

var databaseName = "siglog"
var mongoDbClient *mongo.Client
var mongoDb *mongo.Database
// TODO: what a heck context.Context is?; add context
var usersColl *mongo.Collection

var collections map[string]*mongo.Collection

func getColl(collName string) (*mongo.Collection, error) {
    if _, ok := collections[collName]; !ok {
	coll := mongoDb.Collection(collName)
	if coll == nil {
	    return nil, errors.New(fmt.Sprintf("WTF is '%s' collection?", collName))
	}
	
	collections[collName] = coll
    }
    return collections[collName], nil
}


// CreateUser returns userID if User creaion was successfull, or error if failure.
func (*MongoDbController) CreateUser(user *models.User) (string, error) {
    if user == nil {
	return "", errors.New("I can't create new user cuz there is no user, bastard.") 
    }

    coll, err := getColl("users")
    if err != nil {
	return "", errors.New(fmt.Sprintf("I can't create new user. %s", err))
    }

    res, err := coll.InsertOne(context.TODO(), user)
    if err != nil {
	return "", errors.New(fmt.Sprintf("I can't create new user. %s", err))
    }
    
    insertedId := fmt.Sprintf("%v", res.InsertedID)

    return insertedId, nil
}

func (*MongoDbController) ReadUserByUsername(username string) (*models.User, error) {
    coll, err := getColl("users")
    if err != nil {
	return nil, errors.New(fmt.Sprintf("I can't read user. %s", err))
    }
    
    var user *models.User
    err = coll.FindOne(context.TODO(), bson.M{"username": username}).Decode(user)
    if err != nil {
	if err == mongo.ErrNoDocuments {
	    return nil, errors.New("No user found with given username.")
	} else {
	    return nil, errors.New(fmt.Sprintf("Can't read user by username. WTF error. %s", err))
	}
    }

    return user, nil
}

func ConnectToMongoDb() error {
    env := utils.GetEnv()

    // wtf, this thing doesn't care what to connect to
    mongoDbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(env.ConnString))
    if err != nil {
	return errors.New(fmt.Sprintf("Can't get mongo client. %s", err))
    }

    mongoDb = mongoDbClient.Database(databaseName)
    
    return nil
}

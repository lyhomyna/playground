package mongoDb

import (
	"context"
	"errors"
	"fmt"
	"qqweq/siglog/model/models"
	"qqweq/siglog/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbDao struct {}

var databaseName = "siglog"
var mongoDbClient *mongo.Client
var mongoDb *mongo.Database
// TODO: what a heck context.Context is?; add context
var usersColl *mongo.Collection

var collections = map[string]*mongo.Collection{}

func getColl(collName string) (*mongo.Collection) {
    if _, ok := collections[collName]; !ok {
	coll := mongoDb.Collection(collName)

	// and if collName collection don't exist?
	// but I don't know how to handle this, cuz Collection()
	// never returns nil
	collections[collName] = coll
    }
    return collections[collName]
}

func ConnectToMongoDb() error {
    env := utils.GetEnv()

    // wtf, this thing doesn't care what to connect to
    mongoDbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(env.ConnString))
    if err != nil {
	return fmt.Errorf("Couldn't get mongo client. %s", err)
    }

    if err := mongoDbClient.Ping(context.TODO(), nil); err != nil {
	return fmt.Errorf("Couldn't get access to db. %w", err)
    }

    mongoDb = mongoDbClient.Database(databaseName)
    
    return nil
}


// CreateUser returns userID if User creaion was successfull, or error if failure.
func (*MongoDbDao) CreateUser(user *models.User) (string, error) {
    if user == nil {
	return "", errors.New("I can't create new user cuz there is no user.") 
    }

    coll := getColl("users")

    res, err := coll.InsertOne(context.TODO(), user)
    if err != nil {
	return "", fmt.Errorf("I can't create new user. %w", err)
    }
    
    insertedId, _:= res.InsertedID.(primitive.ObjectID)

    return insertedId.Hex(), nil
}

func (*MongoDbDao) ReadUserByUsername(username string) (*models.User, error) {
    coll := getColl("users")
    
    var user *models.User
    filter := bson.M{"username": username}

    err := coll.FindOne(context.TODO(), filter).Decode(&user)
    if err != nil {
	if err == mongo.ErrNoDocuments {
	    return nil, errors.New("No user found with given username.")
	} else {
	    return nil, fmt.Errorf("Couldn't get user by username. WTF error. %w", err)
	}
    }

    return user, nil
}

func (*MongoDbDao) DeleteUser(user *models.User) error {
    if user == nil {
	return errors.New("No user provided to delete.")
    }

    coll := getColl("users")

    res, err := coll.DeleteOne(context.TODO(), user)
    if res.DeletedCount != 1 {
	resErr := errors.New("User not deleted.")
    
	if err != nil {
	    resErr = fmt.Errorf("%w. %w", resErr, err)
	}

	return resErr 
    }

    return nil
}

func (*MongoDbDao) CreateSession(username string) (string, error) {
    sessionId := uuid.NewString() 
    
    coll := getColl("sessions")

    filter := bson.M{ 
	"session-id":  sessionId, 
	"username": username,
    }
    _, err := coll.InsertOne(context.TODO(), filter)
    if err != nil {
	return "", fmt.Errorf("Couldn't create new session. %w", err)
    }

    return sessionId, nil
}

func (*MongoDbDao) DeleteSession(sessionId string) error {
    coll := getColl("sessions")
    
    filter := bson.M { "session-id": sessionId }
    _, err := coll.DeleteOne(context.TODO(), filter)
    if err != nil {
	return fmt.Errorf("Couldn't delete session. %w", err)
    }
    
    return nil
}

func (*MongoDbDao) UsernameFromSessionId(sessionId string) (string, error) {
    coll := getColl("sessions")   

    var findResObj struct {
	Username string `bson:"username"`
    }
    filter := bson.M{ "session-id": sessionId }
    opts := options.FindOne().SetProjection(bson.M{"username": 1, "_id": 0})

    err := coll.FindOne(context.TODO(), filter, opts).Decode(&findResObj)
    if err != nil {
	if err == mongo.ErrNoDocuments {
	    return "", fmt.Errorf("There is no username associated with session '%s'. %w", sessionId, err)
	}
	return "", errors.New(fmt.Sprintf("Error fetching username for session %s. %s", sessionId, err))
    }

    return findResObj.Username, nil
}

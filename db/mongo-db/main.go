package main

import (
	"context"
	"encoding/json"
	"log"
	"mongodbtest/models"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mdClient *mongo.Client
var ctx context.Context

func main() {
    env := readEnv()
    ctx = context.Background()

    var err error
    mdClient, err = mongo.Connect(ctx, options.Client().ApplyURI(env.ConnString))
    if err != nil {
	log.Fatalf("Can't connect to db. %s", err)
    }
    defer func() {
	if err := mdClient.Disconnect(ctx); err != nil {
	    log.Panic(err)
	}
    } ()

    db := mdClient.Database("testdb")
    coll := db.Collection("testcollection")
    _, err = coll.InsertOne(ctx, bson.D{
	{"name", "superuser"},
    })
    if err != nil {
	log.Printf("InsertOne failure. %s", err)
    }
}

func readEnv() *models.EnvVars {
    envFile, err := os.Open("./variables.env")
    if err != nil {
	log.Panicf("Can't open file. %s", err)
    }
    defer envFile.Close()

    var envVars models.EnvVars
    err = json.NewDecoder(envFile).Decode(&envVars)
    if err != nil {
	log.Panicf("Can't decode. %s", err)
    }

    return &envVars
}

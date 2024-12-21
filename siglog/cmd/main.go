package main

import (
	"context"
	"log"
	"path/filepath"
	"qqweq/siglog/api"

	"github.com/joho/godotenv"
)

// For loading env variables
var envFilePath = filepath.Join("..", ".env")

func init() {
    if err := godotenv.Load(envFilePath); err != nil {
	panic(err)
    }
}

func main() {
    s := api.SiglogServer{}

    // Channel for consuming errors
    ec := make(chan error, 1)
    go func() {
	ec <- s.Run(context.Background())
    }()
    
    err := <- ec
    if err != nil {
	log.Printf("Server terminated by error: %s", err)
    }
}


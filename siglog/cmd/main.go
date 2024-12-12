package main

import (
	"context"
	"log"
	"path/filepath"
	"qqweq/siglog/api"

	"github.com/joho/godotenv"
)

var envFilePath = filepath.Join("..", ".env")

func init() {
    if err := godotenv.Load(envFilePath); err != nil {
	panic(err)
    }
}

func main() {
    s := api.Server{}
    
    // error channel
    ec := make(chan error, 1)
    go func() {
	ec <- s.Run(context.Background())
    }()
    
    err := <- ec
    if err != nil {
	log.Printf("Server terminated by error: %s", err)
    }

    log.Println("end of main.")
}


package main

import (
	"context"
	"log"
	"qqweq/siglog/api"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
    s := api.Server{}
    
    ec := make(chan error, 1)
    go func() {
	ec <- s.Run(context.Background())
    }()
    
    err := <- ec
    if err != nil {
	log.Printf("Server terminated by error: %s", err)
    }
}


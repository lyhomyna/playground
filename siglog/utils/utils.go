package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// TODO: refactor
func GetEnv() *EnvModel {
    wd, err := os.Getwd()
    if err != nil {
	log.Panicf("Hello, genious! Cant get working directory. %s", err)
    }
	
    envFilePath := filepath.Join(wd, "variables.env")

    f, err := os.Open(envFilePath)
    if err != nil { log.Fatalf("Can't open env file. %s", err) }
    defer f.Close()
    
    var env EnvModel
    if err := json.NewDecoder(f).Decode(&env); err != nil {
	log.Fatalf("Can't decode env object. %s", err)
    }
    return &env
}


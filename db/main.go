package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

type Config struct {
    Db struct {
	User 	 string `yaml:"user"`	
	Password string `yaml:"password"`
	Host 	 string `yaml:"host"`
	Name	 string `yaml:"name"`
	Port	 int 	`yaml:"port"`
    }
}

func main() {
    config, err := loadConfig()
    if err != nil {
	log.Fatal(err)
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Db.User, config.Db.Password, config.Db.Host, config.Db.Port, config.Db.Name)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
	log.Fatal("Failed to connect to database.")
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
	log.Fatal("Failed to ping database.")
    }

    log.Println("Successfully conected to db.")
}

func loadConfig() (*Config, error) {
    cfg, err := os.ReadFile("./config.yml")
    if err != nil {
	log.Println(err)
	return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(cfg, &config); err != nil {
	log.Println(err)
	return nil, err
    }

    return &config, nil
}

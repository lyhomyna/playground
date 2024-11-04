package main

import (
	"database/sql"
	"errors"
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
    log.SetFlags(log.LstdFlags)

    db, err := getDb()
    if err != nil {
	log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Query("CREATE TABLE IF NOT EXISTS users(id INT NOT NULL AUTO_INCREMENT, name VARCHAR(20), PRIMARY KEY (id));")
    if err != nil {
        log.Fatalf("Failed CREATE TABLE. %+v", err)
    }

    _, err = db.Query("INSERT INTO users (name) VALUES ('Sam'), ('Mike'), ('Marek');")
    if err != nil {
        log.Fatalf("Failed INSERT values. %+v", err)
    }

    rows, err := db.Query("SELECT * FROM users;")
    if err != nil {
        log.Fatalf("Failed SELECT. %+v", err)
    }

    for rows.Next() {
	var id int
	var name string

	err = rows.Scan(&id, &name)
	if err != nil {
	    log.Fatal(err)
	}

	log.Printf("Id: %d\tName: %s", id, name)
    }

    _, err = db.Query("TRUNCATE TABLE users")
    if err != nil {
        log.Fatalf("Failed TRUNCATE TABLE. %+v", err)
    }
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

func getDb() (*sql.DB, error) {
    config, err := loadConfig()
    if err != nil {
	return nil, errors.New("Failed to load config.") 
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Db.User, config.Db.Password, config.Db.Host, config.Db.Port, config.Db.Name)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
	return nil, errors.New("Failed to connect to db.") 
    }

    if err := db.Ping(); err != nil {
	return nil, errors.New("Failed to ping db.")
    }

    log.Println("Successfully conected to db.")
    return db, nil
}

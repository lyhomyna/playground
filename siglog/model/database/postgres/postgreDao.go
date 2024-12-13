package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"qqweq/siglog/model/models"
)

type PostgresDao struct {
    db *pgx.Conn
    ctx context.Context
}

var (
    dao *PostgresDao

    // To create postgre DAO only once even if multiple goroutines
    // will call GetDao function. 
    daoOnce sync.Once
)

func GetDao(ctx context.Context) (*PostgresDao, error) {
    var err error
    if dao == nil {
	if ctx == nil { ctx = context.Background() }
	daoOnce.Do(func() {
	    var conn *pgx.Conn
	    conn, err = connectToDb(ctx)
	    if err == nil {
		log.Println("Postgres DAO created.")
		dao = &PostgresDao{ db: conn, ctx: ctx }
	    }
	})
    }
    
    return dao, err
}

// Associated with USER functions 

func (pd *PostgresDao) CreateUser(user *models.User) (string, error) {
    _, err := pd.db.Exec(pd.ctx, "INSERT INTO users (username, password, firstname, lastname, role) VALUES ($1, $2, $3, $4, $5);", user.Username, user.Password, user.Firstname, user.Lastname, user.Role)
    if err != nil {
	return "", err 
    }
 
    row := pd.db.QueryRow(pd.ctx, "SELECT id FROM users WHERE username=$1;", user.Username)
    var id string
    err = row.Scan(&id)
    if err != nil {
	return "", errors.New("WTF, user wasn't created.")
    }

    return id, nil
}
func (pd *PostgresDao) ReadUserByUsername(username string) (*models.User, error) {
    row := pd.db.QueryRow(pd.ctx, "SELECT * FROM users WHERE username=$1;", username)
    var user *models.User
    err := row.Scan(user)
    if err != nil {
	return nil, fmt.Errorf("Cannot read user by username. %w", err)
    }
    
    return user, nil
}
func (pd *PostgresDao) DeleteUser(user *models.User) error {
    _, err := pd.db.Exec(pd.ctx, "DELETE FROM users WHERE username=$1;", user.Username)
    if err != nil {
	err = fmt.Errorf("Cannot delete user %s. %w", user.Username, err)
    }
    return err
} 

// Associated with SESSIONS functions

func (pd *PostgresDao) CreateSession(username string) (string, error) {
    sessionId := uuid.NewString()
    _, err := pd.db.Exec(pd.ctx, "INSERT INTO sessions (sessionId, username) VALUES ($1, $2);", sessionId, username)
    if err != nil {
	return "", fmt.Errorf("Couldn't create session. %w", err)
    }
    return sessionId, nil
}
func (pd *PostgresDao) DeleteSession(sessionId string) error {
    execRes, err := pd.db.Exec(pd.ctx, "DELETE FROM sessions WHERE sessionId=$1;")

    switch {
	case err != nil:
	    break;
	    case execRes.RowsAffected() != 1: // FIXME: magic number
	    err = errors.New("Nothing was deleted from sessions table.")
    }

    return err
}
func (pd *PostgresDao) UsernameFromSessionId(sessionId string) (string, error) {
    row := pd.db.QueryRow(pd.ctx, "SELECT * FROM sessions WHERE sessionId=$1", sessionId)

    var record struct {
	sId string
	username string
    }
    err := row.Scan(&record)
    if err != nil {
	return "", fmt.Errorf("Couldnt get username from session %s. %w", sessionId, err)
    }

    return record.username, nil
}

func connectToDb(ctx context.Context) (*pgx.Conn, error)  { 
    connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

    log.Println(connString)

    conn, err := pgx.Connect(ctx, connString)
    if err != nil {
	err = fmt.Errorf("Failed to connect to the db. %w", err)
	return nil, err 
    }

    if err := conn.Ping(ctx); err != nil {
	err = fmt.Errorf("Failed to connect to the db. %w", err)
	return nil, err 
    }

    return conn, nil
}
